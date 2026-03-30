package service_test

import (
	"context"
	"errors"
	"testing"

	"triangleoflove/backend/internal/service"
)

type stubPinger struct{ err error }

func (s stubPinger) PingContext(ctx context.Context) error { return s.err }

func TestHealthService_GivenPingerFails_ThenCheckReturnsError(t *testing.T) {
	svc := service.NewHealthService(stubPinger{err: errors.New("connection refused")})
	if err := svc.Check(context.Background()); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestHealthService_GivenPingerSucceeds_ThenCheckReturnsNil(t *testing.T) {
	svc := service.NewHealthService(stubPinger{err: nil})
	if err := svc.Check(context.Background()); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}
