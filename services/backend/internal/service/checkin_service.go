package service

import (
	"context"

	"triangleoflove/backend/internal/domain"
)

// SessionRepo is the storage interface required by SessionService.
type SessionRepo interface {
	FindByAccountAndDate(ctx context.Context, accountID string) (domain.Checkin, error)
	Save(ctx context.Context, accountID string, c domain.Checkin) (domain.Checkin, error)
}

// SessionService orchestrates loading and saving today's reflection session.
type SessionService struct {
	repo SessionRepo
}

func NewSessionService(repo SessionRepo) *SessionService {
	return &SessionService{repo: repo}
}

// GetToday returns today's session for accountID, or domain.ErrNotFound if none.
func (s *SessionService) GetToday(ctx context.Context, accountID string) (domain.Checkin, error) {
	return s.repo.FindByAccountAndDate(ctx, accountID)
}

// SaveToday saves today's session for accountID and returns the saved record.
func (s *SessionService) SaveToday(ctx context.Context, accountID string, c domain.Checkin) (domain.Checkin, error) {
	return s.repo.Save(ctx, accountID, c)
}
