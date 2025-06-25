package errors

import (
	"net/http"
)

type Error interface {
	Error() string
	Code() int
}

type BasicError struct {
	message string
	code    int
}

func (e *BasicError) Error() string {
	return e.message
}

func (e *BasicError) Code() int {
	return e.code
}

type notFoundError struct {
	BasicError
}

func NewNotFoundError(message string) *notFoundError {
	return &notFoundError{
		BasicError: BasicError{
			message: message,
			code:    http.StatusNotFound,
		},
	}
}

type unauthorizedError struct {
	BasicError
}

func NewUnauthorizedError() *unauthorizedError {
	return &unauthorizedError{
		BasicError: BasicError{
			message: "user is not authorized",
			code:    http.StatusUnauthorized,
		},
	}
}

type forbiddenError struct {
	BasicError
}

func NewForbiddenError() *unauthorizedError {
	return &unauthorizedError{
		BasicError: BasicError{
			message: "access denied",
			code:    http.StatusForbidden,
		},
	}
}

type validationError struct {
	BasicError
}

func NewValidationError(message string) *validationError {
	return &validationError{
		BasicError: BasicError{
			message: message,
			code:    http.StatusBadRequest,
		},
	}
}

type databaseError struct {
	BasicError
}

func NewDatabaseError(err error) *databaseError {
	return &databaseError{
		BasicError: BasicError{
			message: err.Error(),
			code:    http.StatusInternalServerError,
		},
	}
}

type internalServerError struct {
	BasicError
}

func NewInternalServerError(err error) *internalServerError {
	return &internalServerError{
		BasicError: BasicError{
			message: err.Error(),
			code:    http.StatusInternalServerError,
		},
	}
}
