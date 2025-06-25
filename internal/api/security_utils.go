package api

import (
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

const TOKEN_COOKIE_KEY string = "access_token"

type ContextKey string

const UserIdKey ContextKey = "userId"

func GetAuthUserId(r *http.Request) (string, error) {
	id, ok := r.Context().Value(UserIdKey).(string)
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

func GetTokenFromCookie(authCookie *http.Cookie) string {
	log.Printf("%-15s ==>️ Validating for Authorization header...", "Authentication")

	tokenAuth := authCookie.Value
	if tokenAuth == "" {
		log.Printf("%-15s ==> No access token found.", "Authentication")
		return ""
	}

	log.Printf("%-15s ==> Access token found!", "Authentication")
	return tokenAuth
}

func GetAuthCookie(r *http.Request) (*http.Cookie, error) {
	log.Printf("%-15s ==>️ Extract Authorization cookie...", "Authentication")

	authCookie, err := r.Cookie(TOKEN_COOKIE_KEY)
	if err != nil {
		log.Printf("%-15s ==> No access token cookie found!", "Authentication")
		return nil, err
	}

	return authCookie, nil
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

func CreateJwt(secret []byte, id string) (string, error) {
	log.Printf("%-15s ==> Starting JWT token creation...", "Authentication")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    id,
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	log.Printf("%-15s ==> Signing JWT token...", "Authentication")

	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Printf("%-15s ==> Error signing JWT token: %v", "Authentication", err)
		return "", errors.NewUnauthorizedError()
	}

	log.Printf("%-15s ==> JWT token created successfully!", "Authentication")
	return signedToken, nil
}

func ValidateJWT(t, secret string) (token *jwt.Token, err error) {
	fail := func() (*jwt.Token, error) { return nil, errors.NewUnauthorizedError() }

	log.Printf("%-15s ==> Validating JWT token...", "Authentication")

	token, err = jwt.Parse(t, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("%-15s ==> Unexpected signing method: %v", "Authentication", t.Header["alg"])
			return fail()
		}

		log.Printf("%-15s ==> Comparing secret...", "Authentication")

		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("%-15s ==> JWT validation failed!", "Authentication")
		return fail()
	} else {
		log.Printf("%-15s ==> JWT token validated successfully!", "Authentication")
	}
	return
}
