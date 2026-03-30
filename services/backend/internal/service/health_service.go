package service

import "context"

// Pinger is satisfied by *sql.DB.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// HealthService checks DB connectivity.
type HealthService struct {
	db Pinger
}

// NewHealthService returns a HealthService backed by the given Pinger.
func NewHealthService(db Pinger) *HealthService {
	return &HealthService{db: db}
}

// Check returns an error if the DB is unreachable.
func (s *HealthService) Check(ctx context.Context) error {
	return s.db.PingContext(ctx)
}
