package db

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type subscriptionRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

// SubscriptionRepository defines the interface for subscription data operations
type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, followerId, followingId string) error
	DeleteSubscription(ctx context.Context, followerId, followingId string) error
	GetFollowers(ctx context.Context, followingId string, limit int) ([]app.Subscription, error)
	GetFollowing(ctx context.Context, followerId string, limit int) ([]app.Subscription, error)
	IsFollowing(ctx context.Context, followerId, followingId string) (bool, error)
	CountFollowers(ctx context.Context, followingId string) (int, error)
	CountFollowing(ctx context.Context, followerId string) (int, error)
	GetMutualFollowings(ctx context.Context, userId1, userId2 string) ([]app.Subscription, error)
}

// CreateSubscription creates a new follow relationship
func (sr *subscriptionRepository) CreateSubscription(ctx context.Context, followerId, followingId string) error {
	if followerId == "" {
		return errors.NewValidationError("follower ID is required")
	}

	if followingId == "" {
		return errors.NewValidationError("following ID is required")
	}

	if followerId == followingId {
		return errors.NewValidationError("cannot follow yourself")
	}

	// Check if already following
	isFollowing, err := sr.IsFollowing(ctx, followerId, followingId)
	if err != nil {
		return err
	}

	if isFollowing {
		return errors.NewValidationError("already following this user")
	}

	now := time.Now()

	// Insert into subscriptions table (following -> followers)
	query := `INSERT INTO mingle.subscriptions (follower_id, following_id, created_at)
			  VALUES (?, ?, ?)`

	err = sr.session.Query(query, followerId, followingId, now).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to create subscription",
			"follower_id", followerId,
			"following_id", followingId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Insert into followers table (reverse lookup)
	followersQuery := `INSERT INTO mingle.followers (following_id, follower_id, created_at)
					   VALUES (?, ?, ?)`

	err = sr.session.Query(followersQuery, followingId, followerId, now).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to create follower entry",
			"follower_id", followerId,
			"following_id", followingId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Info("Subscription created successfully",
		"follower_id", followerId,
		"following_id", followingId,
	)

	return nil
}

// DeleteSubscription removes a follow relationship
func (sr *subscriptionRepository) DeleteSubscription(ctx context.Context, followerId, followingId string) error {
	if followerId == "" {
		return errors.NewValidationError("follower ID is required")
	}

	if followingId == "" {
		return errors.NewValidationError("following ID is required")
	}

	// Check if currently following
	isFollowing, err := sr.IsFollowing(ctx, followerId, followingId)
	if err != nil {
		return err
	}

	if !isFollowing {
		return errors.NewNotFoundError("subscription not found")
	}

	// Delete from subscriptions table
	query := `DELETE FROM mingle.subscriptions WHERE follower_id = ? AND following_id = ?`
	err = sr.session.Query(query, followerId, followingId).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to delete subscription",
			"follower_id", followerId,
			"following_id", followingId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Delete from followers table
	followersQuery := `DELETE FROM mingle.followers WHERE following_id = ? AND follower_id = ?`
	err = sr.session.Query(followersQuery, followingId, followerId).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to delete follower entry",
			"follower_id", followerId,
			"following_id", followingId,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Info("Subscription deleted successfully",
		"follower_id", followerId,
		"following_id", followingId,
	)

	return nil
}

