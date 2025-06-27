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

func TestCommentRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	repo := db.NewCommentRepository(testDB.Session)

	t.Run("Save Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content := "This is a test comment"

		// When
		comment, err := repo.Save(ctx, authorID, postID, content)

		// Then
		assert.NoError(t, err)
		assert.NotEmpty(t, comment.ID)
		assert.Equal(t, authorID, comment.AuthorID)
		assert.Equal(t, postID, comment.PostID)
		assert.Equal(t, content, comment.Content)
		assert.False(t, comment.CreatedAt.IsZero())
		assert.False(t, comment.UpdatedAt.IsZero())
	})

	t.Run("Save AuthorID Empty", func(t *testing.T) {
		// Given
		authorID := ""
		postID := uuid.New().String()
		content := "This is a test comment"

		// When
		_, err = repo.Save(ctx, authorID, postID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "author ID is required")
	})

	t.Run("Save PostID Empty", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := ""
		content := "This is a test comment"

		// When
		_, err = repo.Save(ctx, authorID, postID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("Save Content Empty", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content := ""

		// When
		_, err = repo.Save(ctx, authorID, postID, content)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "content is required")
	})

	t.Run("SaveComment Success", func(t *testing.T) {
		// Given
		comment := &app.Comment{
			ID:       uuid.New().String(),
			AuthorID: "author123",
			PostID:   uuid.New().String(),
			Content:  "This is a test comment",
		}

		// When
		err = repo.SaveComment(ctx, comment)

		// Then
		assert.NoError(t, err)
		assert.False(t, comment.CreatedAt.IsZero())
		assert.False(t, comment.UpdatedAt.IsZero())
	})

	t.Run("SaveComment Comment Nil", func(t *testing.T) {
		// When
		err = repo.SaveComment(ctx, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment cannot be nil")
	})

	t.Run("GetById Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content := "This is a test comment"

		// Save comment first
		savedComment, err := repo.Save(ctx, authorID, postID, content)
		require.NoError(t, err)

		// When
		comment, err := repo.GetById(ctx, savedComment.ID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, savedComment.ID, comment.ID)
		assert.Equal(t, authorID, comment.AuthorID)
		assert.Equal(t, postID, comment.PostID)
		assert.Equal(t, content, comment.Content)
		assert.Equal(t, savedComment.CreatedAt.Unix(), comment.CreatedAt.Unix())
	})

	t.Run("GetById NotFound", func(t *testing.T) {
		// Given
		commentID := uuid.New().String()

		// When
		_, err = repo.GetById(ctx, commentID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment not found")
	})

	t.Run("GetById CommentID Empty", func(t *testing.T) {
		// When
		_, err = repo.GetById(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment ID is required")
	})

	t.Run("GetByPost Success", func(t *testing.T) {
		// Given
		authorID1 := "author1"
		authorID2 := "author2"
		postID := uuid.New().String()
		content1 := "First comment"
		content2 := "Second comment"
		content3 := "Third comment"

		// Save comments with different timestamps
		_, err = repo.Save(ctx, authorID1, postID, content1)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = repo.Save(ctx, authorID2, postID, content2)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)

		_, err = repo.Save(ctx, authorID1, postID, content3)
		require.NoError(t, err)

		// When
		comments, err := repo.GetByPost(ctx, postID, 10)

		// Then
		assert.NoError(t, err)
		assert.Len(t, comments, 3)

		// Comments should be ordered by created_at DESC
		assert.Equal(t, content3, comments[0].Content)
		assert.Equal(t, content2, comments[1].Content)
		assert.Equal(t, content1, comments[2].Content)
	})

	t.Run("GetByPost Limit", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()

		// Save 10 comments
		for i := 0; i < 10; i++ {
			_, err := repo.Save(ctx, authorID, postID, "Comment "+string(rune('0'+i)))
			require.NoError(t, err)
			time.Sleep(10 * time.Millisecond)
		}

		// When
		comments, err := repo.GetByPost(ctx, postID, 5)

		// Then
		assert.NoError(t, err)
		assert.Len(t, comments, 5)
	})

	t.Run("GetByPost Empty", func(t *testing.T) {
		// Given
		postID := uuid.New().String()

		// When
		comments, err := repo.GetByPost(ctx, postID, 10)

		// Then
		assert.NoError(t, err)
		assert.Empty(t, comments)
	})

	t.Run("GetByPost PostID Empty", func(t *testing.T) {
		// When
		_, err = repo.GetByPost(ctx, "", 10)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("GetAll Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content1 := "First comment"
		content2 := "Second comment"

		// Save comments
		_, err = repo.Save(ctx, authorID, postID, content1)
		require.NoError(t, err)

		_, err = repo.Save(ctx, authorID, postID, content2)
		require.NoError(t, err)

		// When (using legacy method)
		comments, err := repo.GetAll(ctx, postID)

		// Then
		assert.NoError(t, err)
		assert.Len(t, comments, 2)
	})

	t.Run("Update Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		originalContent := "Original content"
		newContent := "Updated content"

		// Save comment first
		savedComment, err := repo.Save(ctx, authorID, postID, originalContent)
		require.NoError(t, err)

		// When
		updatedComment, err := repo.Update(ctx, savedComment.ID, newContent)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, savedComment.ID, updatedComment.ID)
		assert.Equal(t, newContent, updatedComment.Content)
		assert.True(t, updatedComment.UpdatedAt.After(savedComment.CreatedAt))

		// Verify by getting the comment
		retrievedComment, err := repo.GetById(ctx, savedComment.ID)
		require.NoError(t, err)
		assert.Equal(t, newContent, retrievedComment.Content)
	})

	t.Run("Update CommentNotFound", func(t *testing.T) {
		// Given
		commentID := uuid.New().String()
		newContent := "Updated content"

		// When
		_, err = repo.Update(ctx, commentID, newContent)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment not found")
	})

	t.Run("Update CommentID Empty", func(t *testing.T) {
		// When
		_, err = repo.Update(ctx, "", "Updated content")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment ID is required")
	})

	t.Run("Update Content Empty", func(t *testing.T) {
		// Given
		commentID := uuid.New().String()

		// When
		_, err = repo.Update(ctx, commentID, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "content is required")
	})

	t.Run("Delete Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content := "This is a test comment"

		// Save comment first
		savedComment, err := repo.Save(ctx, authorID, postID, content)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, authorID, savedComment.ID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(ctx, savedComment.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment not found")
	})

	t.Run("Delete CommentNotFound", func(t *testing.T) {
		// Given
		userID := "user123"
		commentID := uuid.New().String()

		// When
		err = repo.Delete(ctx, userID, commentID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment not found")
	})

	t.Run("Delete UnauthorizedUser", func(t *testing.T) {
		// Given
		authorID := "author123"
		otherUserID := "other456"
		postID := uuid.New().String()
		content := "This is a test comment"

		// Save comment first
		savedComment, err := repo.Save(ctx, authorID, postID, content)
		require.NoError(t, err)

		// When - try to delete with different user
		err = repo.Delete(ctx, otherUserID, savedComment.ID)

		// Then
		assert.Error(t, err)
		// This would typically be a forbidden error
	})

	t.Run("Delete UserID Empty", func(t *testing.T) {
		// When
		err = repo.Delete(ctx, "", "comment123")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Delete CommentID Empty", func(t *testing.T) {
		// When
		err = repo.Delete(ctx, "user123", "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment ID is required")
	})

	t.Run("Exists True", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		content := "This is a test comment"

		// Save comment first
		savedComment, err := repo.Save(ctx, authorID, postID, content)
		require.NoError(t, err)

		// When
		exists, err := repo.Exists(ctx, savedComment.ID)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists False", func(t *testing.T) {
		// Given
		commentID := uuid.New().String()

		// When
		exists, err := repo.Exists(ctx, commentID)

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists CommentID Empty", func(t *testing.T) {
		// When
		_, err = repo.Exists(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "comment ID is required")
	})

	t.Run("CountByPost Success", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()

		// Save multiple comments
		for i := 0; i < 5; i++ {
			_, err := repo.Save(ctx, authorID, postID, "Comment "+string(rune('0'+i)))
			require.NoError(t, err)
		}

		// When
		count, err := repo.CountByPost(ctx, postID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 5, count)
	})

	t.Run("CountByPost Empty", func(t *testing.T) {
		// Given
		postID := uuid.New().String()

		// When
		count, err := repo.CountByPost(ctx, postID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("CountByPost PostID Empty", func(t *testing.T) {
		// When
		_, err = repo.CountByPost(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "post ID is required")
	})

	t.Run("ConcurrentOperations", func(t *testing.T) {
		runner := NewTestRunner()

		// Given
		authorID1 := "author1"
		authorID2 := "author2"
		postID := uuid.New().String()
		content1 := "First comment"
		content2 := "Second comment"

		// When - concurrent saves
		runner.RunConcurrently(t,
			func() error {
				_, err := repo.Save(ctx, authorID1, postID, content1)
				return err
			},
			func() error {
				_, err := repo.Save(ctx, authorID2, postID, content2)
				return err
			},
		)

		// Then - verify both comments exist
		comments, err := repo.GetByPost(ctx, postID, 10)
		assert.NoError(t, err)
		assert.Len(t, comments, 2)
	})

	t.Run("LargeContentHandling", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		largeContent := make([]byte, 5000)

		// Fill with valid characters
		for i := range largeContent {
			largeContent[i] = 'A'
		}

		// When
		comment, err := repo.Save(ctx, authorID, postID, string(largeContent))

		// Then
		assert.NoError(t, err)
		assert.Equal(t, string(largeContent), comment.Content)

		// Verify by retrieving
		retrievedComment, err := repo.GetById(ctx, comment.ID)
		assert.NoError(t, err)
		assert.Equal(t, string(largeContent), retrievedComment.Content)
	})

	t.Run("UpdateAndDeleteSequence", func(t *testing.T) {
		// Given
		authorID := "author123"
		postID := uuid.New().String()
		originalContent := "Original content"
		updatedContent := "Updated content"

		// Save comment
		savedComment, err := repo.Save(ctx, authorID, postID, originalContent)
		require.NoError(t, err)

		// Update comment
		updatedComment, err := repo.Update(ctx, savedComment.ID, updatedContent)
		require.NoError(t, err)
		assert.Equal(t, updatedContent, updatedComment.Content)

		// Delete comment
		err = repo.Delete(ctx, authorID, savedComment.ID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(ctx, savedComment.ID)
		assert.Error(t, err)
	})
}
