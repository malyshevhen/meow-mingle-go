package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	repo := db.NewPostRepository(testDB.Session)

	// Helper function to clean the database before each test
	setupTest := func(t *testing.T) {
		err := testDB.Clean(ctx)
		require.NoError(t, err, "Failed to clean test database")
	}

	t.Run("Save Success", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content := "This is a test post"

		// When
		post, err := repo.Save(ctx, authorID, content)

		// Then
		assert.NoError(t, err)
		assert.NotEmpty(t, post.ID)
		assert.Equal(t, authorID, post.AuthorID)
		assert.Equal(t, content, post.Content)
		assert.False(t, post.CreatedAt.IsZero())
		assert.False(t, post.UpdatedAt.IsZero())
	})

	t.Run("Save Validation Error Empty AuthorID", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := ""
		content := "This is a test post"

		// When
		_, err := repo.Save(ctx, authorID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("Save Validation Error Empty Content", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content := ""

		// When
		_, err := repo.Save(ctx, authorID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "content is required")
	})

	t.Run("SavePost Success", func(t *testing.T) {
		setupTest(t)
		// Given
		post := &app.Post{
			ID:       uuid.New().String(),
			AuthorID: "author123",
			Content:  "This is a test post",
		}

		// When
		err := repo.SavePost(ctx, post)

		// Then
		assert.NoError(t, err)
		assert.False(t, post.CreatedAt.IsZero())
		assert.False(t, post.UpdatedAt.IsZero())
	})

	t.Run("SavePost Validation Error Nil Post", func(t *testing.T) {
		setupTest(t)
		// When
		err := repo.SavePost(ctx, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post cannot be nil")
	})

	t.Run("Get Success", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content := "This is a test post"

		// Save post first
		savedPost, err := repo.Save(ctx, authorID, content)
		require.NoError(t, err)

		// When
		post, err := repo.Get(ctx, savedPost.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, savedPost.ID, post.ID)
		assert.Equal(t, authorID, post.AuthorID)
		assert.Equal(t, content, post.Content)
		assert.Equal(t, savedPost.CreatedAt.Unix(), post.CreatedAt.Unix())
	})

	t.Run("Get Not Found", func(t *testing.T) {
		setupTest(t)
		// Given
		postID := uuid.New().String()

		// When
		_, err := repo.Get(ctx, postID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post not found")
	})

	t.Run("Get Validation Error Empty PostID", func(t *testing.T) {
		setupTest(t)
		// When
		_, err := repo.Get(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("List Success", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content1 := "First post"
		content2 := "Second post"
		content3 := "Third post"

		// Save posts
		_, err := repo.Save(ctx, authorID, content1)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond) // Ensure different timestamps

		_, err = repo.Save(ctx, authorID, content2)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = repo.Save(ctx, authorID, content3)
		require.NoError(t, err)

		// When
		posts, err := repo.List(ctx, authorID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, posts, 3)

		// Posts should be ordered by created_at DESC
		assert.Equal(t, content3, posts[0].Content)
		assert.Equal(t, content2, posts[1].Content)
		assert.Equal(t, content1, posts[2].Content)
	})

	t.Run("List Empty Result", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "nonexistent"

		// When
		posts, err := repo.List(ctx, authorID)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("GetByAuthor With Limit", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"

		// Save 10 posts
		for i := range 10 {
			_, err := repo.Save(ctx, authorID, "Post "+string(rune('0'+i)))
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}

		// When
		posts, err := repo.GetByAuthor(ctx, authorID, 5)

		// Then
		assert.NoError(t, err)
		assert.Len(t, posts, 5)
	})

	t.Run("GetByAuthor Validation Error Empty AuthorID", func(t *testing.T) {
		setupTest(t)
		// When
		_, err := repo.GetByAuthor(ctx, "", 10)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("Update Success", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		originalContent := "Original content"
		newContent := "Updated content"

		// Save post first
		savedPost, err := repo.Save(ctx, authorID, originalContent)
		require.NoError(t, err)

		// When
		updatedPost, err := repo.Update(ctx, savedPost.ID, newContent)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, savedPost.ID, updatedPost.ID)
		assert.Equal(t, newContent, updatedPost.Content)
		assert.True(t, updatedPost.UpdatedAt.After(savedPost.CreatedAt))

		// Verify by getting the post
		retrievedPost, err := repo.Get(ctx, savedPost.ID)
		require.NoError(t, err)
		assert.Equal(t, newContent, retrievedPost.Content)
	})

	t.Run("Update Post Not Found", func(t *testing.T) {
		setupTest(t)
		// Given
		postID := uuid.New().String()
		newContent := "Updated content"

		// When
		_, err := repo.Update(ctx, postID, newContent)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post not found")
	})

	t.Run("Update Validation Error Empty PostID", func(t *testing.T) {
		setupTest(t)
		// When
		_, err := repo.Update(ctx, "", "Updated content")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("Update Validation Error Empty Content", func(t *testing.T) {
		setupTest(t)
		// Given
		postID := uuid.New().String()

		// When
		_, err := repo.Update(ctx, postID, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "content is required")
	})

	t.Run("Delete Success", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content := "This is a test post"

		// Save post first
		savedPost, err := repo.Save(ctx, authorID, content)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, savedPost.ID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.Get(ctx, savedPost.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post not found")
	})

	t.Run("Delete Post Not Found", func(t *testing.T) {
		setupTest(t)
		// Given
		postID := uuid.New().String()

		// When
		err := repo.Delete(ctx, postID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post not found")
	})

	t.Run("Delete Validation Error Empty PostID", func(t *testing.T) {
		setupTest(t)
		// When
		err := repo.Delete(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("Exists True", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		content := "This is a test post"

		// Save post first
		savedPost, err := repo.Save(ctx, authorID, content)
		require.NoError(t, err)

		// When
		exists, err := repo.Exists(ctx, savedPost.ID)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists False", func(t *testing.T) {
		setupTest(t)
		// Given
		postID := uuid.New().String()

		// When
		exists, err := repo.Exists(ctx, postID)

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists Validation Error Empty PostID", func(t *testing.T) {
		setupTest(t)
		// When
		_, err := repo.Exists(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("Feed Empty Result", func(t *testing.T) {
		setupTest(t)
		// Given
		userID := "user123"

		// When
		posts, err := repo.Feed(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, posts)
	})

	t.Run("Feed Validation Error Empty UserID", func(t *testing.T) {
		setupTest(t)
		// When
		_, err := repo.Feed(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		setupTest(t)
		runner := NewTestRunner()

		// Given
		authorID := "author123"
		content1 := "First post"
		content2 := "Second post"

		// When - concurrent saves
		runner.RunConcurrently(t,
			func() error {
				_, err := repo.Save(ctx, authorID, content1)
				return err
			},
			func() error {
				_, err := repo.Save(ctx, authorID, content2)
				return err
			},
		)

		// Then
		posts, err := repo.List(ctx, authorID)
		assert.NoError(t, err)
		assert.Len(t, posts, 2)
	})

	t.Run("SavePost Duplicate ID", func(t *testing.T) {
		t.Skip()

		setupTest(t)
		// Given
		postID := uuid.New().String()
		post1 := &app.Post{
			ID:       postID,
			AuthorID: "author1",
			Content:  "First post",
		}
		post2 := &app.Post{
			ID:       postID,
			AuthorID: "author2",
			Content:  "Second post",
		}

		// Save first post
		err := repo.SavePost(ctx, post1)
		require.NoError(t, err)

		// When - try to save with same ID
		err = repo.SavePost(ctx, post2)

		// Then - should fail due to primary key constraint
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Primary key already exists")
	})

	t.Run("Large Content Handling", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		largeContent := make([]byte, 10000) // 10KB content

		// Fill with valid characters
		for i := range largeContent {
			largeContent[i] = 'A'
		}

		// When
		post, err := repo.Save(ctx, authorID, string(largeContent))

		// Then
		assert.NoError(t, err)
		assert.Equal(t, string(largeContent), post.Content)

		// Verify by retrieving
		retrievedPost, err := repo.Get(ctx, post.ID)
		assert.NoError(t, err)
		assert.Equal(t, string(largeContent), retrievedPost.Content)
	})

	t.Run("Update And Delete Sequence", func(t *testing.T) {
		setupTest(t)
		// Given
		authorID := "author123"
		originalContent := "Original content"
		updatedContent := "Updated content"

		// Save post
		savedPost, err := repo.Save(ctx, authorID, originalContent)
		require.NoError(t, err)

		// Update post
		updatedPost, err := repo.Update(ctx, savedPost.ID, updatedContent)
		require.NoError(t, err)
		assert.Equal(t, updatedContent, updatedPost.Content)

		// Delete post
		err = repo.Delete(ctx, savedPost.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.Get(ctx, savedPost.ID)
		assert.Error(t, err)
	})

	t.Run("Multiple Authors Scenario", func(t *testing.T) {
		setupTest(t)
		// Given
		author1 := "author1"
		author2 := "author2"

		// Save posts from different authors
		_, err := repo.Save(ctx, author1, "Author 1 Post 1")
		require.NoError(t, err)

		_, err = repo.Save(ctx, author2, "Author 2 Post 1")
		require.NoError(t, err)

		_, err = repo.Save(ctx, author1, "Author 1 Post 2")
		require.NoError(t, err)

		// When
		author1Posts, err := repo.List(ctx, author1)
		require.NoError(t, err)

		author2Posts, err := repo.List(ctx, author2)
		require.NoError(t, err)

		// Then
		assert.Len(t, author1Posts, 2)
		assert.Len(t, author2Posts, 1)

		// Verify posts belong to correct authors
		for _, post := range author1Posts {
			assert.Equal(t, author1, post.AuthorID)
		}
		for _, post := range author2Posts {
			assert.Equal(t, author2, post.AuthorID)
		}
	})
}
