package repository

import (
	"context"
	"database/sql"
)

type RoundtripResult struct {
	ID    int64  `json:"id"`
	Value string `json:"value"`
}

type RoundtripRepository struct {
	db *sql.DB
}

func NewRoundtripRepository(db *sql.DB) *RoundtripRepository {
	return &RoundtripRepository{db: db}
}

func (r *RoundtripRepository) InsertAndReturn(ctx context.Context, value string) (RoundtripResult, error) {
	if _, err := r.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS demo_roundtrip (
			id BIGSERIAL PRIMARY KEY,
			value TEXT NOT NULL
		)
	`); err != nil {
		return RoundtripResult{}, err
	}

	result := RoundtripResult{}
	err := r.db.QueryRowContext(ctx,
		"INSERT INTO demo_roundtrip (value) VALUES ($1) RETURNING id, value",
		value,
	).Scan(&result.ID, &result.Value)

	if err != nil {
		return RoundtripResult{}, err
	}

	return result, nil
}
