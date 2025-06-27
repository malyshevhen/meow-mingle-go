package db

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type postRepository struct {
	session *gocql.Session
}

// Save implements post.repository.
func (pr *postRepository) Save(ctx context.Context, authorId, content string) (post app.Post, err error) {
	panic("unimplemented")
}

// Get implements post.repository.
func (pr *postRepository) Get(ctx context.Context, postId string) (post app.Post, err error) {
	panic("unimplemented")
}

// Feed implements post.repository.
func (pr *postRepository) Feed(ctx context.Context, userId string) (feed []app.Post, err error) {
	panic("unimplemented")
}

// List implements post.repository.
func (pr *postRepository) List(ctx context.Context, profileId string) (posts []app.Post, err error) {
	panic("unimplemented")
}

// Update implements post.repository.
func (pr *postRepository) Update(ctx context.Context, postId, content string) (post app.Post, err error) {
	panic("unimplemented")
}

// Delete implements post.repository.
func (pr *postRepository) Delete(ctx context.Context, postId string) error {
	panic("unimplemented")
}

func NewPostRepository(session *gocql.Session) *postRepository {
	return &postRepository{
		session: session,
	}
}
