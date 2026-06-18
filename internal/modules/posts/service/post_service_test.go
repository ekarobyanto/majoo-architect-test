package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	authDomain "github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/modules/posts/service"
	"github.com/user/simple-blog/models"
)

type mockPostRepository struct {
	mock.Mock
}

func (m *mockPostRepository) Create(ctx context.Context, post *models.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *mockPostRepository) GetByID(ctx context.Context, id string) (*models.Post, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *mockPostRepository) GetPaginated(ctx context.Context, page, limit int) ([]models.Post, int64, error) {
	args := m.Called(ctx, page, limit)
	if args.Get(0) == nil {
		return nil, 0, args.Error(2)
	}
	return args.Get(0).([]models.Post), args.Get(1).(int64), args.Error(2)
}

func (m *mockPostRepository) Update(ctx context.Context, post *models.Post) error {
	args := m.Called(ctx, post)
	return args.Error(0)
}

func (m *mockPostRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockTransactor struct {
	mock.Mock
}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func TestPostService_Create(t *testing.T) {
	repo := new(mockPostRepository)
	tx := new(mockTransactor)
	svc := service.NewPostService(repo, tx)

	ctx := context.Background()
	authorID := "author-1"
	req := domain.CreatePostRequest{Title: "Title", Content: "Content"}

	repo.On("Create", ctx, mock.AnythingOfType("*models.Post")).Return(nil)

	post, err := svc.Create(ctx, authorID, req)
	assert.NoError(t, err)
	assert.NotNil(t, post)
	assert.Equal(t, "Title", post.Title)
	assert.Equal(t, authorID, post.AuthorID)
	repo.AssertExpectations(t)
}

func TestPostService_Update_Success(t *testing.T) {
	repo := new(mockPostRepository)
	tx := new(mockTransactor)
	svc := service.NewPostService(repo, tx)

	ctx := context.Background()
	user := &authDomain.UserContext{ID: "author-1", Roles: []string{"user"}}
	req := domain.UpdatePostRequest{Title: "Updated Title"}

	existingPost := &models.Post{ID: "post-1", AuthorID: "author-1"}
	repo.On("GetByID", ctx, "post-1").Return(existingPost, nil)
	repo.On("Update", ctx, mock.AnythingOfType("*models.Post")).Return(nil)

	post, err := svc.Update(ctx, "post-1", user, req)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Title", post.Title)
	repo.AssertExpectations(t)
}

func TestPostService_Update_Forbidden(t *testing.T) {
	repo := new(mockPostRepository)
	tx := new(mockTransactor)
	svc := service.NewPostService(repo, tx)

	ctx := context.Background()
	user := &authDomain.UserContext{ID: "other-user", Roles: []string{"user"}}
	req := domain.UpdatePostRequest{Title: "Updated Title"}

	existingPost := &models.Post{ID: "post-1", AuthorID: "author-1"}
	repo.On("GetByID", ctx, "post-1").Return(existingPost, nil)

	post, err := svc.Update(ctx, "post-1", user, req)
	assert.Error(t, err)
	assert.Nil(t, post)
	assert.Contains(t, err.Error(), "You do not have permission")
	repo.AssertExpectations(t)
}

func TestPostService_Delete_Success(t *testing.T) {
	repo := new(mockPostRepository)
	tx := new(mockTransactor)
	svc := service.NewPostService(repo, tx)

	ctx := context.Background()
	user := &authDomain.UserContext{ID: "admin-1", Roles: []string{"admin"}}

	existingPost := &models.Post{ID: "post-1", AuthorID: "author-1"}
	repo.On("GetByID", ctx, "post-1").Return(existingPost, nil)
	repo.On("Delete", ctx, "post-1").Return(nil)

	err := svc.Delete(ctx, "post-1", user)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestPostService_GetPaginated(t *testing.T) {
	repo := new(mockPostRepository)
	tx := new(mockTransactor)
	svc := service.NewPostService(repo, tx)

	ctx := context.Background()
	query := domain.PaginationQuery{Page: 1, Limit: 10}

	posts := []models.Post{{ID: "1"}, {ID: "2"}}
	repo.On("GetPaginated", ctx, 1, 10).Return(posts, int64(20), nil)

	resp, err := svc.GetPaginated(ctx, query)
	assert.NoError(t, err)
	assert.Equal(t, int64(20), resp.Total)
	assert.Equal(t, 2, resp.TotalPages)
	repo.AssertExpectations(t)
}
