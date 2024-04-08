package main

import "errors"

type BasicError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *BasicError) Error() string {
	return e.Message
}

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")
