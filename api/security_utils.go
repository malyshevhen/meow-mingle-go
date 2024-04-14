package api

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/config"
	"github.com/malyshEvhen/meow_mingle/errors"
	"golang.org/x/crypto/bcrypt"
)

func getAuthUserId(r *http.Request) (int64, error) {
	tokenString := getTokenFromRequest(r)

	token, err := validateJWT(tokenString)
	if err != nil {
		log.Printf("%-15s ==> 😢 Authentication failed: Invalid JWT token", "Authentication")
		return 0, errors.NewUnauthorizedError()
	}

	claims := token.Claims.(jwt.MapClaims)
	id := claims["userId"].(string)
	numId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%-15s ==> 😕 Failed to convert user Id to integer", "Authentication")
		return 0, errors.NewUnauthorizedError()
	}

	log.Printf("%-15s ==> 🎉 User Id converted to integer successfully! ID: %d\n", "Authentication", numId)
	return int64(numId), nil
}

func getTokenFromRequest(r *http.Request) string {
	log.Printf("%-15s ==> 🕵️ Validating for Authorization header...", "Authentication")

	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		log.Printf("%-15s ==> 🎉 Authorization header found!", "Authentication")
		return tokenAuth
	}

	log.Printf("%-15s ==> 😢 No Authorization header found.", "Authentication")
	return ""
}

func hashPwd(s string) (string, error) {
	log.Printf("%-15s ==> 🌈 Starting password hashing...", "Authentication")

	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("%-15s ==> 😱 Error generating password hash: %v", "Authentication", err)
		return "", errors.NewInternalServerError(err)
	}

	log.Printf("%-15s ==> ✨ Password hashed successfully!", "Authentication")
	return string(hash), nil
}

func createJwt(secret []byte, id int64) (string, error) {
	log.Printf("%-15s ==> 🌟 Starting JWT token creation...", "Authentication")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(int(id)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	log.Printf("%-15s ==> 🔏 Signing JWT token...", "Authentication")

	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error signing JWT token: %v", "Authentication", err)
		return "", errors.NewUnauthorizedError()
	}

	log.Printf("%-15s ==> ✅ JWT token created successfully!", "Authentication")
	return signedToken, nil
}

func validateJWT(t string) (token *jwt.Token, err error) {
	var (
		secret = config.Envs.JWTSecret
		fail   = func() (*jwt.Token, error) { return nil, errors.NewUnauthorizedError() }
	)

	log.Printf("%-15s ==> 🕵 Validating JWT token...", "Authentication")

	token, err = jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("%-15s ==> ❌ Unexpected signing method: %v", "Authentication", t.Header["alg"])
			return fail()
		}

		log.Printf("%-15s ==> 🔑 Comparing secret...", "Authentication")

		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("%-15s ==> 🚨 JWT validation failed!", "Authentication")
		return fail()
	} else {
		log.Printf("%-15s ==> ✅ JWT token validated successfully!", "Authentication")
	}
	return
}
