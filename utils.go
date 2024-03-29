package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}

	// 📝 Add descriptive logging with emojis
	log.Printf("WriteJson ===> 📝 Responded with JSON status %d and payload: %+v", status, v)
}
