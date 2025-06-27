package integration

import (
	"context"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type SubscriptionRepositoryTestSuite struct {
	suite.Suite
	testDB *TestDatabase
	repo   db.SubscriptionRepository
	ctx    context.Context
}

func (suite *SubscriptionRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	testDB, err := NewTestDatabase(suite.ctx)
	require.NoError(suite.T(), err, "Failed to create test database")

	suite.testDB = testDB
	suite.repo = db.NewSubscriptionRepository(testDB.Session)
}

func (suite *SubscriptionRepositoryTestSuite) TearDownSuite() {
	if suite.testDB != nil {
		suite.testDB.Close(suite.ctx)
	}
}

func (suite *SubscriptionRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.testDB.Clean(suite.ctx)
	require.NoError(suite.T(), err, "Failed to clean test database")
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_Success() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// When
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify subscription exists
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, followerID, followingID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing)
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_ValidationError_EmptyFollowerID() {
	// Given
	followerID := ""
	followingID := "user2"

	// When
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "follower ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_ValidationError_EmptyFollowingID() {
	// Given
	followerID := "user1"
	followingID := ""

	// When
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "following ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_ValidationError_SelfFollow() {
	// Given
	userID := "user1"

	// When
	err := suite.repo.CreateSubscription(suite.ctx, userID, userID)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "cannot follow yourself")
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateSubscription_AlreadyFollowing() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// Create subscription first
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)

	// When - try to create again
	err = suite.repo.CreateSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "already following this user")
}

func (suite *SubscriptionRepositoryTestSuite) TestDeleteSubscription_Success() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// Create subscription first
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.DeleteSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify subscription no longer exists
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, followerID, followingID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isFollowing)
}

func (suite *SubscriptionRepositoryTestSuite) TestDeleteSubscription_NotFound() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// When
	err := suite.repo.DeleteSubscription(suite.ctx, followerID, followingID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "subscription not found")
}

func (suite *SubscriptionRepositoryTestSuite) TestDeleteSubscription_ValidationError_EmptyFollowerID() {
	// When
	err := suite.repo.DeleteSubscription(suite.ctx, "", "user2")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "follower ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestDeleteSubscription_ValidationError_EmptyFollowingID() {
	// When
	err := suite.repo.DeleteSubscription(suite.ctx, "user1", "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "following ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowers_Success() {
	// Given
	followingID := "user1"
	follower1ID := "user2"
	follower2ID := "user3"
	follower3ID := "user4"

	// Create subscriptions
	err := suite.repo.CreateSubscription(suite.ctx, follower1ID, followingID)
	require.NoError(suite.T(), err)

	err = suite.repo.CreateSubscription(suite.ctx, follower2ID, followingID)
	require.NoError(suite.T(), err)

	err = suite.repo.CreateSubscription(suite.ctx, follower3ID, followingID)
	require.NoError(suite.T(), err)

	// When
	followers, err := suite.repo.GetFollowers(suite.ctx, followingID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), followers, 3)

	// Verify all followers are present
	followerIDs := make(map[string]bool)
	for _, subscription := range followers {
		followerIDs[subscription.FollowerID] = true
		assert.Equal(suite.T(), followingID, subscription.FollowingID)
		assert.False(suite.T(), subscription.CreatedAt.IsZero())
	}

	assert.True(suite.T(), followerIDs[follower1ID])
	assert.True(suite.T(), followerIDs[follower2ID])
	assert.True(suite.T(), followerIDs[follower3ID])
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowers_WithLimit() {
	// Given
	followingID := "user1"

	// Create 10 followers
	for i := 0; i < 10; i++ {
		followerID := "follower" + string(rune('0'+i))
		err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
		require.NoError(suite.T(), err)
	}

	// When
	followers, err := suite.repo.GetFollowers(suite.ctx, followingID, 5)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), followers, 5)
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowers_EmptyResult() {
	// Given
	followingID := "user1"

	// When
	followers, err := suite.repo.GetFollowers(suite.ctx, followingID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), followers)
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowers_ValidationError_EmptyFollowingID() {
	// When
	_, err := suite.repo.GetFollowers(suite.ctx, "", 10)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "following ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowing_Success() {
	// Given
	followerID := "user1"
	following1ID := "user2"
	following2ID := "user3"
	following3ID := "user4"

	// Create subscriptions
	err := suite.repo.CreateSubscription(suite.ctx, followerID, following1ID)
	require.NoError(suite.T(), err)

	err = suite.repo.CreateSubscription(suite.ctx, followerID, following2ID)
	require.NoError(suite.T(), err)

	err = suite.repo.CreateSubscription(suite.ctx, followerID, following3ID)
	require.NoError(suite.T(), err)

	// When
	following, err := suite.repo.GetFollowing(suite.ctx, followerID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), following, 3)

	// Verify all following are present
	followingIDs := make(map[string]bool)
	for _, subscription := range following {
		followingIDs[subscription.FollowingID] = true
		assert.Equal(suite.T(), followerID, subscription.FollowerID)
		assert.False(suite.T(), subscription.CreatedAt.IsZero())
	}

	assert.True(suite.T(), followingIDs[following1ID])
	assert.True(suite.T(), followingIDs[following2ID])
	assert.True(suite.T(), followingIDs[following3ID])
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowing_WithLimit() {
	// Given
	followerID := "user1"

	// Create 10 following relationships
	for i := 0; i < 10; i++ {
		followingID := "following" + string(rune('0'+i))
		err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
		require.NoError(suite.T(), err)
	}

	// When
	following, err := suite.repo.GetFollowing(suite.ctx, followerID, 5)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), following, 5)
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowing_EmptyResult() {
	// Given
	followerID := "user1"

	// When
	following, err := suite.repo.GetFollowing(suite.ctx, followerID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), following)
}

