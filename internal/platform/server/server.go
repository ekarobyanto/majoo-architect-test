package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/jmoiron/sqlx"
	"github.com/user/go-backend-boilerplate/config"
	"github.com/user/go-backend-boilerplate/internal/modules/health/handler"
	"github.com/user/go-backend-boilerplate/internal/modules/health/repository"
	"github.com/user/go-backend-boilerplate/internal/modules/health/service"
)

// Server holds the fiber app and dependencies
type Server struct {
	App *fiber.App
	Cfg *config.Config
	DB  *sqlx.DB
}

// NewServer initializes a new fiber server with basic middlewares
func NewServer(cfg *config.Config, db *sqlx.DB) *Server {
	app := fiber.New(fiber.Config{
		AppName: "Go Backend Boilerplate",
	})

	// Middlewares
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New())

	// Initialize Health Module
	healthRepo := repository.NewHealthRepository(db)
	healthSvc := service.NewHealthService(healthRepo)
	handler.RegisterRoutes(app, healthSvc)

	return &Server{
		App: app,
		Cfg: cfg,
		DB:  db,
	}
}

// Start starts the fiber server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.Cfg.Port)
	if s.Cfg.Port == "" {
		addr = ":8080"
	}

	log.Printf("Server starting on %s", addr)
	return s.App.Listen(addr)
}

// Shutdown gracefully shuts down the server with a context
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")
	return s.App.ShutdownWithContext(ctx)
}
