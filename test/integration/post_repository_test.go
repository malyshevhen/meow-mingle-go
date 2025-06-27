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

type PostRepositoryTestSuite struct {
	suite.Suite
	testDB *TestDatabase
	repo   db.PostRepository
	ctx    context.Context
}

func (suite *PostRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	testDB, err := NewTestDatabase(suite.ctx)
	require.NoError(suite.T(), err, "Failed to create test database")

	suite.testDB = testDB
	suite.repo = db.NewPostRepository(testDB.Session)
}

func (suite *PostRepositoryTestSuite) TearDownSuite() {
	if suite.testDB != nil {
		suite.testDB.Close(suite.ctx)
	}
}

func (suite *PostRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.testDB.Clean(suite.ctx)
	require.NoError(suite.T(), err, "Failed to clean test database")
}

func (suite *PostRepositoryTestSuite) TestSave_Success() {
	// Given
	authorID := "author123"
	content := "This is a test post"

	// When
	post, err := suite.repo.Save(suite.ctx, authorID, content)

	// Then
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), post.ID)
	assert.Equal(suite.T(), authorID, post.AuthorID)
	assert.Equal(suite.T(), content, post.Content)
	assert.False(suite.T(), post.CreatedAt.IsZero())
	assert.False(suite.T(), post.UpdatedAt.IsZero())
}

func (suite *PostRepositoryTestSuite) TestSave_ValidationError_EmptyAuthorID() {
	// Given
	authorID := ""
	content := "This is a test post"

	// When
	_, err := suite.repo.Save(suite.ctx, authorID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *PostRepositoryTestSuite) TestSave_ValidationError_EmptyContent() {
	// Given
	authorID := "author123"
	content := ""

	// When
	_, err := suite.repo.Save(suite.ctx, authorID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "content is required")
}

func (suite *PostRepositoryTestSuite) TestSavePost_Success() {
	// Given
	post := &app.Post{
		ID:       uuid.New().String(),
		AuthorID: "author123",
		Content:  "This is a test post",
	}

	// When
	err := suite.repo.SavePost(suite.ctx, post)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), post.CreatedAt.IsZero())
	assert.False(suite.T(), post.UpdatedAt.IsZero())
}

func (suite *PostRepositoryTestSuite) TestSavePost_ValidationError_NilPost() {
	// When
	err := suite.repo.SavePost(suite.ctx, nil)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post cannot be nil")
}

func (suite *PostRepositoryTestSuite) TestGet_Success() {
	// Given
	authorID := "author123"
	content := "This is a test post"

	// Save post first
	savedPost, err := suite.repo.Save(suite.ctx, authorID, content)
	require.NoError(suite.T(), err)

	// When
	post, err := suite.repo.Get(suite.ctx, savedPost.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), savedPost.ID, post.ID)
	assert.Equal(suite.T(), authorID, post.AuthorID)
	assert.Equal(suite.T(), content, post.Content)
	assert.Equal(suite.T(), savedPost.CreatedAt.Unix(), post.CreatedAt.Unix())
}

func (suite *PostRepositoryTestSuite) TestGet_NotFound() {
	// Given
	postID := uuid.New().String()

	// When
	_, err := suite.repo.Get(suite.ctx, postID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "post not found")
}

func (suite *PostRepositoryTestSuite) TestGet_ValidationError_EmptyPostID() {
	// When
	_, err := suite.repo.Get(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *PostRepositoryTestSuite) TestList_Success() {
	// Given
	authorID := "author123"
	content1 := "First post"
	content2 := "Second post"
	content3 := "Third post"

	// Save posts
	_, err := suite.repo.Save(suite.ctx, authorID, content1)
	require.NoError(suite.T(), err)
	time.Sleep(10 * time.Millisecond) // Ensure different timestamps

	_, err = suite.repo.Save(suite.ctx, authorID, content2)
	require.NoError(suite.T(), err)
	time.Sleep(10 * time.Millisecond)

	_, err = suite.repo.Save(suite.ctx, authorID, content3)
	require.NoError(suite.T(), err)

	// When
	posts, err := suite.repo.List(suite.ctx, authorID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 3)

	// Posts should be ordered by created_at DESC
	assert.Equal(suite.T(), content3, posts[0].Content)
	assert.Equal(suite.T(), content2, posts[1].Content)
	assert.Equal(suite.T(), content1, posts[2].Content)
}

func (suite *PostRepositoryTestSuite) TestList_EmptyResult() {
	// Given
	authorID := "nonexistent"

	// When
	posts, err := suite.repo.List(suite.ctx, authorID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), posts)
}

