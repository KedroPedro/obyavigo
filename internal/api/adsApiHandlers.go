package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
)

func parseIntParam(r *http.Request, key string, defaultVal, maxVal int) int {
	if val := r.URL.Query().Get(key); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil && parsed > 0 {
			if maxVal > 0 && parsed > maxVal {
				return maxVal
			}
			return parsed
		}
	}
	return defaultVal
}

func parsePriceParam(r *http.Request, key string) *int {
	if val := r.URL.Query().Get(key); val != "" {
		if price, err := strconv.Atoi(val); err == nil && price >= 0 {
			return &price
		}
	}
	return nil
}

func (h *Handlers) GetAdsAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			filters := models.AdSearchFilters{
				Page:        parseIntParam(r, "page", 1, 0),
				Limit:       parseIntParam(r, "limit", 20, 100),
				Category:    r.URL.Query().Get("category"),
				Subcategory: r.URL.Query().Get("subcategory"),
				Region:      r.URL.Query().Get("region"),
				Location:    r.URL.Query().Get("location"),
				Condition:   r.URL.Query().Get("condition"),
				SearchQuery: r.URL.Query().Get("q"),
				SortBy:      r.URL.Query().Get("sort"),
				MinPrice:    parsePriceParam(r, "min_price"),
				MaxPrice:    parsePriceParam(r, "max_price"),
			}

			if filters.MinPrice != nil && filters.MaxPrice != nil && *filters.MinPrice > *filters.MaxPrice {
				sendToClient(w, http.StatusBadRequest, "min_price cannot be greater than max_price")
				return
			}

			ads, err := h.db.Psql.SearchAds(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error while searching ads") {
				return
			}

			totalCount, err := h.db.Psql.SearchAdsCount(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error while getting ads count") {
				return
			}

			response := map[string]interface{}{
				"success":    true,
				"data":       ads,
				"totalCount": totalCount,
				"page":       filters.Page,
				"limit":      filters.Limit,
				"totalPages": (totalCount + filters.Limit - 1) / filters.Limit,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) GetAdByIDAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			adData, err := h.db.Psql.GetAdInfo(&adID)
			if adData == nil || err != nil {
				sendToClient(w, http.StatusNotFound, "ad not found")
				return
			}

			userID, _ := userIDFromCtx(r)
			role := "user"
			if userID != nil {
				role, err = h.db.Psql.GetUserRole(userID)
				if handleError(w, err, http.StatusInternalServerError, "error while trying to get user role") {
					return
				}
			}

			if userID == nil || (*userID != adData.UserID && role == "user") {
				if adData.AdStatus != "public" {
					sendToClient(w, http.StatusNotFound, "ad not found")
					return
				}
			}

			
			imageIDs, err := h.db.Psql.GetAdImageIDs(&adID)
			if err != nil {
				slog.Error("error getting ad images", slog.String("ad_id", adID.String()), slog.String("error", err.Error()))
				imageIDs = []string{}
			}
			adData.Images = imageIDs

			response := map[string]interface{}{
				"success": true,
				"data":    adData,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) GetUserAdsAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ads, err := h.db.Psql.GetUserAds(userID)
			if handleError(w, err, http.StatusInternalServerError, "error while getting user ads") {
				return
			}

			response := map[string]interface{}{
				"success": true,
				"data":    ads,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}

func (h *Handlers) DeleteAdAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			
			isOwner, err := h.db.Psql.CheckAdOwnership(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error checking ad ownership") {
				return
			}

			if !isOwner {
				sendToClient(w, http.StatusForbidden, "you are not the owner of this ad")
				return
			}

<<<<<<< HEAD
			if err := h.db.Mongo.DeleteAdImages(r.Context(), adID.String()); err != nil {
				slog.Error("error deleting ad images", slog.String("ad_id", adID.String()), slog.String("error", err.Error()))
			}

=======
			
			if err := h.db.Mongo.DeleteAdImages(r.Context(), adID.String()); err != nil {
				
				slog.Error("error deleting ad images", slog.String("ad_id", adID.String()), slog.String("error", err.Error()))
			}

			
