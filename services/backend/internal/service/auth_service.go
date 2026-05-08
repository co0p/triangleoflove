package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"triangleoflove/backend/internal/auth"
	"triangleoflove/backend/internal/domain"
)

func newUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ErrInvalidCredentials is returned when email/password do not match.
var ErrInvalidCredentials = errors.New("invalid credentials")

// ErrAccountInactive is returned when a valid but inactive account attempts to log in.
var ErrAccountInactive = errors.New("account inactive")

// ErrPasswordTooShort is returned when a registration password is shorter than 8 characters.
var ErrPasswordTooShort = errors.New("password too short")

// ErrEmailAlreadyExists is returned when a registration email is already taken.
var ErrEmailAlreadyExists = errors.New("email already exists")

// AccountRepo is the storage interface required by AuthService.
type AccountRepo interface {
	FindByEmail(ctx context.Context, email string) (domain.Account, error)
	FindByID(ctx context.Context, id string) (domain.Account, error)
	SaveHashedPassword(ctx context.Context, id string, hashedPassword string) error
	Register(ctx context.Context, account domain.Account) error
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
	Email     string `json:"email"`
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

	if !account.IsActive {
		return LoginResult{}, ErrAccountInactive
	}

	token, err := auth.SignToken(account.ID, string(account.Role))
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{Token: token}, nil
}

// Register creates a new account with role user and active status.
func (s *AuthService) Register(ctx context.Context, email, password, firstName string) error {
	if len(password) < 8 {
		return ErrPasswordTooShort
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account := domain.Account{
		ID:             newUUID(),
		Email:          email,
		HashedPassword: string(hashed),
		FirstName:      firstName,
		Role:           domain.RoleUser,
		IsActive:       true,
	}

	if err := s.accounts.Register(ctx, account); err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			return ErrEmailAlreadyExists
		}
		return err
	}
	return nil
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

	return ProfileResult{FirstName: account.FirstName, Email: account.Email}, nil
}

// ChangePassword verifies the current password and replaces it with a new bcrypt hash.
func (s *AuthService) ChangePassword(ctx context.Context, accountID, currentPassword, newPassword string) error {
	account, err := s.accounts.FindByID(ctx, accountID)
	if errors.Is(err, domain.ErrNotFound) {
		return ErrInvalidCredentials
	}
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.HashedPassword), []byte(currentPassword)); err != nil {
		return ErrInvalidCredentials
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.accounts.SaveHashedPassword(ctx, accountID, string(newHash))
}
