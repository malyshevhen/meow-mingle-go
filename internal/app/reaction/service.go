package reaction

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	Save(ctx context.Context, targetID, authorID, content string) error
	Delete(ctx context.Context, targetID, authorID string) error
}

type service struct {
	reactionRepo repository
}

// Add implements app.ReactionService.
func (s *service) Add(ctx context.Context, reaction *app.Reaction) error {
	panic("unimplemented")
}

// Remove implements app.ReactionService.
func (s *service) Remove(ctx context.Context, reactionID string) error {
	panic("unimplemented")
}

func NewService(reactionRepo repository) app.ReactionService {
	return &service{
		reactionRepo: reactionRepo,
	}
}
