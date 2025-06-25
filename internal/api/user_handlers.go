package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func handleRegistration(userRepo db.IUserRepository, secret string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		uForm, err := readReqBody[UserRegistrationForm](r)
		if err != nil {
			log.Printf("%-15s ==> Error reading request body: %v\n", "User Handler", err)
			return err
		}

		user := db.CreateUserParams{
			Email:     uForm.Email,
			FirstName: uForm.FirstName,
			LastName:  uForm.LastName,
			Password:  uForm.Password,
		}

		log.Printf("%-15s ==> Hashing password...", "User Handler")

		hashedPwd, err := HashPwd(user.Password)
		if err != nil {
			log.Printf("%-15s ==> Error hashing password: %v\n", "User Handler", err)
			return err
		}

		user.Password = hashedPwd

		log.Printf("%-15s ==> Creating user in database...\n", "User Handler")

		savedUser, err := userRepo.CreateUser(ctx, user)
		if err != nil {
			log.Printf("%-15s ==> Error creating user: %v\n", "User Handler", err)
			return err
		}

		log.Printf("%-15s ==> Creating auth token...\n", "User Handler")

		secret := []byte(secret)
		token, err := CreateJwt(secret, savedUser.ID)
		if err != nil {
			log.Printf("%-15s ==> Error generating JWT token: %s\n", "User Handler", err)
			return errors.NewValidationError("error create token")
		}

		log.Printf("%-15s ==> Setting auth cookie..\n", "User Handler.")

		http.SetCookie(w, &http.Cookie{
			Name:     TOKEN_COOKIE_KEY,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Duration(TOKEN_EXPIRATION_TIME) * time.Hour),
			Secure:   true,
			HttpOnly: true,
		})

		log.Printf("%-15s ==> User created successfully!\n", "User Handler")

		return WriteJson(w, http.StatusCreated, savedUser)
	}
}

func handleGetUser(userRepo db.IUserRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := parseIdParam(r)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID parameter: '%s' Error: %v", id, err)
			return errors.NewValidationError(msg)
		}

		log.Printf("User ID is %s\n", id)

		authUserID, err := GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if id != authUserID {
			log.Printf(
				"%-15s ==> User with ID: %s have no permissions to access account with ID: %s\n",
				"User Handler",
				authUserID,
				id,
			)
			return errors.NewForbiddenError()
		}

		log.Printf("%-15s ==> Searching for user with Id:%s\n", "User Handler", id)

		savedUser, err := userRepo.GetUserById(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> User not found for Id:%s\n", "User Handler", id)
			return err
		}

		log.Printf("%-15s ==> Found user: %s\n", "User Handler", savedUser.ID)

		return WriteJson(w, http.StatusOK, savedUser)
	}
}

func handleSubscribe(userRepo db.IUserRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		authUserID, err := GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := userRepo.CreateSubscription(ctx, db.CreateSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleUnsubscribe(userRepo db.IUserRepository) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		authUserID, err := GetAuthUserId(r)
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

		return WriteJson(w, http.StatusNoContent, nil)
	}
}
