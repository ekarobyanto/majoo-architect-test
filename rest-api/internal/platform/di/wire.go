//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth"
	"github.com/user/simple-blog/internal/modules/comments"
	"github.com/user/simple-blog/internal/modules/health"
	"github.com/user/simple-blog/internal/modules/posts"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/internal/platform/server"
)

func InitializeServer(cfg *config.Config, db *sqlx.DB) *server.Server {
	wire.Build(
		health.ProviderSet,
		auth.ProviderSet,
		posts.ProviderSet,
		comments.ProviderSet,
		database.NewTransactor,
		server.NewServer,
	)
	return nil
}
