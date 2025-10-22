package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handlers) CreateAd() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to get user id from context") {
				return
			}

			err = r.ParseMultipartForm(70 << 20) // 70 MB
			if handleError(w, err, http.StatusInternalServerError, "error while parsing multipart form") {
				return
			}

			adData := models.AdTemplate{
				UserId:       *userID,
				CategoryName: r.FormValue("categoryName"),
				LocationName: r.FormValue("locationName"),
				Title:        r.FormValue("title"),
				Description:  r.FormValue("desc"),
				ContactPhone: r.FormValue("phone"),
				Condition:    r.FormValue("condition"), // или задайте по умолчанию
			}

			if priceStr := r.FormValue("price"); priceStr != "" {
				if price, err := strconv.Atoi(priceStr); err == nil {
					adData.Price = price
				} else {
					sendToClient(w, http.StatusBadRequest, "invalid price format")
					return
				}
			}

			// Получаем файлы
			files := r.MultipartForm.File["images"]
			if len(files) == 0 {
				sendToClient(w, http.StatusBadRequest, "no images uploaded")
				return
			}

			if handleError(w, h.db.Psql.GetCreateAdDependencies(&adData), http.StatusInternalServerError, "error while trying to get dependencies for create new ad") {
				return
			}

			adId, err := h.db.Psql.CreateAd(&adData)
			if handleError(w, err, http.StatusInternalServerError, "error while trying to create new ad") {
				return
			}

			ids, err := h.db.Mongo.UploadImages(r.Context(), files, adId.String())
			if handleError(w, err, http.StatusInternalServerError, "error while trying to upload photos") {
				return
			}

			if handleError(w, h.db.Psql.InsertImages(userID, adId, ids), http.StatusInternalServerError, "error while trying to insert images") {
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": true,
				"adId":    adId.String(),
				"url":     "/ads/" + adId.String(),
			})
		},
	)
}
