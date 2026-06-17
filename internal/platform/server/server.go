package server

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/config"
	authHandler "github.com/user/simple-blog/internal/modules/auth/handler"
	healthHandler "github.com/user/simple-blog/internal/modules/health/handler"
)

// Server holds the fiber app and dependencies
type Server struct {
	App           *fiber.App
	Cfg           *config.Config
	DB            *sqlx.DB
	HealthHandler *healthHandler.HealthHandler
	AuthHandler   *authHandler.AuthHandler
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
