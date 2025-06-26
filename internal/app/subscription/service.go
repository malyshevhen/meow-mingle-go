package subscription

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	subscriptionRepo app.SubscriptionRepository
}

// CreateSubscription implements app.SubscriptionService.
func (s *service) CreateSubscription(ctx context.Context, followingId string) error {
	panic("unimplemented")
}

// DeleteSubscription implements app.SubscriptionService.
func (s *service) DeleteSubscription(ctx context.Context, subscriptionId string) error {
	panic("unimplemented")
}

// ListSubscriptionsByFollowerId implements app.SubscriptionService.
func (s *service) ListSubscriptionsByFollowerId(ctx context.Context, followerId string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

// ListSubscriptionsByFollowingId implements app.SubscriptionService.
func (s *service) ListSubscriptionsByFollowingId(ctx context.Context, followingId string) (subscriptions []*app.Subscription, err error) {
	panic("unimplemented")
}

func NewService(subscriptionRepo app.SubscriptionRepository) app.SubscriptionService {
	return &service{
		subscriptionRepo: subscriptionRepo,
	}
}