>>>>>>> c04cbe9c777b551f29c288a2c9d239c0b97177a5
			if err := h.db.Psql.DeleteAd(&adID); err != nil {
				sendToClient(w, http.StatusInternalServerError, "error deleting ad")
				return
			}

			sendToClient(w, http.StatusOK, "ad deleted successfully")
		},
	)
}

type UpdateAdRequest struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	Condition    string `json:"condition"`
	ContactPhone string `json:"contact_phone"`
}

func (h *Handlers) UpdateAdAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			token := r.PathValue("token")
			adID, err := parseUUID(token)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			
			isOwner, err := h.db.Psql.CheckAdOwnership(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error checking ad ownership") {
				return
			}

			if !isOwner {
				sendToClient(w, http.StatusForbidden, "you are not the owner of this ad")
				return
			}

			var req UpdateAdRequest
			if err := json.NewDecoder(r.Body).Decode(&req); handleError(w, err, http.StatusBadRequest, "invalid request body") {
				return
			}

			if req.Title == "" {
				sendToClient(w, http.StatusBadRequest, "title is required")
				return
			}

			if req.Price < 0 {
				sendToClient(w, http.StatusBadRequest, "price must be positive")
				return
			}

			
			if err := h.db.Psql.UpdateAd(&adID, req.Title, req.Description, req.Price, req.Condition, req.ContactPhone); err != nil {
				sendToClient(w, http.StatusInternalServerError, "error updating ad")
				return
			}

			sendToClient(w, http.StatusOK, "ad updated successfully")
		},
	)
}

func (h *Handlers) DeleteAdImageAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			adToken := r.PathValue("token")
			adID, err := parseUUID(adToken)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			imageID := r.PathValue("imageId")
			if imageID == "" {
				sendToClient(w, http.StatusBadRequest, "image id is required")
				return
			}

			
			isOwner, err := h.db.Psql.CheckAdOwnership(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error checking ad ownership") {
				return
			}

			if !isOwner {
				sendToClient(w, http.StatusForbidden, "you are not the owner of this ad")
				return
			}

			
			if err := h.db.Mongo.DeleteImageByID(r.Context(), imageID); err != nil {
				slog.Error("error deleting image from MongoDB", slog.String("image_id", imageID), slog.String("error", err.Error()))
				sendToClient(w, http.StatusInternalServerError, "error deleting image")
				return
			}

			
			if err := h.db.Psql.DeleteAdImage(imageID); err != nil {
				slog.Error("error deleting image record from PostgreSQL", slog.String("image_id", imageID), slog.String("error", err.Error()))
				
			}

			sendToClient(w, http.StatusOK, "image deleted successfully")
		},
	)
}

func (h *Handlers) UploadAdImagesAPI() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendToClient(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			adToken := r.PathValue("token")
			adID, err := parseUUID(adToken)
			if err != nil {
				sendToClient(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			
			isOwner, err := h.db.Psql.CheckAdOwnership(userID, &adID)
			if handleError(w, err, http.StatusInternalServerError, "error checking ad ownership") {
				return
			}

			if !isOwner {
				sendToClient(w, http.StatusForbidden, "you are not the owner of this ad")
				return
			}

			if err := r.ParseMultipartForm(10 << 20); err != nil {
				sendToClient(w, http.StatusBadRequest, "failed to parse form")
				return
			}

			files := r.MultipartForm.File["images"]
			if len(files) == 0 {
				sendToClient(w, http.StatusBadRequest, "no images provided")
				return
			}

			
			imageIDs, err := h.db.Mongo.UploadImages(r.Context(), files, adID.String())
			if err != nil {
				slog.Error("error uploading images", slog.String("ad_id", adID.String()), slog.String("error", err.Error()))
				sendToClient(w, http.StatusInternalServerError, "error uploading images")
				return
			}

			
			if err := h.db.Psql.InsertImages(userID, &adID, imageIDs); err != nil {
				slog.Error("error inserting image records", slog.String("ad_id", adID.String()), slog.String("error", err.Error()))
				
				for _, imgID := range imageIDs {
					h.db.Mongo.DeleteImageByID(r.Context(), imgID)
				}
				sendToClient(w, http.StatusInternalServerError, "error saving image records")
				return
			}

			response := map[string]interface{}{
				"success":   true,
				"message":   "images uploaded successfully",
				"image_ids": imageIDs,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}
