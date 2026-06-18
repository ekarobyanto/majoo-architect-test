package middleware_test

import (
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/middleware"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthMiddleware(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Middleware Suite")
}

var _ = Describe("AuthMiddleware", func() {
	var app *fiber.App
	var cfg *config.Config

	BeforeEach(func() {
		app = fiber.New()
		cfg = &config.Config{
			Auth: config.AuthConfig{
				JWTSecret: "secret",
			},
		}

		app.Get("/protected", middleware.JWTAuth(cfg), func(c *fiber.Ctx) error {
			return c.SendString("ok")
		})
	})

	Context("Missing Authorization header", func() {
		It("should return 401", func() {
			req := httptest.NewRequest("GET", "/protected", nil)
			resp, _ := app.Test(req)
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		})
	})

	Context("Invalid Token", func() {
		It("should return 401 for bad format", func() {
			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer invalid-token")
			resp, _ := app.Test(req)
			Expect(resp.StatusCode).To(Equal(fiber.StatusUnauthorized))
		})
	})

	Context("Valid Token", func() {
		It("should return 200 and set context", func() {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":   "123",
				"email": "test@example.com",
				"exp":   time.Now().Add(time.Hour).Unix(),
			})
			tokenString, _ := token.SignedString([]byte(cfg.Auth.JWTSecret))

			req := httptest.NewRequest("GET", "/protected", nil)
			req.Header.Set("Authorization", "Bearer "+tokenString)
			resp, _ := app.Test(req)
			Expect(resp.StatusCode).To(Equal(fiber.StatusOK))
		})
	})
})
