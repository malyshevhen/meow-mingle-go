package post

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	postRepo app.PostRepository
}

// Create implements app.PostService.
func (s *service) Create(ctx context.Context, post *app.Post) error {
	panic("unimplemented")
}

// Delete implements app.PostService.
func (s *service) Delete(ctx context.Context, postId string) error {
	panic("unimplemented")
}

// Feed implements app.PostService.
func (s *service) Feed(ctx context.Context, userId string) (feed []*app.Post, err error) {
	panic("unimplemented")
}

// Get implements app.PostService.
func (s *service) Get(ctx context.Context, id string) (post *app.Post, err error) {
	panic("unimplemented")
}

// List implements app.PostService.
func (s *service) List(ctx context.Context, authorId string) (posts []*app.Post, err error) {
	panic("unimplemented")
}

// Update implements app.PostService.
func (s *service) Update(ctx context.Context, postId, content string) error {
	panic("unimplemented")
}

func NewService(postRepo app.PostRepository) app.PostService {
	return &service{
		postRepo: postRepo,
	}
}
