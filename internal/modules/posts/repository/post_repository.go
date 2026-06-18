package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/internal/modules/posts/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/models"
)

type postRepository struct {
	db *sqlx.DB
}

func NewPostRepository(db *sqlx.DB) domain.PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(ctx context.Context, post *models.Post) error {
	query := `INSERT INTO posts (id, author_id, title, content) VALUES (:id, :author_id, :title, :content) RETURNING created_at, updated_at`
	rows, err := sqlx.NamedQueryContext(ctx, database.GetQueryer(ctx, r.db), query, post)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(post)
	}
	return nil
}

func (r *postRepository) GetByID(ctx context.Context, id string) (*models.Post, error) {
	var post models.Post
	query := `SELECT id, author_id, title, content, created_at, updated_at FROM posts WHERE id = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &post, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &post, err
}

func (r *postRepository) GetPaginated(ctx context.Context, page, limit int) ([]models.Post, int64, error) {
	var total int64
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &total, `SELECT COUNT(*) FROM posts`)
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	var posts []models.Post
	query := `SELECT id, author_id, title, content, created_at, updated_at FROM posts ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	err = sqlx.SelectContext(ctx, database.GetQueryer(ctx, r.db), &posts, query, limit, offset)
	return posts, total, err
}

func (r *postRepository) Update(ctx context.Context, post *models.Post) error {
	query := `UPDATE posts SET title = :title, content = :content, updated_at = :updated_at WHERE id = :id`
	post.UpdatedAt = time.Now()
	_, err := sqlx.NamedExecContext(ctx, database.GetQueryer(ctx, r.db), query, post)
	return err
}

func (r *postRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM posts WHERE id = $1`
	_, err := database.GetQueryer(ctx, r.db).ExecContext(ctx, query, id)
	return err
}
