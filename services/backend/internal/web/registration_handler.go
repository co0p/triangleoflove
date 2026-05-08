package web

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"triangleoflove/backend/internal/service"
)

// RegistrationService is the interface required by RegistrationHandler.
type RegistrationService interface {
	Register(ctx context.Context, email, password, firstName string) error
}

// RegistrationHandler handles POST /api/v1/register.
type RegistrationHandler struct {
	svc RegistrationService
}

func NewRegistrationHandler(svc RegistrationService) http.Handler {
	h := &RegistrationHandler{svc: svc}
	return http.HandlerFunc(h.handle)
}

func (h *RegistrationHandler) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var body struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		FirstName string `json:"firstName"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request"})
		return
	}

	if body.Email == "" || body.Password == "" || body.FirstName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email, password, and firstName are required"})
		return
	}

	err := h.svc.Register(r.Context(), body.Email, body.Password, body.FirstName)
	if errors.Is(err, service.ErrEmailAlreadyExists) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "email already exists"})
		return
	}
	if errors.Is(err, service.ErrPasswordTooShort) {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "password must be at least 8 characters"})
		return
	}
	if err != nil {
		log.Printf("register failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	w.WriteHeader(http.StatusCreated)
}