func (suite *SubscriptionRepositoryTestSuite) TestGetFollowing_ValidationError_EmptyFollowerID() {
	// When
	_, err := suite.repo.GetFollowing(suite.ctx, "", 10)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "follower ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestIsFollowing_True() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// Create subscription
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)

	// When
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, followerID, followingID)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing)
}

func (suite *SubscriptionRepositoryTestSuite) TestIsFollowing_False() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// When
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, followerID, followingID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isFollowing)
}

func (suite *SubscriptionRepositoryTestSuite) TestIsFollowing_ValidationError_EmptyFollowerID() {
	// When
	_, err := suite.repo.IsFollowing(suite.ctx, "", "user2")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "follower ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestIsFollowing_ValidationError_EmptyFollowingID() {
	// When
	_, err := suite.repo.IsFollowing(suite.ctx, "user1", "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "following ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowers_Success() {
	// Given
	followingID := "user1"

	// Create multiple followers
	for i := 0; i < 5; i++ {
		followerID := "follower" + string(rune('0'+i))
		err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
		require.NoError(suite.T(), err)
	}

	// When
	count, err := suite.repo.CountFollowers(suite.ctx, followingID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 5, count)
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowers_EmptyResult() {
	// Given
	followingID := "user1"

	// When
	count, err := suite.repo.CountFollowers(suite.ctx, followingID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowers_ValidationError_EmptyFollowingID() {
	// When
	_, err := suite.repo.CountFollowers(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "following ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowing_Success() {
	// Given
	followerID := "user1"

	// Create multiple following relationships
	for i := 0; i < 5; i++ {
		followingID := "following" + string(rune('0'+i))
		err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
		require.NoError(suite.T(), err)
	}

	// When
	count, err := suite.repo.CountFollowing(suite.ctx, followerID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 5, count)
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowing_EmptyResult() {
	// Given
	followerID := "user1"

	// When
	count, err := suite.repo.CountFollowing(suite.ctx, followerID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *SubscriptionRepositoryTestSuite) TestCountFollowing_ValidationError_EmptyFollowerID() {
	// When
	_, err := suite.repo.CountFollowing(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "follower ID is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestGetMutualFollowings_Success() {
	// Given
	user1 := "user1"
	user2 := "user2"
	commonUser1 := "common1"
	commonUser2 := "common2"
	onlyUser1 := "only1"
	onlyUser2 := "only2"

	// Create followings for user1
	err := suite.repo.CreateSubscription(suite.ctx, user1, commonUser1)
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, user1, commonUser2)
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, user1, onlyUser1)
	require.NoError(suite.T(), err)

	// Create followings for user2
	err = suite.repo.CreateSubscription(suite.ctx, user2, commonUser1)
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, user2, commonUser2)
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, user2, onlyUser2)
	require.NoError(suite.T(), err)

	// When
	mutualFollowings, err := suite.repo.GetMutualFollowings(suite.ctx, user1, user2)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), mutualFollowings, 2)

	// Verify mutual followings
	mutualIDs := make(map[string]bool)
	for _, subscription := range mutualFollowings {
		mutualIDs[subscription.FollowingID] = true
	}

	assert.True(suite.T(), mutualIDs[commonUser1])
	assert.True(suite.T(), mutualIDs[commonUser2])
	assert.False(suite.T(), mutualIDs[onlyUser1])
	assert.False(suite.T(), mutualIDs[onlyUser2])
}

func (suite *SubscriptionRepositoryTestSuite) TestGetMutualFollowings_EmptyResult() {
	// Given
	user1 := "user1"
	user2 := "user2"

	// When
	mutualFollowings, err := suite.repo.GetMutualFollowings(suite.ctx, user1, user2)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), mutualFollowings)
}

func (suite *SubscriptionRepositoryTestSuite) TestGetMutualFollowings_ValidationError_EmptyUser1() {
	// When
	_, err := suite.repo.GetMutualFollowings(suite.ctx, "", "user2")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID 1 is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestGetMutualFollowings_ValidationError_EmptyUser2() {
	// When
	_, err := suite.repo.GetMutualFollowings(suite.ctx, "user1", "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID 2 is required")
}

func (suite *SubscriptionRepositoryTestSuite) TestBidirectionalFollowing() {
	// Given
	user1 := "user1"
	user2 := "user2"

	// When - create bidirectional following
	err := suite.repo.CreateSubscription(suite.ctx, user1, user2)
	require.NoError(suite.T(), err)

	err = suite.repo.CreateSubscription(suite.ctx, user2, user1)
	require.NoError(suite.T(), err)

	// Then
	// Verify user1 follows user2
	isFollowing1, err := suite.repo.IsFollowing(suite.ctx, user1, user2)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing1)

	// Verify user2 follows user1
	isFollowing2, err := suite.repo.IsFollowing(suite.ctx, user2, user1)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing2)

	// Verify counts
	user1Followers, err := suite.repo.CountFollowers(suite.ctx, user1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, user1Followers)

	user2Followers, err := suite.repo.CountFollowers(suite.ctx, user2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, user2Followers)
}

