package app

import (
	"context"
	"time"
)

type Subscription struct {
	ID          string    `json:"id"`
	FollowerID  string    `json:"follower_id"`
	FollowingID string    `json:"following_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, followingId string) error
	ListSubscriptionsByFollowerId(ctx context.Context, followerId string) (subscriptions []*Subscription, err error)
	ListSubscriptionsByFollowingId(ctx context.Context, followingId string) (subscriptions []*Subscription, err error)
	DeleteSubscription(ctx context.Context, followingId string) error
}

type CreateSubscriptionParams struct {
	ID          string `json:"id"`
	FollowerID  string `json:"follower_id" validate:"required"`
	FollowingID string `json:"following_id" validate:"required"`
}

type DeleteSubscriptionParams struct {
	ID          string `json:"id"`
	FollowerID  string `json:"follower_id"`
	FollowingID string `json:"following_id"`
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, params CreateSubscriptionParams) error
	DeleteSubscription(ctx context.Context, params DeleteSubscriptionParams) error
}
