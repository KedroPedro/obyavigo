package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func (h *Handlers) NoDirListing(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/") {
				h.sendNotFound(w)
				return
			}
			next.ServeHTTP(w, r)
		},
	)
}

func (h *Handlers) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("auth_token")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			token := cookie.Value

			id, err := h.jwt.ParseJWTToken(token)
			if err != nil {
				slog.Error("error while parsing token", slog.String("error", err.Error()))
				next.ServeHTTP(w, r)
				return
			}

			
			status, err := h.db.Psql.GetUserStatus(id)
			if err == nil && status == "banned" {
				http.SetCookie(w, &http.Cookie{
					Name:     "auth_token",
					Value:    "",
					Path:     "/",
					Expires:  time.Unix(0, 0),
					MaxAge:   -1,
					HttpOnly: true,
					Secure:   true,
				})
				next.ServeHTTP(w, r)
				return
			}

			ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
			defer cancel()
			ctx = context.WithValue(ctx, "userID", id.String())

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		},
	)
}
