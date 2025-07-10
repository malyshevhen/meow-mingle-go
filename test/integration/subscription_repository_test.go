package integration

import (
	"context"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscriptionRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err, "Failed to create test database")
	defer testDB.Close(ctx)

	repo := db.NewSubscriptionRepository(testDB.Session)

	t.Run("CreateSubscription", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"
			followingID := "user2"

			// When
			err := repo.CreateSubscription(ctx, followerID, followingID)

			// Then
			assert.NoError(t, err)

			// Verify subscription exists
			isFollowing, err := repo.IsFollowing(ctx, followerID, followingID)
			assert.NoError(t, err)
			assert.True(t, isFollowing)
		})

		t.Run("ValidationError_EmptyFollowerID", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := ""
			followingID := "user2"

			// When
			err := repo.CreateSubscription(ctx, followerID, followingID)

			// Then
			assert.Error(t, err)
			assert.Equal(t, "follower ID is required", err.Error())
		})
	})

	t.Run("GetFollowers", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followingID := "user1"
			follower1ID := "user2"
			follower2ID := "user3"
			follower3ID := "user4"

			// Create subscriptions
			err := repo.CreateSubscription(ctx, follower1ID, followingID)
			require.NoError(t, err)

			err = repo.CreateSubscription(ctx, follower2ID, followingID)
			require.NoError(t, err)

			err = repo.CreateSubscription(ctx, follower3ID, followingID)
			require.NoError(t, err)

			// When
			followers, err := repo.GetFollowers(ctx, followingID, 10)

			// Then
			assert.NoError(t, err)
			assert.Len(t, followers, 3)

			// Verify all followers are present
			followerIDs := make(map[string]bool)
			for _, subscription := range followers {
				followerIDs[subscription.FollowerID] = true
				assert.Equal(t, followingID, subscription.FollowingID)
				assert.False(t, subscription.CreatedAt.IsZero())
			}

			assert.True(t, followerIDs[follower1ID])
			assert.True(t, followerIDs[follower2ID])
			assert.True(t, followerIDs[follower3ID])
		})

		t.Run("WithLimit", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followingID := "user1"

			// Create 10 followers
			for i := range 10 {
				followerID := "follower" + string(rune('0'+i))
				err := repo.CreateSubscription(ctx, followerID, followingID)
				require.NoError(t, err)
			}

			// When
			followers, err := repo.GetFollowers(ctx, followingID, 5)

			// Then
			assert.NoError(t, err)
			assert.Len(t, followers, 5)
		})

		t.Run("EmptyResult", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followingID := "user1"

			// When
			followers, err := repo.GetFollowers(ctx, followingID, 10)

			// Then
			assert.NoError(t, err)
			assert.Empty(t, followers)
		})

		t.Run("ValidationError_EmptyFollowingID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.GetFollowers(ctx, "", 10)

			// Then
			assert.Error(t, err)
			assert.Equal(t, "following ID is required", err.Error())
		})
	})

	t.Run("GetFollowing", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"
			following1ID := "user2"
			following2ID := "user3"
			following3ID := "user4"

			// Create subscriptions
			err := repo.CreateSubscription(ctx, followerID, following1ID)
			require.NoError(t, err)

			err = repo.CreateSubscription(ctx, followerID, following2ID)
			require.NoError(t, err)

			err = repo.CreateSubscription(ctx, followerID, following3ID)
			require.NoError(t, err)

			// When
			following, err := repo.GetFollowing(ctx, followerID, 10)

			// Then
			assert.NoError(t, err)
			assert.Len(t, following, 3)

			// Verify all following are present
			followingIDs := make(map[string]bool)
			for _, subscription := range following {
				followingIDs[subscription.FollowingID] = true
				assert.Equal(t, followerID, subscription.FollowerID)
				assert.False(t, subscription.CreatedAt.IsZero())
			}

			assert.True(t, followingIDs[following1ID])
			assert.True(t, followingIDs[following2ID])
			assert.True(t, followingIDs[following3ID])
		})

		t.Run("WithLimit", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"

			// Create 10 following relationships
			for i := 0; i < 10; i++ {
				followingID := "following" + string(rune('0'+i))
				err := repo.CreateSubscription(ctx, followerID, followingID)
				require.NoError(t, err)
			}

			// When
			following, err := repo.GetFollowing(ctx, followerID, 5)

			// Then
			assert.NoError(t, err)
			assert.Len(t, following, 5)
		})

		t.Run("EmptyResult", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"

			// When
			following, err := repo.GetFollowing(ctx, followerID, 10)

			// Then
			assert.NoError(t, err)
			assert.Empty(t, following)
		})

		t.Run("ValidationError_EmptyFollowerID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.GetFollowing(ctx, "", 10)

			// Then
			assert.Error(t, err)
			assert.Equal(t, "follower ID is required", err.Error())
		})
	})

	t.Run("IsFollowing", func(t *testing.T) {
		t.Run("True", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"
			followingID := "user2"

			// Create subscription
			err := repo.CreateSubscription(ctx, followerID, followingID)
			require.NoError(t, err)

			// When
			isFollowing, err := repo.IsFollowing(ctx, followerID, followingID)

			// Then
			assert.NoError(t, err)
			assert.True(t, isFollowing)
		})

		t.Run("False", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"
			followingID := "user2"

			// When
			isFollowing, err := repo.IsFollowing(ctx, followerID, followingID)

			// Then
			assert.NoError(t, err)
			assert.False(t, isFollowing)
		})

		t.Run("ValidationError_EmptyFollowerID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.IsFollowing(ctx, "", "user2")

			// Then
			assert.Error(t, err)
			assert.Equal(t, "follower ID is required", err.Error())
		})

		t.Run("ValidationError_EmptyFollowingID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.IsFollowing(ctx, "user1", "")

			// Then
			assert.Error(t, err)
			assert.Equal(t, "following ID is required", err.Error())
		})
	})

	t.Run("CountFollowers", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followingID := "user1"

			// Create multiple followers
			for i := 0; i < 5; i++ {
				followerID := "follower" + string(rune('0'+i))
				err := repo.CreateSubscription(ctx, followerID, followingID)
				require.NoError(t, err)
			}

			// When
			count, err := repo.CountFollowers(ctx, followingID)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 5, count)
		})

		t.Run("EmptyResult", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followingID := "user1"

			// When
			count, err := repo.CountFollowers(ctx, followingID)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})

		t.Run("ValidationError_EmptyFollowingID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.CountFollowers(ctx, "")

			// Then
			assert.Error(t, err)
			assert.Equal(t, "following ID is required", err.Error())
		})
	})

	t.Run("CountFollowing", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"

			// Create multiple following relationships
			for i := 0; i < 5; i++ {
				followingID := "following" + string(rune('0'+i))
				err := repo.CreateSubscription(ctx, followerID, followingID)
				require.NoError(t, err)
			}

			// When
			count, err := repo.CountFollowing(ctx, followerID)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 5, count)
		})

		t.Run("EmptyResult", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			followerID := "user1"

			// When
			count, err := repo.CountFollowing(ctx, followerID)

			// Then
			assert.NoError(t, err)
			assert.Equal(t, 0, count)
		})

		t.Run("ValidationError_EmptyFollowerID", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.CountFollowing(ctx, "")

			// Then
			assert.Error(t, err)
			assert.Equal(t, "follower ID is required", err.Error())
		})
	})

	t.Run("GetMutualFollowings", func(t *testing.T) {
		t.Run("Success", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			user1 := "user1"
			user2 := "user2"
			commonUser1 := "common1"
			commonUser2 := "common2"
			onlyUser1 := "only1"
			onlyUser2 := "only2"

			// Create followings for user1
			err := repo.CreateSubscription(ctx, user1, commonUser1)
			require.NoError(t, err)
			err = repo.CreateSubscription(ctx, user1, commonUser2)
			require.NoError(t, err)
			err = repo.CreateSubscription(ctx, user1, onlyUser1)
			require.NoError(t, err)

			// Create followings for user2
			err = repo.CreateSubscription(ctx, user2, commonUser1)
			require.NoError(t, err)
			err = repo.CreateSubscription(ctx, user2, commonUser2)
			require.NoError(t, err)
			err = repo.CreateSubscription(ctx, user2, onlyUser2)
			require.NoError(t, err)

			// When
			mutualFollowings, err := repo.GetMutualFollowings(ctx, user1, user2)

			// Then
			assert.NoError(t, err)
			assert.Len(t, mutualFollowings, 2)

			// Verify mutual followings
			mutualIDs := make(map[string]bool)
			for _, subscription := range mutualFollowings {
				mutualIDs[subscription.FollowingID] = true
			}

			assert.True(t, mutualIDs[commonUser1])
			assert.True(t, mutualIDs[commonUser2])
			assert.False(t, mutualIDs[onlyUser1])
			assert.False(t, mutualIDs[onlyUser2])
		})

		t.Run("EmptyResult", func(t *testing.T) {
			testDB.Clean(ctx)
			// Given
			user1 := "user1"
			user2 := "user2"

			// When
			mutualFollowings, err := repo.GetMutualFollowings(ctx, user1, user2)

			// Then
			assert.NoError(t, err)
			assert.Empty(t, mutualFollowings)
		})

		t.Run("ValidationError_EmptyUser1", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.GetMutualFollowings(ctx, "", "user2")

			// Then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user ID 1 is required")
		})

		t.Run("ValidationError_EmptyUser2", func(t *testing.T) {
			testDB.Clean(ctx)
			// When
			_, err := repo.GetMutualFollowings(ctx, "user1", "")

			// Then
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "user ID 2 is required")
		})
	})

	t.Run("BidirectionalFollowing", func(t *testing.T) {
		testDB.Clean(ctx)
		// Given
		user1 := "user1"
		user2 := "user2"

		// When - create bidirectional following
		err := repo.CreateSubscription(ctx, user1, user2)
		require.NoError(t, err)

		err = repo.CreateSubscription(ctx, user2, user1)
		require.NoError(t, err)

		// Then
		// Verify user1 follows user2
		isFollowing1, err := repo.IsFollowing(ctx, user1, user2)
		assert.NoError(t, err)
		assert.True(t, isFollowing1)

		// Verify user2 follows user1
		isFollowing2, err := repo.IsFollowing(ctx, user2, user1)
		assert.NoError(t, err)
		assert.True(t, isFollowing2)

		// Verify counts
		user1Followers, err := repo.CountFollowers(ctx, user1)
		assert.NoError(t, err)
		assert.Equal(t, 1, user1Followers)

		user2Followers, err := repo.CountFollowers(ctx, user2)
		assert.NoError(t, err)
		assert.Equal(t, 1, user2Followers)
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		testDB.Clean(ctx)
		// Given
		follower1 := "user1"
		follower2 := "user2"
		following := "user3"

		// When - concurrent subscriptions
		done1 := make(chan error, 1)
		done2 := make(chan error, 1)

		go func() {
			err := repo.CreateSubscription(ctx, follower1, following)
			done1 <- err
		}()

		go func() {
			err := repo.CreateSubscription(ctx, follower2, following)
			done2 <- err
		}()

		// Then
		err1 := <-done1
		err2 := <-done2

		assert.NoError(t, err1)
		assert.NoError(t, err2)

		// Verify both subscriptions exist
		count, err := repo.CountFollowers(ctx, following)
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("ComplexFollowingNetwork", func(t *testing.T) {
		testDB.Clean(ctx)
		// Given - create a network of users
		users := []string{"user1", "user2", "user3", "user4", "user5"}

		// Create a complex following network
		// user1 follows user2, user3
		err := repo.CreateSubscription(ctx, users[0], users[1])
		require.NoError(t, err)
		err = repo.CreateSubscription(ctx, users[0], users[2])
		require.NoError(t, err)

		// user2 follows user3, user4
		err = repo.CreateSubscription(ctx, users[1], users[2])
		require.NoError(t, err)
		err = repo.CreateSubscription(ctx, users[1], users[3])
		require.NoError(t, err)

		// user3 follows user4, user5
		err = repo.CreateSubscription(ctx, users[2], users[3])
		require.NoError(t, err)
		err = repo.CreateSubscription(ctx, users[2], users[4])
		require.NoError(t, err)

		// When - verify the network
		// Check user1's following
		user1Following, err := repo.GetFollowing(ctx, users[0], 10)
		require.NoError(t, err)
		assert.Len(t, user1Following, 2)

		// Check user3's followers
		user3Followers, err := repo.GetFollowers(ctx, users[2], 10)
		require.NoError(t, err)
		assert.Len(t, user3Followers, 2)

		// Check mutual followings between user1 and user2
		mutualFollowings, err := repo.GetMutualFollowings(ctx, users[0], users[1])
		require.NoError(t, err)
		assert.Len(t, mutualFollowings, 1) // Both follow user3

		// Then - verify specific relationships
		isFollowing, err := repo.IsFollowing(ctx, users[0], users[2])
		assert.NoError(t, err)
		assert.True(t, isFollowing)

		isFollowing, err = repo.IsFollowing(ctx, users[2], users[0])
		assert.NoError(t, err)
		assert.False(t, isFollowing) // user3 doesn't follow user1
	})

	t.Run("CreateAndDeleteSequence", func(t *testing.T) {
		testDB.Clean(ctx)
		// Given
		followerID := "user1"
		followingID := "user2"

		// Create subscription
		err := repo.CreateSubscription(ctx, followerID, followingID)
		require.NoError(t, err)

		// Verify creation
		isFollowing, err := repo.IsFollowing(ctx, followerID, followingID)
		require.NoError(t, err)
		assert.True(t, isFollowing)

		// Delete subscription
		err = repo.DeleteSubscription(ctx, followerID, followingID)
		require.NoError(t, err)

		// Verify deletion
		isFollowing, err = repo.IsFollowing(ctx, followerID, followingID)
		require.NoError(t, err)
		assert.False(t, isFollowing)

		// Try to delete again (should fail)
		err = repo.DeleteSubscription(ctx, followerID, followingID)
		assert.Error(t, err)
	})
}
