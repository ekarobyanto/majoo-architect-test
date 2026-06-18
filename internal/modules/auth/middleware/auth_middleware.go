package middleware

import (
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.Unauthorized("missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.Unauthorized("invalid authorization header format")
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return errors.Unauthorized("invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.Unauthorized("invalid token claims")
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		rolesInterface, ok := claims["roles"].([]interface{})
		var roles []string
		if ok {
			for _, r := range rolesInterface {
				if roleName, ok := r.(string); ok {
					roles = append(roles, roleName)
				}
			}
		}

		c.Locals("user", &domain.UserContext{
			ID:    userID,
			Email: email,
			Roles: roles,
		})

		return c.Next()
	}
}
