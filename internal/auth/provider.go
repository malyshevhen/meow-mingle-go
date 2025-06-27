package auth

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (user *User, err error)
}

type Provider struct {
	userRepo UserRepository
	secret   string
}

func (ai *Provider) Basic(h api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		ctx := r.Context()

		email, password, err := getCredentials(r)
		if err != nil {
			return err
		}

		user, err := ai.userRepo.GetByEmail(ctx, email)
		if err != nil {
			log.Printf("%-15s ==> Authentication failed: User Id not found. Error: %v", "Auth", err)
			return errors.NewUnauthorizedError()
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return errors.NewUnauthorizedError()
		}

		rCtx := context.WithValue(r.Context(), UserIdKey, user.ID)
		r = r.WithContext(rCtx)

		log.Printf("%-15s ==> User %s authenticate successfully", "Auth", email)

		return h(w, r)
	}
}

func NewProvider(userRepo UserRepository) *Provider {
	return &Provider{
		userRepo: userRepo,
	}
}

func getCredentials(r *http.Request) (email, password string, err error) {
	authHeader := r.Header.Get("Authorization")

	encodedCredsStr, ok := strings.CutPrefix(authHeader, "Basic ")
	if !ok {
		return "", "", errors.NewUnauthorizedError()
	}

	decodedCredBytes, err := base64.StdEncoding.DecodeString(encodedCredsStr)
	if err != nil {
		return "", "", errors.NewUnauthorizedError()
	}

	creds := strings.Split(string(decodedCredBytes), ":")
	email = creds[0]
	password = creds[1]

	return email, password, nil
}
