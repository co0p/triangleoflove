package web

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"triangleoflove/backend/internal/domain"
)

// InsightsWeeklyService is the interface required by InsightsWeeklyHandler.
type InsightsWeeklyService interface {
	GetWeekly(ctx context.Context, accountID string) ([]domain.WeeklyInsight, error)
	GetWindow(ctx context.Context, accountID string, past int) ([]domain.WeeklyInsight, error)
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

	accountID := CallerFromContext(r.Context()).ID

	pastStr := r.URL.Query().Get("past")
	if pastStr == "" {
		weekly, err := h.svc.GetWeekly(r.Context(), accountID)
		if err != nil {
			log.Printf("get weekly insights failed: %v", err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
			return
		}
		writeJSON(w, http.StatusOK, weekly)
		return
	}

	past, err := strconv.Atoi(pastStr)
	if err != nil || past < 1 || past > 100 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "past must be an integer between 1 and 100"})
		return
	}

	window, err := h.svc.GetWindow(r.Context(), accountID, past)
	if err != nil {
		log.Printf("get window insights failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, window)
}
