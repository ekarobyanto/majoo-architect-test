package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	"github.com/user/go-backend-boilerplate/config"
	_ "github.com/user/go-backend-boilerplate/docs/swagger"
	authHandler "github.com/user/go-backend-boilerplate/internal/modules/auth/handler"
	authRouter "github.com/user/go-backend-boilerplate/internal/modules/auth/router"
	healthHandler "github.com/user/go-backend-boilerplate/internal/modules/health/handler"
	healthRouter "github.com/user/go-backend-boilerplate/internal/modules/health/router"
	"github.com/user/go-backend-boilerplate/internal/platform/errors"
)

// NewServer initializes a new fiber server with basic middlewares and routes
func NewServer(
	cfg *config.Config,
	db *sqlx.DB,
	healthHdl *healthHandler.HealthHandler,
	authHdl *authHandler.AuthHandler,
) *Server {
	app := fiber.New(fiber.Config{
		AppName:      "Go Backend Boilerplate",
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

	return &Server{
		App:           app,
		Cfg:           cfg,
		DB:            db,
		HealthHandler: healthHdl,
		AuthHandler:   authHdl,
	}
}
