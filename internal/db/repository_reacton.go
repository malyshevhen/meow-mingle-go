package db

import (
	"context"

	"github.com/gocql/gocql"
)

type reactionRepository struct {
	session *gocql.Session
}

// Save implements reaction.repository.
func (r *reactionRepository) Save(ctx context.Context, targetId, authorId, content string) error {
	panic("unimplemented")
}

// GetAll implements reaction.repository.
func (r *reactionRepository) Delete(ctx context.Context, targetId, authorId string) error {
	panic("unimplemented")
}

func NewReactionRepository(session *gocql.Session) *reactionRepository {
	return &reactionRepository{
		session: session,
	}
}