func (suite *PostRepositoryTestSuite) TestGetByAuthor_WithLimit() {
	// Given
	authorID := "author123"

	// Save 10 posts
	for i := 0; i < 10; i++ {
		_, err := suite.repo.Save(suite.ctx, authorID, "Post "+string(rune('0'+i)))
		require.NoError(suite.T(), err)
		time.Sleep(10 * time.Millisecond)
	}

	// When
	posts, err := suite.repo.GetByAuthor(suite.ctx, authorID, 5)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 5)
}

func (suite *PostRepositoryTestSuite) TestGetByAuthor_ValidationError_EmptyAuthorID() {
	// When
	_, err := suite.repo.GetByAuthor(suite.ctx, "", 10)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *PostRepositoryTestSuite) TestUpdate_Success() {
	// Given
	authorID := "author123"
	originalContent := "Original content"
	newContent := "Updated content"

	// Save post first
	savedPost, err := suite.repo.Save(suite.ctx, authorID, originalContent)
	require.NoError(suite.T(), err)

	// When
	updatedPost, err := suite.repo.Update(suite.ctx, savedPost.ID, newContent)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), savedPost.ID, updatedPost.ID)
	assert.Equal(suite.T(), newContent, updatedPost.Content)
	assert.True(suite.T(), updatedPost.UpdatedAt.After(savedPost.CreatedAt))

	// Verify by getting the post
	retrievedPost, err := suite.repo.Get(suite.ctx, savedPost.ID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), newContent, retrievedPost.Content)
}

func (suite *PostRepositoryTestSuite) TestUpdate_PostNotFound() {
	// Given
	postID := uuid.New().String()
	newContent := "Updated content"

	// When
	_, err := suite.repo.Update(suite.ctx, postID, newContent)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "post not found")
}

func (suite *PostRepositoryTestSuite) TestUpdate_ValidationError_EmptyPostID() {
	// When
	_, err := suite.repo.Update(suite.ctx, "", "Updated content")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *PostRepositoryTestSuite) TestUpdate_ValidationError_EmptyContent() {
	// Given
	postID := uuid.New().String()

	// When
	_, err := suite.repo.Update(suite.ctx, postID, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "content is required")
}

func (suite *PostRepositoryTestSuite) TestDelete_Success() {
	// Given
	authorID := "author123"
	content := "This is a test post"

	// Save post first
	savedPost, err := suite.repo.Save(suite.ctx, authorID, content)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.Delete(suite.ctx, savedPost.ID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.Get(suite.ctx, savedPost.ID)
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
}

func (suite *PostRepositoryTestSuite) TestDelete_PostNotFound() {
	// Given
	postID := uuid.New().String()

	// When
	err := suite.repo.Delete(suite.ctx, postID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "post not found")
}

func (suite *PostRepositoryTestSuite) TestDelete_ValidationError_EmptyPostID() {
	// When
	err := suite.repo.Delete(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *PostRepositoryTestSuite) TestExists_True() {
	// Given
	authorID := "author123"
	content := "This is a test post"

	// Save post first
	savedPost, err := suite.repo.Save(suite.ctx, authorID, content)
	require.NoError(suite.T(), err)

	// When
	exists, err := suite.repo.Exists(suite.ctx, savedPost.ID)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *PostRepositoryTestSuite) TestExists_False() {
	// Given
	postID := uuid.New().String()

	// When
	exists, err := suite.repo.Exists(suite.ctx, postID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *PostRepositoryTestSuite) TestExists_ValidationError_EmptyPostID() {
	// When
	_, err := suite.repo.Exists(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "post ID is required")
}

func (suite *PostRepositoryTestSuite) TestFeed_EmptyResult() {
	// Given
	userID := "user123"

	// When
	posts, err := suite.repo.Feed(suite.ctx, userID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), posts)
}

func (suite *PostRepositoryTestSuite) TestFeed_ValidationError_EmptyUserID() {
	// When
	_, err := suite.repo.Feed(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *PostRepositoryTestSuite) TestConcurrentOperations() {
	// Given
	authorID := "author123"
	content1 := "First post"
	content2 := "Second post"

	// When - concurrent saves
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		_, err := suite.repo.Save(suite.ctx, authorID, content1)
		done1 <- err
	}()

	go func() {
		_, err := suite.repo.Save(suite.ctx, authorID, content2)
		done2 <- err
	}()

	// Then
	err1 := <-done1
	err2 := <-done2

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Verify both posts exist
	posts, err := suite.repo.List(suite.ctx, authorID)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), posts, 2)
}

