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

// FindByAccountID returns the couple summary for the given account, or domain.ErrNotFound when unpaired.
func (r *CoupleRepository) FindByAccountID(ctx context.Context, accountID string) (domain.CoupleSummary, error) {
	var s domain.CoupleSummary
	err := r.db.QueryRowContext(ctx, `
		SELECT a.first_name, c.formed_on
		FROM couples c
		JOIN accounts a ON a.id = CASE
			WHEN c.account_id_a = $1 THEN c.account_id_b
			ELSE c.account_id_a
		END
		WHERE c.account_id_a = $1 OR c.account_id_b = $1
	`, accountID).Scan(&s.PartnerFirstName, &s.FormedOn)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.CoupleSummary{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.CoupleSummary{}, err
	}
	return s, nil
}
