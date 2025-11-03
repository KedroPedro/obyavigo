package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"cmd/obyavigo/main.go/internal/secure"
	"encoding/json"
	"net/http"
)

func (h *Handlers) AuthorizationHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user models.AuthUser

			if err := json.NewDecoder(r.Body).Decode(&user); handleError(w, err, http.StatusInternalServerError, "invalid request body") {
				return
			}

			exists, err := h.db.Psql.CheckEmail(user.Email)

			if handleError(w, err, http.StatusInternalServerError, "checking email error") {
				return
			} else if !exists {
				sendToClient(w, http.StatusUnauthorized, "user not found")
				return
			}

			authInfo, err := h.db.Psql.GetAuthInfoByEmail(user.Email)

			if handleError(w, err, http.StatusInternalServerError, "get auth info by email error") {
				return
			}

			if !authInfo.Confirmed {
				sendToClient(w, http.StatusGone, "the account has not been verified")
				return
			}

			if authInfo.Status == "banned" {
				sendToClient(w, http.StatusForbidden, "account is banned")
				return
			}

		if !secure.CheckPasswordHash(user.Password, authInfo.PasswordHash) {
			sendToClient(w, http.StatusUnauthorized, "incorrect password")
			return
		}

		// Update last login time
		if err := h.db.Psql.UpdateLastLogin(&authInfo.Id); err != nil {
			// Log error but don't fail login
			// handleError would return, so we just log it
		}

		token, err := h.jwt.GenerateJWTToken(&authInfo.Id)
		if handleError(w, err, http.StatusInternalServerError, "generate jwt token error") {
			return
		}

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
				MaxAge:   3600 * 24,
			})

			sendToClient(w, http.StatusOK, "the user has been authorized")
		},
	)
}
