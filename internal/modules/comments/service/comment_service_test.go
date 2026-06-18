package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/modules/comments/service"
	postsDomain "github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/models"
)

type mockCommentRepository struct {
	mock.Mock
}

func (m *mockCommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *mockCommentRepository) GetByID(ctx context.Context, id string) (*models.Comment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Comment), args.Error(1)
}

func (m *mockCommentRepository) Update(ctx context.Context, comment *models.Comment) error {
	args := m.Called(ctx, comment)
	return args.Error(0)
}

func (m *mockCommentRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockPostService struct {
	mock.Mock
}

func (m *mockPostService) Create(ctx context.Context, authorID string, req postsDomain.CreatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, authorID, req)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetByID(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) GetPaginated(ctx context.Context, query postsDomain.PaginationQuery) (*postsDomain.PaginatedPostResponse, error) {
	args := m.Called(ctx, query)
	return args.Get(0).(*postsDomain.PaginatedPostResponse), args.Error(1)
}

func (m *mockPostService) Update(ctx context.Context, id string, user *authDomain.UserContext, req postsDomain.UpdatePostRequest) (*models.Post, error) {
	args := m.Called(ctx, id, user, req)
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostService) Delete(ctx context.Context, id string, user *authDomain.UserContext) error {
	args := m.Called(ctx, id, user)
	return args.Error(0)
}

func TestCommentService_Create(t *testing.T) {
	repo := new(mockCommentRepository)
	postSvc := new(mockPostService)
	svc := service.NewCommentService(repo, postSvc)

	ctx := context.Background()
	postID := "post-1"
	authorID := "user-1"
	req := domain.CreateCommentRequest{Content: "Great post!"}

	postSvc.On("GetByID", ctx, postID).Return(&models.Post{ID: postID}, nil)
	repo.On("Create", ctx, mock.AnythingOfType("*models.Comment")).Return(nil)

	comment, err := svc.Create(ctx, postID, authorID, req)
	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, "Great post!", comment.Content)
	repo.AssertExpectations(t)
}

func TestCommentService_Update_Success(t *testing.T) {
	repo := new(mockCommentRepository)
	postSvc := new(mockPostService)
	svc := service.NewCommentService(repo, postSvc)

	ctx := context.Background()
	user := &authDomain.UserContext{ID: "user-1", Roles: []string{"user"}}
	req := domain.UpdateCommentRequest{Content: "Updated content"}

	existingComment := &models.Comment{ID: "c1", AuthorID: "user-1", Content: "Old"}
	repo.On("GetByID", ctx, "c1").Return(existingComment, nil)
	repo.On("Update", ctx, existingComment).Return(nil)

	comment, err := svc.Update(ctx, "c1", user, req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated content", comment.Content)
}

func TestCommentService_Delete_Forbidden(t *testing.T) {
	repo := new(mockCommentRepository)
	postSvc := new(mockPostService)
	svc := service.NewCommentService(repo, postSvc)

	ctx := context.Background()
	user := &authDomain.UserContext{ID: "other-user", Roles: []string{"user"}}

	existingComment := &models.Comment{ID: "c1", AuthorID: "owner-user"}
	repo.On("GetByID", ctx, "c1").Return(existingComment, nil)

	err := svc.Delete(ctx, "c1", user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "You do not have permission")
}