// GetFollowers retrieves followers for a user
func (sr *subscriptionRepository) GetFollowers(ctx context.Context, followingId string, limit int) ([]app.Subscription, error) {
	if followingId == "" {
		return nil, errors.NewValidationError("following ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var subscriptions []app.Subscription

	query := `SELECT follower_id, created_at FROM mingle.followers
			  WHERE following_id = ? LIMIT ?`

	iter := sr.session.Query(query, followingId, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var followerId string
	var createdAt time.Time

	for iter.Scan(&followerId, &createdAt) {
		subscriptions = append(subscriptions, app.Subscription{
			FollowerID:  followerId,
			FollowingID: followingId,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to get followers",
			"following_id", followingId,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Debug("Followers retrieved successfully",
		"following_id", followingId,
		"followers_count", len(subscriptions),
	)

	return subscriptions, nil
}

// GetFollowing retrieves users that a user is following
func (sr *subscriptionRepository) GetFollowing(ctx context.Context, followerId string, limit int) ([]app.Subscription, error) {
	if followerId == "" {
		return nil, errors.NewValidationError("follower ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var subscriptions []app.Subscription

	query := `SELECT following_id, created_at FROM mingle.subscriptions
			  WHERE follower_id = ? LIMIT ?`

	iter := sr.session.Query(query, followerId, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var followingId string
	var createdAt time.Time

	for iter.Scan(&followingId, &createdAt) {
		subscriptions = append(subscriptions, app.Subscription{
			FollowerID:  followerId,
			FollowingID: followingId,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to get following",
			"follower_id", followerId,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Debug("Following retrieved successfully",
		"follower_id", followerId,
		"following_count", len(subscriptions),
	)

	return subscriptions, nil
}

// IsFollowing checks if a user is following another user
func (sr *subscriptionRepository) IsFollowing(ctx context.Context, followerId, followingId string) (bool, error) {
	if followerId == "" {
		return false, errors.NewValidationError("follower ID is required")
	}

	if followingId == "" {
		return false, errors.NewValidationError("following ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.subscriptions WHERE follower_id = ? AND following_id = ?`

	err := sr.session.Query(query, followerId, followingId).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to check following status",
			"follower_id", followerId,
			"following_id", followingId,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

// CountFollowers counts the number of followers for a user
func (sr *subscriptionRepository) CountFollowers(ctx context.Context, followingId string) (int, error) {
	if followingId == "" {
		return 0, errors.NewValidationError("following ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.followers WHERE following_id = ?`

	err := sr.session.Query(query, followingId).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to count followers",
			"following_id", followingId,
			"error", err.Error(),
		)
		return 0, errors.NewDatabaseError(err)
	}

	return count, nil
}

// CountFollowing counts the number of users a user is following
func (sr *subscriptionRepository) CountFollowing(ctx context.Context, followerId string) (int, error) {
	if followerId == "" {
		return 0, errors.NewValidationError("follower ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.subscriptions WHERE follower_id = ?`

	err := sr.session.Query(query, followerId).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to count following",
			"follower_id", followerId,
			"error", err.Error(),
		)
		return 0, errors.NewDatabaseError(err)
	}

	return count, nil
}

// GetMutualFollowings finds mutual followings between two users
func (sr *subscriptionRepository) GetMutualFollowings(ctx context.Context, userId1, userId2 string) ([]app.Subscription, error) {
	if userId1 == "" {
		return nil, errors.NewValidationError("user ID 1 is required")
	}

	if userId2 == "" {
		return nil, errors.NewValidationError("user ID 2 is required")
	}

	var mutualFollowings []app.Subscription

	// Get all users that userId1 is following
	user1Following, err := sr.GetFollowing(ctx, userId1, 1000) // Large limit for comparison
	if err != nil {
		return nil, err
	}

	// Get all users that userId2 is following
	user2Following, err := sr.GetFollowing(ctx, userId2, 1000) // Large limit for comparison
	if err != nil {
		return nil, err
	}

	// Find intersections
	user2FollowingMap := make(map[string]app.Subscription)
	for _, sub := range user2Following {
		user2FollowingMap[sub.FollowingID] = sub
	}

	for _, sub := range user1Following {
		if _, exists := user2FollowingMap[sub.FollowingID]; exists {
			mutualFollowings = append(mutualFollowings, sub)
		}
	}

	sr.logger.WithComponent("subscription-repository").Debug("Mutual followings retrieved successfully",
		"user_id_1", userId1,
		"user_id_2", userId2,
		"mutual_count", len(mutualFollowings),
	)

	return mutualFollowings, nil
}

func NewSubscriptionRepository(session *gocql.Session) SubscriptionRepository {
	return &subscriptionRepository{
		session: session,
		logger:  logger.GetLogger(),
	}
}
