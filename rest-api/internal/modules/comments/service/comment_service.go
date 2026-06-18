package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/user/simple-blog/internal/modules/auth/authorization"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
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

	if !authorization.IsOwnerOrAdmin(user, comment.AuthorID) {
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

	if !authorization.IsOwnerOrAdmin(user, comment.AuthorID) {
		return errors.Forbidden("You do not have permission to delete this comment")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.Internal("Failed to delete comment")
	}
	return nil
}
