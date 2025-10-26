package api

import (
	"cmd/obyavigo/main.go/internal/secure"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type PasswordChangeRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type AccountDeleteRequest struct {
	Password string `json:"password"`
}

func (h *Handlers) ChangePasswordHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			var req PasswordChangeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusBadRequest, "invalid request body") {
				return
			}

			userData, err := h.db.Psql.GetUserData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting user data") {
				return
			}

			authInfo, err := h.db.Psql.GetAuthInfoByEmail(userData.Email)
			if handleError(w, err, http.StatusInternalServerError, "error while getting auth info") {
				return
			}

			if !secure.CheckPasswordHash(req.OldPassword, authInfo.PasswordHash) {
				sendToClient(w, http.StatusUnauthorized, "incorrect old password")
				return
			}

			newPasswordHash, err := secure.HashPassword(req.NewPassword)
			if handleError(w, err, http.StatusInternalServerError, "error while hashing new password") {
				return
			}

			err = h.db.Psql.UpdatePassword(userID, newPasswordHash)
			if handleError(w, err, http.StatusInternalServerError, "error while updating password") {
				return
			}

			sendToClient(w, http.StatusOK, "password updated successfully")
		},
	)
}

func (h *Handlers) DeleteAccountHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				http.Redirect(w, r, "/auth/", http.StatusPermanentRedirect)
				return
			}

			var req AccountDeleteRequest
			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusBadRequest, "invalid request body") {
				return
			}

			userData, err := h.db.Psql.GetUserData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting user data") {
				return
			}

			authInfo, err := h.db.Psql.GetAuthInfoByEmail(userData.Email)
			if handleError(w, err, http.StatusInternalServerError, "error while getting auth info") {
				return
			}

			if !secure.CheckPasswordHash(req.Password, authInfo.PasswordHash) {
				sendToClient(w, http.StatusUnauthorized, "incorrect password")
				return
			}

			err = h.db.Psql.DeleteAccount(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while deleting account") {
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "auth_token",
				Value:    "",
				Path:     "/",
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   true,
			})

			slog.Info("account deleted", slog.String("user_id", userID.String()))
			sendToClient(w, http.StatusOK, "account deleted successfully")
		},
	)
}
