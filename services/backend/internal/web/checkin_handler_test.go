package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/repository"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

type mockCheckinRepo struct {
	saved repository.Checkin
}

func (m *mockCheckinRepo) FindToday(_ context.Context, _ string) (repository.Checkin, error) {
	return repository.Checkin{}, repository.ErrCheckinNotFound
}

func (m *mockCheckinRepo) Save(_ context.Context, _ string, c repository.Checkin) (repository.Checkin, error) {
	m.saved = c
	return c, nil
}

func TestCheckinHandler_GivenValidPUT_ThenSavesBody(t *testing.T) {
	val := 3
	payload := repository.Checkin{FeltClose: &val}

	mock := &mockCheckinRepo{}
	svc := service.NewCheckinService(mock)
	handler := web.Middleware(web.NewCheckinHandler(svc))

	body, _ := json.Marshal(payload)
	token, _ := auth.SignToken("account-123")

	req := httptest.NewRequest(http.MethodPut, "/api/v1/checkins/today", bytes.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if mock.saved.FeltClose == nil || *mock.saved.FeltClose != 3 {
		t.Fatalf("expected saved FeltClose=3, got %+v", mock.saved.FeltClose)
	}
}
