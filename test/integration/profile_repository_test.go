package integration

import (
	"context"
	"testing"
	"time"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProfileRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()
	testDB, err := NewSimpleTestDatabase(ctx)
	require.NoError(t, err)
	defer testDB.Close(ctx)

	repo := db.NewProfileRepository(testDB.Session)
	dataBuilder := NewTestDataBuilder()
	runner := NewTestRunner() // For concurrent operations

	t.Run("Save Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("save-success")

		// When
		profile, err := repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, testData.UserID, profile.UserID)
		assert.Equal(t, testData.Email, profile.Email)
		assert.Equal(t, testData.FirstName, profile.FirstName)
		assert.Equal(t, testData.LastName, profile.LastName)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())
	})

	t.Run("Save Validation Error - Empty UserID", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("validation-empty-user-id")
		testData.UserID = ""

		// When
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Save Validation Error - Empty Email", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("validation-empty-email")
		testData.Email = ""

		// When
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("SaveProfile Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("save-profile-success")
		profile := &app.Profile{
			UserID:    testData.UserID,
			Email:     testData.Email,
			FirstName: testData.FirstName,
			LastName:  testData.LastName,
		}

		// When
		err = repo.SaveProfile(ctx, profile)

		// Then
		assert.NoError(t, err)
		assert.False(t, profile.CreatedAt.IsZero())
		assert.False(t, profile.UpdatedAt.IsZero())
	})

	t.Run("SaveProfile Validation Error - Nil Profile", func(t *testing.T) {
		// When
		err = repo.SaveProfile(ctx, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile cannot be nil")
	})

	t.Run("GetById Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("get-by-id-success")

		// Save profile first
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// When
		profile, err := repo.GetById(ctx, testData.UserID)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, testData.UserID, profile.UserID)
		assert.Equal(t, testData.Email, profile.Email)
		assert.Equal(t, testData.FirstName, profile.FirstName)
		assert.Equal(t, testData.LastName, profile.LastName)
	})

	t.Run("GetById Not Found", func(t *testing.T) {
		// Given
		userID := dataBuilder.UserID("nonexistent-get-by-id")

		// When
		_, err = repo.GetById(ctx, userID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("GetById Validation Error - Empty UserID", func(t *testing.T) {
		// When
		_, err = repo.GetById(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("GetByEmail Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("get-by-email-success")

		// Save profile first
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// When
		profile, err := repo.GetByEmail(ctx, testData.Email)

		// Then
		assert.NoError(t, err)
		assert.Equal(t, testData.UserID, profile.UserID)
		assert.Equal(t, testData.Email, profile.Email)
		assert.Equal(t, testData.FirstName, profile.FirstName)
		assert.Equal(t, testData.LastName, profile.LastName)
	})

	t.Run("GetByEmail Not Found", func(t *testing.T) {
		// Given
		email := dataBuilder.Email("nonexistent-get-by-email")

		// When
		_, err = repo.GetByEmail(ctx, email)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("GetByEmail Validation Error - Empty Email", func(t *testing.T) {
		// When
		_, err = repo.GetByEmail(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email is required")
	})

	t.Run("Update Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("update-success")

		// Save profile first
		originalProfile, err := repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// Modify profile
		originalProfile.FirstName = "UpdatedFirstName"
		originalProfile.Email = dataBuilder.Email("updated-email")
		time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt is different

		// When
		err = repo.Update(ctx, &originalProfile)

		// Then
		assert.NoError(t, err)

		// Verify update
		updatedProfile, err := repo.GetById(ctx, testData.UserID)
		require.NoError(t, err)
		assert.Equal(t, "UpdatedFirstName", updatedProfile.FirstName)
		assert.Equal(t, dataBuilder.Email("updated-email"), updatedProfile.Email)
		assert.True(t, updatedProfile.UpdatedAt.After(originalProfile.CreatedAt))
	})

	t.Run("Update Profile Not Found", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("update-profile-not-found")
		profile := &app.Profile{
			UserID:    testData.UserID,
			Email:     testData.Email,
			FirstName: testData.FirstName,
			LastName:  testData.LastName,
		}

		// When
		err = repo.Update(ctx, profile)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Update Validation Error - Nil Profile", func(t *testing.T) {
		// When
		err = repo.Update(ctx, nil)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile cannot be nil")
	})

	t.Run("Update Validation Error - Empty UserID", func(t *testing.T) {
		// Given
		profile := &app.Profile{
			UserID:    "",
			Email:     dataBuilder.Email("test-update-empty-user-id"),
			FirstName: "Test",
			LastName:  "User",
		}

		// When
		err = repo.Update(ctx, profile)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Delete Success", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("delete-success")

		// Save profile first
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// When
		err = repo.Delete(ctx, testData.UserID)

		// Then
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(ctx, testData.UserID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Delete Profile Not Found", func(t *testing.T) {
		// Given
		userID := dataBuilder.UserID("nonexistent-delete")

		// When
		err = repo.Delete(ctx, userID)

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Delete Validation Error - Empty UserID", func(t *testing.T) {
		// When
		err = repo.Delete(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Exists True", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("exists-true")

		// Save profile first
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// When
		exists, err := repo.Exists(ctx, testData.UserID)

		// Then
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Exists False", func(t *testing.T) {
		// Given
		userID := dataBuilder.UserID("exists-false")

		// When
		exists, err := repo.Exists(ctx, userID)

		// Then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Exists Validation Error - Empty UserID", func(t *testing.T) {
		// When
		_, err = repo.Exists(ctx, "")

		// Then
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user ID is required")
	})

	t.Run("Concurrent Operations", func(t *testing.T) {
		// Given
		testData1 := dataBuilder.CreateTestProfile("concurrent1")
		testData2 := dataBuilder.CreateTestProfile("concurrent2")

		// When - concurrent saves using optimized runner
		runner.RunConcurrently(t,
			func() error {
				_, err := repo.Save(ctx, testData1.UserID, testData1.Email, testData1.FirstName, testData1.LastName)
				return err
			},
			func() error {
				_, err := repo.Save(ctx, testData2.UserID, testData2.Email, testData2.FirstName, testData2.LastName)
				return err
			},
		)

		// Then - verify both profiles exist
		profile1, err := repo.GetById(ctx, testData1.UserID)
		assert.NoError(t, err)
		assert.Equal(t, testData1.UserID, profile1.UserID)

		profile2, err := repo.GetById(ctx, testData2.UserID)
		assert.NoError(t, err)
		assert.Equal(t, testData2.UserID, profile2.UserID)
	})

	t.Run("SaveProfile Duplicate UserID", func(t *testing.T) {
		t.Skip()
		// Given
		testData := dataBuilder.CreateTestProfile("duplicate-user-id")

		// Save first profile
		_, err = repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// When - try to save with same userID but different email
		_, err = repo.Save(ctx, testData.UserID, dataBuilder.Email("different-duplicate-email"), "Different", "User")

		// Then - should fail due to primary key constraint
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate key value violates unique constraint")
	})

	t.Run("Large Data Handling", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("large-data")

		// Create large field data
		largeFirstName := make([]byte, 1000)
		largeLastName := make([]byte, 1000)
		for i := range largeFirstName {
			largeFirstName[i] = 'A'
			largeLastName[i] = 'B'
		}

		// When
		profile, err := repo.Save(ctx, testData.UserID, testData.Email, string(largeFirstName), string(largeLastName))

		// Then
		assert.NoError(t, err)
		assert.Equal(t, string(largeFirstName), profile.FirstName)
		assert.Equal(t, string(largeLastName), profile.LastName)

		// Verify by retrieving
		retrievedProfile, err := repo.GetById(ctx, profile.UserID)
		assert.NoError(t, err)
		assert.Equal(t, string(largeFirstName), retrievedProfile.FirstName)
		assert.Equal(t, string(largeLastName), retrievedProfile.LastName)
	})

	t.Run("Create Update Delete Sequence", func(t *testing.T) {
		// Given
		testData := dataBuilder.CreateTestProfile("crud-sequence")

		// Create profile
		savedProfile, err := repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
		require.NoError(t, err)

		// Update profile
		savedProfile.FirstName = "UpdatedName"
		err = repo.Update(ctx, &savedProfile)
		require.NoError(t, err)

		// Verify update
		updatedProfile, err := repo.GetById(ctx, testData.UserID)
		require.NoError(t, err)
		assert.Equal(t, "UpdatedName", updatedProfile.FirstName)

		// Delete profile
		err = repo.Delete(ctx, testData.UserID)
		require.NoError(t, err)

		// Verify deletion
		_, err = repo.GetById(ctx, testData.UserID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "profile not found")
	})

	t.Run("Bulk Operations", func(t *testing.T) {
		// Given
		numProfiles := 5
		profiles := make([]TestProfile, numProfiles)
		for i := range numProfiles {
			profiles[i] = dataBuilder.CreateTestProfile(string(rune('b' + i)))
		}

		// When - save all profiles
		for _, testData := range profiles {
			_, err := repo.Save(ctx, testData.UserID, testData.Email, testData.FirstName, testData.LastName)
			require.NoError(t, err)
		}

		// Then - verify all profiles exist
		for _, testData := range profiles {
			profile, err := repo.GetById(ctx, testData.UserID)
			assert.NoError(t, err)
			assert.Equal(t, testData.UserID, profile.UserID)
			assert.Equal(t, testData.Email, profile.Email)
		}

		// Cleanup - delete all profiles
		for _, testData := range profiles {
			err := repo.Delete(ctx, testData.UserID)
			assert.NoError(t, err)
		}
	})
}
