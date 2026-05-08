package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

// NewCheckinHandler returns an http.Handler for GET/PUT /api/v1/checkins/today.
// The handler expects CallerAccount to be set in context (i.e. wrapped by Middleware).
func NewCheckinHandler(svc *service.CheckinService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accountID := CallerFromContext(r.Context()).ID

		switch r.Method {
		case http.MethodGet:
			c, err := svc.GetToday(r.Context(), accountID)
			if errors.Is(err, domain.ErrNotFound) {
				writeJSON(w, http.StatusNotFound, map[string]string{"error": "no check-in for today"})
				return
			}
			if err != nil {
				log.Printf("get checkin failed: %v", err)
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				return
			}
			writeJSON(w, http.StatusOK, c)

		case http.MethodPut:
			var body domain.Checkin
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
				return
			}
			saved, err := svc.SaveToday(r.Context(), accountID, body)
			if err != nil {
				log.Printf("save checkin failed: %v", err)
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
				return
			}
			writeJSON(w, http.StatusOK, saved)

		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})
}
