package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/errors"
)

// RequireRoles restricts access to users possessing at least one of the specified roles
func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*domain.UserContext)
		if !ok || user == nil {
			return errors.Unauthorized("Authentication required")
		}

		hasRole := false
		for _, userRole := range user.Roles {
			for _, allowedRole := range allowedRoles {
				if userRole == allowedRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			return errors.Forbidden("Insufficient permissions")
		}

		return c.Next()
	}
}
