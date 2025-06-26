package app

import (
	"context"
	"time"
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

type PostService interface {
	Create(ctx context.Context, post *Post) error
	Get(ctx context.Context, id string) (post *Post, err error)
	Feed(ctx context.Context, userId string) (feed []*Post, err error)
	List(ctx context.Context, authorId string) (posts []*Post, err error)
	Update(ctx context.Context, postId, content string) error
	Delete(ctx context.Context, postId string) error
}

type PostRepository interface {
	Save(ctx context.Context, authorId, content string) (post Post, err error)
	Get(ctx context.Context, id string) (post Post, err error)
	Feed(ctx context.Context, userId string) (feed []Post, err error)
	List(ctx context.Context, profileId string) (posts []Post, err error)
	Update(ctx context.Context, postId, content string) (post Post, err error)
	Delete(ctx context.Context, postId string) error
}
