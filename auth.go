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
			log.Println("failed to authenticate token")
			WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: "permission denied"})
			return
		}

		if !token.Valid {
			log.Println("failed to authenticate token")
			WriteJson(w, http.StatusUnauthorized, ErrorResponse{Error: "permission denied"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		id := claims["userId"].(string)

		_, err = store.GetUserById(id)
		if err != nil {
			WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Id is required"})
			return
		}

		handlerFunc(w, r)
	}
}

func GetAuthUserId(t string) (int64, error) {
	token, err := validateJWT(t)
	if err != nil {
		return 0, err
	}

	claims := token.Claims.(jwt.MapClaims)
	id := claims["userId"].(int64)

	return id, nil
}

func validateJWT(t string) (*jwt.Token, error) {
	secret := Envs.JWTSecret

	return jwt.Parse(t, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func HashPwd(s string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CreateJwt(secret []byte, id int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    strconv.Itoa(int(id)),
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(),
	})

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
