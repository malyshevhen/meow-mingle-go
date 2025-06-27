package integration

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ReactionRepositoryTestSuite struct {
	suite.Suite
	testDB *TestDatabase
	repo   db.ReactionRepository
	ctx    context.Context
}

func (suite *ReactionRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	testDB, err := NewTestDatabase(suite.ctx)
	require.NoError(suite.T(), err, "Failed to create test database")

	suite.testDB = testDB
	suite.repo = db.NewReactionRepository(testDB.Session)
}

func (suite *ReactionRepositoryTestSuite) TearDownSuite() {
	if suite.testDB != nil {
		suite.testDB.Close(suite.ctx)
	}
}

func (suite *ReactionRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.testDB.Clean(suite.ctx)
	require.NoError(suite.T(), err, "Failed to clean test database")
}

func (suite *ReactionRepositoryTestSuite) TestSave_Success() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"
	content := "like"

	// When
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)

	// Then
	assert.NoError(suite.T(), err)

	// Verify reaction exists
	exists, err := suite.repo.Exists(suite.ctx, targetID, authorID)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *ReactionRepositoryTestSuite) TestSave_ValidationError_EmptyTargetID() {
	// Given
	targetID := ""
	authorID := "author123"
	content := "like"

	// When
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestSave_ValidationError_EmptyAuthorID() {
	// Given
	targetID := uuid.New().String()
	authorID := ""
	content := "like"

	// When
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestSave_ValidationError_EmptyContent() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"
	content := ""

	// When
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "reaction content is required")
}

func (suite *ReactionRepositoryTestSuite) TestSaveReaction_Success() {
	// Given
	reaction := &app.Reaction{
		TargetID: uuid.New().String(),
		AuthorID: "author123",
		Content:  "love",
	}

	// When
	err := suite.repo.SaveReaction(suite.ctx, reaction)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), reaction.CreatedAt.IsZero())
	assert.False(suite.T(), reaction.UpdatedAt.IsZero())
}

func (suite *ReactionRepositoryTestSuite) TestSaveReaction_ValidationError_NilReaction() {
	// When
	err := suite.repo.SaveReaction(suite.ctx, nil)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "reaction cannot be nil")
}

func (suite *ReactionRepositoryTestSuite) TestDelete_Success() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"
	content := "like"

	// Save reaction first
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.Delete(suite.ctx, targetID, authorID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify reaction no longer exists
	exists, err := suite.repo.Exists(suite.ctx, targetID, authorID)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *ReactionRepositoryTestSuite) TestDelete_NotFound() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"

	// When
	err := suite.repo.Delete(suite.ctx, targetID, authorID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "reaction not found")
}

