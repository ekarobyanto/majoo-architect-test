# Authentication Login Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the user login functionality, verifying credentials and returning a JWT token for subsequent authenticated requests.

**Architecture:** We will implement the login flow across the existing domain, repository, service, and handler layers. We will introduce JWT generation in the service layer and configure the necessary secrets via environment variables. Tests will follow the established pattern: testify/mock for unit tests and Ginkgo for handler/integration tests.

**Tech Stack:** Go 1.25, Fiber v2, PostgreSQL (sqlx), bcrypt, golang-jwt/jwt/v5, Ginkgo/Gomega, testify/mock.

---

### Task 1: Add JWT Dependency and Configuration

**Files:**
- Modify: `go.mod`
- Modify: `config/config.go`
- Modify: `.env.example`

- [ ] **Step 1: Install JWT dependency**

Run: `go get github.com/golang-jwt/jwt/v5`
Expected: go.mod is updated.

- [ ] **Step 2: Update Configuration Structure**

Modify `config/config.go` to add JWT configuration:

```go
package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Port          string   `mapstructure:"PORT"`
	JWTSecret     string   `mapstructure:"JWT_SECRET"`
	JWTExpiration int      `mapstructure:"JWT_EXPIRATION_HOURS"`
	DB            DBConfig `mapstructure:",squash"`
}

// LoadConfig loads configuration from .env file or environment variables
func LoadConfig() (*Config, error) {
	v := viper.New()

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Bind general env vars
	v.BindEnv("PORT")
	v.BindEnv("JWT_SECRET")
	v.BindEnv("JWT_EXPIRATION_HOURS")

	// Set defaults
	v.SetDefault("JWT_EXPIRATION_HOURS", 24)

	// Bind module-specific env vars
	BindDBEnv(v)

	v.SetConfigFile(".env")
	v.SetConfigType("env")

	// Ignore error if config file not found, fallback to env vars
	_ = v.ReadInConfig()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
```

- [ ] **Step 3: Update .env.example**

Add JWT variables to `.env.example`:

```env
# Application
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_HOURS=24
```

- [ ] **Step 4: Commit**

```bash
git add go.mod go.sum config/config.go .env.example
git commit -m "feat(auth): add JWT dependency and configuration"
```

---

### Task 2: Implement Repository Method

**Files:**
- Modify: `internal/modules/auth/repository/auth_repository.go`

*(Note: `GetByEmail` and `GetUserRoles` are already implemented in `auth_repository.go` and defined in `interfaces.go`, so we just need to verify they work as expected for login. No new code is required for the repository in this task, but we will ensure the tests pass).*

- [ ] **Step 1: Verify existing tests**

Run: `go test ./internal/modules/auth/repository/...`
Expected: PASS. If tests don't exist, we will proceed as the implementation is already complete from Task 1.

- [ ] **Step 2: Commit (Empty or documentation update if needed)**

```bash
# No commit necessary if no changes made, but acknowledging the repo layer is ready.
```

---

### Task 3: Implement Service Layer

**Files:**
- Modify: `internal/modules/auth/service/auth_service_test.go`
- Modify: `internal/modules/auth/service/auth_service.go`

- [ ] **Step 1: Write the failing unit tests for Login**

Add to `internal/modules/auth/service/auth_service_test.go`:

```go
func TestAuthService_Login_Success(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{
		JWTSecret:     "secret",
		JWTExpiration: 24,
	}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	// Generate a real bcrypt hash for the mock user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	mockUser := &models.User{
		ID:           "user-123",
		Username:     "testuser",
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	mockRoles := []models.Role{{ID: "role-1", Name: "user"}}

	repo.On("GetByEmail", ctx, req.Email).Return(mockUser, nil)
	repo.On("GetUserRoles", ctx, mockUser.ID).Return(mockRoles, nil)

	resp, err := svc.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotEmpty(t, resp.AccessToken)
	assert.Equal(t, mockUser.ID, resp.User.ID)
	repo.AssertExpectations(t)
}

func TestAuthService_Login_InvalidCredentials(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	repo.On("GetByEmail", ctx, req.Email).Return(nil, nil) // User not found

	resp, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Invalid credentials")
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/modules/auth/service -run TestAuthService_Login`
Expected: FAIL (returns "not implemented")

- [ ] **Step 3: Implement Login Service Logic**

Update `internal/modules/auth/service/auth_service.go` and add the `time` and `github.com/golang-jwt/jwt/v5` imports if needed. Remove the placeholder `fmt.Errorf("not implemented")`.

```go
// Add to imports in auth_service.go
import (
	// ... existing imports
	"time"
	"github.com/golang-jwt/jwt/v5"
)

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	// 1. Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	// 3. Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, errors.Internal("Failed to load user roles")
	}

	// 4. Generate JWT
	var roleNames []string
	for _, role := range roles {
		roleNames = append(roleNames, role.Name)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"roles": roleNames,
		"exp":   time.Now().Add(time.Hour * time.Duration(s.cfg.JWTExpiration)).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.Internal("Failed to generate token")
	}

	// Don't leak the password hash in the response
	responseUser := *user
	responseUser.PasswordHash = ""

	return &domain.LoginResponse{
		AccessToken: tokenString,
		User:        responseUser,
	}, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./internal/modules/auth/service`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/modules/auth/service/
