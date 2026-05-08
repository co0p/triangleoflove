package web_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
	"triangleoflove/backend/internal/web"
)

type mockInsightsWeeklyRepo struct {
	checkins map[string]domain.Checkin
	err      error
}

func (m *mockInsightsWeeklyRepo) FindByAccountAndDate(_ context.Context, _ string, date string) (domain.Checkin, error) {
	if m.err != nil {
		return domain.Checkin{}, m.err
	}
	c, ok := m.checkins[date]
	if !ok {
		return domain.Checkin{}, domain.ErrNotFound
	}
	return c, nil
}

func TestInsightsWeekly_GivenNoAuth_WhenRequested_ThenReturns401(t *testing.T) {
	repo := &mockInsightsWeeklyRepo{checkins: map[string]domain.Checkin{}}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsWeeklyHandler(svc))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rec.Code)
	}
}

func TestInsightsWeekly_GivenAuthenticatedUser_WhenRequested_ThenReturnsSevenDaysOrdered(t *testing.T) {
	// Fix today so the window is deterministic: today=2026-04-28, yesterday=2026-04-27
	fixedNow := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)

	repo := &mockInsightsWeeklyRepo{
		checkins: map[string]domain.Checkin{
			// oldest in window (2026-04-21) has score all 5s → 100
			"20260421": {FeltUnderstood: 5, MeaningfulSharing: 5, CouldCountOnThem: 5, EffortForUs: 5, Desire: 5, Spark: 5},
			// yesterday (2026-04-27) has score all 3s → 50
			"20260427": {FeltUnderstood: 3, MeaningfulSharing: 3, CouldCountOnThem: 3, EffortForUs: 3, Desire: 3, Spark: 3},
		},
	}
	svc := service.NewInsightsService(repo)
	svc.SetClock(func() time.Time { return fixedNow })
	handler := web.Middleware(web.NewInsightsWeeklyHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var result []domain.WeeklyInsight
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(result) != 7 {
		t.Fatalf("expected 7 items, got %d", len(result))
	}

	// Index 0 = oldest (20260421), index 6 = yesterday (20260427)
	if result[0].Date != "20260421" {
		t.Errorf("expected result[0].Date=20260421, got %s", result[0].Date)
	}
	if result[6].Date != "20260427" {
		t.Errorf("expected result[6].Date=20260427, got %s", result[6].Date)
	}
	if result[0].Intimacy != 100 {
		t.Errorf("expected result[0].Intimacy=100, got %d", result[0].Intimacy)
	}
	if result[6].Intimacy != 50 {
		t.Errorf("expected result[6].Intimacy=50, got %d", result[6].Intimacy)
	}
	// A date with no check-in must have all -1
	if result[1].Intimacy != -1 {
		t.Errorf("expected result[1].Intimacy=-1 (no check-in), got %d", result[1].Intimacy)
	}
}

func TestInsightsWeekly_GivenRepoError_WhenRequested_ThenReturns500(t *testing.T) {
	repo := &mockInsightsWeeklyRepo{
		checkins: map[string]domain.Checkin{},
		err:      errors.New("db unavailable"),
	}
	svc := service.NewInsightsService(repo)
	handler := web.Middleware(web.NewInsightsWeeklyHandler(svc))

	token, _ := auth.SignToken("account-123", "user")
	req := httptest.NewRequest(http.MethodGet, "/api/v1/insights", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", rec.Code, rec.Body.String())
	}
}
