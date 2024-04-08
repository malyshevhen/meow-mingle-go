package main

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return &BasicError{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			}
		}
	}

	return nil
}

func Unmarshal[T any](v []byte) (value T, err error) {
	json.Unmarshal(v, &value)
	return
}
