package api

import (
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

func handleCreateProfile(profileService app.ProfileService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("profile_handler")
		ctx := r.Context()

		profileForm, err := readValidBody[CreateProfileForm](r)
		if err != nil {
			logger.WithError(err).Error("Error reading request body")
			return err
		}

		logger.Info("Creating profile in database...")

		profile := app.Profile{
			UserID:    profileForm.UserID,
			Email:     profileForm.Email,
			FirstName: profileForm.FirstName,
			LastName:  profileForm.LastName,
		}

		if err := profileService.Create(ctx, &profile); err != nil {
			logger.WithError(err).Error("Error creating profile")
			return err
		}

		logger.Info("Profile created successfully!")

		return writeJSON(w, http.StatusCreated, profile)
	}
}

func handleGetProfile(profileService app.ProfileService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		logger := logger.GetLogger().WithComponent("profile_handler")
		ctx := r.Context()

		id, err := idPathParam(r)
		if err != nil {
			logger.WithError(err).Error("Error parsing Id parameter")
			return err
		}

		profile, err := profileService.GetByID(ctx, id)
		if err != nil {
			logger.WithError(err).Warn("Profile not found for Id: " + id)
			return err
		}

		logger.Info("Found profile: " + profile.UserID)

		return writeJSON(w, http.StatusOK, profile)
	}
}
