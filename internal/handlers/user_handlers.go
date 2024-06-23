package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

const TOKEN_EXPIRATION_TIME int = 12

type UserHandler struct {
	userRepo db.IUserReposytory
}

func NewUserHandler(userRepo db.IUserReposytory) *UserHandler {
	return &UserHandler{
		userRepo: userRepo,
	}
}

func (uh *UserHandler) HandleCreateUser(cfg config.Config) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		uForm, err := ReadReqBody[UserRegistrationForm](r)
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

		hashedPwd, err := utils.HashPwd(user.Password)
		if err != nil {
			log.Printf("%-15s ==> Error hashing password: %v\n", "User Handler", err)
			return err
		}

		user.Password = hashedPwd

		log.Printf("%-15s ==> Creating user in database...\n", "User Handler")

		savedUser, err := uh.userRepo.CreateUser(ctx, user)
		if err != nil {
			log.Printf("%-15s ==> Error creating user: %v\n", "User Handler", err)
			return err
		}

		log.Printf("%-15s ==> Creating auth token...\n", "User Handler")

		secret := []byte(cfg.JWTSecret)
		token, err := utils.CreateJwt(secret, savedUser.ID)
		if err != nil {
			log.Printf("%-15s ==> Error generating JWT token: %s\n", "User Handler", err)
			return errors.NewValidationError("error create token")
		}

		log.Printf("%-15s ==> Setting auth cookie..\n", "User Handler.")

		http.SetCookie(w, &http.Cookie{
			Name:     utils.TOKEN_COOKIE_KEY,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(time.Duration(TOKEN_EXPIRATION_TIME) * time.Hour),
			Secure:   true,
			HttpOnly: true,
		})

		log.Printf("%-15s ==> User created successfully!\n", "User Handler")

		return utils.WriteJson(w, http.StatusCreated, savedUser)
	}
}

func (uh *UserHandler) HandleGetUser() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := utils.ParseIdParam(r)
		if err != nil {
			msg := fmt.Sprintf("Invalid ID parameter: '%s' Error: %v", id, err)
			return errors.NewValidationError(msg)
		}

		log.Printf("User ID is %s\n", id)

		authUserID, err := utils.GetAuthUserId(r)
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

		savedUser, err := uh.userRepo.GetUserById(ctx, id)
		if err != nil {
			log.Printf("%-15s ==> User not found for Id:%s\n", "User Handler", id)
			return err
		}

		log.Printf("%-15s ==> Found user: %s\n", "User Handler", savedUser.ID)

		return utils.WriteJson(w, http.StatusOK, savedUser)
	}
}

func (uh *UserHandler) HandleSubscribe() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		authUserID, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := uh.userRepo.CreateSubscription(ctx, db.CreateSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}

func (uh *UserHandler) HandleUnsubscribe() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id := mux.Vars(r)["id"]

		authUserID, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := uh.userRepo.DeleteSubscription(ctx, db.DeleteSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return utils.WriteJson(w, http.StatusNoContent, nil)
	}
}
