package repository

import (
	"context"
	"database/sql"
	"errors"

	"triangleoflove/backend/internal/domain"
)

// InsightsRepository reads check-in records for insight score computation.
type InsightsRepository struct {
	db *sql.DB
}

func NewInsightsRepository(db *sql.DB) *InsightsRepository {
	return &InsightsRepository{db: db}
}

// FindByAccountAndDate returns the check-in for the given account and date (YYYYMMDD), or domain.ErrNotFound.
func (r *InsightsRepository) FindByAccountAndDate(ctx context.Context, accountID string, date string) (domain.Checkin, error) {
	var c domain.Checkin
	err := r.db.QueryRowContext(ctx,
		`SELECT felt_understood, meaningful_sharing, could_count_on_them, effort_for_us,
		        desire, spark, mood, note
		 FROM checkins WHERE account_id = $1 AND date = $2`,
		accountID, date,
	).Scan(&c.FeltUnderstood, &c.MeaningfulSharing, &c.CouldCountOnThem, &c.EffortForUs,
		&c.Desire, &c.Spark, &c.Mood, &c.Note)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Checkin{}, domain.ErrNotFound
	}
	return c, err
}
