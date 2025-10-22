package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"cmd/obyavigo/main.go/internal/secure"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

func (h *Handlers) RegistrationHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user models.User
			if err := json.NewDecoder(r.Body).Decode(&user); handleError(w, err, http.StatusBadRequest, "invalid request body") {
				return
			}

			exists, err := h.db.Psql.CheckEmail(user.Email)
			if handleError(w, err, http.StatusInternalServerError, "checking email error") {
				return
			}
			if exists {
				sendToClient(w, http.StatusConflict, "Email is busy")
				return
			}

			user.PasswordHash, err = secure.HashPassword(user.Password)
			if handleError(w, err, http.StatusInternalServerError, "password hashing error") {
				return
			}

			settings, _ := json.Marshal("{}")

			user.RegistrationDate = time.Now()
			user.Role = "user"
			user.Status = "online"
			user.Settings = settings

			id, err := h.db.Psql.CreateNewUser(&user)
			if handleError(w, err, http.StatusInternalServerError, "creating new user error") {
				return
			}

			token, err := h.mail.SendRegConfirm(user.Email)
			if handleError(w, err, http.StatusInternalServerError, "sending registration confirm error") {
				return
			}

			if err := h.db.Psql.CreateEmailConfirmation(id, token); handleError(w, err, http.StatusInternalServerError, "create email confirmation error") {
				return
			}

			sendToClient(w, http.StatusCreated, "User registered successfully")
		},
	)
}

func (h *Handlers) ConfirmRegistrationHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.PathValue("token")

			em, err := h.mail.Cipher.DecryptMail(token)
			if handleError(w, err, http.StatusInternalServerError, "decrypt email error") {
				return
			}

			ok, err := h.db.Psql.IsConfirmationExpired(token)
			if handleError(w, err, http.StatusInternalServerError, "try confirm account request error") {
				return
			}

			if !ok {
				err := h.db.Psql.DeleteAccountByEmail(em)
				if handleError(w, err, http.StatusInternalServerError, "delete account request error") {
					return
				}
				sendToClient(w, http.StatusGone, "Account verification time has expired")
				return
			}

			err = h.db.Psql.ConfirmAccount(token)
			if errors.Is(err, sql.ErrNoRows) {
				h.sendNotFound(w)
				return
			}
			if handleError(w, err, http.StatusInternalServerError, "confirm account request error") {
				return
			}

			http.Redirect(w, r, "/auth", http.StatusOK)
		})
}
