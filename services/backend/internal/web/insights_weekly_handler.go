package web

import (
	"context"
	"log"
	"net/http"

	"triangleoflove/backend/internal/domain"
)

// InsightsWeeklyService is the interface required by InsightsWeeklyHandler.
type InsightsWeeklyService interface {
	GetWeekly(ctx context.Context, accountID string) ([]domain.WeeklyInsight, error)
}

// InsightsWeeklyHandler handles GET /api/v1/insights.
type InsightsWeeklyHandler struct {
	svc InsightsWeeklyService
}

func NewInsightsWeeklyHandler(svc InsightsWeeklyService) http.Handler {
	h := &InsightsWeeklyHandler{svc: svc}
	return http.HandlerFunc(h.handle)
}

func (h *InsightsWeeklyHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	accountID, _ := r.Context().Value(AccountIDKey).(string)

	weekly, err := h.svc.GetWeekly(r.Context(), accountID)
	if err != nil {
		log.Printf("get weekly insights failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, weekly)
}
