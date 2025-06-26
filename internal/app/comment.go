package app

import (
	"context"
	"time"
)

type Comment struct {
	ID        string      `json:"id"`
	AuthorID  string      `json:"author_id"`
	PostID    string      `json:"post_id"`
	Content   string      `json:"content"`
	Reactions []*Reaction `json:"reactions"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type CommentService interface {
	Create(ctx context.Context, comment *Comment) error
	List(ctx context.Context, postID string) (comments []*Comment, err error)
	Update(ctx context.Context, commentId, content string) error
	Delete(ctx context.Context, commentId string) error
}
