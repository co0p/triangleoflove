package web

import (
	"context"
	"errors"
	"log"
	"net/http"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

// InsightsService is the interface required by InsightsHandler.
type InsightsService interface {
	GetByDate(ctx context.Context, accountID string, date string) (domain.DailyInsight, error)
}

// InsightsHandler handles GET /api/v1/insights/{date}.
type InsightsHandler struct {
	svc InsightsService
}

func NewInsightsHandler(svc InsightsService) http.Handler {
	h := &InsightsHandler{svc: svc}
	return http.HandlerFunc(h.handle)
}

func (h *InsightsHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	accountID := CallerFromContext(r.Context()).ID
	date := r.PathValue("date")

	insight, err := h.svc.GetByDate(r.Context(), accountID, date)
	if errors.Is(err, service.ErrInvalidDate) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid date format"})
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no check-in for this date"})
		return
	}
	if err != nil {
		log.Printf("get insights failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, insight)
}
