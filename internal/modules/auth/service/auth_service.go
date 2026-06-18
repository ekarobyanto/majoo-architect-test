package service

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/user/simple-blog/config"
	"github.com/user/simple-blog/internal/modules/auth/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/internal/platform/errors"
	"github.com/user/simple-blog/models"
	"golang.org/x/crypto/bcrypt"
	"time"
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
	// 1. Get user by email
	user, err := s.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	// 2. Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.Unauthorized("Invalid credentials")
	}

	// 3. Get user roles
	roles, err := s.repo.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	// 4. Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"roles": roleNames,
		"exp":   time.Now().Add(time.Duration(s.cfg.JWTExpiration) * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, errors.Internal("Failed to generate token")
	}

	// 5. Return response
	user.PasswordHash = ""
	return &domain.LoginResponse{
		AccessToken: tokenString,
		User:        *user,
	}, nil
}
