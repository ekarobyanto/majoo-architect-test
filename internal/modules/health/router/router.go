package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/go-backend-boilerplate/internal/modules/health/handler"
)

func RegisterRoutes(router fiber.Router, h *handler.HealthHandler) {
	router.Get("/health", h.CheckHealth)
}
