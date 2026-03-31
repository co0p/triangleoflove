package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// ErrCheckinNotFound is returned when no check-in exists for the given account and date.
var ErrCheckinNotFound = errors.New("checkin not found")

// Checkin holds a single daily check-in record.
type Checkin struct {
	FeltClose            int    `json:"felt_close"`
	PositiveEnergy       int    `json:"positive_energy"`
	Supported            int    `json:"supported"`
	CommunicationHealthy int    `json:"communication_healthy"`
	StressLevel          int    `json:"stress_level"`
	Note                 string `json:"note"`
}

// CheckinRepository reads and writes checkin records.
type CheckinRepository struct {
	db *sql.DB
}

func NewCheckinRepository(db *sql.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// FindToday returns today's check-in for accountID, or ErrCheckinNotFound if none exists.
func (r *CheckinRepository) FindToday(ctx context.Context, accountID string) (Checkin, error) {
	today := time.Now().UTC().Format("2006-01-02")
	var c Checkin
	err := r.db.QueryRowContext(ctx,
		`SELECT felt_close, positive_energy, supported, communication_healthy, stress_level, note
		 FROM checkins WHERE account_id = $1 AND date = $2`,
		accountID, today,
	).Scan(&c.FeltClose, &c.PositiveEnergy, &c.Supported, &c.CommunicationHealthy, &c.StressLevel, &c.Note)
	if errors.Is(err, sql.ErrNoRows) {
		return Checkin{}, ErrCheckinNotFound
	}
	return c, err
}

// Upsert creates or replaces today's check-in for accountID.
func (r *CheckinRepository) Upsert(ctx context.Context, accountID string, c Checkin) (Checkin, error) {
	today := time.Now().UTC().Format("2006-01-02")
	var saved Checkin
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO checkins (account_id, date, felt_close, positive_energy, supported, communication_healthy, stress_level, note, saved_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now())
		 ON CONFLICT (account_id, date) DO UPDATE
		   SET felt_close = EXCLUDED.felt_close,
		       positive_energy = EXCLUDED.positive_energy,
		       supported = EXCLUDED.supported,
		       communication_healthy = EXCLUDED.communication_healthy,
		       stress_level = EXCLUDED.stress_level,
		       note = EXCLUDED.note,
		       saved_at = now()
		 RETURNING felt_close, positive_energy, supported, communication_healthy, stress_level, note`,
		accountID, today,
		c.FeltClose, c.PositiveEnergy, c.Supported, c.CommunicationHealthy, c.StressLevel, c.Note,
	).Scan(&saved.FeltClose, &saved.PositiveEnergy, &saved.Supported, &saved.CommunicationHealthy, &saved.StressLevel, &saved.Note)
	return saved, err
}
