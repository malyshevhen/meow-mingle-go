package integration

import (
	"context"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProfileRepositoryTestSuite struct {
	suite.Suite
	testDB *TestDatabase
	repo   db.ProfileRepository
	ctx    context.Context
}

func (suite *ProfileRepositoryTestSuite) SetupSuite() {
	suite.ctx = context.Background()

	testDB, err := NewTestDatabase(suite.ctx)
	require.NoError(suite.T(), err, "Failed to create test database")

	suite.testDB = testDB
	suite.repo = db.NewProfileRepository(testDB.Session)
}

func (suite *ProfileRepositoryTestSuite) TearDownSuite() {
	if suite.testDB != nil {
		suite.testDB.Close(suite.ctx)
	}
}

func (suite *ProfileRepositoryTestSuite) SetupTest() {
	// Clean database before each test
	err := suite.testDB.Clean(suite.ctx)
	require.NoError(suite.T(), err, "Failed to clean test database")
}

func (suite *ProfileRepositoryTestSuite) TestSave_Success() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// When
	profile, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, profile.UserID)
	assert.Equal(suite.T(), email, profile.Email)
	assert.Equal(suite.T(), firstName, profile.FirstName)
	assert.Equal(suite.T(), lastName, profile.LastName)
	assert.False(suite.T(), profile.CreatedAt.IsZero())
	assert.False(suite.T(), profile.UpdatedAt.IsZero())
}

func (suite *ProfileRepositoryTestSuite) TestSave_ValidationError_EmptyUserID() {
	// Given
	userID := ""
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// When
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
}

func (suite *ProfileRepositoryTestSuite) TestSave_ValidationError_EmptyEmail() {
	// Given
	userID := "user123"
	email := ""
	firstName := "John"
	lastName := "Doe"

	// When
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
}

func (suite *ProfileRepositoryTestSuite) TestSaveProfile_Success() {
	// Given
	profile := &app.Profile{
		UserID:    "user123",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// When
	err := suite.repo.SaveProfile(suite.ctx, profile)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), profile.CreatedAt.IsZero())
	assert.False(suite.T(), profile.UpdatedAt.IsZero())
}

func (suite *ProfileRepositoryTestSuite) TestSaveProfile_ValidationError_NilProfile() {
	// When
	err := suite.repo.SaveProfile(suite.ctx, nil)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "profile cannot be nil")
}

func (suite *ProfileRepositoryTestSuite) TestGetById_Success() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Save profile first
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)
	require.NoError(suite.T(), err)

	// When
	profile, err := suite.repo.GetById(suite.ctx, userID)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, profile.UserID)
	assert.Equal(suite.T(), email, profile.Email)
	assert.Equal(suite.T(), firstName, profile.FirstName)
	assert.Equal(suite.T(), lastName, profile.LastName)
}

func (suite *ProfileRepositoryTestSuite) TestGetById_NotFound() {
	// Given
	userID := "nonexistent"

	// When
	_, err := suite.repo.GetById(suite.ctx, userID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "profile not found")
}

func (suite *ProfileRepositoryTestSuite) TestGetById_ValidationError_EmptyUserID() {
	// When
	_, err := suite.repo.GetById(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *ProfileRepositoryTestSuite) TestGetByEmail_Success() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Save profile first
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)
	require.NoError(suite.T(), err)

	// When
	profile, err := suite.repo.GetByEmail(suite.ctx, email)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID, profile.UserID)
	assert.Equal(suite.T(), email, profile.Email)
	assert.Equal(suite.T(), firstName, profile.FirstName)
	assert.Equal(suite.T(), lastName, profile.LastName)
}

func (suite *ProfileRepositoryTestSuite) TestGetByEmail_NotFound() {
	// Given
	email := "nonexistent@example.com"

	// When
	_, err := suite.repo.GetByEmail(suite.ctx, email)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "profile not found")
}

func (suite *ProfileRepositoryTestSuite) TestGetByEmail_ValidationError_EmptyEmail() {
	// When
	_, err := suite.repo.GetByEmail(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "email is required")
}

func (suite *ProfileRepositoryTestSuite) TestUpdate_Success() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Save profile first
	originalProfile, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)
	require.NoError(suite.T(), err)

	// Modify profile
	originalProfile.FirstName = "Jane"
	originalProfile.Email = "jane@example.com"

	// When
	err = suite.repo.Update(suite.ctx, &originalProfile)

	// Then
	assert.NoError(suite.T(), err)

	// Verify update
	updatedProfile, err := suite.repo.GetById(suite.ctx, userID)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Jane", updatedProfile.FirstName)
	assert.Equal(suite.T(), "jane@example.com", updatedProfile.Email)
	assert.True(suite.T(), updatedProfile.UpdatedAt.After(originalProfile.CreatedAt))
}

func (suite *ProfileRepositoryTestSuite) TestUpdate_ProfileNotFound() {
	// Given
	profile := &app.Profile{
		UserID:    "nonexistent",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// When
	err := suite.repo.Update(suite.ctx, profile)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "profile not found")
}

