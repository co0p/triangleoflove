package service

import (
	"context"
	"errors"

	"triangleoflove/backend/internal/domain"
)

// ErrForbidden is returned when the caller does not have the required role.
var ErrForbidden = errors.New("forbidden")

// AdminRepo is the storage interface required by AdminService.
type AdminRepo interface {
	FindAll(ctx context.Context) ([]domain.AccountSummary, error)
	SetActive(ctx context.Context, id string, active bool) error
}

// AdminService handles admin operations with role enforcement.
type AdminService struct {
	repo AdminRepo
}

func NewAdminService(repo AdminRepo) *AdminService {
	return &AdminService{repo: repo}
}

// ListUsers returns all accounts. Only callers with RoleAdmin are permitted.
func (s *AdminService) ListUsers(ctx context.Context, callerRole domain.Role) ([]domain.AccountSummary, error) {
	if callerRole != domain.RoleAdmin {
		return nil, ErrForbidden
	}
	return s.repo.FindAll(ctx)
}

// SetActive activates or deactivates an account. Only callers with RoleAdmin are permitted.
func (s *AdminService) SetActive(ctx context.Context, callerRole domain.Role, targetID string, active bool) error {
	if callerRole != domain.RoleAdmin {
		return ErrForbidden
	}
	err := s.repo.SetActive(ctx, targetID, active)
	if errors.Is(err, domain.ErrNotFound) {
		return domain.ErrNotFound
	}
	return err
}
