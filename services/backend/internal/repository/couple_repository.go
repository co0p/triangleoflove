package repository

import (
	"context"
	"database/sql"
	"errors"

	"triangleoflove/backend/internal/domain"
)

// CoupleRepository manages the couples table.
type CoupleRepository struct {
	db *sql.DB
}

func NewCoupleRepository(db *sql.DB) *CoupleRepository {
	return &CoupleRepository{db: db}
}

// Save inserts a new couple record for the two account IDs.
func (r *CoupleRepository) Save(ctx context.Context, accountIDA, accountIDB string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO couples (account_id_a, account_id_b) VALUES ($1, $2)`,
		accountIDA, accountIDB,
	)
	return err
}

// FindByAccountID returns the couple summary for the given account's active couple,
// or domain.ErrNotFound when unpaired or the couple has ended.
func (r *CoupleRepository) FindByAccountID(ctx context.Context, accountID string) (domain.CoupleSummary, error) {
	var s domain.CoupleSummary
	err := r.db.QueryRowContext(ctx, `
		SELECT a.first_name, c.formed_on
		FROM couples c
		JOIN accounts a ON a.id = CASE
			WHEN c.account_id_a = $1 THEN c.account_id_b
			ELSE c.account_id_a
		END
		WHERE (c.account_id_a = $1 OR c.account_id_b = $1)
		  AND c.ended_on IS NULL
	`, accountID).Scan(&s.PartnerFirstName, &s.FormedOn)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.CoupleSummary{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.CoupleSummary{}, err
	}
	return s, nil
}

// Unpair sets ended_on on the active couple record for accountID.
// Returns domain.ErrNotFound if accountID has no active couple.
func (r *CoupleRepository) Unpair(ctx context.Context, accountID string) error {
	res, err := r.db.ExecContext(ctx, `
		UPDATE couples
		SET ended_on = now()
		WHERE (account_id_a = $1 OR account_id_b = $1)
		  AND ended_on IS NULL
	`, accountID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
