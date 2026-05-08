package web

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

// AdminService is the interface required by AdminHandler.
type AdminService interface {
	ListUsers(ctx context.Context, callerRole domain.Role) ([]domain.AccountSummary, error)
	SetActive(ctx context.Context, callerRole domain.Role, targetID string, active bool) error
}

// AdminHandler handles admin user management endpoints.
type AdminHandler struct {
	svc AdminService
}

func NewAdminHandler(svc AdminService) http.Handler {
	h := &AdminHandler{svc: svc}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/admin/users", h.listUsers)
	mux.HandleFunc("PUT /api/v1/admin/users/{id}/activate", h.activate)
	mux.HandleFunc("PUT /api/v1/admin/users/{id}/deactivate", h.deactivate)
	return mux
}

func (h *AdminHandler) listUsers(w http.ResponseWriter, r *http.Request) {
	caller := CallerFromContext(r.Context())
	users, err := h.svc.ListUsers(r.Context(), caller.Role)
	if errors.Is(err, service.ErrForbidden) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	if err != nil {
		log.Printf("admin list users failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	type userRow struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		Role      string `json:"role"`
		IsActive  bool   `json:"isActive"`
		CreatedAt string `json:"createdAt"`
	}
	rows := make([]userRow, len(users))
	for i, u := range users {
		rows[i] = userRow{
			ID:        u.ID,
			Email:     u.Email,
			FirstName: u.FirstName,
			Role:      string(u.Role),
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
		}
	}
	writeJSON(w, http.StatusOK, rows)
}

func (h *AdminHandler) activate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, true)
}

func (h *AdminHandler) deactivate(w http.ResponseWriter, r *http.Request) {
	h.setActive(w, r, false)
}

func (h *AdminHandler) setActive(w http.ResponseWriter, r *http.Request, active bool) {
	id := strings.TrimPrefix(r.PathValue("id"), "")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing user id"})
		return
	}
	caller := CallerFromContext(r.Context())
	err := h.svc.SetActive(r.Context(), caller.Role, id, active)
	if errors.Is(err, service.ErrForbidden) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	if errors.Is(err, domain.ErrNotFound) {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
		return
	}
	if err != nil {
		log.Printf("admin set active failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
