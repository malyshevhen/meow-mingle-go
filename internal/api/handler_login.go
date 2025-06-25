package api

import (
	"log"
	"net/http"
	"time"

	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

const TOKEN_EXPIRATION_TIME int = 12

func handleLogin(secret string) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		userId, err := GetAuthUserId(r)
		if err != nil {
			log.Printf("%-15s ==> Error getting authenticated user Id %v\n", "Comment Handler", err)
			return err
		}

		secret := []byte(secret)
		token, err := CreateJwt(secret, userId)
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

		return WriteJson(w, http.StatusOK, nil)
	}
}
