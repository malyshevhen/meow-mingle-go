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
	CreateSubscription(ctx context.Context, followerID, followingID string) error
	DeleteSubscription(ctx context.Context, followerID, followingID string) error
	GetFollowers(ctx context.Context, followingID string, limit int) ([]app.Subscription, error)
	GetFollowing(ctx context.Context, followerID string, limit int) ([]app.Subscription, error)
	IsFollowing(ctx context.Context, followerID, followingID string) (bool, error)
	CountFollowers(ctx context.Context, followingID string) (int, error)
	CountFollowing(ctx context.Context, followerID string) (int, error)
	GetMutualFollowings(ctx context.Context, firstUserID, secondUserID string) ([]app.Subscription, error)
}

// CreateSubscription creates a new follow relationship
func (sr *subscriptionRepository) CreateSubscription(ctx context.Context, followerID, followingID string) error {
	if followerID == "" {
		return errors.NewValidationError("follower ID is required")
	}

	if followingID == "" {
		return errors.NewValidationError("following ID is required")
	}

	if followerID == followingID {
		return errors.NewValidationError("cannot follow yourself")
	}

	// Check if already following
	isFollowing, err := sr.IsFollowing(ctx, followerID, followingID)
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

	err = sr.session.Query(query, followerID, followingID, now).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to create subscription",
			"follower_id", followerID,
			"following_id", followingID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Insert into followers table (reverse lookup)
	followersQuery := `INSERT INTO mingle.followers (following_id, follower_id, created_at)
					   VALUES (?, ?, ?)`

	err = sr.session.Query(followersQuery, followingID, followerID, now).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to create follower entry",
			"follower_id", followerID,
			"following_id", followingID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Info("Subscription created successfully",
		"follower_id", followerID,
		"following_id", followingID,
	)

	return nil
}

// DeleteSubscription removes a follow relationship
func (sr *subscriptionRepository) DeleteSubscription(ctx context.Context, followerID, followingID string) error {
	if followerID == "" {
		return errors.NewValidationError("follower ID is required")
	}

	if followingID == "" {
		return errors.NewValidationError("following ID is required")
	}

	// Check if currently following
	isFollowing, err := sr.IsFollowing(ctx, followerID, followingID)
	if err != nil {
		return err
	}

	if !isFollowing {
		return errors.NewNotFoundError("subscription not found")
	}

	// Delete from subscriptions table
	query := `DELETE FROM mingle.subscriptions WHERE follower_id = ? AND following_id = ?`
	err = sr.session.Query(query, followerID, followingID).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to delete subscription",
			"follower_id", followerID,
			"following_id", followingID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	// Delete from followers table
	followersQuery := `DELETE FROM mingle.followers WHERE following_id = ? AND follower_id = ?`
	err = sr.session.Query(followersQuery, followingID, followerID).WithContext(ctx).Exec()
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to delete follower entry",
			"follower_id", followerID,
			"following_id", followingID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Info("Subscription deleted successfully",
		"follower_id", followerID,
		"following_id", followingID,
	)

	return nil
}

// GetFollowers retrieves followers for a user
func (sr *subscriptionRepository) GetFollowers(ctx context.Context, followingID string, limit int) ([]app.Subscription, error) {
	if followingID == "" {
		return nil, errors.NewValidationError("following ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var subscriptions []app.Subscription

	query := `SELECT follower_id, created_at FROM mingle.followers
			  WHERE following_id = ? LIMIT ?`

	iter := sr.session.Query(query, followingID, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var followerID string
	var createdAt time.Time

	for iter.Scan(&followerID, &createdAt) {
		subscriptions = append(subscriptions, app.Subscription{
			FollowerID:  followerID,
			FollowingID: followingID,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to get followers",
			"following_id", followingID,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Debug("Followers retrieved successfully",
		"following_id", followingID,
		"followers_count", len(subscriptions),
	)

	return subscriptions, nil
}

// GetFollowing retrieves users that a user is following
func (sr *subscriptionRepository) GetFollowing(ctx context.Context, followerID string, limit int) ([]app.Subscription, error) {
	if followerID == "" {
		return nil, errors.NewValidationError("follower ID is required")
	}

	if limit <= 0 {
		limit = 50 // Default limit
	}

	var subscriptions []app.Subscription

	query := `SELECT following_id, created_at FROM mingle.subscriptions
			  WHERE follower_id = ? LIMIT ?`

	iter := sr.session.Query(query, followerID, limit).WithContext(ctx).Iter()
	defer iter.Close()

	var followingID string
	var createdAt time.Time

	for iter.Scan(&followingID, &createdAt) {
		subscriptions = append(subscriptions, app.Subscription{
			FollowerID:  followerID,
			FollowingID: followingID,
			CreatedAt:   createdAt,
			UpdatedAt:   createdAt,
		})
	}

	if err := iter.Close(); err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to get following",
			"follower_id", followerID,
			"error", err.Error(),
		)
		return nil, errors.NewDatabaseError(err)
	}

	sr.logger.WithComponent("subscription-repository").Debug("Following retrieved successfully",
		"follower_id", followerID,
		"following_count", len(subscriptions),
	)

	return subscriptions, nil
}

// IsFollowing checks if a user is following another user
func (sr *subscriptionRepository) IsFollowing(ctx context.Context, followerID, followingID string) (bool, error) {
	if followerID == "" {
		return false, errors.NewValidationError("follower ID is required")
	}

	if followingID == "" {
		return false, errors.NewValidationError("following ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.subscriptions WHERE follower_id = ? AND following_id = ?`

	err := sr.session.Query(query, followerID, followingID).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to check following status",
			"follower_id", followerID,
			"following_id", followingID,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

// CountFollowers counts the number of followers for a user
func (sr *subscriptionRepository) CountFollowers(ctx context.Context, followingID string) (int, error) {
	if followingID == "" {
		return 0, errors.NewValidationError("following ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.followers WHERE following_id = ?`

	err := sr.session.Query(query, followingID).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to count followers",
			"following_id", followingID,
			"error", err.Error(),
		)
		return 0, errors.NewDatabaseError(err)
	}

	return count, nil
}

// CountFollowing counts the number of users a user is following
func (sr *subscriptionRepository) CountFollowing(ctx context.Context, followerID string) (int, error) {
	if followerID == "" {
		return 0, errors.NewValidationError("follower ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.subscriptions WHERE follower_id = ?`

	err := sr.session.Query(query, followerID).WithContext(ctx).Scan(&count)
	if err != nil {
		sr.logger.WithComponent("subscription-repository").Error("Failed to count following",
			"follower_id", followerID,
			"error", err.Error(),
		)
		return 0, errors.NewDatabaseError(err)
	}

	return count, nil
}

// GetMutualFollowings finds mutual followings between two users
func (sr *subscriptionRepository) GetMutualFollowings(ctx context.Context, firstUserID, secondUserID string) ([]app.Subscription, error) {
	if firstUserID == "" {
		return nil, errors.NewValidationError("user ID 1 is required")
	}

	if secondUserID == "" {
		return nil, errors.NewValidationError("user ID 2 is required")
	}

	var mutualFollowings []app.Subscription

	// Get all users that userId1 is following
	user1Following, err := sr.GetFollowing(ctx, firstUserID, 1000) // Large limit for comparison
	if err != nil {
		return nil, err
	}

	// Get all users that userId2 is following
	user2Following, err := sr.GetFollowing(ctx, secondUserID, 1000) // Large limit for comparison
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
		"user_id_1", firstUserID,
		"user_id_2", secondUserID,
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
