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
	CreateProfile(ctx context.Context, profile *Profile) error
	GetProfileById(ctx context.Context, profileId string) (user *Profile, err error)
	GetProfileByEmail(ctx context.Context, profileEmail string) (user *Profile, err error)
}

type CreateProfileParams struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type UpdateProfileParams struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type ProfileRepository interface {
	CreateProfile(ctx context.Context, params CreateProfileParams) (user Profile, err error)
	GetProfileById(ctx context.Context, id string) (user Profile, err error)
	GetProfileByEmail(ctx context.Context, email string) (user Profile, err error)
}
