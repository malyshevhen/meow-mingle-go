package app

import (
	"context"
	"time"
)

type Profile struct {
	ID            string          `json:"id"`
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
	GetById(ctx context.Context, profileId string) (user *Profile, err error)
	GetByEmail(ctx context.Context, profileEmail string) (user *Profile, err error)
}

type ProfileRepository interface {
	Save(ctx context.Context, userId, email, firstName, lastName string) (user Profile, err error)
	GetById(ctx context.Context, id string) (user Profile, err error)
	GetByEmail(ctx context.Context, email string) (user Profile, err error)
}
