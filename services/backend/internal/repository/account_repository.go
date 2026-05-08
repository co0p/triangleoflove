package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"

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
	var role string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name, role, is_active, created_at FROM accounts WHERE email = $1`,
		email,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName, &role, &a.IsActive, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, domain.ErrNotFound
	}
	a.Role = domain.Role(role)
	return a, err
}

func (r *AccountRepository) FindByID(ctx context.Context, id string) (domain.Account, error) {
	var a domain.Account
	var role string
	err := r.db.QueryRowContext(ctx,
		`SELECT id, email, hashed_password, first_name, role, is_active, created_at FROM accounts WHERE id = $1`,
		id,
	).Scan(&a.ID, &a.Email, &a.HashedPassword, &a.FirstName, &role, &a.IsActive, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Account{}, domain.ErrNotFound
	}
	a.Role = domain.Role(role)
	return a, err
}

func (r *AccountRepository) SaveHashedPassword(ctx context.Context, id string, hashedPassword string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE accounts SET hashed_password = $1 WHERE id = $2`,
		hashedPassword, id,
	)
	return err
}

func (r *AccountRepository) Register(ctx context.Context, account domain.Account) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO accounts (id, email, hashed_password, first_name, role, is_active)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		account.ID, account.Email, account.HashedPassword, account.FirstName, string(account.Role), account.IsActive,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return domain.ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (r *AccountRepository) FindAll(ctx context.Context) ([]domain.AccountSummary, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, email, first_name, role, is_active, created_at FROM accounts ORDER BY created_at ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []domain.AccountSummary
	for rows.Next() {
		var s domain.AccountSummary
		var role string
		if err := rows.Scan(&s.ID, &s.Email, &s.FirstName, &role, &s.IsActive, &s.CreatedAt); err != nil {
			return nil, err
		}
		s.Role = domain.Role(role)
		summaries = append(summaries, s)
	}
	return summaries, rows.Err()
}

func (r *AccountRepository) SetActive(ctx context.Context, id string, active bool) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE accounts SET is_active = $1 WHERE id = $2`,
		active, id,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return domain.ErrNotFound
	}
	return nil
}
