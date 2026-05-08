package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"triangleoflove/backend/internal/service"
)

// PairingHandler handles all pairing and couple endpoints.
// Each method is registered to its own route in main.go, so no URL dispatch is needed here.
type PairingHandler struct {
	svc *service.PairingService
}

func NewPairingHandler(svc *service.PairingService) *PairingHandler {
	return &PairingHandler{svc: svc}
}

// GetCode handles GET /api/v1/pairing.
func (h *PairingHandler) GetCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	accountID := CallerFromContext(r.Context()).ID
	code, err := h.svc.GetOrCreateCode(r.Context(), accountID)
	if err != nil {
		log.Printf("get or create invite code failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"invite_code": code})
}

// Regenerate handles POST /api/v1/pairing/regenerate.
func (h *PairingHandler) Regenerate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	accountID := CallerFromContext(r.Context()).ID
	code, err := h.svc.RegenerateCode(r.Context(), accountID)
	if err != nil {
		log.Printf("regenerate invite code failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"invite_code": code})
}

// Connect handles POST /api/v1/pairing/connect.
func (h *PairingHandler) Connect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	accountID := CallerFromContext(r.Context()).ID
	var body struct {
		InviteCode string `json:"invite_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.InviteCode == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invite_code required"})
		return
	}
	err := h.svc.Connect(r.Context(), accountID, body.InviteCode)
	if errors.Is(err, service.ErrCodeNotFound) {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "invalid invite code"})
		return
	}
	if errors.Is(err, service.ErrAlreadyPaired) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "already paired"})
		return
	}
	if err != nil {
		log.Printf("connect failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "connected"})
}

// Unpair handles DELETE /api/v1/couples/me.
func (h *PairingHandler) Unpair(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	accountID := CallerFromContext(r.Context()).ID
	err := h.svc.Unpair(r.Context(), accountID)
	if errors.Is(err, service.ErrNotPaired) {
		writeJSON(w, http.StatusConflict, map[string]string{"error": "not paired"})
		return
	}
	if err != nil {
		log.Printf("unpair failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "unpaired"})
}

// GetCoupleStatus handles GET /api/v1/couples/me.
func (h *PairingHandler) GetCoupleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	accountID := CallerFromContext(r.Context()).ID
	status, err := h.svc.GetCoupleStatus(r.Context(), accountID)
	if err != nil {
		log.Printf("get couple status failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}
	if !status.Paired {
		writeJSON(w, http.StatusOK, map[string]any{"paired": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"paired":             true,
		"partner_first_name": status.PartnerFirstName,
		"paired_since":       status.PairedSince.Format(time.RFC3339),
	})
}
