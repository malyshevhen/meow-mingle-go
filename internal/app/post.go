package app

import (
	"context"
	"time"
)

type Post struct {
	ID        string          `json:"id"`
	AuthorID  string          `json:"author_id"`
	Content   string          `json:"content"`
	Comments  []*Comment      `json:"comments"`
	Reactions []*PostReaction `json:"reactions"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type PostReaction struct {
	PostID     string `json:"post_id"`
	ReactionID string `json:"reaction_id"`
}

type PostService interface {
	CreatePost(ctx context.Context, post *Post) error
	GetPost(ctx context.Context, id string) (post *Post, err error)
	GetFeed(ctx context.Context, userId string) (feed []*Post, err error)
	ListPosts(ctx context.Context, authorId string) (posts []*Post, err error)
	UpdatePost(ctx context.Context, postId, content string) error
	DeletePost(ctx context.Context, postId string) error
}

type CreatePostParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
}

type UpdatePostParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId string `json:"author_id"`
}

type PostRepository interface {
	CreatePost(ctx context.Context, params CreatePostParams) (post Post, err error)
	GetPost(ctx context.Context, id string) (post Post, err error)
	GetFeed(ctx context.Context, userId string) (feed []Post, err error)
	ListUserPosts(ctx context.Context, userId string) (posts []Post, err error)
	UpdatePost(ctx context.Context, params UpdatePostParams) (post Post, err error)
	DeletePost(ctx context.Context, userId, postId string) error
}
