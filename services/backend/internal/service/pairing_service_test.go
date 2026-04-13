package service_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"triangleoflove/backend/internal/domain"
	"triangleoflove/backend/internal/service"
)

// mockInviteCodeRepo implements service.InviteCodeRepo.
type mockInviteCodeRepo struct {
	code      string
	partnerID string
	paired    bool
}

func (m *mockInviteCodeRepo) FindInviteCodeByAccountID(_ context.Context, _ string) (domain.InviteCode, error) {
	return domain.InviteCode(m.code), nil
}

func (m *mockInviteCodeRepo) SaveInviteCode(_ context.Context, _ string, code domain.InviteCode) error {
	m.code = string(code)
	return nil
}

func (m *mockInviteCodeRepo) FindByInviteCode(_ context.Context, _ domain.InviteCode) (string, error) {
	if m.partnerID == "" {
		return "", domain.ErrNotFound
	}
	return m.partnerID, nil
}

func (m *mockInviteCodeRepo) ExistsCoupleByAccountID(_ context.Context, _ string) (bool, error) {
	return m.paired, nil
}

// mockCoupleRepo implements service.CoupleRepo.
type mockCoupleRepo struct {
	coupled      bool
	unpairCalled bool
	unpairErr    error
}

func (m *mockCoupleRepo) Save(_ context.Context, _, _ string) error {
	m.coupled = true
	return nil
}

func (m *mockCoupleRepo) FindByAccountID(_ context.Context, _ string) (domain.CoupleSummary, error) {
	return domain.CoupleSummary{}, domain.ErrNotFound
}

func (m *mockCoupleRepo) Unpair(_ context.Context, _ string) error {
	m.unpairCalled = true
	return m.unpairErr
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

func TestUnpairService_GivenPairedCouple_WhenUnpaired_ThenRecordRetainedWithEndedOn(t *testing.T) {
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(&mockInviteCodeRepo{paired: true}, couple)

	err := svc.Unpair(context.Background(), "account-123")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if !couple.unpairCalled {
		t.Fatal("expected repo Unpair to be called")
	}
}

func TestGetCoupleStatus_GivenEndedCouple_WhenCalled_ThenReturnsPairedFalse(t *testing.T) {
	// mockCoupleRepo.FindByAccountID returns ErrNotFound, simulating the
	// ended_on IS NULL filter excluding the ended couple from the result set.
	couple := &mockCoupleRepo{}
	svc := service.NewPairingService(&mockInviteCodeRepo{}, couple)

	status, err := svc.GetCoupleStatus(context.Background(), "account-123")

	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if status.Paired {
		t.Fatal("expected paired=false for ended couple")
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
