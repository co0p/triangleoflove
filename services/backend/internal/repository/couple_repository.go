package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// CoupleSummary holds the partner's display name and the date the couple formed.
type CoupleSummary struct {
	PartnerFirstName string
	FormedOn         time.Time
}

// CoupleRepository manages the couples table.
type CoupleRepository struct {
	db *sql.DB
}

func NewCoupleRepository(db *sql.DB) *CoupleRepository {
	return &CoupleRepository{db: db}
}

// CreateCouple inserts a new couple record for the two account IDs.
func (r *CoupleRepository) CreateCouple(ctx context.Context, accountIDA, accountIDB string) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO couples (account_id_a, account_id_b) VALUES ($1, $2)`,
		accountIDA, accountIDB,
	)
	return err
}

// GetCoupleSummary returns the partner's first name and the date the couple formed.
// Returns (summary, true, nil) when paired, (zero, false, nil) when unpaired.
func (r *CoupleRepository) GetCoupleSummary(ctx context.Context, accountID string) (CoupleSummary, bool, error) {
	var s CoupleSummary
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
		return CoupleSummary{}, false, nil
	}
	if err != nil {
		return CoupleSummary{}, false, err
	}
	return s, true, nil
}
