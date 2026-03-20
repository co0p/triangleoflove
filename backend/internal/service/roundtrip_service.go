package service

import (
	"context"
	"fmt"
	"time"

	"triangleoflove/backend/internal/repository"
)

type RoundtripService struct {
	repo *repository.RoundtripRepository
}

func NewRoundtripService(repo *repository.RoundtripRepository) *RoundtripService {
	return &RoundtripService{repo: repo}
}

func (s *RoundtripService) Execute(ctx context.Context) (repository.RoundtripResult, error) {
	value := fmt.Sprintf("roundtrip-%d", time.Now().UnixNano())
	return s.repo.InsertAndReturn(ctx, value)
}
