package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
)

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

func readCreatePostParams(r *http.Request) (*db.CreatePostParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := Unmarshal[db.CreatePostParams](body)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func readUpdatePostParams(r *http.Request) (*db.UpdatePostParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> Error reading post request %v\n", "Post Handler", err)
		return nil, errors.NewValidationError("parameter ID is not valid")
	}
	defer r.Body.Close()

	p, err := Unmarshal[db.UpdatePostParams](body)
	if err != nil {
		log.Printf("%-15s ==> Error reading post request %v\n", "Post Handler", err)
		return nil, err
	}

	return &p, nil
}

func readCreateCommentParams(r *http.Request) (*db.CreateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[db.CreateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func readUpdateCommentParams(r *http.Request) (*db.UpdateCommentParams, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	c, err := Unmarshal[db.UpdateCommentParams](body)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
