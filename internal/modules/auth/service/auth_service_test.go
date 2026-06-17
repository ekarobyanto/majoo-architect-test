package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/go-backend-boilerplate/config"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/domain"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/service"
	"github.com/user/go-backend-boilerplate/models"
)

type mockAuthRepository struct {
	mock.Mock
}

type mockTransactor struct {
	mock.Mock
}

func (m *mockTransactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func (m *mockAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockAuthRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Role), args.Error(1)
}

func (m *mockAuthRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	args := m.Called(ctx, userID, roleID)
	return args.Error(0)
}

func (m *mockAuthRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *mockAuthRepository) GetUserRoles(ctx context.Context, userID string) ([]models.Role, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]models.Role), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetByUsername", ctx, req.Username).Return(nil, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(nil, nil)
	repo.On("GetRoleByName", ctx, "user").Return(&models.Role{ID: "role-id", Name: "user"}, nil)
	repo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(nil)
	repo.On("AssignRole", ctx, mock.AnythingOfType("string"), "role-id").Return(nil)

	resp, err := svc.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Username, resp.Username)
	assert.Equal(t, req.Email, resp.Email)
	repo.AssertExpectations(t)
}

func TestAuthService_Register_UsernameConflict(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.RegisterRequest{
		Username: "existinguser",
		Email:    "test@example.com",
		Password: "password123",
	}

	repo.On("GetByUsername", ctx, req.Username).Return(&models.User{Username: "existinguser"}, nil)

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Username already taken")
}
