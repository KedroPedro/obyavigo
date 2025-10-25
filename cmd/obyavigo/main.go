package main

import (
	"cmd/obyavigo/main.go/internal/api"
	"cmd/obyavigo/main.go/internal/config"
	"cmd/obyavigo/main.go/internal/database"
	"cmd/obyavigo/main.go/internal/logger"
	"cmd/obyavigo/main.go/internal/mail"
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	cfg := config.New()
	log, h := logger.Init(cfg.Logger, 5000)
	defer h.Close()
	slog.SetDefault(log)
	fmt.Println(cfg)
	db, err := database.Init(cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err := db.Psql.Ping(); err != nil {
		slog.Error(err.Error())
		return
	}
	if err := db.Mongo.Ping(); err != nil {
		slog.Error(err.Error())
		return
	}
	fs := http.FileServer(http.Dir(cfg.Environment.StaticDir))
	mail, err := mail.New(cfg)
	if err != nil {
		slog.Error("email initialization error", slog.String("error", err.Error()))
	}
	handler := api.New(db, cfg, mail)
	mux := http.NewServeMux()

	mux.Handle("POST /api/auth/register/", handler.RegistrationHandler())
	mux.Handle("POST /api/auth/login/", handler.AuthMiddleware(handler.AuthorizationHandler()))
	mux.Handle("POST /api/create-ad/", handler.AuthMiddleware(handler.CreateAd()))
	mux.Handle("POST /api/logout/", handler.UserLogout())
	mux.Handle("GET /api/auth/confirm-email/{token}/", handler.ConfirmRegistrationHandler())
	mux.Handle("GET /", handler.AuthMiddleware(handler.GetMainPage()))
	mux.Handle("GET /profile/", handler.AuthMiddleware(handler.GetProfilePage()))
	mux.Handle("GET /messages/", handler.AuthMiddleware(handler.GetMessagesPage()))
	mux.Handle("GET /create-ad/", handler.AuthMiddleware(handler.GetCreateAdPage()))
	mux.Handle("GET /auth/", handler.AuthMiddleware(handler.GetAuthPage()))
	mux.Handle("GET /ads/", handler.AuthMiddleware(handler.GetAdsPage()))
	mux.Handle("GET /ads/{token}/", handler.AuthMiddleware(handler.GetAdPage()))
	mux.Handle("GET /static/", handler.NoDirListing(http.StripPrefix("/static/", fs)))

	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + cfg.Server.Addr + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently) // 301 редирект
		}))
		if err != nil {
			slog.Error("HTTP redirect server failed", slog.String("error", err.Error()))
		}
	}()

	srv := http.Server{
		Handler:      mux,
		Addr:         cfg.Server.Addr,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}

	log.Info("server is started")
	if err := srv.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		log.Error("error shutdown the server")
	}
	log.Info("server is stopped")
}
