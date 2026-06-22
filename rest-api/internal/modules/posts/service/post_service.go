package service

import (
	"context"
	"math"

	"github.com/google/uuid"
	"github.com/user/simple-blog/internal/modules/auth/authorization"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	commentDomain "github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
)

type postService struct {
	repo        domain.PostRepository
	tx          database.Transactor
	commentRepo commentDomain.CommentRepository
}

func NewPostService(repo domain.PostRepository, tx database.Transactor, commentRepo commentDomain.CommentRepository) domain.PostService {
	return &postService{repo: repo, tx: tx, commentRepo: commentRepo}
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

func (s *postService) GetDetailByID(ctx context.Context, id string) (*domain.PostDetailResponse, error) {
	post, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	comments, err := s.commentRepo.GetByPostID(ctx, id)
	if err != nil {
		return nil, errors.Internal("Failed to fetch post comments")
	}

	return &domain.PostDetailResponse{
		Post:     *post,
		Comments: comments,
	}, nil
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

	if !authorization.IsOwnerOrAdmin(user, post.AuthorID) {
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

	if !authorization.IsOwnerOrAdmin(user, post.AuthorID) {
		return errors.Forbidden("You do not have permission to delete this post")
	}

	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Delete(txCtx, id); err != nil {
			return errors.Internal("Failed to delete post")
		}
		return nil
	})
}
