package service

import (
	"context"

	"triangleoflove/backend/internal/repository"
)

// CheckinRepo is the storage interface required by CheckinService.
type CheckinRepo interface {
	FindToday(ctx context.Context, accountID string) (repository.Checkin, error)
	Save(ctx context.Context, accountID string, c repository.Checkin) (repository.Checkin, error)
}

// CheckinService orchestrates loading and saving daily check-ins.
type CheckinService struct {
	repo CheckinRepo
}

func NewCheckinService(repo CheckinRepo) *CheckinService {
	return &CheckinService{repo: repo}
}

// GetToday returns today's check-in for accountID, or ErrCheckinNotFound if none.
func (s *CheckinService) GetToday(ctx context.Context, accountID string) (repository.Checkin, error) {
	return s.repo.FindToday(ctx, accountID)
}

// SaveToday saves today's check-in for accountID and returns the saved record.
func (s *CheckinService) SaveToday(ctx context.Context, accountID string, c repository.Checkin) (repository.Checkin, error) {
	return s.repo.Save(ctx, accountID, c)
}
