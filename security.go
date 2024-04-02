package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type ContextKey string

const UserKey ContextKey = "user"

type SecurityContextHolder struct {
	contextMap map[string]interface{}
}

type AuthUserService interface {
	GetUserById(id int64) (*UserRequest, error)
}

func InitSecurityContext(authUserService AuthUserService) *SecurityContextHolder {
	sCtx := &SecurityContextHolder{
		contextMap: make(map[string]interface{}),
	}

	sCtx.contextMap["userService"] = authUserService

	return sCtx
}

func (sCtx *SecurityContextHolder) WithJWTAuth(handlerFunc apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id, err := sCtx.GetAuthUserId(r)
		if err != nil {
			return &BasicError{
				Code:    http.StatusUnauthorized,
				Message: "Access denied",
			}
		}

		us, ok := sCtx.contextMap["userService"].(AuthUserService)
		if !ok {
			return &BasicError{
				Code:    http.StatusUnauthorized,
				Message: "Access denied",
			}
		}

		user, err := us.GetUserById(int64(id))
		if err != nil {
			log.Printf("%-15s ==> Authentication failed: User Id not found 🆘", "AuthMW")
			return &BasicError{
				Code:    http.StatusUnauthorized,
				Message: "Access denied",
			}

		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, UserKey, user))
		*r = *req

		log.Printf("%-15s ==> User %d authenticated successfully ✅", "AuthMW", id)
		return handlerFunc(w, r)
	}
}

func (sCtx *SecurityContextHolder) GetAuthUserId(r *http.Request) (int64, error) {
	tokenString := GetTokenFromRequest(r)

	token, err := validateJWT(tokenString)
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

	log.Printf("%-15s ==> 🎉 User Id converted to integer successfully! ID: %d\n", "AuthMW", numId)
	return int64(numId), nil
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
