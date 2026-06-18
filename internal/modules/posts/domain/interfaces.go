package domain

import (
	"context"

	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
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
	Update(ctx context.Context, id string, user *authDomain.UserContext, req UpdatePostRequest) (*models.Post, error)
	Delete(ctx context.Context, id string, user *authDomain.UserContext) error
}