func (suite *ProfileRepositoryTestSuite) TestUpdate_ValidationError_NilProfile() {
	// When
	err := suite.repo.Update(suite.ctx, nil)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "profile cannot be nil")
}

func (suite *ProfileRepositoryTestSuite) TestUpdate_ValidationError_EmptyUserID() {
	// Given
	profile := &app.Profile{
		UserID:    "",
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
	}

	// When
	err := suite.repo.Update(suite.ctx, profile)

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *ProfileRepositoryTestSuite) TestDelete_Success() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Save profile first
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)
	require.NoError(suite.T(), err)

	// When
	err = suite.repo.Delete(suite.ctx, userID)

	// Then
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.GetById(suite.ctx, userID)
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
}

func (suite *ProfileRepositoryTestSuite) TestDelete_ProfileNotFound() {
	// Given
	userID := "nonexistent"

	// When
	err := suite.repo.Delete(suite.ctx, userID)

	// Then
	assert.Error(suite.T(), err)
	var notFoundErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &notFoundErr)
	assert.Contains(suite.T(), err.Error(), "profile not found")
}

func (suite *ProfileRepositoryTestSuite) TestDelete_ValidationError_EmptyUserID() {
	// When
	err := suite.repo.Delete(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *ProfileRepositoryTestSuite) TestExists_True() {
	// Given
	userID := "user123"
	email := "test@example.com"
	firstName := "John"
	lastName := "Doe"

	// Save profile first
	_, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)
	require.NoError(suite.T(), err)

	// When
	exists, err := suite.repo.Exists(suite.ctx, userID)

	// Then
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), exists)
}

func (suite *ProfileRepositoryTestSuite) TestExists_False() {
	// Given
	userID := "nonexistent"

	// When
	exists, err := suite.repo.Exists(suite.ctx, userID)

	// Then
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), exists)
}

func (suite *ProfileRepositoryTestSuite) TestExists_ValidationError_EmptyUserID() {
	// When
	_, err := suite.repo.Exists(suite.ctx, "")

	// Then
	assert.Error(suite.T(), err)
	var validationErr *errors.BasicError
	assert.ErrorAs(suite.T(), err, &validationErr)
	assert.Contains(suite.T(), err.Error(), "user ID is required")
}

func (suite *ProfileRepositoryTestSuite) TestConcurrentOperations() {
	// Given
	userID1 := "user1"
	userID2 := "user2"
	email1 := "user1@example.com"
	email2 := "user2@example.com"

	// When - concurrent saves
	done1 := make(chan error, 1)
	done2 := make(chan error, 1)

	go func() {
		_, err := suite.repo.Save(suite.ctx, userID1, email1, "John", "Doe")
		done1 <- err
	}()

	go func() {
		_, err := suite.repo.Save(suite.ctx, userID2, email2, "Jane", "Smith")
		done2 <- err
	}()

	// Then
	err1 := <-done1
	err2 := <-done2

	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)

	// Verify both profiles exist
	profile1, err := suite.repo.GetById(suite.ctx, userID1)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID1, profile1.UserID)

	profile2, err := suite.repo.GetById(suite.ctx, userID2)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), userID2, profile2.UserID)
}

func (suite *ProfileRepositoryTestSuite) TestSaveProfile_DuplicateEmail() {
	// Given
	userID1 := "user1"
	userID2 := "user2"
	email := "duplicate@example.com"

	// Save first profile
	_, err := suite.repo.Save(suite.ctx, userID1, email, "John", "Doe")
	require.NoError(suite.T(), err)

	// When - try to save with same email
	_, err = suite.repo.Save(suite.ctx, userID2, email, "Jane", "Smith")

	// Then - should succeed as we don't enforce unique email constraint at DB level
	// In a real application, you might want to handle this at the service layer
	assert.NoError(suite.T(), err)
}

func (suite *ProfileRepositoryTestSuite) TestSaveProfile_DuplicateUserID() {
	// Given
	userID := "user123"
	email1 := "test1@example.com"
	email2 := "test2@example.com"

	// Save first profile
	_, err := suite.repo.Save(suite.ctx, userID, email1, "John", "Doe")
	require.NoError(suite.T(), err)

	// When - try to save with same userID
	_, err = suite.repo.Save(suite.ctx, userID, email2, "Jane", "Smith")

	// Then - should fail due to primary key constraint
	assert.Error(suite.T(), err)
}

func (suite *ProfileRepositoryTestSuite) TestLargeDataHandling() {
	// Given
	userID := "user123"
	email := "test@example.com"
	// Test with large names
	firstName := string(make([]byte, 1000)) // Large first name
	lastName := string(make([]byte, 1000))  // Large last name

	// Fill with valid characters
	for i := range firstName {
		firstName = firstName[:i] + "A" + firstName[i+1:]
		lastName = lastName[:i] + "B" + lastName[i+1:]
	}

	// When
	profile, err := suite.repo.Save(suite.ctx, userID, email, firstName, lastName)

	// Then
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), firstName, profile.FirstName)
	assert.Equal(suite.T(), lastName, profile.LastName)
}

func TestProfileRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ProfileRepositoryTestSuite))
}
