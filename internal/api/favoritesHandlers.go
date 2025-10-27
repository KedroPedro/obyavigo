package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func (h *Handlers) AddToFavoritesAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil || userID == nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			err = h.db.Psql.AddToFavorites(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error while adding to favorites") {
				return
			}

			response := map[string]interface{}{
				"success": true,
				"message": "added to favorites",
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) RemoveFromFavoritesAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil || userID == nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			err = h.db.Psql.RemoveFromFavorites(userID, &adID)
			if err == sql.ErrNoRows {
				sendToClient(w, http.StatusNotFound, "favorite not found")
				return
			}
			if handleError(w, err, http.StatusInternalServerError, "error while removing from favorites") {
				return
			}

			response := map[string]interface{}{
				"success": true,
				"message": "removed from favorites",
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) CheckIfFavoriteAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil || userID == nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			isFavorite, err := h.db.Psql.CheckIfFavorite(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error while checking favorite status") {
				return
			}

			response := map[string]interface{}{
				"success":    true,
				"isFavorite": isFavorite,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}
