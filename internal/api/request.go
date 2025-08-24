package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

type valid interface {
	validate() error
}

func readBody[T any](r *http.Request) (target T, readErr error) {
	logger := logger.GetLogger().WithComponent("request_reader")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.WithError(err).Error("Error reading request body")
		readErr = errors.NewValidationError("Invalid request body")
		return
	}
	defer r.Body.Close()

	target, err = unmarshal[T](body)
	if err != nil {
		logger.WithError(err).Error("Error unmarshalling JSON")
		readErr = err
		return
	}

	logger.Info("Validating user payload")

	return
}

func readValidBody[T valid](r *http.Request) (value T, err error) {
	value, err = readBody[T](r)
	if err != nil {
		return
	}

	if err = value.validate(); err != nil {
		return
	}

	return
}

func unmarshal[T any](v []byte) (value T, err error) {
	if err := json.Unmarshal(v, &value); err != nil {
		return value, errors.NewValidationError("error parse JSON payload")
	}
	return
}

func idPathParam(r *http.Request) (string, error) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return "", errors.NewValidationError("Invalid 'ID' parameter")
	}

	return id, nil
}

func IsEmpty[T comparable](object *T) bool {
	return *object == *new(T)
}
