package service

import (
	"context"

	"github.com/user/simple-blog/internal/modules/health/domain"
)

type healthService struct {
	repo domain.HealthRepository
}

func NewHealthService(repo domain.HealthRepository) domain.HealthService {
	return &healthService{
		repo: repo,
	}
}

func (s *healthService) Check(ctx context.Context) (domain.HealthResponse, error) {
	if err := s.repo.Ping(ctx); err != nil {
		return domain.HealthResponse{
			Status:  "DOWN",
			Message: "Database connection failed",
		}, err
	}

	return domain.HealthResponse{
		Status:  "UP",
		Message: "Database connection is healthy",
	}, nil
}
