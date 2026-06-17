package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/modules/auth/service"
	"github.com/user/simple-blog/models"
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

func TestAuthService_Register_EmailConflict(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.RegisterRequest{
		Username: "testuser",
		Email:    "existing@example.com",
		Password: "password123",
	}

	repo.On("GetByUsername", ctx, req.Username).Return(nil, nil)
	repo.On("GetByEmail", ctx, req.Email).Return(&models.User{Email: "existing@example.com"}, nil)

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Email already registered")
}

func TestAuthService_Register_RoleNotFound(t *testing.T) {
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
	repo.On("GetRoleByName", ctx, "user").Return(nil, nil)

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "Default role not found")
}

func TestAuthService_Register_CreateUserError(t *testing.T) {
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
	repo.On("CreateUser", ctx, mock.AnythingOfType("*models.User")).Return(fmt.Errorf("db error"))

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Register_AssignRoleError(t *testing.T) {
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
	repo.On("AssignRole", ctx, mock.AnythingOfType("string"), "role-id").Return(fmt.Errorf("db error"))

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Register_GetByUsernameError(t *testing.T) {
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

	repo.On("GetByUsername", ctx, req.Username).Return(nil, fmt.Errorf("db error"))

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Register_GetByEmailError(t *testing.T) {
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
	repo.On("GetByEmail", ctx, req.Email).Return(nil, fmt.Errorf("db error"))

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Register_GetRoleByNameError(t *testing.T) {
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
	repo.On("GetRoleByName", ctx, "user").Return(nil, fmt.Errorf("db error"))

	resp, err := svc.Register(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "db error", err.Error())
}

func TestAuthService_Login_NotImplemented(t *testing.T) {
	repo := new(mockAuthRepository)
	tx := new(mockTransactor)
	cfg := &config.Config{}
	svc := service.NewAuthService(repo, cfg, tx)

	ctx := context.Background()
	req := domain.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	resp, err := svc.Login(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "not implemented", err.Error())
}
