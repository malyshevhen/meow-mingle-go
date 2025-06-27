package integration

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type CommentRepositoryTestSuite struct {
	suite.Suite
	testDB *TestDatabase
	repo   db.CommentRepository
	ctx    context.Context
}

func (suite *CommentRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	testDB, err := NewTestDatabase(suite.ctx)
	require.NoError(suite.T(), err, "Failed to create test database")

	suite.testDB = testDB
	suite.repo = db.NewCommentRepository(testDB.Session)
}

func (suite *CommentRepositoryTestSuite) TearDownSuite() {
	if suite.testDB != nil {
		suite.testDB.Close(suite.ctx)
	}
}

func (suite *CommentRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.testDB.Clean(suite.ctx)
	require.NoError(suite.T(), err, "Failed to clean test database")
}

func (suite *CommentRepositoryTestSuite) TestSave_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content := "This is a test comment"

	// When
	comment, err := suite.repo.Save(suite.ctx, authorID, postID, content)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), comment.ID)
	assert.Equal(suite.T(), authorID, comment.AuthorID)
	assert.Equal(suite.T(), postID, comment.PostID)
	assert.Equal(suite.T(), content, comment.Content)
	assert.False(suite.T(), comment.CreatedAt.IsZero())
	assert.False(suite.T(), comment.UpdatedAt.IsZero())
}

func (suite *CommentRepositoryTestSuite) TestSave_ValidationError_EmptyAuthorID() {
	// Given
	authorID := ""
	postID := uuid.New().String()
	content := "This is a test comment"

	// When
	_, err := suite.repo.Save(suite.ctx, authorID, postID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *CommentRepositoryTestSuite) TestSave_ValidationError_EmptyPostID() {
	// Given
	authorID := "author123"
	postID := ""
	content := "This is a test comment"

	// When
	_, err := suite.repo.Save(suite.ctx, authorID, postID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *CommentRepositoryTestSuite) TestSave_ValidationError_EmptyContent() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content := ""

	// When
	_, err := suite.repo.Save(suite.ctx, authorID, postID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "content is required")
}

func (suite *CommentRepositoryTestSuite) TestSaveComment_Success() {
	// Given
	comment := &app.Comment{
		ID:       uuid.New().String(),
		AuthorID: "author123",
		PostID:   uuid.New().String(),
		Content:  "This is a test comment",
	}

	// When
	err := suite.repo.SaveComment(suite.ctx, comment)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), comment.CreatedAt.IsZero())
	assert.False(suite.T(), comment.UpdatedAt.IsZero())
}

func (suite *CommentRepositoryTestSuite) TestSaveComment_ValidationError_NilComment() {
	// When
	err := suite.repo.SaveComment(suite.ctx, nil)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "comment cannot be nil")
}

func (suite *CommentRepositoryTestSuite) TestGetById_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content := "This is a test comment"

	// Save comment first
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, content)
	require.NoError(suite.T(), err)

	// When
	comment, err := suite.repo.GetById(suite.ctx, savedComment.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), savedComment.ID, comment.ID)
	assert.Equal(suite.T(), authorID, comment.AuthorID)
	assert.Equal(suite.T(), postID, comment.PostID)
	assert.Equal(suite.T(), content, comment.Content)
	assert.Equal(suite.T(), savedComment.CreatedAt.Unix(), comment.CreatedAt.Unix())
}

func (suite *CommentRepositoryTestSuite) TestGetById_NotFound() {
	// Given
	commentID := uuid.New().String()

	// When
	_, err := suite.repo.GetById(suite.ctx, commentID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "comment not found")
}

func (suite *CommentRepositoryTestSuite) TestGetById_ValidationError_EmptyCommentID() {
	// When
	_, err := suite.repo.GetById(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "comment ID is required")
}

func (suite *CommentRepositoryTestSuite) TestGetByPost_Success() {
	// Given
	authorID1 := "author1"
	authorID2 := "author2"
	postID := uuid.New().String()
	content1 := "First comment"
	content2 := "Second comment"
	content3 := "Third comment"

	// Save comments
	_, err := suite.repo.Save(suite.ctx, authorID1, postID, content1)
	require.NoError(suite.T(), err)
	time.Sleep(10 * time.Millisecond) // Ensure different timestamps

	_, err = suite.repo.Save(suite.ctx, authorID2, postID, content2)
	require.NoError(suite.T(), err)
	time.Sleep(10 * time.Millisecond)

	_, err = suite.repo.Save(suite.ctx, authorID1, postID, content3)
	require.NoError(suite.T(), err)

	// When
	comments, err := suite.repo.GetByPost(suite.ctx, postID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 3)

	// Comments should be ordered by created_at DESC
	assert.Equal(suite.T(), content3, comments[0].Content)
	assert.Equal(suite.T(), content2, comments[1].Content)
	assert.Equal(suite.T(), content1, comments[2].Content)
}

