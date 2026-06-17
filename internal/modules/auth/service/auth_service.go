package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/user/go-backend-boilerplate/config"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/domain"
	"github.com/user/go-backend-boilerplate/internal/platform/database"
	"github.com/user/go-backend-boilerplate/internal/platform/errors"
	"github.com/user/go-backend-boilerplate/models"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo domain.AuthRepository
	cfg  *config.Config
	tx   database.Transactor
}

func NewAuthService(
	repo domain.AuthRepository,
	cfg *config.Config,
	tx database.Transactor,
) domain.AuthService {
	return &authService{
		repo: repo,
		cfg:  cfg,
		tx:   tx,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.RegisterResponse, error) {
	// Check if username exists
	existingUser, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.Conflict("Username already taken")
	}

	// Check if email exists
	existingUser, err = s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.Conflict("Email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Internal("Failed to hash password")
	}

	// Get default role
	role, err := s.repo.GetRoleByName(ctx, "user")
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.Internal("Default role not found")
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	err = s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		err = s.repo.CreateUser(txCtx, user)
		if err != nil {
			return err
		}

		// Assign role
		err = s.repo.AssignRole(txCtx, user.ID, role.ID)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &domain.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (*domain.LoginResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