func (suite *ReactionRepositoryTestSuite) TestDelete_ValidationError_EmptyTargetID() {
	// When
	err := suite.repo.Delete(suite.ctx, "", "author123")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestDelete_ValidationError_EmptyAuthorID() {
	// When
	err := suite.repo.Delete(suite.ctx, uuid.New().String(), "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestGetByTarget_Success() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"
	author1ID := "author1"
	author2ID := "author2"
	author3ID := "author3"

	// Save reactions
	err := suite.repo.Save(suite.ctx, targetID, author1ID, "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, author2ID, "love")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, author3ID, "like")
	require.NoError(suite.T(), err)

	// When
	reactions, err := suite.repo.GetByTarget(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactions, 3)

	// Verify all reactions are present
	authorIDs := make(map[string]bool)
	reactionTypes := make(map[string]int)
	for _, reaction := range reactions {
		authorIDs[reaction.AuthorID] = true
		reactionTypes[reaction.Content]++
		assert.Equal(suite.T(), targetID, reaction.TargetID)
		assert.False(suite.T(), reaction.CreatedAt.IsZero())
	}

	assert.True(suite.T(), authorIDs[author1ID])
	assert.True(suite.T(), authorIDs[author2ID])
	assert.True(suite.T(), authorIDs[author3ID])
	assert.Equal(suite.T(), 2, reactionTypes["like"])
	assert.Equal(suite.T(), 1, reactionTypes["love"])
}

func (suite *ReactionRepositoryTestSuite) TestGetByTarget_EmptyResult() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"

	// When
	reactions, err := suite.repo.GetByTarget(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), reactions)
}

func (suite *ReactionRepositoryTestSuite) TestGetByTarget_ValidationError_EmptyTargetID() {
	// When
	_, err := suite.repo.GetByTarget(suite.ctx, "", "post")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestGetByAuthor_Success() {
	// Given
	authorID := "author123"
	target1ID := uuid.New().String()
	target2ID := uuid.New().String()
	target3ID := uuid.New().String()

	// Save reactions from the same author
	err := suite.repo.Save(suite.ctx, target1ID, authorID, "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, target2ID, authorID, "love")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, target3ID, authorID, "laugh")
	require.NoError(suite.T(), err)

	// When
	reactions, err := suite.repo.GetByAuthor(suite.ctx, authorID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactions, 3)

	// Verify all reactions belong to the author
	targetIDs := make(map[string]bool)
	for _, reaction := range reactions {
		targetIDs[reaction.TargetID] = true
		assert.Equal(suite.T(), authorID, reaction.AuthorID)
		assert.False(suite.T(), reaction.CreatedAt.IsZero())
	}

	assert.True(suite.T(), targetIDs[target1ID])
	assert.True(suite.T(), targetIDs[target2ID])
	assert.True(suite.T(), targetIDs[target3ID])
}

func (suite *ReactionRepositoryTestSuite) TestGetByAuthor_WithLimit() {
	// Given
	authorID := "author123"

	// Save 10 reactions
	for i := 0; i < 10; i++ {
		targetID := uuid.New().String()
		err := suite.repo.Save(suite.ctx, targetID, authorID, "like")
		require.NoError(suite.T(), err)
	}

	// When
	reactions, err := suite.repo.GetByAuthor(suite.ctx, authorID, 5)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactions, 5)
}

func (suite *ReactionRepositoryTestSuite) TestGetByAuthor_EmptyResult() {
	// Given
	authorID := "nonexistent"

	// When
	reactions, err := suite.repo.GetByAuthor(suite.ctx, authorID, 10)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), reactions)
}

func (suite *ReactionRepositoryTestSuite) TestGetByAuthor_ValidationError_EmptyAuthorID() {
	// When
	_, err := suite.repo.GetByAuthor(suite.ctx, "", 10)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestExists_True() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"
	content := "like"

	// Save reaction first
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)
	require.NoError(suite.T(), err)

	// When
	exists, err := suite.repo.Exists(suite.ctx, targetID, authorID)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *ReactionRepositoryTestSuite) TestExists_False() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"

	// When
	exists, err := suite.repo.Exists(suite.ctx, targetID, authorID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *ReactionRepositoryTestSuite) TestExists_ValidationError_EmptyTargetID() {
	// When
	_, err := suite.repo.Exists(suite.ctx, "", "author123")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestExists_ValidationError_EmptyAuthorID() {
	// When
	_, err := suite.repo.Exists(suite.ctx, uuid.New().String(), "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "author ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestCountByTarget_Success() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"

	// Save reactions of different types
	err := suite.repo.Save(suite.ctx, targetID, "author1", "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author2", "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author3", "love")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author4", "laugh")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author5", "like")
	require.NoError(suite.T(), err)

	// When
	counts, err := suite.repo.CountByTarget(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), counts, 3)
	assert.Equal(suite.T(), 3, counts["like"])
	assert.Equal(suite.T(), 1, counts["love"])
	assert.Equal(suite.T(), 1, counts["laugh"])
}

func (suite *ReactionRepositoryTestSuite) TestCountByTarget_EmptyResult() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"

	// When
	counts, err := suite.repo.CountByTarget(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), counts)
}

