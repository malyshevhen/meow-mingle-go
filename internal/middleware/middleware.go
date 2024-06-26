package middleware

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/malyshEvhen/meow_mingle/internal/config"
	"github.com/malyshEvhen/meow_mingle/internal/db"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/types"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

func MiddlewareChain(h types.Handler, m ...types.Middleware) http.HandlerFunc {
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

func LoggerMW(h types.Handler) types.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()

		ww := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		err := h(ww, r)
		log.Printf(
			"%-15s ==> %d %s %s %s",
			"Request",
			ww.status,
			r.Method,
			r.RequestURI,
			time.Since(start),
		)

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

func ErrorHandler(h types.Handler) types.Handler {
	log.Printf("%-15s Apply error handler", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			switch e := err.(type) {
			case errors.Error:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				utils.WriteJson(w, e.Code(), NewErrorResponse(e.Error()))
			default:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				utils.WriteJson(
					w,
					http.StatusInternalServerError,
					NewErrorResponse("Internal error"),
				)
			}
		}

		return nil
	}
}

func WithJWTAuth(store db.IStore, cfg config.Config) types.Middleware {
	ctx := context.Background()

	return func(h types.Handler) types.Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			tokenString := utils.GetTokenFromRequest(r)

			token, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
			if err != nil {
				log.Printf("%-15s ==> Authentication failed. Error: %v", "AuthMW", err)
				return errors.NewUnauthorizedError()
			}

			claims := token.Claims.(jwt.MapClaims)
			id := claims["userId"].(string)
			numId, err := strconv.Atoi(id)
			if err != nil {
				log.Printf(
					"%-15s ==> Failed to convert user Id to integer. Error: %v",
					"AuthMW",
					err,
				)
				return errors.NewUnauthorizedError()
			}

			user, err := store.GetUserTx(ctx, int64(numId))
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

			log.Printf("%-15s ==> User %d authenticated successfully", "AuthMW", numId)
			return h(w, r)
		}
	}
}
