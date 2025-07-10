package subscription

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	CreateSubscription(ctx context.Context, followerID, followingID string) error
	DeleteSubscription(ctx context.Context, followerID, followingID string) error
}

type service struct {
	subscriptionRepo repository
}

// Subscribe implements app.SubscriptionService.
func (s *service) Subscribe(ctx context.Context, followingID string) error {
	panic("unimplemented")
}

// Unsubscribe implements app.SubscriptionService.
func (s *service) Unsubscribe(ctx context.Context, subscriptionID string) error {
	panic("unimplemented")
}

// ListFollowings implements app.SubscriptionService.
func (s *service) ListFollowings(ctx context.Context, followerID string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

// ListFollowers implements app.SubscriptionService.
func (s *service) ListFollowers(ctx context.Context, followingID string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

func NewService(subscriptionRepo repository) app.SubscriptionService {
	return &service{
		subscriptionRepo: subscriptionRepo,
	}
}
