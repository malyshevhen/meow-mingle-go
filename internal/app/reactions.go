package app

import (
	"context"
	"time"
)

type Reaction struct {
	TargetID  string    `json:"target_id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ReactionService interface {
	Add(ctx context.Context, reaction *Reaction) error
	Remove(ctx context.Context, reactionID string) error
}
