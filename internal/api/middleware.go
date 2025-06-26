package api

import (
	"log"
	"net/http"
	"time"

	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func middlewareChain(h api.Handler, m ...api.Middleware) http.HandlerFunc {
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

func loggerMW(h api.Handler) api.Handler {
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

func ErrorHandler(h api.Handler) api.Handler {
	log.Printf("%-15s Apply error handler", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			switch e := err.(type) {
			case errors.Error:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				writeJSON(w, e.Code(), NewErrorResponse(e.Error()))
			default:
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)
				writeJSON(
					w,
					http.StatusInternalServerError,
					NewErrorResponse("Internal error"),
				)
			}
		}

		return nil
	}
}
