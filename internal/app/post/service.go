package post

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	postRepo app.PostRepository
}

// CreatePost implements app.PostService.
func (s *service) CreatePost(ctx context.Context, post *app.Post) error {
	panic("unimplemented")
}

// DeletePost implements app.PostService.
func (s *service) DeletePost(ctx context.Context, postId string) error {
	panic("unimplemented")
}

// GetFeed implements app.PostService.
func (s *service) GetFeed(ctx context.Context, userId string) (feed []*app.Post, err error) {
	panic("unimplemented")
}

// GetPost implements app.PostService.
func (s *service) GetPost(ctx context.Context, id string) (post *app.Post, err error) {
	panic("unimplemented")
}

// ListPosts implements app.PostService.
func (s *service) ListPosts(ctx context.Context, authorId string) (posts []*app.Post, err error) {
	panic("unimplemented")
}

// UpdatePost implements app.PostService.
func (s *service) UpdatePost(ctx context.Context, postId, content string) error {
	panic("unimplemented")
}

func NewService(postRepo app.PostRepository) app.PostService {
	return &service{
		postRepo: postRepo,
	}
}
