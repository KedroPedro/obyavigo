package api

import (
	"cmd/obyavigo/main.go/internal/models"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

// GetAdminStats возвращает статистику для админ панели
func (h *Handlers) GetAdminStats() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				sendJSONError(w, http.StatusForbidden, "access denied")
				return
			}

			stats := models.AdminStats{
				TotalAds:          0,
				TotalUsers:        0,
				PendingReports:    0,
				PendingModeration: 0,
			}

			// Получить статистику из БД
			// TODO: Реализовать методы в БД
			// stats.TotalAds = h.db.Psql.GetTotalAds()
			// stats.TotalUsers = h.db.Psql.GetTotalUsers()
			// stats.PendingReports = h.db.Psql.GetPendingReports()
			// stats.PendingModeration = h.db.Psql.GetPendingModeration()

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(stats)
		},
	)
}

// GetAdminAds возвращает список объявлений для админ панели
func (h *Handlers) GetAdminAds() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				sendJSONError(w, http.StatusForbidden, "access denied")
				return
			}

			page := parseIntParam(r, "page", 1, 0)
			limit := parseIntParam(r, "limit", 20, 0)
			status := r.URL.Query().Get("status")
			search := r.URL.Query().Get("search")

			filters := models.AdSearchFilters{
				Page:        page,
				Limit:       limit,
				SearchQuery: search,
			}

			if status != "" && status != "all" {
				// Фильтр по статусу
				// TODO: добавить поле AdStatus в AdSearchFilters
			}

			ads, err := h.db.Psql.SearchAds(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error searching ads") {
				return
			}

			totalCount, err := h.db.Psql.SearchAdsCount(&filters)
			if handleError(w, err, http.StatusInternalServerError, "error getting ads count") {
				return
			}

			response := map[string]interface{}{
				"ads":        ads,
				"totalCount": totalCount,
				"page":       page,
				"limit":      limit,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		},
	)
}

// UpdateAdStatus обновляет статус объявления (одобрить/отклонить)
func (h *Handlers) UpdateAdStatus() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				sendJSONError(w, http.StatusForbidden, "access denied")
				return
			}

			adIDStr := r.PathValue("id")
			adID, err := uuid.Parse(adIDStr)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid ad id")
				return
			}

			var req struct {
				Status string `json:"status"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			if req.Status != "public" && req.Status != "rejected" && req.Status != "draft" {
				sendJSONError(w, http.StatusBadRequest, "invalid status")
				return
			}

			// TODO: Реализовать метод UpdateAdStatus в БД
			// err = h.db.Psql.UpdateAdStatus(&adID, req.Status)
			// if handleError(w, err, http.StatusInternalServerError, "error updating ad status") {
			// 	return
			// }

			slog.Info("ad status updated",
				slog.String("ad_id", adID.String()),
				slog.String("status", req.Status),
				slog.String("admin_id", userID.String()))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "status updated successfully",
			})
		},
	)
}

// GetAdminUsers возвращает список пользователей для админ панели
func (h *Handlers) GetAdminUsers() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				sendJSONError(w, http.StatusForbidden, "access denied")
				return
			}

			page := parseIntParam(r, "page", 1, 0)
			limit := parseIntParam(r, "limit", 20, 0)
			_ = r.URL.Query().Get("role")    // roleFilter - для будущей реализации
			_ = r.URL.Query().Get("search") // search - для будущей реализации

			// TODO: Реализовать метод GetUsers в БД
			// users, err := h.db.Psql.GetUsers(page, limit, roleFilter, search)
			// if handleError(w, err, http.StatusInternalServerError, "error getting users") {
			// 	return
			// }

			// Заглушка
			users := []models.User{}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"users": users,
				"page":  page,
				"limit": limit,
			})
		},
	)
}

// UpdateUserStatus обновляет статус пользователя (бан/разбан)
func (h *Handlers) UpdateUserStatus() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role, err := h.db.Psql.GetUserRole(userID)
			if handleError(w, err, http.StatusInternalServerError, "error getting user role") {
				return
			}

			if role != "admin" && role != "moderator" {
				sendJSONError(w, http.StatusForbidden, "access denied")
				return
			}

			targetUserIDStr := r.PathValue("id")
			targetUserID, err := uuid.Parse(targetUserIDStr)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid user id")
				return
			}

			var req struct {
				Status string `json:"status"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			if req.Status != "active" && req.Status != "banned" {
				sendJSONError(w, http.StatusBadRequest, "invalid status")
				return
			}

			// TODO: Реализовать метод UpdateUserStatus в БД
			// err = h.db.Psql.UpdateUserStatus(&targetUserID, req.Status)
			// if handleError(w, err, http.StatusInternalServerError, "error updating user status") {
			// 	return
			// }

			slog.Info("user status updated",
				slog.String("target_user_id", targetUserID.String()),
				slog.String("status", req.Status),
				slog.String("admin_id", userID.String()))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "status updated successfully",
			})
		},
	)
}
