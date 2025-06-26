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
	CreateReaction(ctx context.Context, reaction *Reaction) error
	DeleteReaction(ctx context.Context, reactionId string) error
}

type ReactionRepository interface {
	CreateReaction(ctx context.Context, targetId, authorId, content string) error
	DeleteReaction(ctx context.Context, targetId, authorId string) error
}
