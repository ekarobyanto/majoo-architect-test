package domain

import (
	"context"
	"github.com/user/go-backend-boilerplate/models"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetRoleByName(ctx context.Context, name string) (*models.Role, error)
	AssignRole(ctx context.Context, userID, roleID string) error
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]models.Role, error)
}

type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*RegisterResponse, error)
	Login(ctx context.Context, req LoginRequest) (*LoginResponse, error)
}
