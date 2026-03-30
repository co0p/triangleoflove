package repository

import (
	"context"
	"database/sql"
	"errors"
)

// ErrAccountNotFound is returned when no account matches the query.
var ErrAccountNotFound = errors.New("account not found")

// Account holds the fields read from the accounts table.
type Account struct {
	ID             string
	Email          string
	HashedPassword string
	FirstName      string
}

// AccountRepository reads account records from Postgres.
type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (Account, error) {
	var a Account
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name FROM accounts WHERE email = $1`,
		email,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, ErrAccountNotFound
	}
	return a, err
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (Account, error) {
	var a Account
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name FROM accounts WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return Account{}, ErrAccountNotFound
	}
	return a, err
}
