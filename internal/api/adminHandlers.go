package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

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

			stats, err := h.db.Psql.GetAdminStats()
			if handleError(w, err, http.StatusInternalServerError, "error getting admin stats") {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(stats)
		},
	)
}

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

			ads, totalCount, err := h.db.Psql.GetAdminAds(page, limit, status, search)
			if handleError(w, err, http.StatusInternalServerError, "error getting ads") {
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

			err = h.db.Psql.UpdateAdStatus(&adID, req.Status)
			if handleError(w, err, http.StatusInternalServerError, "error updating ad status") {
				return
			}

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
			roleFilter := r.URL.Query().Get("role")
			search := r.URL.Query().Get("search")

			users, err := h.db.Psql.GetUsers(page, limit, roleFilter, search)
			if handleError(w, err, http.StatusInternalServerError, "error getting users") {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"users": users,
				"page":  page,
				"limit": limit,
			})
		},
	)
}

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

			err = h.db.Psql.UpdateUserStatus(&targetUserID, req.Status)
			if handleError(w, err, http.StatusInternalServerError, "error updating user status") {
				return
			}

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

func (h *Handlers) UpdateUserRole() http.Handler {
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

			if role != "admin" {
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
				Role string `json:"role"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			if req.Role != "moderator" && req.Role != "user" && req.Role != "admin" {
				sendJSONError(w, http.StatusBadRequest, "invalid role")
				return
			}

			err = h.db.Psql.UpdateUserRole(&targetUserID, req.Role)
			if handleError(w, err, http.StatusInternalServerError, "error updating user role") {
				return
			}

			slog.Info("user role updated",
				slog.String("target_user_id", targetUserID.String()),
				slog.String("role", req.Role),
				slog.String("admin_id", userID.String()))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "role updated successfully",
			})
		},
	)
}
