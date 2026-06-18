package domain

import (
	"time"
)

type CreateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" validate:"required,min=1"`
}

type CommentResponse struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
