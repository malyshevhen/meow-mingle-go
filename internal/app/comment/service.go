package comment

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	commentRepo app.CommentRepository
}

// Create implements app.CommentService.
func (s *service) Create(ctx context.Context, comment *app.Comment) error {
	panic("unimplemented")
}

// Delete implements app.CommentService.
func (s *service) Delete(ctx context.Context, commentId string) error {
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

func NewService(commentRepo app.CommentRepository) app.CommentService {
	return &service{
		commentRepo: commentRepo,
	}
}
