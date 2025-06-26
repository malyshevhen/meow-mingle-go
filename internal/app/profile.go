package app

import (
	"context"
	"time"
)

type Profile struct {
	UserID        string          `json:"user_id"`
	Email         string          `json:"email"`
	FirstName     string          `json:"first_name"`
	LastName      string          `json:"last_name"`
	Posts         []*Post         `json:"posts"`
	Subscriptions []*Subscription `json:"subscriptions"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

type ProfileService interface {
	Create(ctx context.Context, profile *Profile) error
	GetById(ctx context.Context, userId string) (profile *Profile, err error)
}
