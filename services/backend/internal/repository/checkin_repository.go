package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"triangleoflove/backend/internal/domain"
)

// CheckinRepository reads and writes checkin records.
type CheckinRepository struct {
	db *sql.DB
}

func NewCheckinRepository(db *sql.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// FindByAccountAndDate returns today's check-in for accountID, or domain.ErrNotFound if none exists.
func (r *CheckinRepository) FindByAccountAndDate(ctx context.Context, accountID string) (domain.Checkin, error) {
	today := time.Now().UTC().Format("2006-01-02")
	var c domain.Checkin
	err := r.db.QueryRowContext(ctx,
		`SELECT felt_understood, meaningful_sharing, could_count_on_them, effort_for_us,
		        desire, spark, mood, note
		 FROM checkins WHERE account_id = $1 AND date = $2`,
		accountID, today,
	).Scan(&c.FeltUnderstood, &c.MeaningfulSharing, &c.CouldCountOnThem, &c.EffortForUs,
		&c.Desire, &c.Spark, &c.Mood, &c.Note)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Checkin{}, domain.ErrNotFound
	}
	return c, err
}

// Save creates or replaces today's check-in for accountID.
func (r *CheckinRepository) Save(ctx context.Context, accountID string, c domain.Checkin) (domain.Checkin, error) {
	today := time.Now().UTC().Format("2006-01-02")
	var saved domain.Checkin
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO checkins (account_id, date, felt_understood, meaningful_sharing,
		        could_count_on_them, effort_for_us, desire, spark, mood, note, saved_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now())
		 ON CONFLICT (account_id, date) DO UPDATE
		   SET felt_understood    = EXCLUDED.felt_understood,
		       meaningful_sharing = EXCLUDED.meaningful_sharing,
		       could_count_on_them = EXCLUDED.could_count_on_them,
		       effort_for_us      = EXCLUDED.effort_for_us,
		       desire             = EXCLUDED.desire,
		       spark              = EXCLUDED.spark,
		       mood               = EXCLUDED.mood,
		       note               = EXCLUDED.note,
		       saved_at           = now()
		 RETURNING felt_understood, meaningful_sharing, could_count_on_them, effort_for_us,
		           desire, spark, mood, note`,
		accountID, today,
		c.FeltUnderstood, c.MeaningfulSharing, c.CouldCountOnThem, c.EffortForUs,
		c.Desire, c.Spark, c.Mood, c.Note,
	).Scan(&saved.FeltUnderstood, &saved.MeaningfulSharing, &saved.CouldCountOnThem, &saved.EffortForUs,
		&saved.Desire, &saved.Spark, &saved.Mood, &saved.Note)
	return saved, err
}
