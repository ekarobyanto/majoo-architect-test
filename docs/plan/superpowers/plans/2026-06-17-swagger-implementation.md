# Swagger Integration Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Integrate automated Swagger/OpenAPI documentation using `swag` and serve it via Fiber UI.

**Architecture:** Use `swag` CLI to parse code annotations into a `docs/` package. Serve the interactive UI at `/swagger/*` using `gofiber/swagger` middleware.

**Tech Stack:** Go 1.25+, Fiber v2, swaggo/swag, gofiber/swagger.

---

### Task 1: General API Metadata & Dependencies

**Files:**
- Modify: `cmd/api/main.go`
- Modify: `go.mod` (via `go get`)

- [x] **Step 1: Install dependencies**

Run:
```bash
go get github.com/gofiber/swagger
go get github.com/swaggo/swag/cmd/swag
```

- [x] **Step 2: Add General API Annotations to main.go**

Modify `cmd/api/main.go` to include the following annotations before the `main()` function:

```go
// @title Blog System API
// @version 1.0
// @description This is a sample blog system server.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
    // ... existing code
}
```

- [x] **Step 3: Commit**

```bash
git add go.mod go.sum cmd/api/main.go
git commit -m "feat(swagger): add general api metadata and dependencies"
```

---

### Task 2: Annotate Health Handler

**Files:**
- Modify: `internal/modules/health/handler/health_handler.go`

- [ ] **Step 1: Add annotations to CheckHealth handler**

Find the `CheckHealth` method in `internal/modules/health/handler/health_handler.go` and add these annotations:

```go
// CheckHealth godoc
// @Summary Check service health
// @Description Get the status of the service and database
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} domain.HealthResponse
// @Router /health [get]
func (h *HealthHandler) CheckHealth(c *fiber.Ctx) error {
    // ... existing code
}
```

- [ ] **Step 2: Commit**

```bash
git add internal/modules/health/handler/health_handler.go
git commit -m "feat(swagger): annotate health check handler"
```

---

### Task 3: Register Swagger Route

**Files:**
- Modify: `internal/platform/server/server.go`

- [ ] **Step 1: Import docs and swagger middleware**

Update imports in `internal/platform/server/server.go`:

```go
import (
    "github.com/gofiber/swagger"
    _ "github.com/user/simple-blog/docs" // Import generated docs
    // ... existing imports
)
```

- [ ] **Step 2: Register the /swagger/* route**

In the `NewServer` or route registration section of `internal/platform/server/server.go`:

```go
app.Get("/swagger/*", swagger.HandlerDefault)
```

- [ ] **Step 3: Commit**

```bash
git add internal/platform/server/server.go
git commit -m "feat(swagger): register swagger ui route"
```

---

### Task 4: Automation & Generation

**Files:**
- Modify: `Makefile`
- Create: `docs/` (via `swag init`)

- [ ] **Step 1: Add swagger target to Makefile**

```makefile
## swagger: Generate swagger documentation
swagger:
	swag init -g cmd/api/main.go -o ./docs --parseDependency --parseInternal
```

- [ ] **Step 2: Generate documentation**

Run:
```bash
make swagger
```

- [ ] **Step 3: Verify docs/ folder exists**

Run: `ls -d docs/`
Expected: `docs/` folder contains `docs.go`, `swagger.json`, `swagger.yaml`.

- [ ] **Step 4: Commit**

```bash
git add Makefile docs/
git commit -m "chore(swagger): add make target and generate docs"
```

---

### Task 5: Final Verification

- [ ] **Step 1: Run the application**

Run: `make build && ./bin/api` (ensure DB is up)

- [ ] **Step 2: Verify Swagger UI in browser**

Run: `curl -I http://localhost:3000/swagger/index.html`
Expected: `HTTP/1.1 200 OK`

- [ ] **Step 3: Check Health Check in Swagger JSON**

Run: `grep "Health" docs/swagger.json`
Expected: Matches found for Health tag and summary.
