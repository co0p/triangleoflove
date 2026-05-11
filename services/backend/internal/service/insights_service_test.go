package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

type mockInsightsRepoWeekly struct {
	calls    []string // dates called in order
	checkins map[string]domain.Checkin
	err      error
}

func (m *mockInsightsRepoWeekly) FindByAccountAndDate(_ context.Context, _ string, date string) (domain.Checkin, error) {
	m.calls = append(m.calls, date)
	if m.err != nil {
		return domain.Checkin{}, m.err
	}
	c, ok := m.checkins[date]
	if !ok {
		return domain.Checkin{}, domain.ErrNotFound
	}
	return c, nil
}

func TestInsightsService_GetWeekly_GivenSevenDays_WhenReturned_ThenOrderedOldestToNewest(t *testing.T) {
	// Fix "today" to a known date so the window is deterministic.
	// today = 2026-04-27 → yesterday = 2026-04-26 → window: 2026-04-21 .. 2026-04-27 (yesterday)
	fixedNow := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC) // today is 2026-04-28; yesterday = 2026-04-27

	repo := &mockInsightsRepoWeekly{
		checkins: map[string]domain.Checkin{
			"20260422": {FeltUnderstood: 5, MeaningfulSharing: 5, CouldCountOnThem: 5, EffortForUs: 5, Desire: 5, Spark: 5},
			"20260427": {FeltUnderstood: 3, MeaningfulSharing: 3, CouldCountOnThem: 3, EffortForUs: 3, Desire: 3, Spark: 3},
		},
	}
	svc := service.NewInsightsService(repo)
	svc.SetClock(func() time.Time { return fixedNow })

	result, err := svc.GetWeekly(context.Background(), "account-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 7 {
		t.Fatalf("expected 7 days, got %d", len(result))
	}

	// Index 0 must be the oldest date (6 days before yesterday = 2026-04-21)
	if result[0].Date != "20260421" {
		t.Errorf("expected result[0].Date=20260421, got %s", result[0].Date)
	}

	// Index 6 must be yesterday (2026-04-27)
	if result[6].Date != "20260427" {
		t.Errorf("expected result[6].Date=20260427, got %s", result[6].Date)
	}

	// A date with no check-in must have all -1
	if result[0].Intimacy != -1 || result[0].Commitment != -1 || result[0].Passion != -1 {
		t.Errorf("expected all -1 for missing date, got %+v", result[0])
	}

	// A date with a check-in must have computed scores (all 5s → score = 100)
	if result[1].Intimacy != 100 {
		t.Errorf("expected intimacy=100 for 20260422, got %d", result[1].Intimacy)
	}

	// Yesterday (index 6) must have correct scores (all 3s → score = 50)
	if result[6].Intimacy != 50 {
		t.Errorf("expected intimacy=50 for 20260427, got %d", result[6].Intimacy)
	}
}

func TestInsightsService_GetWeekly_GivenRepoError_WhenRequested_ThenPropagatesError(t *testing.T) {
	fixedNow := time.Date(2026, 4, 28, 0, 0, 0, 0, time.UTC)
	repo := &mockInsightsRepoWeekly{err: errors.New("db error")}
	svc := service.NewInsightsService(repo)
	svc.SetClock(func() time.Time { return fixedNow })

	_, err := svc.GetWeekly(context.Background(), "account-1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInsightsService_GetWindow_GivenPast3_WhenToday_ThenReturns3ItemsWithTodayLast(t *testing.T) {
	// today = 2026-05-11 → window: 20260509, 20260510, 20260511
	fixedNow := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	repo := &mockInsightsRepoWeekly{
		checkins: map[string]domain.Checkin{
			"20260511": {FeltUnderstood: 5, MeaningfulSharing: 5, CouldCountOnThem: 5, EffortForUs: 5, Desire: 5, Spark: 5},
		},
	}
	svc := service.NewInsightsService(repo)
	svc.SetClock(func() time.Time { return fixedNow })

	result, err := svc.GetWindow(context.Background(), "account-1", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 3 {
		t.Fatalf("expected 3 items, got %d", len(result))
	}
	// Index 2 must be today
	if result[2].Date != "20260511" {
		t.Errorf("expected result[2].Date=20260511 (today), got %s", result[2].Date)
	}
	// Index 0 must be 2 days ago
	if result[0].Date != "20260509" {
		t.Errorf("expected result[0].Date=20260509, got %s", result[0].Date)
	}
	// Today has a check-in; score should be 100
	if result[2].Intimacy != 100 {
		t.Errorf("expected result[2].Intimacy=100, got %d", result[2].Intimacy)
	}
	// Days without check-ins must have -1
	if result[0].Intimacy != -1 {
		t.Errorf("expected result[0].Intimacy=-1 (no check-in), got %d", result[0].Intimacy)
	}
}

func TestInsightsService_GetWindow_GivenRepoError_WhenRequested_ThenPropagatesError(t *testing.T) {
	fixedNow := time.Date(2026, 5, 11, 0, 0, 0, 0, time.UTC)
	repo := &mockInsightsRepoWeekly{err: errors.New("db error")}
	svc := service.NewInsightsService(repo)
	svc.SetClock(func() time.Time { return fixedNow })

	_, err := svc.GetWindow(context.Background(), "account-1", 3)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