func (suite *ReactionRepositoryTestSuite) TestCountByTarget_ValidationError_EmptyTargetID() {
	// When
	_, err := suite.repo.CountByTarget(suite.ctx, "", "post")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestGetReactionTypes_Success() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"

	// Save reactions of different types
	err := suite.repo.Save(suite.ctx, targetID, "author1", "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author2", "love")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author3", "laugh")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, targetID, "author4", "like") // Duplicate type
	require.NoError(suite.T(), err)

	// When
	reactionTypes, err := suite.repo.GetReactionTypes(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactionTypes, 3)

	// Convert to map for easier verification
	typeMap := make(map[string]bool)
	for _, reactionType := range reactionTypes {
		typeMap[reactionType] = true
	}

	assert.True(suite.T(), typeMap["like"])
	assert.True(suite.T(), typeMap["love"])
	assert.True(suite.T(), typeMap["laugh"])
}

func (suite *ReactionRepositoryTestSuite) TestGetReactionTypes_EmptyResult() {
	// Given
	targetID := uuid.New().String()
	targetType := "post"

	// When
	reactionTypes, err := suite.repo.GetReactionTypes(suite.ctx, targetID, targetType)

	// Then
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), reactionTypes)
}

func (suite *ReactionRepositoryTestSuite) TestGetReactionTypes_ValidationError_EmptyTargetID() {
	// When
	_, err := suite.repo.GetReactionTypes(suite.ctx, "", "post")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "target ID is required")
}

func (suite *ReactionRepositoryTestSuite) TestConcurrentOperations() {
	// Given
	targetID := uuid.New().String()
	author1ID := "author1"
	author2ID := "author2"

	// When - concurrent saves
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		err := suite.repo.Save(suite.ctx, targetID, author1ID, "like")
		done1 <- err
	}()

	go func() {
		err := suite.repo.Save(suite.ctx, targetID, author2ID, "love")
		done2 <- err
	}()

	// Then
	err1 := <-done1
	err2 := <-done2

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Verify both reactions exist
	reactions, err := suite.repo.GetByTarget(suite.ctx, targetID, "post")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactions, 2)
}

func (suite *ReactionRepositoryTestSuite) TestMultipleTargetsScenario() {
	// Given
	author1ID := "author1"
	target1ID := uuid.New().String()
	target2ID := uuid.New().String()

	// Save reactions for different targets
	err := suite.repo.Save(suite.ctx, target1ID, author1ID, "like")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, target1ID, "author2", "love")
	require.NoError(suite.T(), err)

	err = suite.repo.Save(suite.ctx, target2ID, author1ID, "laugh")
	require.NoError(suite.T(), err)

	// When
	target1Reactions, err := suite.repo.GetByTarget(suite.ctx, target1ID, "post")
	require.NoError(suite.T(), err)

	target2Reactions, err := suite.repo.GetByTarget(suite.ctx, target2ID, "post")
	require.NoError(suite.T(), err)

	author1Reactions, err := suite.repo.GetByAuthor(suite.ctx, author1ID, 10)
	require.NoError(suite.T(), err)

	// Then
	assert.Len(suite.T(), target1Reactions, 2)
	assert.Len(suite.T(), target2Reactions, 1)
	assert.Len(suite.T(), author1Reactions, 2)

	// Verify reactions belong to correct targets
	for _, reaction := range target1Reactions {
		assert.Equal(suite.T(), target1ID, reaction.TargetID)
	}
	for _, reaction := range target2Reactions {
		assert.Equal(suite.T(), target2ID, reaction.TargetID)
	}
	for _, reaction := range author1Reactions {
		assert.Equal(suite.T(), author1ID, reaction.AuthorID)
	}
}

func (suite *ReactionRepositoryTestSuite) TestReactionOverwrite() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"

	// Save initial reaction
	err := suite.repo.Save(suite.ctx, targetID, authorID, "like")
	require.NoError(suite.T(), err)

	// When - save different reaction from same author to same target
	err = suite.repo.Save(suite.ctx, targetID, authorID, "love")

	// Then - should succeed (overwrites previous reaction)
	assert.NoError(suite.T(), err)

	// Verify only one reaction exists
	reactions, err := suite.repo.GetByTarget(suite.ctx, targetID, "post")
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), reactions, 1)
	assert.Equal(suite.T(), "love", reactions[0].Content) // Should be the latest reaction
}

