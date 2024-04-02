package main

import (
	"log"
	"net/http"
)

type BasicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BasicError) Error() string {
	return e.Message
}

func ErrorHandler(h apiHandler) apiHandler {
	log.Printf("%-15s Apply error handler 🕵️", "Error Handler")

	return func(w http.ResponseWriter, r *http.Request) error {
		if err := h(w, r); err != nil {
			if e, ok := err.(*BasicError); ok {
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)

				WriteJson(w, e.Code, e)
			} else {
				log.Printf("%-15s ==> Error: %v", "Error Handler", err)

				WriteJson(
					w,
					http.StatusInternalServerError,
					NewErrorResponse("Internal error"),
				)
			}
		}

		return nil
	}
}
