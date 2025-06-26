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
	Subscribe(ctx context.Context, followingId string) error
	Unsubscribe(ctx context.Context, followingId string) error
	ListFollowings(ctx context.Context, followerId string) (subscriptions []*Subscription, err error)
	ListFollowers(ctx context.Context, followingId string) (subscriptions []*Subscription, err error)
}