func (suite *ReactionRepositoryTestSuite) TestCreateAndDeleteSequence() {
	// Given
	targetID := uuid.New().String()
	authorID := "author123"
	content := "like"

	// Create reaction
	err := suite.repo.Save(suite.ctx, targetID, authorID, content)
	require.NoError(suite.T(), err)

	// Verify creation
	exists, err := suite.repo.Exists(suite.ctx, targetID, authorID)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), exists)

	// Delete reaction
	err = suite.repo.Delete(suite.ctx, targetID, authorID)
	require.NoError(suite.T(), err)

	// Verify deletion
	exists, err = suite.repo.Exists(suite.ctx, targetID, authorID)
	require.NoError(suite.T(), err)
	assert.False(suite.T(), exists)

	// Try to delete again (should fail)
	err = suite.repo.Delete(suite.ctx, targetID, authorID)
	assert.Error(suite.T(), err)
}

func (suite *ReactionRepositoryTestSuite) TestComplexReactionScenario() {
	// Given - create a complex scenario with multiple targets and authors
	post1ID := uuid.New().String()
	post2ID := uuid.New().String()
	comment1ID := uuid.New().String()

	users := []string{"user1", "user2", "user3", "user4", "user5"}
	reactionTypes := []string{"like", "love", "laugh", "angry", "sad"}

	// Create reactions for post1
	for i, user := range users {
		err := suite.repo.Save(suite.ctx, post1ID, user, reactionTypes[i%len(reactionTypes)])
		require.NoError(suite.T(), err)
	}

	// Create reactions for post2 (fewer reactions)
	for i := 0; i < 3; i++ {
		err := suite.repo.Save(suite.ctx, post2ID, users[i], "like")
		require.NoError(suite.T(), err)
	}

	// Create reactions for comment1
	err := suite.repo.Save(suite.ctx, comment1ID, users[0], "love")
	require.NoError(suite.T(), err)
	err = suite.repo.Save(suite.ctx, comment1ID, users[1], "love")
	require.NoError(suite.T(), err)

	// When - analyze the scenario
	post1Reactions, err := suite.repo.GetByTarget(suite.ctx, post1ID, "post")
	require.NoError(suite.T(), err)

	post2Reactions, err := suite.repo.GetByTarget(suite.ctx, post2ID, "post")
	require.NoError(suite.T(), err)

	comment1Reactions, err := suite.repo.GetByTarget(suite.ctx, comment1ID, "comment")
	require.NoError(suite.T(), err)

	post1Counts, err := suite.repo.CountByTarget(suite.ctx, post1ID, "post")
	require.NoError(suite.T(), err)

	user1Reactions, err := suite.repo.GetByAuthor(suite.ctx, users[0], 10)
	require.NoError(suite.T(), err)

	// Then - verify the complex scenario
	assert.Len(suite.T(), post1Reactions, 5)
	assert.Len(suite.T(), post2Reactions, 3)
	assert.Len(suite.T(), comment1Reactions, 2)
	assert.Len(suite.T(), user1Reactions, 3) // user1 reacted to all three targets

	// Verify post1 has diverse reaction types
	assert.True(suite.T(), len(post1Counts) > 1)

	// Verify post2 has uniform reactions
	post2Counts, err := suite.repo.CountByTarget(suite.ctx, post2ID, "post")
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 3, post2Counts["like"])

	// Verify comment1 has specific reaction type
	comment1Types, err := suite.repo.GetReactionTypes(suite.ctx, comment1ID, "comment")
	require.NoError(suite.T(), err)
	assert.Len(suite.T(), comment1Types, 1)
	assert.Equal(suite.T(), "love", comment1Types[0])
}

func TestReactionRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ReactionRepositoryTestSuite))
}
