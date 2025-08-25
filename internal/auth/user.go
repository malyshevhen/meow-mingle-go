package auth

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type ContextKey string

const UserIDKey ContextKey = "userId"

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetAuthUserID(r *http.Request) (string, error) {
	id, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		log.Printf("%-15s ==> Failed to convert user Id to integer", "Authentication")
		return "", errors.NewUnauthorizedError()
	}

	log.Printf(
		"%-15s ==> User Id founded in the request context. ID: %s\n",
		"Authentication",
		id,
	)
	return id, nil
}

func HashPwd(s string) (string, error) {
	log.Printf("%-15s ==> Starting password hashing...", "Authentication")

	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%-15s ==> Error generating password hash: %v", "Authentication", err)
		return "", errors.NewInternalServerError(err)
	}

	log.Printf("%-15s ==> Password hashed successfully!", "Authentication")
	return string(hash), nil
}

func UserID(ctx context.Context) string {
	return ctx.Value(UserIDKey).(string)
}
