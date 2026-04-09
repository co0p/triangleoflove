package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"time"

	"triangleoflove/backend/internal/repository"
)

// ErrCodeNotFound is returned when the submitted invite code does not match any account.
var ErrCodeNotFound = errors.New("invite code not found")

// ErrAlreadyPaired is returned when either party is already in a couple.
var ErrAlreadyPaired = errors.New("already paired")

const codeCharset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// InviteCodeRepo is the storage interface for invite code operations.
type InviteCodeRepo interface {
	GetCode(ctx context.Context, accountID string) (string, error)
	SetCode(ctx context.Context, accountID, code string) error
	FindAccountByCode(ctx context.Context, code string) (string, error)
	IsAccountPaired(ctx context.Context, accountID string) (bool, error)
}

// CoupleRepo is the storage interface for couple operations.
type CoupleRepo interface {
	CreateCouple(ctx context.Context, accountIDA, accountIDB string) error
	GetCoupleSummary(ctx context.Context, accountID string) (repository.CoupleSummary, bool, error)
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
	code, err := s.codes.GetCode(ctx, accountID)
	if err != nil {
		return "", err
	}
	if code != "" {
		return code, nil
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
	partnerID, err := s.codes.FindAccountByCode(ctx, partnerCode)
	if errors.Is(err, repository.ErrCodeNotFound) {
		return ErrCodeNotFound
	}
	if err != nil {
		return err
	}

	if partnerID == submitterID {
		return ErrCodeNotFound
	}

	submitterPaired, err := s.codes.IsAccountPaired(ctx, submitterID)
	if err != nil {
		return err
	}
	if submitterPaired {
		return ErrAlreadyPaired
	}

	partnerPaired, err := s.codes.IsAccountPaired(ctx, partnerID)
	if err != nil {
		return err
	}
	if partnerPaired {
		return ErrAlreadyPaired
	}

	if err := s.couple.CreateCouple(ctx, submitterID, partnerID); err != nil {
		return err
	}

	_, err = s.issueNewCode(ctx, partnerID)
	return err
}

func (s *PairingService) issueNewCode(ctx context.Context, accountID string) (string, error) {
	code, err := generateCode()
	if err != nil {
		return "", err
	}
	if err := s.codes.SetCode(ctx, accountID, code); err != nil {
		return "", err
	}
	return code, nil
}

// GetCoupleStatus returns whether the caller is paired and, if so, their partner's
// first name and the date the couple formed.
func (s *PairingService) GetCoupleStatus(ctx context.Context, accountID string) (CoupleStatus, error) {
	summary, found, err := s.couple.GetCoupleSummary(ctx, accountID)
	if err != nil {
		return CoupleStatus{}, err
	}
	if !found {
		return CoupleStatus{Paired: false}, nil
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
