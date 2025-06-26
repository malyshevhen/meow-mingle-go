package app

import (
	"context"
	"time"
)

type Comment struct {
	ID        string             `json:"id"`
	AuthorID  string             `json:"author_id"`
	PostID    string             `json:"post_id"`
	Content   string             `json:"content"`
	Reactions []*CommentReaction `json:"reactions"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type CommentReaction struct {
	CommentID  string `json:"comment_id"`
	ReactionID string `json:"reaction_id"`
}

type CommentService interface {
	CreateComment(ctx context.Context, comment *Comment) error
	ListPostComments(ctx context.Context, postID string) (comments []*Comment, err error)
	UpdateComment(ctx context.Context, commentId, content string) error
	DeleteComment(ctx context.Context, commentId string) error
}

type CreateCommentParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
	PostID   string `json:"post_id" validate:"required"`
}

type UpdateCommentParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId string `json:"author_id"`
}

type CommentRepository interface {
	CreateComment(ctx context.Context, params CreateCommentParams) (comment Comment, err error)
	ListPostComments(ctx context.Context, id string) (posts []Comment, err error)
	UpdateComment(ctx context.Context, params UpdateCommentParams) (comment Comment, err error)
	DeleteComment(ctx context.Context, userId, commentId string) (err error)
}
