package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"triangleoflove/backend/internal/domain"
)

// ErrCodeNotFound is returned when the submitted invite code does not match any account.
var ErrCodeNotFound = errors.New("invite code not found")

// ErrAlreadyPaired is returned when either party is already in a couple.
var ErrAlreadyPaired = errors.New("already paired")

// ErrNotPaired is returned when an Unpair is attempted but the account has no active couple.
var ErrNotPaired = errors.New("not paired")

const codeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// InviteCodeRepo is the storage interface for invite code operations.
type InviteCodeRepo interface {
	FindInviteCodeByAccountID(ctx context.Context, accountID string) (domain.InviteCode, error)
	SaveInviteCode(ctx context.Context, accountID string, code domain.InviteCode) error
	FindByInviteCode(ctx context.Context, code domain.InviteCode) (string, error)
	ExistsCoupleByAccountID(ctx context.Context, accountID string) (bool, error)
}

// CoupleRepo is the storage interface for couple operations.
type CoupleRepo interface {
	Save(ctx context.Context, accountIDA, accountIDB string) error
	FindByAccountID(ctx context.Context, accountID string) (domain.CoupleSummary, error)
	Unpair(ctx context.Context, accountID string) error
}

// CoupleStatus describes whether the caller is paired and, if so, with whom.
type CoupleStatus struct {
	Paired           bool
	PartnerFirstName string
	PairedSince      time.Time
}

// PairingService orchestrates invite code generation, retrieval, and couple formation.
type PairingService struct {
	codes  InviteCodeRepo
	couple CoupleRepo
}

func NewPairingService(codes InviteCodeRepo, couple CoupleRepo) *PairingService {
	return &PairingService{codes: codes, couple: couple}
}

// GetOrCreateCode returns the caller's current invite code, generating and
// persisting one if none exists yet.
func (s *PairingService) GetOrCreateCode(ctx context.Context, accountID string) (string, error) {
	code, err := s.codes.FindInviteCodeByAccountID(ctx, accountID)
	if err != nil {
		return "", err
	}
	if code != "" {
		return string(code), nil
	}
	return s.issueNewCode(ctx, accountID)
}

// RegenerateCode replaces the caller's invite code with a newly generated one
// and returns it. The previous code is immediately invalid.
func (s *PairingService) RegenerateCode(ctx context.Context, accountID string) (string, error) {
	return s.issueNewCode(ctx, accountID)
}

// Connect forms a couple between the submitter and the account identified by
// partnerCode. Returns ErrCodeNotFound if no account holds that code,
// ErrAlreadyPaired if either party is already in a couple.
// After a successful pairing the partner's invite code is replaced.
func (s *PairingService) Connect(ctx context.Context, submitterID, partnerCode string) error {
	partnerID, err := s.codes.FindByInviteCode(ctx, domain.InviteCode(partnerCode))
	if errors.Is(err, domain.ErrNotFound) {
		return ErrCodeNotFound
	}
	if err != nil {
		return err
	}

	if partnerID == submitterID {
		return ErrCodeNotFound
	}

	submitterPaired, err := s.codes.ExistsCoupleByAccountID(ctx, submitterID)
	if err != nil {
		return err
	}
	if submitterPaired {
		return ErrAlreadyPaired
	}

	partnerPaired, err := s.codes.ExistsCoupleByAccountID(ctx, partnerID)
	if err != nil {
		return err
	}
	if partnerPaired {
		return ErrAlreadyPaired
	}

	if err := s.couple.Save(ctx, submitterID, partnerID); err != nil {
		return err
	}

	_, err = s.issueNewCode(ctx, partnerID)
	return err
}

// Unpair ends the active couple for accountID by setting ended_on on the couple record.
// Returns ErrNotPaired if accountID is not in any active couple.
func (s *PairingService) Unpair(ctx context.Context, accountID string) error {
	err := s.couple.Unpair(ctx, accountID)
	if errors.Is(err, domain.ErrNotFound) {
		return ErrNotPaired
	}
	return err
}

func (s *PairingService) issueNewCode(ctx context.Context, accountID string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	if err := s.codes.SaveInviteCode(ctx, accountID, domain.InviteCode(code)); err != nil {
		return "", err
	}
	return code, nil
}

// GetCoupleStatus returns whether the caller is paired and, if so, their partner's
// first name and the date the couple formed.
func (s *PairingService) GetCoupleStatus(ctx context.Context, accountID string) (CoupleStatus, error) {
	summary, err := s.couple.FindByAccountID(ctx, accountID)
	if errors.Is(err, domain.ErrNotFound) {
		return CoupleStatus{Paired: false}, nil
	}
	if err != nil {
		return CoupleStatus{}, err
	}
	return CoupleStatus{
		Paired:           true,
		PartnerFirstName: summary.PartnerFirstName,
		PairedSince:      summary.FormedOn,
	}, nil
}

func generateCode() (string, error) {
	b := make([]byte, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeCharset))))
		if err != nil {
			return "", err
		}
		b[i] = codeCharset[n.Int64()]
	}
	return string(b), nil
}
