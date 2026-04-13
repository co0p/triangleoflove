package web

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"triangleoflove/backend/internal/service"
)

// ChangePasswordService is the interface required by ChangePasswordHandler.
type ChangePasswordService interface {
	ChangePassword(ctx context.Context, accountID, currentPassword, newPassword string) error
}

// ChangePasswordHandler handles PUT /api/v1/auth/password.
type ChangePasswordHandler struct {
	svc ChangePasswordService
}

func NewChangePasswordHandler(svc ChangePasswordService) http.Handler {
	h := &ChangePasswordHandler{svc: svc}
	return http.HandlerFunc(h.handle)
}

func (h *ChangePasswordHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	accountID, _ := r.Context().Value(AccountIDKey).(string)

	var body struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	err := h.svc.ChangePassword(r.Context(), accountID, body.CurrentPassword, body.NewPassword)
	if errors.Is(err, service.ErrInvalidCredentials) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "current password is incorrect"})
		return
	}
	if err != nil {
		log.Printf("change password failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}
