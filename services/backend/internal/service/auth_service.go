package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
)

// ErrInvalidCredentials is returned when email/password do not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// AccountRepo is the storage interface required by AuthService.
type AccountRepo interface {
	FindByEmail(ctx context.Context, email string) (domain.Account, error)
	FindByID(ctx context.Context, id string) (domain.Account, error)
}

// AuthService handles login and profile retrieval.
type AuthService struct {
	accounts AccountRepo
}

func NewAuthService(accounts AccountRepo) *AuthService {
	return &AuthService{accounts: accounts}
}

// LoginResult is the response payload for a successful login.
type LoginResult struct {
	Token string `json:"token"`
}

// ProfileResult is the response payload for GET /users/me.
type ProfileResult struct {
	FirstName string `json:"firstName"`
}

// Login validates credentials and returns a signed JWT on success.
func (s *AuthService) Login(ctx context.Context, email, password string) (LoginResult, error) {
	account, err := s.accounts.FindByEmail(ctx, email)
	if errors.Is(err, domain.ErrNotFound) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	token, err := auth.SignToken(account.ID)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token}, nil
}

// GetProfile returns the profile for the given account ID.
func (s *AuthService) GetProfile(ctx context.Context, accountID string) (ProfileResult, error) {
	account, err := s.accounts.FindByID(ctx, accountID)
	if errors.Is(err, domain.ErrNotFound) {
		return ProfileResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return ProfileResult{}, err
	}

	return ProfileResult{FirstName: account.FirstName}, nil
}
