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
	account domain.Account
	err     error
}

func (m *mockAccountRepo) FindByEmail(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.err
}

func (m *mockAccountRepo) FindByID(_ context.Context, _ string) (domain.Account, error) {
	return m.account, m.err
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
