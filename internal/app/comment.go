package app

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
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

func NewComment(ctx context.Context, postID, content string) (*Comment, error) {
	if len(content) == 0 {
		return nil, errors.NewValidationError("Comment content is required")
	}

	if len(postID) == 0 {
		return nil, errors.NewValidationError("Post ID is required")
	}

	authorID := auth.UserID(ctx)
	if authorID == "" {
		return nil, errors.NewUnauthorizedError()
	}

	return &Comment{
		ID:       uuid.New().String(),
		AuthorID: authorID,
		PostID:   postID,
		Content:  content,
	}, nil
}

type CommentService interface {
	Add(ctx context.Context, comment *Comment) error
	List(ctx context.Context, postID string) (comments []*Comment, err error)
	Update(ctx context.Context, commentId, content string) error
	Remove(ctx context.Context, commentId string) error
}
