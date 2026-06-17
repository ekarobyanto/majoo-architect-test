package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/user/go-backend-boilerplate/internal/modules/auth/domain"
	"github.com/user/go-backend-boilerplate/internal/platform/database"
	"github.com/user/go-backend-boilerplate/models"
)

type authRepository struct {
	db *sqlx.DB
}

func NewAuthRepository(db *sqlx.DB) domain.AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, email, password_hash) VALUES (:id, :username, :email, :password_hash) RETURNING created_at, updated_at`
	rows, err := sqlx.NamedQueryContext(ctx, database.GetQueryer(ctx, r.db), query, user)
	if err != nil {
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.StructScan(user)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *authRepository) GetRoleByName(ctx context.Context, name string) (*models.Role, error) {
	var role models.Role
	query := `SELECT id, name, created_at FROM roles WHERE name = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &role, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *authRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)`
	_, err := database.GetQueryer(ctx, r.db).ExecContext(ctx, query, userID, roleID)
	return err
}

func (r *authRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &user, query, username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUserRoles(ctx context.Context, userID string) ([]models.Role, error) {
	var roles []models.Role
	query := `
		SELECT r.id, r.name, r.created_at 
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
	`
	err := sqlx.SelectContext(ctx, database.GetQueryer(ctx, r.db), &roles, query, userID)
	if err != nil {
		return nil, err
	}
	return roles, nil
}
