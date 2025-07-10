package post

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	Save(ctx context.Context, authorID, content string) (post app.Post, err error)
	Get(ctx context.Context, id string) (post app.Post, err error)
	Feed(ctx context.Context, userID string) (feed []app.Post, err error)
	List(ctx context.Context, profileID string) (posts []app.Post, err error)
	Update(ctx context.Context, postID, content string) (post app.Post, err error)
	Delete(ctx context.Context, postID string) error
}

type service struct {
	postRepo repository
}

// Create implements app.PostService.
func (s *service) Create(ctx context.Context, post *app.Post) error {
	panic("unimplemented")
}

// Delete implements app.PostService.
func (s *service) Delete(ctx context.Context, postID string) error {
	panic("unimplemented")
}

// Feed implements app.PostService.
func (s *service) Feed(ctx context.Context) (feed []*app.Post, err error) {
	panic("unimplemented")
}

// Get implements app.PostService.
func (s *service) Get(ctx context.Context, id string) (post *app.Post, err error) {
	panic("unimplemented")
}

// List implements app.PostService.
func (s *service) List(ctx context.Context, authorID string) (posts []*app.Post, err error) {
	panic("unimplemented")
}

// Edit implements app.PostService.
func (s *service) Edit(ctx context.Context, postID, content string) error {
	panic("unimplemented")
}

func NewService(postRepo repository) app.PostService {
	return &service{
		postRepo: postRepo,
	}
}
