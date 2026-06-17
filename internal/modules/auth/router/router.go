package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/internal/modules/auth/handler"
)

func RegisterRoutes(router fiber.Router, h *handler.AuthHandler) {
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
}
