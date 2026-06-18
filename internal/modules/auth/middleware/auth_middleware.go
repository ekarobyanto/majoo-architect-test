package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/config"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
