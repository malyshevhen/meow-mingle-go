package api

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
	"github.com/malyshEvhen/meow_mingle/types"
)

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

func ErrorHandler(h Handler) Handler {
	log.Printf("%-15s Apply error handler 🕵️", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			switch e := err.(type) {
			case errors.Error:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				WriteJson(w, e.Code(), types.NewErrorResponse(e.Error()))
			default:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				WriteJson(
					w,
					http.StatusInternalServerError,
					types.NewErrorResponse("Internal error"),
				)
			}
		}

		return nil
	}
}

func WithJWTAuth(store *db.Store, handlerFunc Handler) Middleware {
	ctx := context.Background()

	return func(h Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request) error {
			id, err := GetAuthUserId(r)
			if err != nil {
				return errors.NewUnauthorizedError()
			}

			if _, err = store.GetUser(ctx, int64(id)); err != nil {
				log.Printf("%-15s ==> Authentication failed: User Id not found 🆘", "AuthMW")
				return errors.NewUnauthorizedError()
			}

			log.Printf("%-15s ==> User %d authenticated successfully ✅", "AuthMW", id)
			return handlerFunc(w, r)
		}
	}
}
