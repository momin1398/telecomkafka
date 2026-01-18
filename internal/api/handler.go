// internal/api/handler.go
package api

import (
	"encoding/json"
	"net/http"

	"github.com/momin1398/telecomkafka/internal/model"
	"github.com/momin1398/telecomkafka/internal/service"
)

type Handler struct {
	telemetry *service.TelemetryService
}

func NewHandler(t *service.TelemetryService) *Handler {
	return &Handler{telemetry: t}
}

func (h *Handler) Collect(w http.ResponseWriter, r *http.Request) {
	var event model.TelemetryEvent

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if err := h.telemetry.Process(r.Context(), event); err != nil {
		http.Error(w, "failed to process event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
