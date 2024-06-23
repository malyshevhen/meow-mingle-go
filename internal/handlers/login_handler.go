package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

type LoginHandler struct {
	userRepo *db.IUserReposytory
	cfg      *config.Config
}

func NewLoginHandler(userRepo *db.IUserReposytory, cfg *config.Config) *LoginHandler {
	return &LoginHandler{
		userRepo: userRepo,
	}
}

func (lh *LoginHandler) HahdleLogin() types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, err := utils.GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		secret := []byte(lh.cfg.JWTSecret)
		token, err := utils.CreateJwt(secret, userId)
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

		return utils.WriteJson(w, http.StatusOK, nil)
	}
}
