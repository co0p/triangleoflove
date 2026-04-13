package service

import (
	"context"

	"triangleoflove/backend/internal/domain"
)

// CheckinRepo is the storage interface required by CheckinService.
type CheckinRepo interface {
	FindByAccountAndDate(ctx context.Context, accountID string) (domain.Checkin, error)
	Save(ctx context.Context, accountID string, c domain.Checkin) (domain.Checkin, error)
}

// CheckinService orchestrates loading and saving daily check-ins.
type CheckinService struct {
	repo CheckinRepo
}

func NewCheckinService(repo CheckinRepo) *CheckinService {
	return &CheckinService{repo: repo}
}

// GetToday returns today's check-in for accountID, or domain.ErrNotFound if none.
func (s *CheckinService) GetToday(ctx context.Context, accountID string) (domain.Checkin, error) {
	return s.repo.FindByAccountAndDate(ctx, accountID)
}

// SaveToday saves today's check-in for accountID and returns the saved record.
func (s *CheckinService) SaveToday(ctx context.Context, accountID string, c domain.Checkin) (domain.Checkin, error) {
	return s.repo.Save(ctx, accountID, c)
}
