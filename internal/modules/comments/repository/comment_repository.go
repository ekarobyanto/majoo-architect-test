package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/user/simple-blog/internal/modules/comments/domain"
	"github.com/user/simple-blog/internal/platform/database"
	"github.com/user/simple-blog/models"
)

type commentRepository struct {
	db *sqlx.DB
}

func NewCommentRepository(db *sqlx.DB) domain.CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (id, post_id, author_id, content) VALUES (:id, :post_id, :author_id, :content) RETURNING created_at, updated_at`
	rows, err := sqlx.NamedQueryContext(ctx, database.GetQueryer(ctx, r.db), query, comment)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return rows.StructScan(comment)
	}
	return nil
}

func (r *commentRepository) GetByID(ctx context.Context, id string) (*models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, post_id, author_id, content, created_at, updated_at FROM comments WHERE id = $1`
	err := sqlx.GetContext(ctx, database.GetQueryer(ctx, r.db), &comment, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &comment, err
}

func (r *commentRepository) Update(ctx context.Context, comment *models.Comment) error {
	query := `UPDATE comments SET content = :content, updated_at = :updated_at WHERE id = :id`
	comment.UpdatedAt = time.Now()
	_, err := sqlx.NamedExecContext(ctx, database.GetQueryer(ctx, r.db), query, comment)
	return err
}

func (r *commentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := database.GetQueryer(ctx, r.db).ExecContext(ctx, query, id)
	return err
}
