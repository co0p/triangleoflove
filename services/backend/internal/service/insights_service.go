package service

import (
	"context"
	"errors"
	"regexp"
	"time"

	"triangleoflove/backend/internal/domain"
)

// ErrInvalidDate is returned when the date parameter does not match YYYYMMDD format.
var ErrInvalidDate = errors.New("invalid date format")

var validDate = regexp.MustCompile(`^\d{8}$`)

// InsightsRepo is the read-only storage interface required by InsightsService.
type InsightsRepo interface {
	FindByAccountAndDate(ctx context.Context, accountID string, date string) (domain.Checkin, error)
}

// InsightsService computes daily and weekly insight scores from check-in data.
type InsightsService struct {
	repo  InsightsRepo
	clock func() time.Time
}

func NewInsightsService(repo InsightsRepo) *InsightsService {
	return &InsightsService{repo: repo, clock: time.Now}
}

// SetClock overrides the clock used for window calculation. Intended for tests only.
func (s *InsightsService) SetClock(fn func() time.Time) {
	s.clock = fn
}

// GetByDate returns the daily insight scores for the given account and date.
// Returns ErrInvalidDate if date is not in YYYYMMDD format.
func (s *InsightsService) GetByDate(ctx context.Context, accountID string, date string) (domain.DailyInsight, error) {
	if !validDate.MatchString(date) {
		return domain.DailyInsight{}, ErrInvalidDate
	}
	checkin, err := s.repo.FindByAccountAndDate(ctx, accountID, date)
	if err != nil {
		return domain.DailyInsight{}, err
	}
	return domain.NewDailyInsight(checkin), nil
}

// GetWeekly returns 7 daily insight entries ordered oldest→newest (index 0 = 6 days before
// yesterday, index 6 = yesterday). Missing check-ins produce all -1 scores.
func (s *InsightsService) GetWeekly(ctx context.Context, accountID string) ([]domain.WeeklyInsight, error) {
	yesterday := s.clock().UTC().AddDate(0, 0, -1)
	result := make([]domain.WeeklyInsight, 7)
	for i := 0; i < 7; i++ {
		day := yesterday.AddDate(0, 0, -(6 - i))
		date := day.Format("20060102")
		checkin, err := s.repo.FindByAccountAndDate(ctx, accountID, date)
		if errors.Is(err, domain.ErrNotFound) {
			result[i] = domain.WeeklyInsight{Date: date, Intimacy: -1, Commitment: -1, Passion: -1}
			continue
		}
		if err != nil {
			return nil, err
		}
		di := domain.NewDailyInsight(checkin)
		result[i] = domain.WeeklyInsight{Date: date, Intimacy: di.Intimacy, Commitment: di.Commitment, Passion: di.Passion}
	}
	return result, nil
}