func (suite *PostRepositoryTestSuite) TestSavePost_DuplicateID() {
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
	err := suite.repo.SavePost(suite.ctx, post1)
	require.NoError(suite.T(), err)

	// When - try to save with same ID
	err = suite.repo.SavePost(suite.ctx, post2)

	// Then - should fail due to primary key constraint
	assert.Error(suite.T(), err)
}

func (suite *PostRepositoryTestSuite) TestLargeContentHandling() {
	// Given
	authorID := "author123"
	largeContent := string(make([]byte, 10000)) // 10KB content

	// Fill with valid characters
	for i := range largeContent {
		largeContent = largeContent[:i] + "A" + largeContent[i+1:]
	}

	// When
	post, err := suite.repo.Save(suite.ctx, authorID, largeContent)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), largeContent, post.Content)

	// Verify by retrieving
	retrievedPost, err := suite.repo.Get(suite.ctx, post.ID)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), largeContent, retrievedPost.Content)
}

func (suite *PostRepositoryTestSuite) TestUpdateAndDeleteSequence() {
	// Given
	authorID := "author123"
	originalContent := "Original content"
	updatedContent := "Updated content"

	// Save post
	savedPost, err := suite.repo.Save(suite.ctx, authorID, originalContent)
	require.NoError(suite.T(), err)

	// Update post
	updatedPost, err := suite.repo.Update(suite.ctx, savedPost.ID, updatedContent)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedContent, updatedPost.Content)

	// Delete post
	err = suite.repo.Delete(suite.ctx, savedPost.ID)
	require.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.Get(suite.ctx, savedPost.ID)
	assert.Error(suite.T(), err)
}

func (suite *PostRepositoryTestSuite) TestMultipleAuthorsScenario() {
	// Given
	author1 := "author1"
	author2 := "author2"

	// Save posts from different authors
	_, err := suite.repo.Save(suite.ctx, author1, "Author 1 Post 1")
	require.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, author2, "Author 2 Post 1")
	require.NoError(suite.T(), err)

	_, err = suite.repo.Save(suite.ctx, author1, "Author 1 Post 2")
	require.NoError(suite.T(), err)

	// When
	author1Posts, err := suite.repo.List(suite.ctx, author1)
	require.NoError(suite.T(), err)

	author2Posts, err := suite.repo.List(suite.ctx, author2)
	require.NoError(suite.T(), err)

	// Then
	assert.Len(suite.T(), author1Posts, 2)
	assert.Len(suite.T(), author2Posts, 1)

	// Verify posts belong to correct authors
	for _, post := range author1Posts {
		assert.Equal(suite.T(), author1, post.AuthorID)
	}
	for _, post := range author2Posts {
		assert.Equal(suite.T(), author2, post.AuthorID)
	}
}

func TestPostRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostRepositoryTestSuite))
}
