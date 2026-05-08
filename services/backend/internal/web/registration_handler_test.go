package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

type mockRegistrationSvc struct {
	err error
}

func (m *mockRegistrationSvc) Register(_ context.Context, _, _, _ string) error {
	return m.err
}

func TestRegistrationHandler_GivenDuplicateEmail_WhenPOST_ThenReturns409(t *testing.T) {
	svc := &mockRegistrationSvc{err: service.ErrEmailAlreadyExists}
	handler := web.NewRegistrationHandler(svc)

	body, _ := json.Marshal(map[string]string{
		"email":     "taken@example.com",
		"password":  "securepass",
		"firstName": "User",
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestRegistrationHandler_GivenMissingFields_WhenPOST_ThenReturns400(t *testing.T) {
	svc := &mockRegistrationSvc{}
	handler := web.NewRegistrationHandler(svc)

	body, _ := json.Marshal(map[string]string{
		"email": "user@example.com",
		// password and firstName intentionally absent
	})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}
}
