package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func (h *Handlers) CreateReport() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userID, err := userIDFromCtx(r)
			if err != nil {
				sendJSONError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			var req struct {
				AdId        string `json:"ad_id"`
				ReportType  string `json:"report_type"`
				Description string `json:"description"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			if req.AdId == "" || req.ReportType == "" {
				sendJSONError(w, http.StatusBadRequest, "ad_id and report_type are required")
				return
			}

			adID, err := uuid.Parse(req.AdId)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid ad_id")
				return
			}

			validTypes := map[string]bool{
				"spam":           true,
				"fraud":          true,
				"fake":           true,
				"wrong_category": true,
				"offensive":      true,
				"sold":           true,
				"other":          true,
			}

			if !validTypes[req.ReportType] {
				sendJSONError(w, http.StatusBadRequest, "invalid report type")
				return
			}

			err = h.db.Psql.CreateReport(userID, &adID, req.ReportType, req.Description)
			if handleError(w, err, http.StatusInternalServerError, "error creating report") {
				return
			}

			slog.Info("report created",
				slog.String("user_id", userID.String()),
				slog.String("ad_id", adID.String()),
				slog.String("report_type", req.ReportType))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{
				"message": "report created successfully",
			})
		},
	)
}

func (h *Handlers) GetReports() http.Handler {
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

		reports, err := h.db.Psql.GetReports(page, limit, status)
		if handleError(w, err, http.StatusInternalServerError, "error getting reports") {
			return
		}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"reports": reports,
				"page":    page,
				"limit":   limit,
			})
		},
	)
}

func (h *Handlers) UpdateReportStatus() http.Handler {
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

			reportIDStr := r.PathValue("id")
			reportID, err := uuid.Parse(reportIDStr)
			if err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid report id")
				return
			}

			var req struct {
				Status            string `json:"status"`
				ResolutionComment string `json:"resolution_comment"`
			}

			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				sendJSONError(w, http.StatusBadRequest, "invalid request body")
				return
			}

			if req.Status != "resolved" && req.Status != "rejected" {
				sendJSONError(w, http.StatusBadRequest, "invalid status")
				return
			}

			err = h.db.Psql.UpdateReportStatus(&reportID, req.Status, req.ResolutionComment, userID)
			if handleError(w, err, http.StatusInternalServerError, "error updating report status") {
				return
			}

			slog.Info("report status updated",
				slog.String("report_id", reportID.String()),
				slog.String("status", req.Status),
				slog.String("admin_id", userID.String()))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "report status updated successfully",
			})
		},
	)
}
