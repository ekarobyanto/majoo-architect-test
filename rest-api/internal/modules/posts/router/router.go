package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/config"
	authMiddleware "github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/modules/posts/handler"
)

func RegisterRoutes(router fiber.Router, h *handler.PostHandler, cfg *config.Config) {
	posts := router.Group("/posts")

	// Public routes
	posts.Get("/", h.GetPaginated)
	posts.Get("/:id", h.GetByID)

	// Protected routes
	posts.Use(authMiddleware.JWTAuth(cfg))
	posts.Post("/", h.Create)
	posts.Put("/:id", h.Update)
	posts.Delete("/:id", h.Delete)
}
