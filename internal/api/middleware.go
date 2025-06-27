package api

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/pkg/api"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
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
		requestID := uuid.New().String()
		reqLogger := logger.GetLogger().WithRequest(r.Method, r.URL.Path, requestID)

		// Add request ID to response headers for tracing
		w.Header().Set("X-Request-ID", requestID)

		ww := &wrappedWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		err := h(ww, r)
		duration := time.Since(start)

		reqLogger.LogRequest(
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			ww.status,
			duration.Milliseconds(),
		)

		if err != nil {
			reqLogger.WithError(err).Error("Request processing failed")
		}

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
	errLogger := logger.GetLogger().WithComponent("error_handler")
	errLogger.Debug("Error handler middleware initialized")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			requestLogger := errLogger.WithRequest(r.Method, r.URL.Path, w.Header().Get("X-Request-ID"))

			switch e := err.(type) {
			case errors.Error:
				requestLogger.WithError(err).Warn("API error occurred",
					"error_code", e.Code(),
					"error_type", "api_error",
				)
				writeJSON(w, e.Code(), NewErrorResponse(e.Error()))
			default:
				requestLogger.WithError(err).Error("Internal server error occurred",
					"error_type", "internal_error",
				)
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
