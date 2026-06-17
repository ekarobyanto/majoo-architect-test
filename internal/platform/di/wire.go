//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/user/go-backend-boilerplate/config"
	"github.com/user/go-backend-boilerplate/internal/modules/auth"
	"github.com/user/go-backend-boilerplate/internal/modules/health"
	"github.com/user/go-backend-boilerplate/internal/platform/server"
)

func InitializeServer(cfg *config.Config, db *sqlx.DB) *server.Server {
	wire.Build(
		health.ProviderSet,
		auth.ProviderSet,
		server.NewServer,
	)
	return nil
}