git commit -m "feat(auth): implement login service logic with JWT generation"
```

---

### Task 4: Implement Handler and Route

**Files:**
- Modify: `internal/modules/auth/handler/auth_handler_ginkgo_test.go`
- Modify: `internal/modules/auth/handler/auth_handler.go`
- Modify: `internal/modules/auth/router/router.go`

- [ ] **Step 1: Write Handler Tests (Ginkgo)**

Add to `internal/modules/auth/handler/auth_handler_ginkgo_test.go` inside the `Describe("AuthHandler", ...)` block:

```go
	Describe("Login", func() {
		Context("with valid request", func() {
			It("should return 200 OK with token", func() {
				reqBody := domain.LoginRequest{
					Email:    "test@example.com",
					Password: "password123",
				}
				expectedResp := &domain.LoginResponse{
					AccessToken: "eyJhbGci...",
					User:        models.User{ID: "uuid-123", Email: "test@example.com"},
				}

				mockSvc.On("Login", mock.Anything, reqBody).Return(expectedResp, nil)

				jsonBody, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusOK))

				var fullResp struct {
					Success bool                 `json:"success"`
					Data    domain.LoginResponse `json:"data"`
				}
				body, _ := io.ReadAll(resp.Body)
				json.Unmarshal(body, &fullResp)

				Expect(fullResp.Success).To(BeTrue())
				Expect(fullResp.Data.AccessToken).To(Equal("eyJhbGci..."))
				mockSvc.AssertExpectations(GinkgoT())
			})
		})

		Context("with invalid credentials", func() {
			It("should return 401 Unauthorized", func() {
				reqBody := domain.LoginRequest{
					Email:    "test@example.com",
					Password: "wrong",
				}

				mockSvc.On("Login", mock.Anything, reqBody).Return(nil, errors.Unauthorized("Invalid credentials"))

				jsonBody, _ := json.Marshal(reqBody)
				req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(string(jsonBody)))
				req.Header.Set("Content-Type", "application/json")

				resp, err := app.Test(req)
				Expect(err).NotTo(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
			})
		})
	})
```
*Note: Make sure to add `app.Post("/auth/login", h.Login)` to the `BeforeEach` block in the test file so the route exists during tests.*

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./internal/modules/auth/handler`
Expected: FAIL (missing Login handler method)

- [ ] **Step 3: Implement Handler Method**

Add to `internal/modules/auth/handler/auth_handler.go`:

```go
// Login godoc
// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=domain.LoginResponse}
// @Failure 400 {object} errors.ErrorResponse
// @Failure 401 {object} errors.ErrorResponse
// @Failure 422 {object} errors.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req domain.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid request body")
	}

	if err := validation.Validate(req); err != nil {
		return err
	}

	resp, err := h.svc.Login(c.Context(), req)
	if err != nil {
		return err
	}

	return response.Success(c, http.StatusOK, "Login successful", resp)
}
```

- [ ] **Step 4: Register Route**

Modify `internal/modules/auth/router/router.go`:

```go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/internal/modules/auth/handler"
)

func RegisterRoutes(router fiber.Router, h *handler.AuthHandler) {
	auth := router.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./internal/modules/auth/handler`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/modules/auth/handler/ internal/modules/auth/router/
git commit -m "feat(auth): implement login handler and routing"
```

---

### Task 5: Integration Test and Swagger Update

**Files:**
- Create: `tests/integration/auth_login_integration_test.go`
- Modify: `docs/swagger/` (via make)

- [ ] **Step 1: Write Integration Test**

Create `tests/integration/auth_login_integration_test.go`:

```go
package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/di"
	"github.com/user/simple-blog/internal/platform/server"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Auth Login Integration", func() {
	var (
		cfg  *config.Config
		db   *sqlx.DB
		mock sqlmock.Sqlmock
		srv  *server.Server
	)

	BeforeEach(func() {
		cfg = &config.Config{
			Port: "8080",
			JWTSecret: "test-secret",
			JWTExpiration: 24,
		}
		dbRaw, mockRaw, _ := sqlmock.New()
		db = sqlx.NewDb(dbRaw, "postgres")
		mock = mockRaw
		srv = di.InitializeServer(cfg, db)
	})

	AfterEach(func() {
		db.Close()
	})

	Describe("POST /auth/login", func() {
		It("should login successfully and return token", func() {
			reqBody := domain.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			}

			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
			
			// Mock GetByEmail
			mock.ExpectQuery("SELECT (.+) FROM users WHERE email = \\$1").
				WithArgs(reqBody.Email).
				WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password_hash", "created_at", "updated_at"}).
					AddRow("user-uuid", "testuser", reqBody.Email, string(hashedPassword), time.Now(), time.Now()))

			// Mock GetUserRoles
			mock.ExpectQuery("SELECT r.id, r.name, r.created_at FROM roles r JOIN user_roles ur ON r.id = ur.role_id WHERE ur.user_id = \\$1").
				WithArgs("user-uuid").
				WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
					AddRow("role-uuid", "user", time.Now()))

			body, _ := json.Marshal(reqBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := srv.App.Test(req)
			Expect(err).NotTo(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			var fullResp struct {
				Success bool                 `json:"success"`
				Data    domain.LoginResponse `json:"data"`
			}
			respBody, _ := io.ReadAll(resp.Body)
			json.Unmarshal(respBody, &fullResp)

			Expect(fullResp.Success).To(BeTrue())
			Expect(fullResp.Data.AccessToken).NotTo(BeEmpty())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
	})
})
```

- [ ] **Step 2: Run integration tests**

Run: `go test ./tests/integration`
Expected: PASS

- [ ] **Step 3: Update Swagger Docs**

Run: `make swagger`
Expected: Swagger documentation is generated successfully.

- [ ] **Step 4: Update TASKS.md**

Check off all items under Task 2 in `TASKS.md` and check off `POST /auth/login` under API Endpoints.

- [ ] **Step 5: Commit**

```bash
git add tests/integration/ docs/swagger/ TASKS.md
git commit -m "test(auth): add login integration test and update docs"
```
