package web_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

// --- mock ---

type mockAdminSvc struct {
	users     []domain.AccountSummary
	listErr   error
	activeErr error
}

func (m *mockAdminSvc) ListUsers(_ context.Context, callerRole domain.Role) ([]domain.AccountSummary, error) {
	if callerRole != domain.RoleAdmin {
		return nil, service.ErrForbidden
	}
	return m.users, m.listErr
}

func (m *mockAdminSvc) SetActive(_ context.Context, callerRole domain.Role, _ string, _ bool) error {
	if callerRole != domain.RoleAdmin {
		return service.ErrForbidden
	}
	return m.activeErr
}

// --- helpers ---

func adminToken(t *testing.T) string {
	t.Helper()
	tok, err := auth.SignToken("admin-1", "admin")
	if err != nil {
		t.Fatalf("sign admin token: %v", err)
	}
	return tok
}

func userToken(t *testing.T) string {
	t.Helper()
	tok, err := auth.SignToken("user-1", "user")
	if err != nil {
		t.Fatalf("sign user token: %v", err)
	}
	return tok
}

// --- acceptance tests ---

// TestAdmin_GivenAdminUser_WhenViewingAdminPage_ThenAllUsersListed
func TestAdmin_GivenAdminUser_WhenViewingAdminPage_ThenAllUsersListed(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	svc := &mockAdminSvc{
		users: []domain.AccountSummary{
			{ID: "u-1", Email: "alice@example.com", FirstName: "Alice", Role: domain.RoleUser, IsActive: true, CreatedAt: now},
			{ID: "u-2", Email: "bob@example.com", FirstName: "Bob", Role: domain.RoleUser, IsActive: false, CreatedAt: now},
		},
	}
	handler := web.Middleware(web.NewAdminHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	var got []map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
}

// TestAdmin_GivenNonAdminUser_WhenAccessingAdminPage_ThenRedirectedAndForbidden
func TestAdmin_GivenNonAdminUser_WhenAccessingAdminPage_ThenRedirectedAndForbidden(t *testing.T) {
	svc := &mockAdminSvc{}
	handler := web.Middleware(web.NewAdminHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+userToken(t))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestAdmin_GivenAdminUser_WhenDeactivatingUser_ThenUserCannotLogin
func TestAdmin_GivenAdminUser_WhenDeactivatingUser_ThenUserCannotLogin(t *testing.T) {
	svc := &mockAdminSvc{}
	handler := web.Middleware(web.NewAdminHandler(svc))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/u-1/deactivate", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}

// TestAdmin_GivenAdminUser_WhenReactivatingUser_ThenUserCanLogin
func TestAdmin_GivenAdminUser_WhenReactivatingUser_ThenUserCanLogin(t *testing.T) {
	svc := &mockAdminSvc{}
	handler := web.Middleware(web.NewAdminHandler(svc))

	req := httptest.NewRequest(http.MethodPut, "/api/v1/admin/users/u-1/activate", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken(t))
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", rec.Code, rec.Body.String())
	}
}
