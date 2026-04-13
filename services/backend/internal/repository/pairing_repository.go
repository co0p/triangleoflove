package repository

import (
	"context"
	"database/sql"
	"errors"

	"triangleoflove/backend/internal/domain"
)

// PairingRepository manages invite codes on the accounts table.
type PairingRepository struct {
	db *sql.DB
}

func NewPairingRepository(db *sql.DB) *PairingRepository {
	return &PairingRepository{db: db}
}

// FindInviteCodeByAccountID returns the stored invite code for the given accountID, or "" if none.
func (r *PairingRepository) FindInviteCodeByAccountID(ctx context.Context, accountID string) (domain.InviteCode, error) {
	var code sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT invite_code FROM accounts WHERE id = $1`,
		accountID,
	).Scan(&code)
	if err != nil {
		return "", err
	}
	if !code.Valid {
		return "", nil
	}
	return domain.InviteCode(code.String), nil
}

// SaveInviteCode persists the invite code for the given accountID.
func (r *PairingRepository) SaveInviteCode(ctx context.Context, accountID string, code domain.InviteCode) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE accounts SET invite_code = $1 WHERE id = $2`,
		string(code), accountID,
	)
	return err
}

// FindByInviteCode returns the accountID whose invite code matches, or domain.ErrNotFound.
func (r *PairingRepository) FindByInviteCode(ctx context.Context, code domain.InviteCode) (string, error) {
	var accountID string
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM accounts WHERE invite_code = $1`,
		string(code),
	).Scan(&accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	return accountID, err
}

// ExistsCoupleByAccountID returns true if the accountID appears in any active couple row.
func (r *PairingRepository) ExistsCoupleByAccountID(ctx context.Context, accountID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM couples WHERE (account_id_a = $1 OR account_id_b = $1) AND ended_on IS NULL`,
		accountID,
	).Scan(&count)
	return count > 0, err
}
