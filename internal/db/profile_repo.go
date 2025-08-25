package db

import (
	"context"
	"time"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type profileRepository struct {
	session *gocql.Session
	logger  *logger.Logger
}

// ProfileRepository defines the interface for profile data operations
type ProfileRepository interface {
	Save(ctx context.Context, userID, email, firstName, lastName string) (app.Profile, error)
	SaveProfile(ctx context.Context, profile *app.Profile) error
	GetByID(ctx context.Context, id string) (app.Profile, error)
	GetByEmail(ctx context.Context, email string) (app.Profile, error)
	Update(ctx context.Context, profile *app.Profile) error
	Delete(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (bool, error)
}

// Save creates a new profile with the given parameters
func (pr *profileRepository) Save(ctx context.Context, userID, email, firstName, lastName string) (app.Profile, error) {
	profile := app.Profile{
		UserID:    userID,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := pr.SaveProfile(ctx, &profile); err != nil {
		return app.Profile{}, err
	}

	return profile, nil
}

// SaveProfile saves a complete profile object
func (pr *profileRepository) SaveProfile(ctx context.Context, profile *app.Profile) error {
	if profile == nil {
		return errors.NewValidationError("profile cannot be nil")
	}

	if profile.UserID == "" {
		return errors.NewValidationError("user ID is required")
	}

	if profile.Email == "" {
		return errors.NewValidationError("email is required")
	}

	now := time.Now()
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = now
	}
	profile.UpdatedAt = now

	query := `INSERT INTO mingle.profiles (user_id, email, first_name, last_name, bio, avatar_url, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	bio := ""
	avatarURL := ""

	err := pr.session.Query(query,
		profile.UserID,
		profile.Email,
		profile.FirstName,
		profile.LastName,
		bio,
		avatarURL,
		profile.CreatedAt,
		profile.UpdatedAt,
	).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("profile-repository").Error("Failed to save profile",
			"user_id", profile.UserID,
			"email", profile.Email,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("profile-repository").Info("Profile saved successfully",
		"user_id", profile.UserID,
		"email", profile.Email,
	)

	return nil
}

// GetById retrieves a profile by user ID
func (pr *profileRepository) GetByID(ctx context.Context, id string) (app.Profile, error) {
	if id == "" {
		return app.Profile{}, errors.NewValidationError("user ID is required")
	}

	var profile app.Profile
	var bio, avatarURL string

	query := `SELECT user_id, email, first_name, last_name, bio, avatar_url, created_at, updated_at
			  FROM mingle.profiles WHERE user_id = ?`

	err := pr.session.Query(query, id).WithContext(ctx).Scan(
		&profile.UserID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&bio,
		&avatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			pr.logger.WithComponent("profile-repository").Info("Profile not found",
				"user_id", id,
			)
			return app.Profile{}, errors.NewNotFoundError("profile not found")
		}
		pr.logger.WithComponent("profile-repository").Error("Failed to get profile",
			"user_id", id,
			"error", err.Error(),
		)
		return app.Profile{}, errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("profile-repository").Debug("Profile retrieved successfully",
		"user_id", id,
	)

	return profile, nil
}

// GetByEmail retrieves a profile by email address
func (pr *profileRepository) GetByEmail(ctx context.Context, email string) (app.Profile, error) {
	if email == "" {
		return app.Profile{}, errors.NewValidationError("email is required")
	}

	var profile app.Profile
	var bio, avatarURL string

	query := `SELECT user_id, email, first_name, last_name, bio, avatar_url, created_at, updated_at
			  FROM mingle.profiles WHERE email = ?`

	err := pr.session.Query(query, email).WithContext(ctx).Scan(
		&profile.UserID,
		&profile.Email,
		&profile.FirstName,
		&profile.LastName,
		&bio,
		&avatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if err == gocql.ErrNotFound {
			pr.logger.WithComponent("profile-repository").Info("Profile not found by email",
				"email", email,
			)
			return app.Profile{}, errors.NewNotFoundError("profile not found")
		}
		pr.logger.WithComponent("profile-repository").Error("Failed to get profile by email",
			"email", email,
			"error", err.Error(),
		)
		return app.Profile{}, errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("profile-repository").Debug("Profile retrieved by email successfully",
		"email", email,
		"user_id", profile.UserID,
	)

	return profile, nil
}

// Update updates an existing profile
func (pr *profileRepository) Update(ctx context.Context, profile *app.Profile) error {
	if profile == nil {
		return errors.NewValidationError("profile cannot be nil")
	}

	if profile.UserID == "" {
		return errors.NewValidationError("user ID is required")
	}

	// Check if profile exists
	exists, err := pr.Exists(ctx, profile.UserID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("profile not found")
	}

	profile.UpdatedAt = time.Now()

	query := `UPDATE mingle.profiles
			  SET email = ?, first_name = ?, last_name = ?, updated_at = ?
			  WHERE user_id = ?`

	err = pr.session.Query(query,
		profile.Email,
		profile.FirstName,
		profile.LastName,
		profile.UpdatedAt,
		profile.UserID,
	).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("profile-repository").Error("Failed to update profile",
			"user_id", profile.UserID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("profile-repository").Info("Profile updated successfully",
		"user_id", profile.UserID,
	)

	return nil
}

// Delete removes a profile
func (pr *profileRepository) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.NewValidationError("user ID is required")
	}

	// Check if profile exists
	exists, err := pr.Exists(ctx, userID)
	if err != nil {
		return err
	}

	if !exists {
		return errors.NewNotFoundError("profile not found")
	}

	query := `DELETE FROM mingle.profiles WHERE user_id = ?`

	err = pr.session.Query(query, userID).WithContext(ctx).Exec()
	if err != nil {
		pr.logger.WithComponent("profile-repository").Error("Failed to delete profile",
			"user_id", userID,
			"error", err.Error(),
		)
		return errors.NewDatabaseError(err)
	}

	pr.logger.WithComponent("profile-repository").Info("Profile deleted successfully",
		"user_id", userID,
	)

	return nil
}

// Exists checks if a profile exists
func (pr *profileRepository) Exists(ctx context.Context, userID string) (bool, error) {
	if userID == "" {
		return false, errors.NewValidationError("user ID is required")
	}

	var count int
	query := `SELECT COUNT(*) FROM mingle.profiles WHERE user_id = ?`

	err := pr.session.Query(query, userID).WithContext(ctx).Scan(&count)
	if err != nil {
		pr.logger.WithComponent("profile-repository").Error("Failed to check profile existence",
			"user_id", userID,
			"error", err.Error(),
		)
		return false, errors.NewDatabaseError(err)
	}

	return count > 0, nil
}

func NewProfileRepository(session *gocql.Session) ProfileRepository {
	return &profileRepository{
		session: session,
		logger:  logger.GetLogger(),
	}
}
