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

		
		if err := h.db.Psql.UpdateLastLogin(&authInfo.Id); err != nil {
			
			
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

func (h *Handlers) ForgotPasswordHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req models.PasswordResetRequest

			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusInternalServerError, "invalid request body") {
				return
			}

			
			exists, err := h.db.Psql.CheckEmail(req.Email)
			if handleError(w, err, http.StatusInternalServerError, "checking email error") {
				return
			}

			
			if !exists {
				sendToClient(w, http.StatusOK, "if email exists, reset link will be sent")
				return
			}

			
			token, err := secure.GenerateRandomToken(32)
			if handleError(w, err, http.StatusInternalServerError, "token generation error") {
				return
			}

			
			if err := h.db.Psql.CreatePasswordResetToken(req.Email, token); handleError(w, err, http.StatusInternalServerError, "error creating reset token") {
				return
			}

			
			if err := h.mail.SendPasswordReset(req.Email, token); handleError(w, err, http.StatusInternalServerError, "error sending email") {
				return
			}

			sendToClient(w, http.StatusOK, "reset link sent to email")
		},
	)
}

func (h *Handlers) ResetPasswordHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var req models.ResetPasswordRequest

			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusInternalServerError, "invalid request body") {
				return
			}

			
			resetToken, err := h.db.Psql.GetPasswordResetToken(req.Token)
			if handleError(w, err, http.StatusInternalServerError, "error validating token") {
				return
			}

			if resetToken == nil {
				sendToClient(w, http.StatusBadRequest, "invalid or expired token")
				return
			}

			
			hashedPassword, err := secure.HashPassword(req.NewPassword)
			if handleError(w, err, http.StatusInternalServerError, "error hashing password") {
				return
			}

			
			if err := h.db.Psql.UpdatePassword(&resetToken.UserID, hashedPassword); handleError(w, err, http.StatusInternalServerError, "error updating password") {
				return
			}

			
			if err := h.db.Psql.DeletePasswordResetToken(req.Token); err != nil {
				
			}

			sendToClient(w, http.StatusOK, "password successfully reset")
		},
	)
}
