package db

import (
	"context"

	"github.com/gocql/gocql"
)

type subscriptionRepository struct {
	session *gocql.Session
}

// Save implements subscription.repository.
func (r *subscriptionRepository) CreateSubscription(ctx context.Context, followerId, followingId string) error {
	panic("unimplemented")
}

// DeleteSubscription implements subscription.repository.
func (r *subscriptionRepository) DeleteSubscription(ctx context.Context, followerId, followingId string) error {
	panic("unimplemented")
}

func NewSubscriptionRepository(session *gocql.Session) *subscriptionRepository {
	return &subscriptionRepository{
		session: session,
	}
}
