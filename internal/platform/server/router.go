package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/config"
	_ "github.com/user/simple-blog/docs/swagger"
	authHandler "github.com/user/simple-blog/internal/modules/auth/handler"
	authRouter "github.com/user/simple-blog/internal/modules/auth/router"
	healthHandler "github.com/user/simple-blog/internal/modules/health/handler"
	healthRouter "github.com/user/simple-blog/internal/modules/health/router"
	postHandler "github.com/user/simple-blog/internal/modules/posts/handler"
	postRouter "github.com/user/simple-blog/internal/modules/posts/router"
	"github.com/user/simple-blog/internal/platform/errors"
)

// NewServer initializes a new fiber server with basic middlewares and routes
func NewServer(
	cfg *config.Config,
	db *sqlx.DB,
	healthHdl *healthHandler.HealthHandler,
	authHdl *authHandler.AuthHandler,
	postHdl *postHandler.PostHandler,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "Simple Blog",
		ErrorHandler: errors.GlobalErrorHandler,
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Register Routes
	healthRouter.RegisterRoutes(app, healthHdl)
	authRouter.RegisterRoutes(app, authHdl)
	postRouter.RegisterRoutes(app, postHdl, cfg)

	return &Server{
		App:           app,
		Cfg:           cfg,
		DB:            db,
		HealthHandler: healthHdl,
		AuthHandler:   authHdl,
		PostHandler:   postHdl,
	}
}
