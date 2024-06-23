package middleware

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	userRepo db.IUserReposytory
	cfg      config.Config
}

func NewAuthProvider(userRepo db.IUserReposytory, cfg config.Config) *AuthProvider {
	return &AuthProvider{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (ai *AuthProvider) WithJWTAuth(h types.Handler) types.Handler {
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) error {
		authCookie, err := utils.GetAuthCookie(r)
		if err != nil {
			return errors.NewUnauthorizedError()
		}

		tokenString := utils.GetTokenFromCookie(authCookie)

		token, err := utils.ValidateJWT(tokenString, ai.cfg.JWTSecret)
		if err != nil {
			log.Printf("%-15s ==> Authentication failed. Error: %v", "AuthMW", err)
			return errors.NewUnauthorizedError()
		}

		claims := token.Claims.(jwt.MapClaims)
		id := claims["userId"].(string)

		user, err := ai.userRepo.GetUserById(ctx, id)
		if err != nil {
			log.Printf(
				"%-15s ==> Authentication failed: User Id not found. Error: %v",
				"AuthMW",
				err,
			)
			return errors.NewUnauthorizedError()
		}

		rCtx := context.WithValue(r.Context(), utils.UserIdKey, user.ID)
		r = r.WithContext(rCtx)

		log.Printf("%-15s ==> User %s athenticated successfully", "AuthMW", id)

		http.SetCookie(w, authCookie)

		return h(w, r)
	}
}

func (ai *AuthProvider) WithBasicAuth(h types.Handler) types.Handler {
	ctx := context.Background()

	return func(w http.ResponseWriter, r *http.Request) error {
		authHeader := r.Header.Get("Authentication")

		encodedCredsStr, ok := strings.CutPrefix(authHeader, "Basic ")
		if !ok {
			return errors.NewUnauthorizedError()
		}

		decodedCredBytes, err := base64.StdEncoding.DecodeString(encodedCredsStr)
		if err != nil {
			return errors.NewUnauthorizedError()
		}

		creds := strings.Split(string(decodedCredBytes), ":")
		email := creds[0]
		password := creds[1]

		user, err := ai.userRepo.GetUserByEmail(ctx, email)
		if err != nil {
			log.Printf(
				"%-15s ==> Authentication failed: User Id not found. Error: %v",
				"AuthMW",
				err,
			)
			return errors.NewUnauthorizedError()
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			return errors.NewUnauthorizedError()
		}

		rCtx := context.WithValue(r.Context(), utils.UserIdKey, user.ID)
		r = r.WithContext(rCtx)

		log.Printf("%-15s ==> User %s athenticated successfully", "AuthMW", email)

		return h(w, r)
	}
}
