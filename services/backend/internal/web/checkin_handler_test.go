package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

type mockSessionRepo struct {
	saved domain.Checkin
}

func (m *mockSessionRepo) FindByAccountAndDate(_ context.Context, _ string) (domain.Checkin, error) {
	return domain.Checkin{}, domain.ErrNotFound
}

func (m *mockSessionRepo) Save(_ context.Context, _ string, c domain.Checkin) (domain.Checkin, error) {
	m.saved = c
	return c, nil
}

func TestSessionHandler_GivenNoAuth_WhenGET_ThenReturns401(t *testing.T) {
	svc := service.NewSessionService(&mockSessionRepo{})
	handler := web.Middleware(web.NewSessionHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/today", nil)
	// deliberately no Authorization header
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestSessionHandler_GivenNoEntry_WhenGET_ThenReturns404(t *testing.T) {
	svc := service.NewSessionService(&mockSessionRepo{})
	handler := web.Middleware(web.NewSessionHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/sessions/today", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}

	var payload map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v", err)
	}
	if payload["error"] != "no session for today" {
		t.Fatalf("expected error 'no session for today', got %q", payload["error"])
	}
}

func TestSessionHandler_GivenMalformedBody_WhenPUT_ThenReturns400(t *testing.T) {
	svc := service.NewSessionService(&mockSessionRepo{})
	handler := web.Middleware(web.NewSessionHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodPut, "/api/v1/sessions/today", bytes.NewReader([]byte("not-json")))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestSessionHandler_GivenValidPUT_ThenSavesBody(t *testing.T) {
	payload := domain.Checkin{FeltUnderstood: 3, MeaningfulSharing: 2, Mood: 4}

	mock := &mockSessionRepo{}
	svc := service.NewSessionService(mock)
	handler := web.Middleware(web.NewSessionHandler(svc))

	body, _ := json.Marshal(payload)
	token, _ := auth.SignToken("account-123", "user")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/sessions/today", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if mock.saved.FeltUnderstood != 3 {
		t.Fatalf("expected saved FeltUnderstood=3, got %d", mock.saved.FeltUnderstood)
	}
}
