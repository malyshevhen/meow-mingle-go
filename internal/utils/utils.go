package utils

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/malyshEvhen/meow_mingle/internal/errors"
)

type ContextKey string

const UserIdKey ContextKey = "userId"

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return err
		}
	}

	return nil
}

func Unmarshal[T any](v []byte) (value T, err error) {
	if err := json.Unmarshal(v, &value); err != nil {
		return value, errors.NewValidationError("error parse JSON payload")
	}
	return
}

func Validate(s interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(s); err != nil {
		return errors.NewValidationError(err.Error())
	}
	return nil
}

func ParseIdParam(r *http.Request) (int64, error) {
	id := r.PathValue("id")

	numId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("%-15s ==> Error parsing Id parameter %v\n", "Post Handler", err)
		return 0, errors.NewValidationError("Error parsing Id parameter")
	}

	return int64(numId), nil
}
