package api

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/config"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
)

func handleCreateUser(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("%-15s ==> Error reading request body: %v\n", "User Handler", err)
			return err
		}
		defer r.Body.Close()

		user, err := Unmarshal[db.CreateUserParams](body)
		if err != nil {
			log.Printf("%-15s ==> Error unmarshal JSON: %v\n", "User Handler", err)
			return err
		}

		log.Printf("%-15s ==> Validating user payload: %s\n", "User Handler", user)

		if err := Validate(user); err != nil {
			return err
		}

		log.Printf("%-15s ==> Hashing password...", "User Handler")

		hashedPwd, err := hashPwd(user.Password)
		if err != nil {
			log.Printf("%-15s ==> Error hashing password: %v\n", "User Handler", err)
			return err
		}

		user.Password = hashedPwd

		log.Printf("%-15s ==> Creating user in database...\n", "User Handler")

		savedUser, err := store.CreateUserTx(ctx, user)
		if err != nil {
			log.Printf("%-15s ==> Error creating user: %v\n", "User Handler", err)
			return err
		}

		log.Printf("%-15s ==> Creating auth token...\n", "User Handler")

		secret := []byte(config.Envs.JWTSecret)
		token, err := createJwt(secret, savedUser.ID)
		if err != nil {
			log.Printf("%-15s ==> Error generating JWT token: %s\n", "User Handler", err)
			return errors.NewValidationError("error create token")
		}

		log.Printf("%-15s ==> Setting auth cookie..\n", "User Handler.")

		http.SetCookie(w, &http.Cookie{
			Name:  "Authorization",
			Value: token,
		})

		log.Printf("%-15s ==> User created successfully!\n", "User Handler")

		return WriteJson(w, http.StatusCreated, savedUser)
	}
}

func handleGetUser(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		authUserID, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if id != authUserID {
			log.Printf("%-15s ==> User with ID: %d have no permissions to access account with ID: %d\n", "User Handler", authUserID, id)
			return errors.NewForbiddenError()
		}

		log.Printf("%-15s ==> Searching for user with Id:%d\n", "User Handler", id)

		savedUser, err := store.GetUserTx(ctx, int64(id))
		if err != nil {
			log.Printf("%-15s ==> User not found for Id:%d\n", "User Handler", id)
			return err
		}

		log.Printf("%-15s ==> Found user: %d\n", "User Handler", savedUser.ID)

		return WriteJson(w, http.StatusOK, savedUser)
	}
}

func handleSubscribe(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		authUserID, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := store.CreateSubscriptionTx(ctx, db.CreateSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleUnsubscribe(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		authUserID, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		if err := store.DeleteSubscriptionTx(ctx, db.DeleteSubscriptionParams{
			UserID:         authUserID,
			SubscriptionID: id,
		}); err != nil {
			return err
		}

		return WriteJson(w, http.StatusNoContent, nil)
	}
}

func handleOwnersFeed(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		authUserID, err := getAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> No authenticated user found", "User Handler")
			return err
		}

		feed, err := store.GetFeed(ctx, authUserID)
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, feed)
	}
}

func handleUsersFeed(store db.IStore) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := context.Background()

		id, err := ParseIdParam(r)
		if err != nil {
			return errors.NewValidationError("ID parameter is invalid")
		}

		feed, err := store.GetFeed(ctx, id)
		if err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, feed)
	}
}
