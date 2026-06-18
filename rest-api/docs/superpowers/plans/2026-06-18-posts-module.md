# Posts Module Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Posts module to allow users to create, read, update, and delete blog posts with proper authorization and pagination.

**Architecture:** We will implement the domain interfaces, repository (using `sqlx`), service (enforcing RBAC/ownership), and handler (using Fiber). The repository will support pagination.

**Tech Stack:** Go, Fiber v2, PostgreSQL (sqlx), Ginkgo/Gomega, testify/mock.

---

### Task 1: Domain and Repository

**Files:**
- Create: `internal/modules/posts/domain/dto.go`
- Create: `internal/modules/posts/domain/interfaces.go`
- Create: `internal/modules/posts/repository/post_repository.go`
- Create: `internal/modules/posts/repository/post_repository_test.go`
- Create: `internal/modules/posts/repository/repository_suite_test.go`

- [ ] **Step 1: Define Domain Interfaces and DTOs**

Create `internal/modules/posts/domain/dto.go`:

```go
package domain

import (
	"time"
	"github.com/user/simple-blog/models"
)

type CreatePostRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=255"`
	Content string `json:"content" validate:"required"`
}

type UpdatePostRequest struct {
	Title   string `json:"title" validate:"omitempty,min=3,max=255"`
	Content string `json:"content" validate:"omitempty"`
}

type PostResponse struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PaginationQuery struct {
	Page  int `query:"page" validate:"omitempty,min=1"`
	Limit int `query:"limit" validate:"omitempty,min=1,max=100"`
}

type PaginatedPostResponse struct {
	Data       []models.Post `json:"data"`
	Total      int64         `json:"total"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
	TotalPages int           `json:"total_pages"`
}
```

Create `internal/modules/posts/domain/interfaces.go`:

```go
package domain

import (
	"context"
	"github.com/user/simple-blog/models"
)

type PostRepository interface {
	Create(ctx context.Context, post *models.Post) error
	GetByID(ctx context.Context, id string) (*models.Post, error)
	GetPaginated(ctx context.Context, page, limit int) ([]models.Post, int64, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id string) error
}

type PostService interface {
	Create(ctx context.Context, authorID string, req CreatePostRequest) (*models.Post, error)
	GetByID(ctx context.Context, id string) (*models.Post, error)
	GetPaginated(ctx context.Context, query PaginationQuery) (*PaginatedPostResponse, error)
	Update(ctx context.Context, id string, user *UserContext, req UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, id string, user *UserContext) error
}
```

- [ ] **Step 2: Write failing unit tests for Repository**

Create `internal/modules/posts/repository/repository_suite_test.go` and `post_repository_test.go` using `sqlmock`. Ensure coverage for Create, GetByID, GetPaginated, Update, Delete. *(Omitted full mock code for brevity, but implement standard CRUD mocks)*.

- [ ] **Step 3: Implement PostRepository**

Create `internal/modules/posts/repository/post_repository.go` using `sqlx`.

```go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/models"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (id, author_id, title, content) VALUES (:id, :author_id, :title, :content) RETURNING created_at, updated_at`
	rows, err := sqlx.NamedQueryContext(ctx, database.GetQueryer(ctx, r.db), query, post)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(post)
	}
	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	query := `SELECT id, author_id, title, content, created_at, updated_at FROM posts WHERE id = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &post, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &post, err
}

func (r *postRepository) GetPaginated(ctx context.Context, page, limit int) ([]models.Post, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &total, `SELECT COUNT(*) FROM posts`)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var posts []models.Post
	query := `SELECT id, author_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err = sqlx.SelectContext(ctx, database.GetQueryer(ctx, r.db), &posts, query, limit, offset)
	return posts, total, err
}

func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	query := `UPDATE posts SET title = :title, content = :content, updated_at = :updated_at WHERE id = :id`
	post.UpdatedAt = time.Now()
	_, err := sqlx.NamedExecContext(ctx, database.GetQueryer(ctx, r.db), query, post)
	return err
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := database.GetQueryer(ctx, r.db).ExecContext(ctx, query, id)
	return err
}
```

- [ ] **Step 4: Run tests to verify**

Run: `go test ./internal/modules/posts/repository/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/modules/posts/domain/ internal/modules/posts/repository/
git commit -m "feat(posts): implement domain and repository for posts"
```

---

### Task 2: Implement Service Layer

**Files:**
- Create: `internal/modules/posts/service/post_service.go`
- Create: `internal/modules/posts/service/post_service_test.go`

- [ ] **Step 1: Write unit tests for PostService**

Create `internal/modules/posts/service/post_service_test.go` testing Create, GetByID, GetPaginated, Update (including ownership checks), and Delete.
*Ensure to test `errors.Forbidden` on Update/Delete if user is not owner/admin.*

- [ ] **Step 2: Implement PostService**

Create `internal/modules/posts/service/post_service.go`.
*Use `middleware.IsOwnerOrAdmin` to verify permissions.*
*Use `database.Transactor` for the `Delete` method to satisfy the "Implement delete transaction" requirement.*

```go
package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
)

