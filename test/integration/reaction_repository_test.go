package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReactionRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	repo := db.NewReactionRepository(testDB.Session)
	runner := NewTestRunner() // For concurrent operations

	t.Run("Save Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		targetID := uuid.New().String()
		authorID := "author123"
		content := "like"

		// When
		err = repo.Save(ctx, targetID, authorID, content)

		// Then
		assert.NoError(t, err)

		// Verify reaction exists
		exists, err := repo.Exists(ctx, targetID, authorID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Save Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := ""
		authorID := "author123"
		content := "like"

		// When
		err = repo.Save(ctx, targetID, authorID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("Save Validation Error - Empty AuthorID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := ""
		content := "like"

		// When
		err = repo.Save(ctx, targetID, authorID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("Save Validation Error - Empty Content", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"
		content := ""

		// When
		err = repo.Save(ctx, targetID, authorID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reaction content is required")
	})

	t.Run("SaveReaction Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		reaction := &app.Reaction{
			TargetID: uuid.New().String(),
			AuthorID: "author123",
			Content:  "love",
		}

		// When
		err = repo.SaveReaction(ctx, reaction)

		// Then
		assert.NoError(t, err)
		assert.False(t, reaction.CreatedAt.IsZero())
		assert.False(t, reaction.UpdatedAt.IsZero())
	})

	t.Run("SaveReaction Validation Error - Nil Reaction", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		err = repo.SaveReaction(ctx, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reaction cannot be nil")
	})

	t.Run("Delete Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"
		content := "like"

		// Save reaction first
		err = repo.Save(ctx, targetID, authorID, content)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, targetID, authorID)

		// Then
		assert.NoError(t, err)

		// Verify reaction no longer exists
		exists, err := repo.Exists(ctx, targetID, authorID)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Delete Not Found", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"

		// When
		err = repo.Delete(ctx, targetID, authorID) // Provide a dummy reaction type for non-existent deletion

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reaction not found")
	})

	t.Run("Delete Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		err = repo.Delete(ctx, "", "author123")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("Delete Validation Error - Empty AuthorID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		err = repo.Delete(ctx, uuid.New().String(), "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("GetByTarget Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"
		author1ID := "author1"
		author2ID := "author2"
		author3ID := "author3"

		// Save reactions
		err = repo.Save(ctx, targetID, author1ID, "like")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, author2ID, "love")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, author3ID, "like")
		require.NoError(t, err)

		// When
		reactions, err := repo.GetByTarget(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Len(t, reactions, 3)

		// Verify all reactions are present
		authorIDs := make(map[string]bool)
		reactionTypes := make(map[string]int)
		for _, reaction := range reactions {
			authorIDs[reaction.AuthorID] = true
			reactionTypes[reaction.Content]++
			assert.Equal(t, targetID, reaction.TargetID)
			assert.False(t, reaction.CreatedAt.IsZero())
		}

		assert.True(t, authorIDs[author1ID])
		assert.True(t, authorIDs[author2ID])
		assert.True(t, authorIDs[author3ID])
		assert.Equal(t, 2, reactionTypes["like"])
		assert.Equal(t, 1, reactionTypes["love"])
	})

	t.Run("GetByTarget Empty Result", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"

		// When
		reactions, err := repo.GetByTarget(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, reactions)
	})

	t.Run("GetByTarget Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.GetByTarget(ctx, "", "post")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("GetByAuthor Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		authorID := "author123"
		target1ID := uuid.New().String()
		target2ID := uuid.New().String()
		target3ID := uuid.New().String()

		// Save reactions from the same author
		err = repo.Save(ctx, target1ID, authorID, "like")
		require.NoError(t, err)

		err = repo.Save(ctx, target2ID, authorID, "love")
		require.NoError(t, err)

		err = repo.Save(ctx, target3ID, authorID, "laugh")
		require.NoError(t, err)

		// When
		reactions, err := repo.GetByAuthor(ctx, authorID, 10)

		// Then
		assert.NoError(t, err)
		assert.Len(t, reactions, 3)

		// Verify all reactions belong to the author
		targetIDs := make(map[string]bool)
		for _, reaction := range reactions {
			targetIDs[reaction.TargetID] = true
			assert.Equal(t, authorID, reaction.AuthorID)
			assert.False(t, reaction.CreatedAt.IsZero())
		}

		assert.True(t, targetIDs[target1ID])
		assert.True(t, targetIDs[target2ID])
		assert.True(t, targetIDs[target3ID])
	})

	t.Run("GetByAuthor With Limit", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		authorID := "author123"

		// Save 10 reactions
		for range 10 {
			targetID := uuid.New().String()
			err = repo.Save(ctx, targetID, authorID, "like")
			require.NoError(t, err)
		}

		// When
		reactions, err := repo.GetByAuthor(ctx, authorID, 5)

		// Then
		assert.NoError(t, err)
		assert.Len(t, reactions, 5)
	})

	t.Run("GetByAuthor Empty Result", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		authorID := "nonexistent"

		// When
		reactions, err := repo.GetByAuthor(ctx, authorID, 10)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, reactions)
	})

	t.Run("GetByAuthor Validation Error - Empty AuthorID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.GetByAuthor(ctx, "", 10)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("Exists True", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"
		content := "like"

		// Save reaction first
		err = repo.Save(ctx, targetID, authorID, content)
		require.NoError(t, err)

		// When
		exists, err := repo.Exists(ctx, targetID, authorID)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists False", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"

		// When
		exists, err := repo.Exists(ctx, targetID, authorID)

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.Exists(ctx, "", "author123")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("Exists Validation Error - Empty AuthorID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.Exists(ctx, uuid.New().String(), "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("CountByTarget Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"

		// Save reactions of different types
		err = repo.Save(ctx, targetID, "author1", "like")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author2", "like")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author3", "love")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author4", "laugh")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author5", "like")
		require.NoError(t, err)

		// When
		counts, err := repo.CountByTarget(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Len(t, counts, 3)
		assert.Equal(t, 3, counts["like"])
		assert.Equal(t, 1, counts["love"])
		assert.Equal(t, 1, counts["laugh"])
	})

	t.Run("CountByTarget Empty Result", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"

		// When
		counts, err := repo.CountByTarget(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, counts)
	})

	t.Run("CountByTarget Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.CountByTarget(ctx, "", "post")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("GetReactionTypes Success", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"

		// Save reactions of different types
		err = repo.Save(ctx, targetID, "author1", "like")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author2", "love")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author3", "laugh")
		require.NoError(t, err)

		err = repo.Save(ctx, targetID, "author4", "like") // Duplicate type
		require.NoError(t, err)

		// When
		reactionTypes, err := repo.GetReactionTypes(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Len(t, reactionTypes, 3)

		// Convert to map for easier verification
		typeMap := make(map[string]bool)
		for _, reactionType := range reactionTypes {
			typeMap[reactionType] = true
		}

		assert.True(t, typeMap["like"])
		assert.True(t, typeMap["love"])
		assert.True(t, typeMap["laugh"])
	})

	t.Run("GetReactionTypes Empty Result", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		targetType := "post"

		// When
		reactionTypes, err := repo.GetReactionTypes(ctx, targetID, targetType)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, reactionTypes)
	})

	t.Run("GetReactionTypes Validation Error - Empty TargetID", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// When
		_, err = repo.GetReactionTypes(ctx, "", "post")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target ID is required")
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		author1ID := "author1"
		author2ID := "author2"

		// When - concurrent saves
		runner.RunConcurrently(t,
			func() error {
				err := repo.Save(ctx, targetID, author1ID, "like")
				return err
			},
			func() error {
				err := repo.Save(ctx, targetID, author2ID, "love")
				return err
			},
		)

		// Then
		reactions, err := repo.GetByTarget(ctx, targetID, "post")
		assert.NoError(t, err)
		assert.Len(t, reactions, 2)
	})

	t.Run("Multiple Targets Scenario", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		author1ID := "author1"
		target1ID := uuid.New().String()
		target2ID := uuid.New().String()

		// Save reactions for different targets
		err = repo.Save(ctx, target1ID, author1ID, "like")
		require.NoError(t, err)

		err = repo.Save(ctx, target1ID, "author2", "love")
		require.NoError(t, err)

		err = repo.Save(ctx, target2ID, author1ID, "laugh")
		require.NoError(t, err)

		// When
		target1Reactions, err := repo.GetByTarget(ctx, target1ID, "post")
		require.NoError(t, err)

		target2Reactions, err := repo.GetByTarget(ctx, target2ID, "post")
		require.NoError(t, err)

		author1Reactions, err := repo.GetByAuthor(ctx, author1ID, 10)
		require.NoError(t, err)

		// Then
		assert.Len(t, target1Reactions, 2)
		assert.Len(t, target2Reactions, 1)
		assert.Len(t, author1Reactions, 2)

		// Verify reactions belong to correct targets
		for _, reaction := range target1Reactions {
			assert.Equal(t, target1ID, reaction.TargetID)
		}
		for _, reaction := range target2Reactions {
			assert.Equal(t, target2ID, reaction.TargetID)
		}
		for _, reaction := range author1Reactions {
			assert.Equal(t, author1ID, reaction.AuthorID)
		}
	})

	t.Run("Reaction Overwrite", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"

		// Save initial reaction
		err = repo.Save(ctx, targetID, authorID, "like")
		require.NoError(t, err)

		// When - save different reaction from same author to same target
		err = repo.Save(ctx, targetID, authorID, "love")

		// Then - should succeed (overwrites previous reaction)
		assert.NoError(t, err)

		// Verify only one reaction exists
		reactions, err := repo.GetByTarget(ctx, targetID, "post")
		assert.NoError(t, err)
		assert.Len(t, reactions, 1)
		assert.Equal(t, "love", reactions[0].Content) // Should be the latest reaction
	})

	t.Run("Create and Delete Sequence", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given
		targetID := uuid.New().String()
		authorID := "author123"
		content := "like"

		// Create reaction
		err = repo.Save(ctx, targetID, authorID, content)
		require.NoError(t, err)

		// Verify creation
		exists, err := repo.Exists(ctx, targetID, authorID)
		require.NoError(t, err)
		assert.True(t, exists)

		// Delete reaction
		err = repo.Delete(ctx, targetID, authorID)
		require.NoError(t, err)

		// Verify deletion
		exists, err = repo.Exists(ctx, targetID, authorID)
		require.NoError(t, err)
		assert.False(t, exists)

		// Try to delete again (should fail)
		err = repo.Delete(ctx, targetID, authorID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "reaction not found")
	})

	t.Run("Complex Reaction Scenario", func(t *testing.T) {
		err = testDB.Clean(ctx)
		require.NoError(t, err)
		// Given - create a complex scenario with multiple targets and authors
		post1ID := uuid.New().String()
		post2ID := uuid.New().String()
		comment1ID := uuid.New().String()

		users := []string{"user1", "user2", "user3", "user4", "user5"}
		reactionTypes := []string{"like", "love", "laugh", "angry", "sad"}

		// Create reactions for post1
		for i, user := range users {
			err = repo.Save(ctx, post1ID, user, reactionTypes[i%len(reactionTypes)])
			require.NoError(t, err)
		}

		// Create reactions for post2 (fewer reactions)
		for i := range 3 {
			err = repo.Save(ctx, post2ID, users[i], "like")
			require.NoError(t, err)
		}

		// Create reactions for comment1
		err = repo.Save(ctx, comment1ID, users[0], "love")
		require.NoError(t, err)
		err = repo.Save(ctx, comment1ID, users[1], "love")
		require.NoError(t, err)

		// When - analyze the scenario
		post1Reactions, err := repo.GetByTarget(ctx, post1ID, "post")
		require.NoError(t, err)

		post2Reactions, err := repo.GetByTarget(ctx, post2ID, "post")
		require.NoError(t, err)

		comment1Reactions, err := repo.GetByTarget(ctx, comment1ID, "comment")
		require.NoError(t, err)

		post1Counts, err := repo.CountByTarget(ctx, post1ID, "post")
		require.NoError(t, err)

		user1Reactions, err := repo.GetByAuthor(ctx, users[0], 10)
		require.NoError(t, err)

		// Then - verify the complex scenario
		assert.Len(t, post1Reactions, 5)
		assert.Len(t, post2Reactions, 3)
		assert.Len(t, comment1Reactions, 2)
		assert.Len(t, user1Reactions, 3) // user1 reacted to all three targets

		// Verify post1 has diverse reaction types
		assert.True(t, len(post1Counts) > 1)

		// Verify post2 has uniform reactions
		post2Counts, err := repo.CountByTarget(ctx, post2ID, "post")
		require.NoError(t, err)
		assert.Equal(t, 3, post2Counts["like"])

		// Verify comment1 has specific reaction type
		comment1Types, err := repo.GetReactionTypes(ctx, comment1ID, "comment")
		require.NoError(t, err)
		assert.Len(t, comment1Types, 1)
		assert.Equal(t, "love", comment1Types[0])
	})
}
