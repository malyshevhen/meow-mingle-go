package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

type Post struct {
	ID        string      `json:"id"`
	AuthorID  string      `json:"author_id"`
	Content   string      `json:"content"`
	Comments  []*Comment  `json:"comments"`
	Reactions []*Reaction `json:"reactions"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

func NewPost(ctx context.Context, content string) (*Post, error) {
	if len(content) == 0 {
		return nil, errors.NewValidationError("Post content is required")
	}

	authorID := auth.UserID(ctx)
	if authorID == "" {
		return nil, errors.NewUnauthorizedError()
	}

	return &Post{
		ID:       uuid.New().String(),
		AuthorID: authorID,
		Content:  content,
	}, nil
}

type PostService interface {
	Create(ctx context.Context, post *Post) error
	Get(ctx context.Context, id string) (post *Post, err error)
	Feed(ctx context.Context) (feed []*Post, err error)
	List(ctx context.Context, authorID string) (posts []*Post, err error)
	Edit(ctx context.Context, postID, content string) error
	Delete(ctx context.Context, postID string) error
}
