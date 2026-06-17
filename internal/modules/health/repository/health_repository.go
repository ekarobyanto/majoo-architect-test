package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/internal/modules/health/domain"
)

type healthRepository struct {
	db *sqlx.DB
}

func NewHealthRepository(db *sqlx.DB) domain.HealthRepository {
	return &healthRepository{
		db: db,
	}
}

func (r *healthRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}
