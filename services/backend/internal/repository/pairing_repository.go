package repository

import (
	"context"
	"database/sql"
	"errors"
)

// ErrCodeNotFound is returned when no account holds the given invite code.
var ErrCodeNotFound = errors.New("invite code not found")

// PairingRepository manages invite codes on the accounts table.
type PairingRepository struct {
	db *sql.DB
}

func NewPairingRepository(db *sql.DB) *PairingRepository {
	return &PairingRepository{db: db}
}

// GetCode returns the stored invite_code for the given accountID, or "" if null.
func (r *PairingRepository) GetCode(ctx context.Context, accountID string) (string, error) {
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
	return code.String, nil
}

// SetCode persists invite_code for the given accountID.
func (r *PairingRepository) SetCode(ctx context.Context, accountID, code string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE accounts SET invite_code = $1 WHERE id = $2`,
		code, accountID,
	)
	return err
}

// FindAccountByCode returns the accountID whose invite_code matches, or ErrCodeNotFound.
func (r *PairingRepository) FindAccountByCode(ctx context.Context, code string) (string, error) {
	var accountID string
	err := r.db.QueryRowContext(ctx,
		`SELECT id FROM accounts WHERE invite_code = $1`,
		code,
	).Scan(&accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrCodeNotFound
	}
	return accountID, err
}

// IsAccountPaired returns true if the accountID appears in any couple row.
func (r *PairingRepository) IsAccountPaired(ctx context.Context, accountID string) (bool, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM couples WHERE account_id_a = $1 OR account_id_b = $1`,
		accountID,
	).Scan(&count)
	return count > 0, err
}
