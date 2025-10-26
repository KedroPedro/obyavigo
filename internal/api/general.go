package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
)

func sendToClient(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{
		"message": message,
	})
}

func (h *Handlers) sendNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	err := h.tmpl.ExecuteTemplate(w, "not-found.html", nil)
	if err != nil {
		sendToClient(w, 404, "Internal server error")
	}
}

func handleError(w http.ResponseWriter, err error, code int, logMsg string) bool {
	if err != nil {
		sendToClient(w, code, http.StatusText(code))
		slog.Error(logMsg, slog.String("error", err.Error()))
		return true
	}
	return false
}

func userIDFromCtx(r *http.Request) (*uuid.UUID, error) {
	id := r.Context().Value("userID")
	strID, ok := id.(string)
	if !ok {
		return nil, fmt.Errorf("error while trying to convert any value to a string")
	}
	userID, err := uuid.Parse(strID)
	if err != nil {
		return nil, fmt.Errorf("error while trying to convert any value to a string")
	}
	return &userID, nil
}

func parseUUID(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("error while parsing uuid: %w", err)
	}
	return id, nil
}
