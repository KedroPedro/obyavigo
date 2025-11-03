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

type ProfileUpdateRequest struct {
	Username    string  `json:"username"`
	PhoneNumber *string `json:"phone_number"`
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
				sendToClient(w, http.StatusUnauthorized, "Неверный старый пароль")
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

			sendToClient(w, http.StatusOK, "Пароль успешно обновлён")
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

			
			if err := h.db.Mongo.DeleteUserImages(r.Context(), userID.String()); err != nil {
				slog.Error("error deleting user images", slog.String("user_id", userID.String()), slog.String("error", err.Error()))
				
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
			sendToClient(w, http.StatusOK, "Аккаунт удалён")
		},
	)
}

func (h *Handlers) UpdateProfileHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			var req ProfileUpdateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusBadRequest, "invalid request body") {
				return
			}

			if req.Username == "" {
				sendToClient(w, http.StatusBadRequest, "Имя не может быть пустым")
				return
			}

			err = h.db.Psql.UpdateUserProfile(userID, req.Username, req.PhoneNumber)
			if handleError(w, err, http.StatusInternalServerError, "error while updating profile") {
				return
			}

			slog.Info("profile updated", slog.String("user_id", userID.String()))
			sendToClient(w, http.StatusOK, "Профиль успешно обновлён")
		},
	)
}

func (h *Handlers) GetUserProfileAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userData, err := h.db.Psql.GetUserData(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user data") {
				return
			}

			userRole, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

		avatarID, _ := h.db.Psql.GetUserAvatar(userID)

		response := map[string]interface{}{
			"username":            userData.Username,
			"email":               userData.Email,
			"phone_number":        userData.PhoneNumber,
			"role":                userRole,
			"profile_picture_id": avatarID,
		}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) UploadAvatarHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			if err := r.ParseMultipartForm(5 << 20); err != nil {
				sendToClient(w, http.StatusBadRequest, "Файл слишком большой")
				return
			}

			file, _, err := r.FormFile("avatar")
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "Не удалось получить файл")
				return
			}
			defer file.Close()

			fileHeader := r.MultipartForm.File["avatar"][0]

			
			oldAvatarID, _ := h.db.Psql.GetUserAvatar(userID)
			if oldAvatarID != nil {
				_ = h.db.Mongo.DeleteUserAvatar(r.Context(), userID.String())
			}

			
			avatarID, err := h.db.Mongo.UploadUserAvatar(r.Context(), fileHeader, userID.String())
			if handleError(w, err, http.StatusInternalServerError, "error uploading avatar") {
				return
			}

			
			if err := h.db.Psql.UpdateUserAvatar(userID, avatarID); err != nil {
				sendToClient(w, http.StatusInternalServerError, "Ошибка при обновлении аватара")
				return
			}

		slog.Info("avatar uploaded", slog.String("user_id", userID.String()), slog.String("avatar_id", avatarID))
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message":  "Аватар успешно загружен",
			"image_id": avatarID,
		})
		},
	)
}
