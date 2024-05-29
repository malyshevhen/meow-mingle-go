package api

import (
	"io"
	"log"
	"net/http"

	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
)

func ReadReqBody[T any](r *http.Request) (target T, readErr error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("%-15s ==> Error reading request body: %v\n", "User Handler", readErr)
		readErr = errors.NewValidationError("Invalid request body")
		return
	}
	defer r.Body.Close()

	target, err = utils.Unmarshal[T](body)
	if err != nil {
		log.Printf("%-15s ==> Error unmarshal JSON: %v\n", "User Handler", readErr)
		readErr = err
		return
	}

	log.Printf("%-15s ==> Validating user payload: %v\n", "User Handler", target)

	if err := utils.Validate(target); err != nil {
		readErr = err
		return
	}

	return
}

func Map[T, S any](source S, mapper func(S) T) (target T) {
	return mapper(source)
}
