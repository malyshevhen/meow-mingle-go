package comment

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	Save(ctx context.Context, authorID, postID, content string) (comment app.Comment, err error)
	GetAll(ctx context.Context, id string) (posts []app.Comment, err error)
	Update(ctx context.Context, commentID, content string) (comment app.Comment, err error)
	Delete(ctx context.Context, userID, commentID string) (err error)
}

type service struct {
	commentRepo repository
}

// Add implements app.CommentService.
func (s *service) Add(ctx context.Context, comment *app.Comment) error {
	panic("unimplemented")
}

// Remove implements app.CommentService.
func (s *service) Remove(ctx context.Context, commentID string) error {
	panic("unimplemented")
}

// List implements app.CommentService.
func (s *service) List(ctx context.Context, postID string) (comments []*app.Comment, err error) {
	panic("unimplemented")
}

// Update implements app.CommentService.
func (s *service) Update(ctx context.Context, id, content string) error {
	panic("unimplemented")
}

func NewService(commentRepo repository) app.CommentService {
	return &service{
		commentRepo: commentRepo,
	}
}