func (suite *CommentRepositoryTestSuite) TestGetByPost_WithLimit() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()

	// Save 10 comments
	for i := 0; i < 10; i++ {
		_, err := suite.repo.Save(suite.ctx, authorID, postID, "Comment "+string(rune('0'+i)))
		require.NoError(suite.T(), err)
		time.Sleep(10 * time.Millisecond)
	}

	// When
	comments, err := suite.repo.GetByPost(suite.ctx, postID, 5)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 5)
}

func (suite *CommentRepositoryTestSuite) TestGetByPost_EmptyResult() {
	// Given
	postID := uuid.New().String()

	// When
	comments, err := suite.repo.GetByPost(suite.ctx, postID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), comments)
}

func (suite *CommentRepositoryTestSuite) TestGetByPost_ValidationError_EmptyPostID() {
	// When
	_, err := suite.repo.GetByPost(suite.ctx, "", 10)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *CommentRepositoryTestSuite) TestGetAll_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content1 := "First comment"
	content2 := "Second comment"

	// Save comments
	_, err := suite.repo.Save(suite.ctx, authorID, postID, content1)
	require.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, authorID, postID, content2)
	require.NoError(suite.T(), err)

	// When (using legacy method)
	comments, err := suite.repo.GetAll(suite.ctx, postID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 2)
}

func (suite *CommentRepositoryTestSuite) TestUpdate_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	originalContent := "Original content"
	newContent := "Updated content"

	// Save comment first
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, originalContent)
	require.NoError(suite.T(), err)

	// When
	updatedComment, err := suite.repo.Update(suite.ctx, savedComment.ID, newContent)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), savedComment.ID, updatedComment.ID)
	assert.Equal(suite.T(), newContent, updatedComment.Content)
	assert.True(suite.T(), updatedComment.UpdatedAt.After(savedComment.CreatedAt))

	// Verify by getting the comment
	retrievedComment, err := suite.repo.GetById(suite.ctx, savedComment.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), newContent, retrievedComment.Content)
}

func (suite *CommentRepositoryTestSuite) TestUpdate_CommentNotFound() {
	// Given
	commentID := uuid.New().String()
	newContent := "Updated content"

	// When
	_, err := suite.repo.Update(suite.ctx, commentID, newContent)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "comment not found")
}

func (suite *CommentRepositoryTestSuite) TestUpdate_ValidationError_EmptyCommentID() {
	// When
	_, err := suite.repo.Update(suite.ctx, "", "Updated content")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "comment ID is required")
}

func (suite *CommentRepositoryTestSuite) TestUpdate_ValidationError_EmptyContent() {
	// Given
	commentID := uuid.New().String()

	// When
	_, err := suite.repo.Update(suite.ctx, commentID, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "content is required")
}

func (suite *CommentRepositoryTestSuite) TestDelete_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content := "This is a test comment"

	// Save comment first
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, content)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.Delete(suite.ctx, authorID, savedComment.ID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.GetById(suite.ctx, savedComment.ID)
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
}

func (suite *CommentRepositoryTestSuite) TestDelete_CommentNotFound() {
	// Given
	userID := "user123"
	commentID := uuid.New().String()

	// When
	err := suite.repo.Delete(suite.ctx, userID, commentID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "comment not found")
}

func (suite *CommentRepositoryTestSuite) TestDelete_UnauthorizedUser() {
	// Given
	authorID := "author123"
	otherUserID := "other456"
	postID := uuid.New().String()
	content := "This is a test comment"

	// Save comment first
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, content)
	require.NoError(suite.T(), err)

	// When - try to delete with different user
	err = suite.repo.Delete(suite.ctx, otherUserID, savedComment.ID)

	// Then
	assert.Error(suite.T(), err)
	var forbiddenErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &forbiddenErr)
}

