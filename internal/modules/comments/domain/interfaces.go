package domain

import (
	"context"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/models"
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
