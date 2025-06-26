package subscription

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	CreateSubscription(ctx context.Context, followerId, followingId string) error
	DeleteSubscription(ctx context.Context, followerId, followingId string) error
}

type service struct {
	subscriptionRepo repository
}

// Subscribe implements app.SubscriptionService.
func (s *service) Subscribe(ctx context.Context, followingId string) error {
	panic("unimplemented")
}

// Unsubscribe implements app.SubscriptionService.
func (s *service) Unsubscribe(ctx context.Context, subscriptionId string) error {
	panic("unimplemented")
}

// ListFollowings implements app.SubscriptionService.
func (s *service) ListFollowings(ctx context.Context, followerId string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

// ListFollowers implements app.SubscriptionService.
func (s *service) ListFollowers(ctx context.Context, followingId string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

func NewService(subscriptionRepo repository) app.SubscriptionService {
	return &service{
		subscriptionRepo: subscriptionRepo,
	}
}
