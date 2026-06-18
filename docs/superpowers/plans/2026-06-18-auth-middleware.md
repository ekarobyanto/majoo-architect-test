# Auth Middleware Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement JWT middleware to protect routes.

**Architecture:** Middleware extracts Bearer token, validates via `golang-jwt/jwt/v5`, stores `domain.UserContext` in Fiber context locals.

**Tech Stack:** Go, Fiber, JWT.

---

### Task 1: Scaffolding and Failing Tests

**Files:**
- Create: `internal/modules/auth/middleware/auth_middleware_test.go`
- Create: `internal/modules/auth/middleware/auth_middleware.go`

- [ ] **Step 1: Create empty middleware**

```go
package middleware

import (
	"github.com/gofiber/fiber/v2"
	"majoo-architect-test/config"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Next()
	}
}
```

- [ ] **Step 2: Write failing tests**

```go
package middleware_test

import (
	"majoo-architect-test/config"
	"majoo-architect-test/internal/modules/auth/middleware"
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
```

- [ ] **Step 3: Run test to verify it fails**

Run: `go test ./internal/modules/auth/middleware/...`
Expected: FAIL (missing 401s)

- [ ] **Step 4: Commit**

```bash
git add internal/modules/auth/middleware/
git commit -m "test(auth): add failing tests for JWT middleware"
```

### Task 2: Implement JWTAuth Middleware

**Files:**
- Modify: `internal/modules/auth/middleware/auth_middleware.go`

- [ ] **Step 1: Implement extraction and validation**

```go
package middleware

import (
	"majoo-architect-test/config"
	"majoo-architect-test/internal/modules/auth/domain"
	"majoo-architect-test/internal/platform/errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(cfg *config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.NewUnauthorizedError("missing authorization header")
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return errors.NewUnauthorizedError("invalid authorization header format")
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(cfg.Auth.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return errors.NewUnauthorizedError("invalid or expired token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return errors.NewUnauthorizedError("invalid token claims")
		}

		userID, _ := claims["sub"].(string)
		email, _ := claims["email"].(string)

		c.Locals("user", &domain.UserContext{
			ID:    userID,
			Email: email,
		})

		return c.Next()
	}
}
```

- [ ] **Step 2: Run test to verify it passes**

Run: `go test ./internal/modules/auth/middleware/...`
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add internal/modules/auth/middleware/auth_middleware.go
git commit -m "feat(auth): implement JWT authentication middleware"
```

### Task 3: Wrap Up

- [ ] **Step 1: Update TASKS.md**
- [ ] **Step 2: Verify overall suite**

Run: `go test ./...`
