package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := GetTokenFromRequest(r)

		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("%-15s ==> Authentication failed: Invalid JWT token 🚫", "AuthMW")
			WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Permission denied. Invalid JWT token."))
			return
		}

		if !token.Valid {
			log.Printf("%-15s ==> Authentication failed: JWT token not valid ❌", "AuthMW")
			WriteJson(w, http.StatusUnauthorized, NewErrorResponse("Permission denied. JWT token not valid."))
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		stringId := claims["userId"].(string)

		id, err := strconv.Atoi(stringId)
		if err != nil {
			return
		}

		if _, err := store.GetUserById(int64(id)); err != nil {
			log.Printf("%-15s ==> Authentication failed: User Id not found 🆘", "AuthMW")
			WriteJson(w, http.StatusBadRequest, NewErrorResponse("User Id not found."))
			return
		}

		log.Printf("%-15s ==> User %s authenticated successfully ✅", "AuthMW", id)
		handlerFunc(w, r)
	}
}

func GetAuthUserId(t string) (int64, error) {
	token, err := validateJWT(t)
	if err != nil {
		log.Printf("%-15s ==> 😢 Authentication failed: Invalid JWT token", "AuthMW")
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := claims["userId"].(string)
	numId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%-15s ==> 😕 Failed to convert user Id to integer", "AuthMW")
		return 0, nil
	}

	log.Printf("%-15s ==> 🎉 User Id converted to integer successfully!", "AuthMW")
	return int64(numId), nil
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	log.Printf("%-15s ==> 🕵 Validating JWT token...", "AuthMW")

	token, err := jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("%-15s ==> ❌ Unexpected signing method!", "AuthMW")
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		log.Printf("%-15s ==> 🔑 Comparing secret...", "AuthMW")
		return []byte(secret), nil
	})

	if err != nil {
		log.Printf("%-15s ==> 🚨 JWT validation failed!", "AuthMW")
	} else {
		log.Printf("%-15s ==> ✅ JWT token validated successfully!", "AuthMW")
	}

	return token, err
}

func GetTokenFromRequest(r *http.Request) string {
	log.Printf("%-15s ==> 🕵️ Validating for Authorization header...", "AuthMW")

	tokenAuth := r.Header.Get("Authorization")

	if tokenAuth != "" {
		log.Printf("%-15s ==> 🎉 Authorization header found!", "AuthMW")
		return tokenAuth
	}

	log.Printf("%-15s ==> 😢 No Authorization header found.", "AuthMW")
	return ""
}

func HashPwd(s string) (string, error) {
	log.Printf("%-15s ==> 🌈 Starting password hashing...", "AuthMW")

	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("%-15s ==> 😱 Error generating password hash: %v", "AuthMW", err)
		return "", err
	}

	log.Printf("%-15s ==> ✨ Password hashed successfully!", "AuthMW")
	return string(hash), nil
}

func CreateJwt(secret []byte, id int64) (string, error) {
	log.Printf("%-15s ==> 🌟 Starting JWT token creation...", "AuthMW")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(int(id)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	log.Printf("%-15s ==> 🔏 Signing JWT token...", "AuthMW")
	signedToken, err := token.SignedString(secret)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error signing JWT token: %v", "AuthMW", err)
		return "", err
	}

	log.Printf("%-15s ==> ✅ JWT token created successfully!", "AuthMW")
	return signedToken, nil
}
