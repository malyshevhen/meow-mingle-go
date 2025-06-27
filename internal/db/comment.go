package db

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type commentRepository struct {
	session *gocql.Session
}

// Save implements comment.repository.
func (cr *commentRepository) Save(ctx context.Context, authorId, postId, content string) (comment app.Comment, err error) {
	panic("unimplemented")
}

// GetAll implements comment.repository.
func (cr *commentRepository) GetAll(ctx context.Context, id string) (posts []app.Comment, err error) {
	panic("unimplemented")
}

// Update implements comment.repository.
func (cr *commentRepository) Update(ctx context.Context, commentId, content string) (comment app.Comment, err error) {
	panic("unimplemented")
}

// Delete implements comment.repository.
func (cr *commentRepository) Delete(ctx context.Context, userId, commentId string) (err error) {
	panic("unimplemented")
}

func NewCommentRepository(session *gocql.Session) *commentRepository {
	return &commentRepository{
		session: session,
	}
}
