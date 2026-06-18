package domain

import "context"

type HealthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type HealthService interface {
	Check(ctx context.Context) (HealthResponse, error)
}

type HealthRepository interface {
	Ping(ctx context.Context) error
}
