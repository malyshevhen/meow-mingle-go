package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
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
		tokenString := utils.GetTokenFromRequest(r)

		token, err := utils.ValidateJWT(tokenString, ai.cfg.JWTSecret)
		if err != nil {
			log.Printf("%-15s ==> Authentication failed. Error: %v", "AuthMW", err)
			return errors.NewUnauthorizedError()
		}

		claims := token.Claims.(jwt.MapClaims)
		id := claims["userId"].(string)

		user, err := ai.userRepo.GetUser(ctx, id)
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
		return h(w, r)
	}
}
