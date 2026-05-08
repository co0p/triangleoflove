package web_test

import (
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

type mockInsightsRepo struct {
	checkin domain.Checkin
	err     error
}

func (m *mockInsightsRepo) FindByAccountAndDate(_ context.Context, _ string, _ string) (domain.Checkin, error) {
	return m.checkin, m.err
}

func TestInsights_GivenNoAuth_WhenNavigated_ThenRedirectsToLogin(t *testing.T) {
	repo := &mockInsightsRepo{}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights/20260415", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestInsights_GivenCheckinExists_WhenRequested_ThenReturnsThreeScores(t *testing.T) {
	repo := &mockInsightsRepo{
		checkin: domain.Checkin{
			FeltUnderstood:    4,
			MeaningfulSharing: 4,
			CouldCountOnThem:  3,
			EffortForUs:       5,
			Desire:            2,
			Spark:             4,
		},
	}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights/20260415", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("date", "20260415")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result domain.DailyInsight
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Both metrics are 4: avg=4, normalized=(4-1)/4*100=75
	if result.Intimacy != 75 {
		t.Errorf("expected intimacy=75, got %d", result.Intimacy)
	}
	// Metrics are 3 and 5: avg=4, normalized=(4-1)/4*100=75
	if result.Commitment != 75 {
		t.Errorf("expected commitment=75, got %d", result.Commitment)
	}
	// Metrics are 2 and 4: avg=3, normalized=(3-1)/4*100=50
	if result.Passion != 50 {
		t.Errorf("expected passion=50, got %d", result.Passion)
	}
}

func TestInsights_GivenUnpairedUser_WhenRequested_ThenReturnsScores(t *testing.T) {
	// This test proves that the insights endpoint does not check pairing status.
	// Any authenticated user (paired or not) gets scores.
	repo := &mockInsightsRepo{
		checkin: domain.Checkin{
			FeltUnderstood:    5,
			MeaningfulSharing: 5,
			CouldCountOnThem:  5,
			EffortForUs:       5,
			Desire:            5,
			Spark:             5,
		},
	}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsHandler(svc))

	token, _ := auth.SignToken("unpaired-account-456", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights/20260415", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("date", "20260415")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result domain.DailyInsight
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// All metrics are 5: normalized=(5-1)/4*100=100
	if result.Intimacy != 100 {
		t.Errorf("expected intimacy=100, got %d", result.Intimacy)
	}
	if result.Commitment != 100 {
		t.Errorf("expected commitment=100, got %d", result.Commitment)
	}
	if result.Passion != 100 {
		t.Errorf("expected passion=100, got %d", result.Passion)
	}
}

func TestInsights_GivenNoCheckin_WhenRequested_ThenNoDataReturned(t *testing.T) {
	repo := &mockInsightsRepo{err: domain.ErrNotFound}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights/20260416", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("date", "20260416")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["error"] != "no check-in for this date" {
		t.Fatalf("expected no-data error message, got %q", body["error"])
	}
}

func TestInsights_GivenInvalidDate_WhenRequested_ThenErrorReturned(t *testing.T) {
	repo := &mockInsightsRepo{}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights/2026-04-16", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.SetPathValue("date", "2026-04-16")

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", rec.Code, rec.Body.String())
	}

	var body map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body["error"] != "invalid date format" {
		t.Fatalf("expected invalid-date error message, got %q", body["error"])
	}
}
