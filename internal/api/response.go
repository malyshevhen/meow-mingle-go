package api

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if v != nil {
		if err := json.NewEncoder(w).Encode(v); err != nil {
			return err
		}
	}

	return nil
}
