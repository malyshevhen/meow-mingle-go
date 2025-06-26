package app

import (
	"context"
	"time"
)

type Reaction struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReactionService interface {
	CreateReaction(ctx context.Context, reaction *Reaction) error
	ListReactionsByIDs(ctx context.Context, reactionIds []string) (reactions []*Reaction, err error)
	DeleteReaction(ctx context.Context, reactionId string) error
}

type CreateReactionParams struct {
	ID       string `json:"id"`
	AuthorID string `json:"author_id" validate:"required"`
	Content  string `json:"content" validate:"required"`
}

type DeleteReactionParams struct {
	ID         string `json:"id"`
	AuthorID   string `json:"author_id"`
	ReactionID string `json:"reaction_id"`
}

type ReactionRepository interface {
	CreateReaction(ctx context.Context, params CreateReactionParams) error
	DeleteReaction(ctx context.Context, params DeleteReactionParams) error
}
