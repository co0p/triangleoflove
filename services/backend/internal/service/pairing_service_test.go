package service_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"triangleoflove/backend/internal/repository"
	"triangleoflove/backend/internal/service"
)

// mockInviteCodeRepo implements service.InviteCodeRepo.
type mockInviteCodeRepo struct {
	code      string
	partnerID string
	paired    bool
}

func (m *mockInviteCodeRepo) GetCode(_ context.Context, _ string) (string, error) {
	return m.code, nil
}

func (m *mockInviteCodeRepo) SetCode(_ context.Context, _ string, code string) error {
	m.code = code
	return nil
}

func (m *mockInviteCodeRepo) FindAccountByCode(_ context.Context, _ string) (string, error) {
	if m.partnerID == "" {
		return "", repository.ErrCodeNotFound
	}
	return m.partnerID, nil
}

func (m *mockInviteCodeRepo) IsAccountPaired(_ context.Context, _ string) (bool, error) {
	return m.paired, nil
}

// mockCoupleRepo implements service.CoupleRepo.
type mockCoupleRepo struct {
	coupled bool
}

func (m *mockCoupleRepo) CreateCouple(_ context.Context, _, _ string) error {
	m.coupled = true
	return nil
}

func (m *mockCoupleRepo) GetCoupleSummary(_ context.Context, _ string) (repository.CoupleSummary, bool, error) {
	return repository.CoupleSummary{}, false, nil
}

func TestPairingService_Connect_GivenInvalidCode_ThenReturnsErrCodeNotFound(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: ""}
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(codes, couple)

	err := svc.Connect(context.Background(), "submitter-id", "BADCOD")

	if !errors.Is(err, service.ErrCodeNotFound) {
		t.Fatalf("expected ErrCodeNotFound, got %v", err)
	}
	if couple.coupled {
		t.Fatal("expected no couple to be created")
	}
}

func TestPairingService_Connect_GivenSubmitterAlreadyPaired_ThenReturnsErrAlreadyPaired(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: "partner-id", paired: true}
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(codes, couple)

	err := svc.Connect(context.Background(), "submitter-id", "VALIDC")

	if !errors.Is(err, service.ErrAlreadyPaired) {
		t.Fatalf("expected ErrAlreadyPaired, got %v", err)
	}
	if couple.coupled {
		t.Fatal("expected no couple to be created")
	}
}

func TestPairingService_Connect_GivenValidCode_ThenCoupleCreated(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: "partner-id", paired: false}
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(codes, couple)

	err := svc.Connect(context.Background(), "submitter-id", "VALIDC")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !couple.coupled {
		t.Fatal("expected couple to be created")
	}
}

func TestPairingService_Connect_GivenOwnCode_ThenReturnsErrCodeNotFound(t *testing.T) {
	codes := &mockInviteCodeRepo{partnerID: "submitter-id"}
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(codes, couple)

	err := svc.Connect(context.Background(), "submitter-id", "OWNcod")

	if !errors.Is(err, service.ErrCodeNotFound) {
		t.Fatalf("expected ErrCodeNotFound for self-pairing, got %v", err)
	}
}

func TestPairingService_GetOrCreateCode_WhenNoCodeExists_ThenGenerates6CharAlphanumericCode(t *testing.T) {
	codes := &mockInviteCodeRepo{code: ""}
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(codes, couple)

	code, err := svc.GetOrCreateCode(context.Background(), "account-id")

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !regexp.MustCompile(`^[A-Z0-9]{6}$`).MatchString(code) {
		t.Fatalf("expected 6-char uppercase alphanumeric code, got %q", code)
	}
}
