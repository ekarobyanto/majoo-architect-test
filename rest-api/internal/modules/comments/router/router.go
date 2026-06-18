package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/config"
	authMiddleware "github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/modules/comments/handler"
)

func RegisterRoutes(router fiber.Router, h *handler.CommentHandler, cfg *config.Config) {
	// Protected routes
	auth := authMiddleware.JWTAuth(cfg)

	// POST /posts/:id/comments
	router.Post("/posts/:id/comments", auth, h.Create)

	comments := router.Group("/comments", auth)
	comments.Put("/:id", h.Update)
	comments.Delete("/:id", h.Delete)
}
