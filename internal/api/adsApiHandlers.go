package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"encoding/json"
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
				Location:    r.URL.Query().Get("location"),
				Condition:   r.URL.Query().Get("condition"),
				SearchQuery: r.URL.Query().Get("q"),
				SortBy:      r.URL.Query().Get("sort"),
				MinPrice:    parsePriceParam(r, "min_price"),
				MaxPrice:    parsePriceParam(r, "max_price"),
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

			response := map[string]interface{}{
				"success": true,
				"data":    adData,
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(response)
		},
	)
}
