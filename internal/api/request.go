package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/malyshEvhen/meow_mingle/pkg/errors"
)

func readBody[T any](r *http.Request) (target T, readErr error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> Error reading request body: %v\n", "User Handler", readErr)
		readErr = errors.NewValidationError("Invalid request body")
		return
	}
	defer r.Body.Close()

	target, err = unmarshal[T](body)
	if err != nil {
		log.Printf("%-15s ==> Error unmarshal JSON: %v\n", "User Handler", readErr)
		readErr = err
		return
	}

	log.Printf("%-15s ==> Validating user payload: %v\n", "User Handler", target)

	if err := validate(target); err != nil {
		readErr = err
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

func validate(s any) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(s); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

func iaPathParam(r *http.Request) (string, error) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return "", errors.NewValidationError("Invalid 'ID' parameter")
	}

	return id, nil
}

func IsEmpty[T comparable](object *T) bool {
	return *object == *new(T)
}
