package middleware_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/platform/errors"
)

var _ = Describe("RoleMiddleware", func() {
	var app *fiber.App

	BeforeEach(func() {
		app = fiber.New(fiber.Config{
			ErrorHandler: errors.GlobalErrorHandler,
		})

		// Mock authentication by injecting a user
		app.Use(func(c *fiber.Ctx) error {
			role := c.Get("X-Mock-Role")
			if role != "" {
				c.Locals("user", &domain.UserContext{
					ID:    "user-1",
					Roles: []string{role},
				})
			}
			return c.Next()
		})

		app.Get("/admin-only", middleware.RequireRoles("admin"), func(c *fiber.Ctx) error {
			return c.SendString("admin access")
		})
		
		app.Get("/editor-or-admin", middleware.RequireRoles("admin", "editor"), func(c *fiber.Ctx) error {
			return c.SendString("editor access")
		})
	})

	It("should return 401 if user is not in context", func() {
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		resp, _ := app.Test(req)
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("should return 403 if user lacks required role", func() {
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		req.Header.Set("X-Mock-Role", "user")
		resp, _ := app.Test(req)
		Expect(resp.StatusCode).To(Equal(http.StatusForbidden))
	})

	It("should return 200 if user has exactly the required role", func() {
		req := httptest.NewRequest(http.MethodGet, "/admin-only", nil)
		req.Header.Set("X-Mock-Role", "admin")
		resp, _ := app.Test(req)
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("should return 200 if user has one of the required roles", func() {
		req := httptest.NewRequest(http.MethodGet, "/editor-or-admin", nil)
		req.Header.Set("X-Mock-Role", "editor")
		resp, _ := app.Test(req)
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})
})
