package service

import (
	"context"

	"triangleoflove/backend/internal/repository"
)

// CheckinService orchestrates loading and saving daily check-ins.
type CheckinService struct {
	repo *repository.CheckinRepository
}

func NewCheckinService(repo *repository.CheckinRepository) *CheckinService {
	return &CheckinService{repo: repo}
}

// GetToday returns today's check-in for accountID, or ErrCheckinNotFound if none.
func (s *CheckinService) GetToday(ctx context.Context, accountID string) (repository.Checkin, error) {
	return s.repo.FindToday(ctx, accountID)
}

// SaveToday upserts today's check-in for accountID and returns the saved record.
func (s *CheckinService) SaveToday(ctx context.Context, accountID string, c repository.Checkin) (repository.Checkin, error) {
	return s.repo.Upsert(ctx, accountID, c)
}
