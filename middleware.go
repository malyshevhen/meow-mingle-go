package main

import (
	"context"
	"log"
	"net/http"
	"time"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

type Middleware func(h apiHandler) apiHandler

func MiddlewareChain(h apiHandler, m ...Middleware) http.HandlerFunc {
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

func LoggerMiddleware(h apiHandler) apiHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		start := time.Now()

		ww := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		err := h(ww, r)
		log.Println(ww.status, r.Method, r.URL.Path, time.Since(start))

		return err
	}
}

func ErrorHandler(h apiHandler) apiHandler {
	log.Printf("%-15s Apply error handler 🕵️", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			if e, ok := err.(*BasicError); ok {
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)

				WriteJson(w, e.Code, e)
			} else {
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

func WithJWTAuth(store *db.Store, handlerFunc apiHandler) Middleware {
	return func(h apiHandler) apiHandler {
		return func(w http.ResponseWriter, r *http.Request) error {
			ctx := context.Background()

			id, err := GetAuthUserId(r)
			if err != nil {
				return &BasicError{
					Code:    http.StatusUnauthorized,
					Message: "Access denied",
				}
			}

			if _, err = store.GetUser(ctx, int64(id)); err != nil {
				log.Printf("%-15s ==> Authentication failed: User Id not found 🆘", "AuthMW")
				return &BasicError{
					Code:    http.StatusUnauthorized,
					Message: "Access denied",
				}

			}

			log.Printf("%-15s ==> User %d authenticated successfully ✅", "AuthMW", id)
			return handlerFunc(w, r)
		}
	}
}
