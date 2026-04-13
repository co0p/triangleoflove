package repository

import (
	"context"
	"database/sql"
	"errors"

	"triangleoflove/backend/internal/domain"
)

// AccountRepository reads account records from Postgres.
type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name FROM accounts WHERE email = $1`,
		email,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, domain.ErrNotFound
	}
	return a, err
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (domain.Account, error) {
	var a domain.Account
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name FROM accounts WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, domain.ErrNotFound
	}
	return a, err
}