type postService struct {
	repo domain.PostRepository
	tx   database.Transactor
}

func NewPostService(repo domain.PostRepository, tx database.Transactor) domain.PostService {
	return &postService{repo: repo, tx: tx}
}

func (s *postService) Create(ctx context.Context, authorID string, req domain.CreatePostRequest) (*models.Post, error) {
	post := &models.Post{
		ID:       uuid.New().String(),
		AuthorID: authorID,
		Title:    req.Title,
		Content:  req.Content,
	}
	if err := s.repo.Create(ctx, post); err != nil {
		return nil, errors.Internal("Failed to create post")
	}
	return post, nil
}

func (s *postService) GetByID(ctx context.Context, id string) (*models.Post, error) {
	post, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Internal("Failed to fetch post")
	}
	if post == nil {
		return nil, errors.NotFound("Post not found")
	}
	return post, nil
}

func (s *postService) GetPaginated(ctx context.Context, query domain.PaginationQuery) (*domain.PaginatedPostResponse, error) {
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	posts, total, err := s.repo.GetPaginated(ctx, page, limit)
	if err != nil {
		return nil, errors.Internal("Failed to fetch posts")
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &domain.PaginatedPostResponse{
		Data:       posts,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (s *postService) Update(ctx context.Context, id string, user *authDomain.UserContext, req domain.UpdatePostRequest) (*models.Post, error) {
	post, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !middleware.IsOwnerOrAdmin(user, post.AuthorID) {
		return nil, errors.Forbidden("You do not have permission to update this post")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := s.repo.Update(ctx, post); err != nil {
		return nil, errors.Internal("Failed to update post")
	}

	return post, nil
}

func (s *postService) Delete(ctx context.Context, id string, user *authDomain.UserContext) error {
	post, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if !middleware.IsOwnerOrAdmin(user, post.AuthorID) {
		return errors.Forbidden("You do not have permission to delete this post")
	}

	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		// In a real scenario, you might also manually delete tags/comments if DB doesn't cascade
		if err := s.repo.Delete(txCtx, id); err != nil {
			return errors.Internal("Failed to delete post")
		}
		return nil
	})
}
```

- [ ] **Step 3: Run tests to verify**

Run: `go test ./internal/modules/posts/service/...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/modules/posts/service/
git commit -m "feat(posts): implement post service with authorization and pagination"
```

---

### Task 3: Implement Handler and Route

**Files:**
- Create: `internal/modules/posts/handler/post_handler.go`
- Create: `internal/modules/posts/handler/handler_suite_test.go`
- Create: `internal/modules/posts/handler/post_handler_ginkgo_test.go`
- Create: `internal/modules/posts/router/router.go`

- [ ] **Step 1: Implement Handler**

Create `internal/modules/posts/handler/post_handler.go` with Swagger annotations for POST, GET (list), GET (detail), PUT, DELETE. Handle extracting user context for authenticated routes.

- [ ] **Step 2: Write Ginkgo Tests for Handler**

Implement tests for the handler endpoints ensuring 400s on bad input, 401/403s on bad auth, and 200/201s on success.

- [ ] **Step 3: Implement Router**

Create `internal/modules/posts/router/router.go`:

```go
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
```

- [ ] **Step 4: Run tests to verify**

Run: `go test ./internal/modules/posts/handler/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/modules/posts/handler/ internal/modules/posts/router/
git commit -m "feat(posts): implement handler and routing for posts"
```

---

### Task 4: Integration and Final Steps

**Files:**
- Modify: `cmd/api/main.go`
- Modify: `internal/platform/server/server.go`
- Create: `tests/integration/posts_integration_test.go`
- Modify: `docs/swagger/` (via make)
- Modify: `TASKS.md`

- [ ] **Step 1: Wiring DI and Server**
Create a `provider.go` in `internal/modules/posts/` and wire the module in `internal/platform/di/wire.go`. Ensure `RegisterRoutes` is called in `server.go`.

- [ ] **Step 2: Write Integration Test**
Add tests in `tests/integration/posts_integration_test.go` covering creation, updating (by owner vs non-owner), and fetching.

- [ ] **Step 3: Update Swagger**
Run `make swagger`.

- [ ] **Step 4: Update TASKS.md**
Check off Task 5 items.

- [ ] **Step 5: Commit**
```bash
git add .
git commit -m "test(posts): wire posts module and add integration tests"
```
