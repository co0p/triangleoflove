package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

// mockInviteCodeRepo implements service.InviteCodeRepo.
type mockInviteCodeRepo struct {
	code      string
	partnerID string
	paired    bool
}

func (m *mockInviteCodeRepo) FindInviteCodeByAccountID(_ context.Context, _ string) (domain.InviteCode, error) {
	return domain.InviteCode(m.code), nil
}

func (m *mockInviteCodeRepo) SaveInviteCode(_ context.Context, _ string, code domain.InviteCode) error {
	m.code = string(code)
	return nil
}

func (m *mockInviteCodeRepo) FindByInviteCode(_ context.Context, _ domain.InviteCode) (string, error) {
	if m.partnerID == "" {
		return "", domain.ErrNotFound
	}
	return m.partnerID, nil
}

func (m *mockInviteCodeRepo) ExistsCoupleByAccountID(_ context.Context, _ string) (bool, error) {
	return m.paired, nil
}

// mockCoupleRepo implements service.CoupleRepo.
type mockCoupleRepo struct {
	coupled   bool
	summary   *domain.CoupleSummary
	unpairErr error
}

func (m *mockCoupleRepo) Save(_ context.Context, _, _ string) error {
	m.coupled = true
	return nil
}

func (m *mockCoupleRepo) FindByAccountID(_ context.Context, _ string) (domain.CoupleSummary, error) {
	if m.summary == nil {
		return domain.CoupleSummary{}, domain.ErrNotFound
	}
	return *m.summary, nil
}

func (m *mockCoupleRepo) Unpair(_ context.Context, _ string) error {
	return m.unpairErr
}

func TestPairingHandler_GivenNoStoredCode_WhenGET_ThenReturns6CharCode(t *testing.T) {
	codes := &mockInviteCodeRepo{code: ""}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.GetCode))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pairing", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	code, ok := body["invite_code"]
	if !ok {
		t.Fatal("expected invite_code in response")
	}
	if !regexp.MustCompile(`^[A-Z0-9]{6}$`).MatchString(code) {
		t.Fatalf("expected 6-char alphanumeric code, got %q", code)
	}
}

func TestPairingHandler_GivenStoredCode_WhenGET_ThenReturnsSameCode(t *testing.T) {
	codes := &mockInviteCodeRepo{code: "ABC123"}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.GetCode))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/pairing", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	if body["invite_code"] != "ABC123" {
		t.Fatalf("expected ABC123, got %q", body["invite_code"])
	}
}

func TestPairingHandler_GivenStoredCode_WhenPOSTRegenerate_ThenReturnsNewCode(t *testing.T) {
	codes := &mockInviteCodeRepo{code: "OLD123"}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Regenerate))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pairing/regenerate", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	json.NewDecoder(rec.Body).Decode(&body)
	newCode := body["invite_code"]

	if newCode == "OLD123" {
		t.Fatal("expected a new code, got the old one")
	}
	if !regexp.MustCompile(`^[A-Z0-9]{6}$`).MatchString(newCode) {
		t.Fatalf("expected 6-char alphanumeric code, got %q", newCode)
	}
}

func TestPairingHandler_GivenStoredCode_WhenPOSTRegenerate_ThenOldCodeReplaced(t *testing.T) {
	codes := &mockInviteCodeRepo{code: "OLD123"}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Regenerate))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pairing/regenerate", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if codes.code == "OLD123" {
		t.Fatal("expected stored code to be replaced after regenerate")
	}
}

func TestPairingHandler_GivenInvalidCode_WhenPOSTConnect_ThenReturns422(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: ""}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Connect))

	token, _ := auth.SignToken("account-123", "user")
	payload, _ := json.Marshal(map[string]string{"invite_code": "BADCOD"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pairing/connect", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPairingHandler_GivenAlreadyPaired_WhenPOSTConnect_ThenReturns409(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: "partner-id", paired: true}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Connect))

	token, _ := auth.SignToken("account-123", "user")
	payload, _ := json.Marshal(map[string]string{"invite_code": "VALIDC"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pairing/connect", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPairingHandler_GivenValidCode_WhenPOSTConnect_ThenReturns200(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: "partner-id", paired: false}
	svc := service.NewPairingService(codes, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Connect))

	token, _ := auth.SignToken("account-123", "user")
	payload, _ := json.Marshal(map[string]string{"invite_code": "VALIDC"})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/pairing/connect", bytes.NewReader(payload))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUnpairHandler_GivenNoAuth_WhenDELETE_ThenReturns401(t *testing.T) {
	svc := service.NewPairingService(&mockInviteCodeRepo{}, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Unpair))

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/couples/me", nil)
	// No Authorization header.
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestUnpairHandler_GivenNotPaired_WhenDELETE_ThenReturns409(t *testing.T) {
	couple := &mockCoupleRepo{unpairErr: service.ErrNotPaired}
	svc := service.NewPairingService(&mockInviteCodeRepo{}, couple)
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.Unpair))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/couples/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestPairingHandler_GivenNotPaired_WhenGETCouplesMe_ThenReturnsPairedFalse(t *testing.T) {
	svc := service.NewPairingService(&mockInviteCodeRepo{}, &mockCoupleRepo{})
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.GetCoupleStatus))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/couples/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	if body["paired"] != false {
		t.Fatalf("expected paired=false, got %v", body["paired"])
	}
}

func TestPairingHandler_GivenPaired_WhenGETCouplesMe_ThenReturnsPartnerName(t *testing.T) {
	coupleRepo := &mockCoupleRepo{
		summary: &domain.CoupleSummary{
			PartnerFirstName: "Jordan",
			FormedOn:         time.Date(2026, 4, 9, 10, 0, 0, 0, time.UTC),
		},
	}
	svc := service.NewPairingService(&mockInviteCodeRepo{}, coupleRepo)
	ph := web.NewPairingHandler(svc)
	handler := web.Middleware(http.HandlerFunc(ph.GetCoupleStatus))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/couples/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var body map[string]any
	json.NewDecoder(rec.Body).Decode(&body)
	if body["paired"] != true {
		t.Fatalf("expected paired=true, got %v", body["paired"])
	}
	if body["partner_first_name"] != "Jordan" {
		t.Fatalf("expected partner_first_name=Jordan, got %v", body["partner_first_name"])
	}
	if body["paired_since"] == nil {
		t.Fatal("expected paired_since in response")
	}
}
