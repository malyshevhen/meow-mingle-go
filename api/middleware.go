package api

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
)

type ContextKey string

const UserIdKey ContextKey = "userId"

type Middleware func(h Handler) Handler

func MiddlewareChain(h Handler, m ...Middleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(m) < 1 {
			h(w, r)
		}

		wrapped := h

		for i := len(m) - 1; i >= 0; i-- {
			wrapped = m[i](wrapped)
		}

		wrapped(w, r)
	}
}

type wrappedWriter struct {
	http.ResponseWriter
	status int
}

func (ww *wrappedWriter) WriteHeader(code int) {
	ww.status = code
	ww.ResponseWriter.WriteHeader(code)
}

func LoggerMiddleware(h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()

		ww := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		err := h(ww, r)
		log.Printf("%-15s ==> %d %s %s %s", "Request", ww.status, r.Method, r.RequestURI, time.Since(start))

		return err
	}
}

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now(),
	}
}

func ErrorHandler(h Handler) Handler {
	log.Printf("%-15s Apply error handler", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			switch e := err.(type) {
			case errors.Error:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				WriteJson(w, e.Code(), NewErrorResponse(e.Error()))
			default:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				WriteJson(
					w,
					http.StatusInternalServerError,
					NewErrorResponse("Internal error"),
				)
			}
		}

		return nil
	}
}

func WithJWTAuth(store db.IStore, handlerFunc Handler) Middleware {
	ctx := context.Background()

	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			tokenString := getTokenFromRequest(r)

			token, err := validateJWT(tokenString)
			if err != nil {
				log.Printf("%-15s ==> Authentication failed: Invalid JWT token", "AuthMW")
				return errors.NewUnauthorizedError()
			}

			claims := token.Claims.(jwt.MapClaims)
			id := claims["userId"].(string)
			numId, err := strconv.Atoi(id)
			if err != nil {
				log.Printf("%-15s ==> Failed to convert user Id to integer", "AuthMW")
				return errors.NewUnauthorizedError()
			}

			user, err := store.GetUserTx(ctx, int64(numId))
			if err != nil {
				log.Printf("%-15s ==> Authentication failed: User Id not found", "AuthMW")
				return errors.NewUnauthorizedError()
			}

			ctx = context.WithValue(r.Context(), UserIdKey, user.ID)
			r = r.WithContext(ctx)

			log.Printf("%-15s ==> User %d authenticated successfully", "AuthMW", numId)
			return handlerFunc(w, r)
		}
	}
}
