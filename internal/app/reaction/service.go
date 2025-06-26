package reaction

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	reactionRepo app.ReactionRepository
}

// CreateReaction implements app.ReactionService.
func (s *service) CreateReaction(ctx context.Context, reaction *app.Reaction) error {
	panic("unimplemented")
}

// DeleteReaction implements app.ReactionService.
func (s *service) DeleteReaction(ctx context.Context, reactionId string) error {
	panic("unimplemented")
}

// ListReactionsByIDs implements app.ReactionService.
func (s *service) ListReactionsByIDs(ctx context.Context, reactionIds []string) (reactions []*app.Reaction, err error) {
	panic("unimplemented")
}

func NewService(reactionRepo app.ReactionRepository) app.ReactionService {
	return &service{
		reactionRepo: reactionRepo,
	}
}
