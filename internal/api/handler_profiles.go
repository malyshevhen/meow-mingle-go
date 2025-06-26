package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/auth"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleCreateProfile(userRepo db.IProfileRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		profileForm, err := readReqBody[CreateProfileForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading request body: %v\n", "User Handler", err)
			return err
		}

		user := db.CreateProfileParams{
			UserID:    profileForm.UserID,
			Email:     profileForm.Email,
			FirstName: profileForm.FirstName,
			LastName:  profileForm.LastName,
		}

		log.Printf("%-15s ==> Creating profile in database...\n", "Profile Handler")

		savedUser, err := userRepo.CreateProfile(ctx, user)
		if err != nil {
			log.Printf("%-15s ==> Error creating profile: %v\n", "Profile Handler", err)
			return err
		}

		log.Printf("%-15s ==> Profile created successfully!\n", "Profile Handler")

		return writeJson(w, http.StatusCreated, savedUser)
	}
}

func handleGetProfile(profileRepo db.IProfileRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id, err := parseIdParam(r)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID parameter: '%s' Error: %v", id, err)
			return errors.NewValidationError(msg)
		}

		profile, err := profileRepo.GetProfileById(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> Profile not found for Id:%s\n", "User Handler", id)
			return err
		}

		log.Printf("%-15s ==> Found profile: %s\n", "profile Handler", profile.ID)

		return writeJson(w, http.StatusOK, profile)
	}
}

func handleSubscribe(profileRepo db.IProfileRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		authUserID, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := profileRepo.CreateSubscription(ctx, db.CreateSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return writeJson(w, http.StatusNoContent, nil)
	}
}

func handleUnsubscribe(userRepo db.IProfileRepository) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		id := mux.Vars(r)["id"]

		authUserID, err := auth.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := userRepo.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return writeJson(w, http.StatusNoContent, nil)
	}
}
