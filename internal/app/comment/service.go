package comment

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	commentRepo app.CommentRepository
}

// CreateComment implements app.CommentService.
func (s *service) CreateComment(ctx context.Context, comment *app.Comment) error {
	panic("unimplemented")
}

// DeleteComment implements app.CommentService.
func (s *service) DeleteComment(ctx context.Context, commentId string) error {
	panic("unimplemented")
}

// ListPostComments implements app.CommentService.
func (s *service) ListPostComments(ctx context.Context, postID string) (comments []*app.Comment, err error) {
	panic("unimplemented")
}

// UpdateComment implements app.CommentService.
func (s *service) UpdateComment(ctx context.Context, id, content string) error {
	panic("unimplemented")
}

func NewService(commentRepo app.CommentRepository) app.CommentService {
	return &service{
		commentRepo: commentRepo,
	}
}
