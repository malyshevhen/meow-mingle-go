package integration

import (
	"context"
	"testing"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleProfileRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	// Create test database
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err, "Failed to create test database")
	defer testDB.Close(ctx)

	// Create repository
	repo := db.NewProfileRepository(testDB.Session)

	t.Run("Save_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "user123"
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"

		// When
		profile, err := repo.Save(ctx, userID, email, firstName, lastName)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, email, profile.Email)
		assert.Equal(t, firstName, profile.FirstName)
		assert.Equal(t, lastName, profile.LastName)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())
	})

	t.Run("Save_ValidationError_EmptyUserID", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := ""
		email := "test@example.com"
		firstName := "John"
		lastName := "Doe"

		// When
		_, err = repo.Save(ctx, userID, email, firstName, lastName)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("GetById_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "user456"
		email := "john@example.com"
		firstName := "John"
		lastName := "Smith"

		// Save profile first
		_, err = repo.Save(ctx, userID, email, firstName, lastName)
		require.NoError(t, err)

		// When
		profile, err := repo.GetById(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, email, profile.Email)
		assert.Equal(t, firstName, profile.FirstName)
		assert.Equal(t, lastName, profile.LastName)
	})

	t.Run("GetById_NotFound", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "nonexistent"

		// When
		_, err = repo.GetById(ctx, userID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Update_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "user789"
		email := "original@example.com"
		firstName := "Original"
		lastName := "Name"

		// Save profile first
		originalProfile, err := repo.Save(ctx, userID, email, firstName, lastName)
		require.NoError(t, err)

		// Modify profile
		originalProfile.FirstName = "Updated"
		originalProfile.Email = "updated@example.com"

		// When
		err = repo.Update(ctx, &originalProfile)

		// Then
		assert.NoError(t, err)

		// Verify update
		updatedProfile, err := repo.GetById(ctx, userID)
		require.NoError(t, err)
		assert.Equal(t, "Updated", updatedProfile.FirstName)
		assert.Equal(t, "updated@example.com", updatedProfile.Email)
		assert.True(t, updatedProfile.UpdatedAt.After(originalProfile.CreatedAt))
	})

	t.Run("Delete_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "user999"
		email := "delete@example.com"
		firstName := "Delete"
		lastName := "Me"

		// Save profile first
		_, err = repo.Save(ctx, userID, email, firstName, lastName)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, userID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(ctx, userID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Exists_True", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "exists123"
		email := "exists@example.com"
		firstName := "Exists"
		lastName := "User"

		// Save profile first
		_, err = repo.Save(ctx, userID, email, firstName, lastName)
		require.NoError(t, err)

		// When
		exists, err := repo.Exists(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists_False", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "doesnotexist"

		// When
		exists, err := repo.Exists(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("SaveProfile_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		profile := &app.Profile{
			UserID:    "profile123",
			Email:     "profile@example.com",
			FirstName: "Profile",
			LastName:  "Test",
		}

		// When
		err = repo.SaveProfile(ctx, profile)

		// Then
		assert.NoError(t, err)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())

		// Verify by getting the profile
		savedProfile, err := repo.GetById(ctx, profile.UserID)
		require.NoError(t, err)
		assert.Equal(t, profile.UserID, savedProfile.UserID)
		assert.Equal(t, profile.Email, savedProfile.Email)
	})

	t.Run("GetByEmail_Success", func(t *testing.T) {
		// Clean database before test
		err := testDB.Clean(ctx)
		require.NoError(t, err)

		// Given
		userID := "email123"
		email := "unique@example.com"
		firstName := "Email"
		lastName := "Test"

		// Save profile first
		_, err = repo.Save(ctx, userID, email, firstName, lastName)
		require.NoError(t, err)

		// When
		profile, err := repo.GetByEmail(ctx, email)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, userID, profile.UserID)
		assert.Equal(t, email, profile.Email)
		assert.Equal(t, firstName, profile.FirstName)
		assert.Equal(t, lastName, profile.LastName)
	})
}
