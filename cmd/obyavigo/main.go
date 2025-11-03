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
	mux.Handle("POST /api/profile/change-password/", handler.AuthMiddleware(handler.ChangePasswordHandler()))
	mux.Handle("POST /api/profile/update/", handler.AuthMiddleware(handler.UpdateProfileHandler()))
	mux.Handle("POST /api/profile/upload-avatar/", handler.AuthMiddleware(handler.UploadAvatarHandler()))
	mux.Handle("POST /api/profile/delete-account/", handler.AuthMiddleware(handler.DeleteAccountHandler()))
	mux.Handle("GET /api/auth/confirm-email/{token}/", handler.ConfirmRegistrationHandler())
	mux.Handle("GET /api/ads/", handler.AuthMiddleware(handler.GetAdsAPI()))
	mux.Handle("GET /api/ads/{token}/", handler.AuthMiddleware(handler.GetAdByIDAPI()))
	mux.Handle("GET /api/user/ads/", handler.AuthMiddleware(handler.GetUserAdsAPI()))
	mux.Handle("DELETE /api/ads/{token}/", handler.AuthMiddleware(handler.DeleteAdAPI()))
	mux.Handle("PUT /api/ads/{token}/", handler.AuthMiddleware(handler.UpdateAdAPI()))
	mux.Handle("DELETE /api/ads/{token}/images/{imageId}/", handler.AuthMiddleware(handler.DeleteAdImageAPI()))
	mux.Handle("POST /api/ads/{token}/images/", handler.AuthMiddleware(handler.UploadAdImagesAPI()))
	mux.Handle("GET /api/user/profile/", handler.AuthMiddleware(handler.GetUserProfileAPI()))
	mux.Handle("POST /api/favorites/{token}/", handler.AuthMiddleware(handler.AddToFavoritesAPI()))
	mux.Handle("DELETE /api/favorites/{token}/", handler.AuthMiddleware(handler.RemoveFromFavoritesAPI()))
	mux.Handle("GET /api/favorites/check/{token}/", handler.AuthMiddleware(handler.CheckIfFavoriteAPI()))
	mux.Handle("GET /api/images/{id}/", handler.GetImageByID())
	mux.Handle("GET /api/avatars/{id}/", handler.GetAvatarByID())
	mux.Handle("GET /api/admin/stats/", handler.AuthMiddleware(handler.GetAdminStats()))
	mux.Handle("GET /api/admin/ads/", handler.AuthMiddleware(handler.GetAdminAds()))
	mux.Handle("POST /api/admin/ads/{id}/status/", handler.AuthMiddleware(handler.UpdateAdStatus()))
	mux.Handle("GET /api/admin/users/", handler.AuthMiddleware(handler.GetAdminUsers()))
	mux.Handle("POST /api/admin/users/{id}/status/", handler.AuthMiddleware(handler.UpdateUserStatus()))
	mux.Handle("POST /api/admin/users/{id}/role/", handler.AuthMiddleware(handler.UpdateUserRole()))
	mux.Handle("POST /api/reports/", handler.AuthMiddleware(handler.CreateReport()))
	mux.Handle("GET /api/admin/reports/", handler.AuthMiddleware(handler.GetReports()))
	mux.Handle("POST /api/admin/reports/{id}/status/", handler.AuthMiddleware(handler.UpdateReportStatus()))
	mux.Handle("GET /", handler.AuthMiddleware(handler.GetMainPage()))
	mux.Handle("GET /profile/", handler.AuthMiddleware(handler.GetProfilePage()))
	mux.Handle("GET /messages/", handler.AuthMiddleware(handler.GetMessagesPage()))
	mux.Handle("GET /create-ad/", handler.AuthMiddleware(handler.GetCreateAdPage()))
	mux.Handle("GET /auth/", handler.AuthMiddleware(handler.GetAuthPage()))
	mux.Handle("GET /ads/", handler.AuthMiddleware(handler.GetAdsPage()))
	mux.Handle("GET /ads/{token}/", handler.AuthMiddleware(handler.GetAdPage()))
	mux.Handle("GET /edit-ad/{token}/", handler.AuthMiddleware(handler.GetEditAdPage()))
	mux.Handle("GET /admin-panel/", handler.AuthMiddleware(handler.GetAdminPanelPage()))
	mux.Handle("GET /liked-ads/", handler.AuthMiddleware(handler.GetLikedAdsPage()))
	mux.Handle("GET /static/", handler.NoDirListing(http.StripPrefix("/static/", fs)))

	go func() {
		err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			target := "https://" + cfg.Server.Addr + r.URL.RequestURI()
			http.Redirect(w, r, target, http.StatusMovedPermanently)
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

	slog.Info("server is started")
	if err := srv.ListenAndServeTLS("./localhost.crt", "./localhost.key"); err != nil {
		slog.Error("error shutdown the server", slog.String("error", err.Error()))
	}
	slog.Info("server is stopped")
}
