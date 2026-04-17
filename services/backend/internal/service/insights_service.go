package service

import (
	"context"
	"errors"
	"regexp"

	"triangleoflove/backend/internal/domain"
)

// ErrInvalidDate is returned when the date parameter does not match YYYYMMDD format.
var ErrInvalidDate = errors.New("invalid date format")

var validDate = regexp.MustCompile(`^\d{8}$`)

// InsightsRepo is the read-only storage interface required by InsightsService.
type InsightsRepo interface {
	FindByAccountAndDate(ctx context.Context, accountID string, date string) (domain.Checkin, error)
}

// InsightsService computes daily insight scores from check-in data.
type InsightsService struct {
	repo InsightsRepo
}

func NewInsightsService(repo InsightsRepo) *InsightsService {
	return &InsightsService{repo: repo}
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
