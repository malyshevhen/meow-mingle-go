package comment

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type CommentRepository interface {
	Save(ctx context.Context, authorId, postId, content string) (comment app.Comment, err error)
	GetAll(ctx context.Context, id string) (posts []app.Comment, err error)
	Update(ctx context.Context, commentId, content string) (comment app.Comment, err error)
	Delete(ctx context.Context, userId, commentId string) (err error)
}

type service struct {
	commentRepo CommentRepository
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

func NewService(commentRepo CommentRepository) app.CommentService {
	return &service{
		commentRepo: commentRepo,
	}
}
