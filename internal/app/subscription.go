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
	Subscribe(ctx context.Context, followingID string) error
	Unsubscribe(ctx context.Context, followingID string) error
	ListFollowings(ctx context.Context, followerID string) (subscriptions []*Subscription, err error)
	ListFollowers(ctx context.Context, followingID string) (subscriptions []*Subscription, err error)
}
