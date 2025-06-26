package app

import (
	"context"
	"time"
)

type Subscription struct {
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

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, followerId, followingId string) error
	DeleteSubscription(ctx context.Context, followerId, followingId string) error
}