func (suite *SubscriptionRepositoryTestSuite) TestConcurrentOperations() {
	// Given
	follower1 := "user1"
	follower2 := "user2"
	following := "user3"

	// When - concurrent subscriptions
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		err := suite.repo.CreateSubscription(suite.ctx, follower1, following)
		done1 <- err
	}()

	go func() {
		err := suite.repo.CreateSubscription(suite.ctx, follower2, following)
		done2 <- err
	}()

	// Then
	err1 := <-done1
	err2 := <-done2

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Verify both subscriptions exist
	count, err := suite.repo.CountFollowers(suite.ctx, following)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count)
}

func (suite *SubscriptionRepositoryTestSuite) TestComplexFollowingNetwork() {
	// Given - create a network of users
	users := []string{"user1", "user2", "user3", "user4", "user5"}

	// Create a complex following network
	// user1 follows user2, user3
	err := suite.repo.CreateSubscription(suite.ctx, users[0], users[1])
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, users[0], users[2])
	require.NoError(suite.T(), err)

	// user2 follows user3, user4
	err = suite.repo.CreateSubscription(suite.ctx, users[1], users[2])
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, users[1], users[3])
	require.NoError(suite.T(), err)

	// user3 follows user4, user5
	err = suite.repo.CreateSubscription(suite.ctx, users[2], users[3])
	require.NoError(suite.T(), err)
	err = suite.repo.CreateSubscription(suite.ctx, users[2], users[4])
	require.NoError(suite.T(), err)

	// When - verify the network
	// Check user1's following
	user1Following, err := suite.repo.GetFollowing(suite.ctx, users[0], 10)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), user1Following, 2)

	// Check user3's followers
	user3Followers, err := suite.repo.GetFollowers(suite.ctx, users[2], 10)
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), user3Followers, 2)

	// Check mutual followings between user1 and user2
	mutualFollowings, err := suite.repo.GetMutualFollowings(suite.ctx, users[0], users[1])
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), mutualFollowings, 1) // Both follow user3

	// Then - verify specific relationships
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, users[0], users[2])
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing)

	isFollowing, err = suite.repo.IsFollowing(suite.ctx, users[2], users[0])
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), isFollowing) // user3 doesn't follow user1
}

func (suite *SubscriptionRepositoryTestSuite) TestCreateAndDeleteSequence() {
	// Given
	followerID := "user1"
	followingID := "user2"

	// Create subscription
	err := suite.repo.CreateSubscription(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)

	// Verify creation
	isFollowing, err := suite.repo.IsFollowing(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), isFollowing)

	// Delete subscription
	err = suite.repo.DeleteSubscription(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)

	// Verify deletion
	isFollowing, err = suite.repo.IsFollowing(suite.ctx, followerID, followingID)
	require.NoError(suite.T(), err)
	assert.False(suite.T(), isFollowing)

	// Try to delete again (should fail)
	err = suite.repo.DeleteSubscription(suite.ctx, followerID, followingID)
	assert.Error(suite.T(), err)
}

func TestSubscriptionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionRepositoryTestSuite))
}
