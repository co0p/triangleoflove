package service_test

import (
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

type mockAccountRepo struct {
	account     domain.Account
	err         error
	registerErr error
	registered  domain.Account
}

func (m *mockAccountRepo) FindByEmail(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.err
}

func (m *mockAccountRepo) FindByID(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.err
}

func (m *mockAccountRepo) SaveHashedPassword(_ context.Context, _ string, _ string) error {
	return nil
}

func (m *mockAccountRepo) Register(_ context.Context, account domain.Account) error {
	m.registered = account
	return m.registerErr
}

func hashedPassword(t *testing.T, plain string) string {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.MinCost)
	if err != nil {
		t.Fatalf("bcrypt: %v", err)
	}
	return string(h)
}

func TestAuthService_GivenCorrectPassword_WhenLogin_ThenReturnsToken(t *testing.T) {
	repo := &mockAccountRepo{
		account: domain.Account{
			ID:             "acc-1",
			Email:          "alice@example.com",
			HashedPassword: hashedPassword(t, "correct-pass"),
			FirstName:      "Alice",
			Role:           domain.RoleUser,
			IsActive:       true,
		},
	}
	svc := service.NewAuthService(repo)

	result, err := svc.Login(context.Background(), "alice@example.com", "correct-pass")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Token == "" {
		t.Fatal("expected a non-empty token")
	}
}

func TestAuthService_GivenWrongPassword_WhenLogin_ThenReturnsError(t *testing.T) {
	repo := &mockAccountRepo{
		account: domain.Account{
			ID:             "acc-1",
			Email:          "alice@example.com",
			HashedPassword: hashedPassword(t, "correct-pass"),
			FirstName:      "Alice",
		},
	}
	svc := service.NewAuthService(repo)

	_, err := svc.Login(context.Background(), "alice@example.com", "wrong-pass")

	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_GivenUnknownEmail_WhenLogin_ThenReturnsError(t *testing.T) {
	repo := &mockAccountRepo{err: domain.ErrNotFound}
	svc := service.NewAuthService(repo)

	_, err := svc.Login(context.Background(), "nobody@example.com", "any-pass")

	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthService_GivenShortPassword_WhenRegister_ThenValidationError(t *testing.T) {
	repo := &mockAccountRepo{}
	svc := service.NewAuthService(repo)

	err := svc.Register(context.Background(), "new@example.com", "short", "Tester")

	if !errors.Is(err, service.ErrPasswordTooShort) {
		t.Fatalf("expected ErrPasswordTooShort, got %v", err)
	}
}

func TestAuthService_GivenValidRegistration_WhenRegister_ThenCreatesInactiveAccount(t *testing.T) {
	repo := &mockAccountRepo{}
	svc := service.NewAuthService(repo)

	err := svc.Register(context.Background(), "new@example.com", "securepass", "Tester")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.registered.Email != "new@example.com" {
		t.Fatalf("expected registered email to be preserved, got %q", repo.registered.Email)
	}
	if repo.registered.IsActive {
		t.Fatal("expected newly registered account to be inactive")
	}
}

func TestAuthService_GivenInactiveAccount_WhenLogin_ThenRejectedWithClearError(t *testing.T) {
	repo := &mockAccountRepo{
		account: domain.Account{
			ID:             "acc-2",
			Email:          "inactive@example.com",
			HashedPassword: hashedPassword(t, "correct-pass"),
			FirstName:      "Inactive",
			Role:           domain.RoleUser,
			IsActive:       false,
		},
	}
	svc := service.NewAuthService(repo)

	_, err := svc.Login(context.Background(), "inactive@example.com", "correct-pass")

	if !errors.Is(err, service.ErrAccountInactive) {
		t.Fatalf("expected ErrAccountInactive, got %v", err)
	}
}

func TestLogin_GivenInvalidCredentials_WhenLoginAttempted_ThenGenericError(t *testing.T) {
	// wrong password for active account → ErrInvalidCredentials (not ErrAccountInactive)
	repoActive := &mockAccountRepo{
		account: domain.Account{
			ID:             "acc-1",
			Email:          "alice@example.com",
			HashedPassword: hashedPassword(t, "correct-pass"),
			FirstName:      "Alice",
			Role:           domain.RoleUser,
			IsActive:       true,
		},
	}
	svc := service.NewAuthService(repoActive)
	_, err := svc.Login(context.Background(), "alice@example.com", "wrong-pass")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("wrong password: expected ErrInvalidCredentials, got %v", err)
	}

	// unknown email → ErrInvalidCredentials (not ErrAccountInactive)
	repoMissing := &mockAccountRepo{err: domain.ErrNotFound}
	svc2 := service.NewAuthService(repoMissing)
	_, err = svc2.Login(context.Background(), "nobody@example.com", "any-pass")
	if !errors.Is(err, service.ErrInvalidCredentials) {
		t.Fatalf("unknown email: expected ErrInvalidCredentials, got %v", err)
	}
}
