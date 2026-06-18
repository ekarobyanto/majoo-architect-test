# Comments Module Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the Comments module to allow users to add comments to posts, update them, and delete them, with appropriate authorization checks.

**Architecture:** We will implement domain interfaces, repository (sqlx), service (RBAC/ownership and post verification), and handler (Fiber) for comments. We will inject the `PostService` into the `CommentService` to verify a post exists before commenting.

**Tech Stack:** Go, Fiber v2, PostgreSQL (sqlx), Ginkgo/Gomega, testify/mock.

---

### Task 1: Domain and Repository

**Files:**
- Create: `internal/modules/comments/domain/dto.go`
- Create: `internal/modules/comments/domain/interfaces.go`
- Create: `internal/modules/comments/repository/comment_repository.go`
- Create: `internal/modules/comments/repository/comment_repository_test.go`
- Create: `internal/modules/comments/repository/repository_suite_test.go`

- [ ] **Step 1: Define Domain Interfaces and DTOs**

Create `internal/modules/comments/domain/dto.go`:

```go
package domain

import (
	"time"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

Create `internal/modules/comments/domain/interfaces.go`:

```go
package domain

import (
	"context"
	"github.com/user/simple-blog/models"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
)

type CommentRepository interface {
	Create(ctx context.Context, comment *models.Comment) error
	GetByID(ctx context.Context, id string) (*models.Comment, error)
	Update(ctx context.Context, comment *models.Comment) error
	Delete(ctx context.Context, id string) error
}

type CommentService interface {
	Create(ctx context.Context, postID, authorID string, req CreateCommentRequest) (*models.Comment, error)
	Update(ctx context.Context, id string, user *authDomain.UserContext, req UpdateCommentRequest) (*models.Comment, error)
	Delete(ctx context.Context, id string, user *authDomain.UserContext) error
}
```

- [ ] **Step 2: Write failing unit tests for Repository**

Create `internal/modules/comments/repository/repository_suite_test.go` and `comment_repository_test.go` using `sqlmock`. Ensure basic coverage for Create, GetByID, Update, Delete.

- [ ] **Step 3: Implement CommentRepository**

Create `internal/modules/comments/repository/comment_repository.go` using `sqlx`.

```go
package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/models"
)

type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (id, post_id, author_id, content) VALUES (:id, :post_id, :author_id, :content) RETURNING created_at, updated_at`
	rows, err := sqlx.NamedQueryContext(ctx, database.GetQueryer(ctx, r.db), query, comment)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(comment)
	}
	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, post_id, author_id, content, created_at, updated_at FROM comments WHERE id = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &comment, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &comment, err
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `UPDATE comments SET content = :content, updated_at = :updated_at WHERE id = :id`
	comment.UpdatedAt = time.Now()
	_, err := sqlx.NamedExecContext(ctx, database.GetQueryer(ctx, r.db), query, comment)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := database.GetQueryer(ctx, r.db).ExecContext(ctx, query, id)
	return err
}
```

- [ ] **Step 4: Run tests to verify**

Run: `go test ./internal/modules/comments/repository/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/modules/comments/domain/ internal/modules/comments/repository/
git commit -m "feat(comments): implement domain and repository for comments"
```

---

### Task 2: Implement Service Layer

**Files:**
- Create: `internal/modules/comments/service/comment_service.go`
- Create: `internal/modules/comments/service/comment_service_test.go`

- [ ] **Step 1: Write unit tests for CommentService**

Create `internal/modules/comments/service/comment_service_test.go` testing Create, Update, Delete. Mock both `CommentRepository` and `postsDomain.PostService`. Ensure ownership tests work properly.

- [ ] **Step 2: Implement CommentService**

Create `internal/modules/comments/service/comment_service.go`. Inject `postsDomain.PostService` to verify the post exists during creation.

