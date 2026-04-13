package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

type mockChangePasswordAccountRepo struct {
	account   domain.Account
	findErr   error
	saveErr   error
	savedHash string
}

func (m *mockChangePasswordAccountRepo) FindByEmail(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.findErr
}

func (m *mockChangePasswordAccountRepo) FindByID(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.findErr
}

func (m *mockChangePasswordAccountRepo) SaveHashedPassword(_ context.Context, _ string, hashedPassword string) error {
	m.savedHash = hashedPassword
	return m.saveErr
}

func hashedPasswordForTest(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(h)
}

func TestChangePasswordHandler_GivenNoAuth_WhenPUT_ThenReturns401(t *testing.T) {
	repo := &mockChangePasswordAccountRepo{}
	svc := service.NewAuthService(repo)
	handler := web.Middleware(web.NewChangePasswordHandler(svc))

	body, _ := json.Marshal(map[string]string{
		"current_password": "old",
		"new_password":     "new",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/password", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestChangePasswordHandler_GivenWrongCurrentPassword_WhenPUT_ThenReturns409(t *testing.T) {
	repo := &mockChangePasswordAccountRepo{
		account: domain.Account{
			ID:             "acc-1",
			Email:          "alice@example.com",
			HashedPassword: hashedPasswordForTest(t, "correct-pass"),
			FirstName:      "Alice",
		},
	}
	svc := service.NewAuthService(repo)
	handler := web.Middleware(web.NewChangePasswordHandler(svc))

	token, err := auth.SignToken("acc-1")
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}

	body, _ := json.Marshal(map[string]string{
		"current_password": "wrong-pass",
		"new_password":     "brand-new",
	})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/auth/password", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}
