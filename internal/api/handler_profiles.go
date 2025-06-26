package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/app"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreateProfile(profileService app.ProfileService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		profileForm, err := readBody[CreateProfileForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading request body: %v\n", "Profile Handler", err)
			return err
		}

		log.Printf("%-15s ==> Creating profile in database...\n", "Profile Handler")

		profile := app.Profile{
			UserID:    profileForm.UserID,
			Email:     profileForm.Email,
			FirstName: profileForm.FirstName,
			LastName:  profileForm.LastName,
		}

		if err := profileService.Create(ctx, &profile); err != nil {
			log.Printf("%-15s ==> Error creating profile: %v\n", "Profile Handler", err)
			return err
		}

		log.Printf("%-15s ==> Profile created successfully!\n", "Profile Handler")

		return writeJSON(w, http.StatusCreated, profile)
	}
}

func handleGetProfile(profileService app.ProfileService) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := iaPathParam(r)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID parameter: '%s' Error: %v", id, err)
			return errors.NewValidationError(msg)
		}

		profile, err := profileService.GetById(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Profile not found for Id:%s\n", "Profile Handler", id)
			return err
		}

		log.Printf("%-15s ==> Found profile: %s\n", "Profile Handler", profile.ID)

		return writeJSON(w, http.StatusOK, profile)
	}
}