```go
package service

import (
	"context"

	"github.com/google/uuid"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	postsDomain "github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
)

type commentService struct {
	repo    domain.CommentRepository
	postSvc postsDomain.PostService
}

func NewCommentService(repo domain.CommentRepository, postSvc postsDomain.PostService) domain.CommentService {
	return &commentService{repo: repo, postSvc: postSvc}
}

func (s *commentService) Create(ctx context.Context, postID, authorID string, req domain.CreateCommentRequest) (*models.Comment, error) {
	// Verify post exists
	_, err := s.postSvc.GetByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	comment := &models.Comment{
		ID:       uuid.New().String(),
		PostID:   postID,
		AuthorID: authorID,
		Content:  req.Content,
	}

	if err := s.repo.Create(ctx, comment); err != nil {
		return nil, errors.Internal("Failed to create comment")
	}

	return comment, nil
}

func (s *commentService) Update(ctx context.Context, id string, user *authDomain.UserContext, req domain.UpdateCommentRequest) (*models.Comment, error) {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Internal("Failed to fetch comment")
	}
	if comment == nil {
		return nil, errors.NotFound("Comment not found")
	}

	if !middleware.IsOwnerOrAdmin(user, comment.AuthorID) {
		return nil, errors.Forbidden("You do not have permission to update this comment")
	}

	comment.Content = req.Content
	if err := s.repo.Update(ctx, comment); err != nil {
		return nil, errors.Internal("Failed to update comment")
	}

	return comment, nil
}

func (s *commentService) Delete(ctx context.Context, id string, user *authDomain.UserContext) error {
	comment, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.Internal("Failed to fetch comment")
	}
	if comment == nil {
		return errors.NotFound("Comment not found")
	}

	if !middleware.IsOwnerOrAdmin(user, comment.AuthorID) {
		return errors.Forbidden("You do not have permission to delete this comment")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.Internal("Failed to delete comment")
	}
	return nil
}
```

- [ ] **Step 3: Run tests to verify**

Run: `go test ./internal/modules/comments/service/...`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add internal/modules/comments/service/
git commit -m "feat(comments): implement comment service with post verification and auth"
```

---

### Task 3: Implement Handler and Route

**Files:**
- Create: `internal/modules/comments/handler/comment_handler.go`
- Create: `internal/modules/comments/handler/handler_suite_test.go`
- Create: `internal/modules/comments/handler/comment_handler_ginkgo_test.go`
- Create: `internal/modules/comments/router/router.go`

- [ ] **Step 1: Implement Handler**

Create `internal/modules/comments/handler/comment_handler.go`. 
It should have endpoints for: `Create` (taking `post_id` from params), `Update`, `Delete`.

- [ ] **Step 2: Write Ginkgo Tests for Handler**

Implement tests in `internal/modules/comments/handler/comment_handler_ginkgo_test.go`.

- [ ] **Step 3: Implement Router**

Create `internal/modules/comments/router/router.go`:

```go
package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/user/simple-blog/config"
	authMiddleware "github.com/user/simple-blog/internal/modules/auth/middleware"
	"github.com/user/simple-blog/internal/modules/comments/handler"
)

func RegisterRoutes(app *fiber.App, h *handler.CommentHandler, cfg *config.Config) {
	// Protected routes
	auth := authMiddleware.JWTAuth(cfg)

	// POST /posts/:id/comments
	app.Post("/posts/:id/comments", auth, h.Create)

	comments := app.Group("/comments", auth)
	comments.Put("/:id", h.Update)
	comments.Delete("/:id", h.Delete)
}
```

- [ ] **Step 4: Run tests to verify**

Run: `go test ./internal/modules/comments/handler/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/modules/comments/handler/ internal/modules/comments/router/
git commit -m "feat(comments): implement handler and routing for comments"
```

---

### Task 4: Integration and Final Steps

**Files:**
- Modify: `internal/modules/comments/provider.go`
- Modify: `internal/platform/di/wire.go`
- Modify: `internal/platform/server/server.go`
- Modify: `internal/platform/server/router.go`
- Create: `tests/integration/comments_integration_test.go`
- Modify: `docs/swagger/` (via make)
- Modify: `TASKS.md`

- [ ] **Step 1: Wiring DI and Server**
Create a `provider.go` in `internal/modules/comments/` and wire it in `internal/platform/di/wire.go`. Ensure `RegisterRoutes` is called in `router.go`.

- [ ] **Step 2: Write Integration Test**
Add tests in `tests/integration/comments_integration_test.go`.

- [ ] **Step 3: Update Swagger**
Run `make swagger`.

- [ ] **Step 4: Update TASKS.md**
Check off Task 6 items.

- [ ] **Step 5: Commit**
```bash
git add .
git commit -m "test(comments): wire comments module and add integration tests"
```