func (suite *CommentRepositoryTestSuite) TestDelete_ValidationError_EmptyUserID() {
	// When
	err := suite.repo.Delete(suite.ctx, "", "comment123")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *CommentRepositoryTestSuite) TestDelete_ValidationError_EmptyCommentID() {
	// When
	err := suite.repo.Delete(suite.ctx, "user123", "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "comment ID is required")
}

func (suite *CommentRepositoryTestSuite) TestExists_True() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	content := "This is a test comment"

	// Save comment first
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, content)
	require.NoError(suite.T(), err)

	// When
	exists, err := suite.repo.Exists(suite.ctx, savedComment.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *CommentRepositoryTestSuite) TestExists_False() {
	// Given
	commentID := uuid.New().String()

	// When
	exists, err := suite.repo.Exists(suite.ctx, commentID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *CommentRepositoryTestSuite) TestExists_ValidationError_EmptyCommentID() {
	// When
	_, err := suite.repo.Exists(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "comment ID is required")
}

func (suite *CommentRepositoryTestSuite) TestCountByPost_Success() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()

	// Save multiple comments
	for i := 0; i < 5; i++ {
		_, err := suite.repo.Save(suite.ctx, authorID, postID, "Comment "+string(rune('0'+i)))
		require.NoError(suite.T(), err)
	}

	// When
	count, err := suite.repo.CountByPost(suite.ctx, postID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 5, count)
}

func (suite *CommentRepositoryTestSuite) TestCountByPost_EmptyResult() {
	// Given
	postID := uuid.New().String()

	// When
	count, err := suite.repo.CountByPost(suite.ctx, postID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *CommentRepositoryTestSuite) TestCountByPost_ValidationError_EmptyPostID() {
	// When
	_, err := suite.repo.CountByPost(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *CommentRepositoryTestSuite) TestConcurrentOperations() {
	// Given
	authorID1 := "author1"
	authorID2 := "author2"
	postID := uuid.New().String()
	content1 := "First comment"
	content2 := "Second comment"

	// When - concurrent saves
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		_, err := suite.repo.Save(suite.ctx, authorID1, postID, content1)
		done1 <- err
	}()

	go func() {
		_, err := suite.repo.Save(suite.ctx, authorID2, postID, content2)
		done2 <- err
	}()

	// Then
	err1 := <-done1
	err2 := <-done2

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Verify both comments exist
	comments, err := suite.repo.GetByPost(suite.ctx, postID, 10)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), comments, 2)
}

func (suite *CommentRepositoryTestSuite) TestMultiplePostsScenario() {
	// Given
	authorID := "author123"
	postID1 := uuid.New().String()
	postID2 := uuid.New().String()

	// Save comments for different posts
	_, err := suite.repo.Save(suite.ctx, authorID, postID1, "Comment for post 1")
	require.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, authorID, postID1, "Another comment for post 1")
	require.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, authorID, postID2, "Comment for post 2")
	require.NoError(suite.T(), err)

	// When
	post1Comments, err := suite.repo.GetByPost(suite.ctx, postID1, 10)
	require.NoError(suite.T(), err)

	post2Comments, err := suite.repo.GetByPost(suite.ctx, postID2, 10)
	require.NoError(suite.T(), err)

	// Then
	assert.Len(suite.T(), post1Comments, 2)
	assert.Len(suite.T(), post2Comments, 1)

	// Verify comments belong to correct posts
	for _, comment := range post1Comments {
		assert.Equal(suite.T(), postID1, comment.PostID)
	}
	for _, comment := range post2Comments {
		assert.Equal(suite.T(), postID2, comment.PostID)
	}
}

func (suite *CommentRepositoryTestSuite) TestLargeContentHandling() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	largeContent := string(make([]byte, 5000)) // 5KB content

	// Fill with valid characters
	for i := range largeContent {
		largeContent = largeContent[:i] + "A" + largeContent[i+1:]
	}

	// When
	comment, err := suite.repo.Save(suite.ctx, authorID, postID, largeContent)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), largeContent, comment.Content)

	// Verify by retrieving
	retrievedComment, err := suite.repo.GetById(suite.ctx, comment.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), largeContent, retrievedComment.Content)
}

func (suite *CommentRepositoryTestSuite) TestUpdateAndDeleteSequence() {
	// Given
	authorID := "author123"
	postID := uuid.New().String()
	originalContent := "Original content"
	updatedContent := "Updated content"

	// Save comment
	savedComment, err := suite.repo.Save(suite.ctx, authorID, postID, originalContent)
	require.NoError(suite.T(), err)

	// Update comment
	updatedComment, err := suite.repo.Update(suite.ctx, savedComment.ID, updatedContent)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedContent, updatedComment.Content)

	// Delete comment
	err = suite.repo.Delete(suite.ctx, authorID, savedComment.ID)
	require.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.GetById(suite.ctx, savedComment.ID)
	assert.Error(suite.T(), err)
}

func TestCommentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentRepositoryTestSuite))
}
