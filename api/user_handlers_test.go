package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
	"github.com/malyshEvhen/meow_mingle/errors"
)

func TestHandleCreateUser(t *testing.T) {

	user := db.User{
		ID:        1,
		Email:     "john@doe.com",
		FirstName: "John",
		LastName:  "Doe",
		Password:  "password",
		CreatedAt: time.Time{},
	}

	validParams := db.CreateUserParams{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Password:  user.Password,
	}

	invalidParams := db.CreateUserParams{
		Email:     "john@doe.com",
		FirstName: "",
		LastName:  "",
		Password:  "",
	}

	t.Run("should create a new user", func(t *testing.T) {
		store := &db.MockStore{}
		store.SetUser(user)

		req, err := http.NewRequest("POST", "/users", reqBodyOf(validParams))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MiddlewareChain(handleCreateUser(store), LoggerMiddleware, ErrorHandler))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusCreated)
		}
	})

	t.Run("should return 400 if user already exists", func(t *testing.T) {
		store := &db.MockStore{}
		store.SetUser(user)
		store.SetError(errors.NewValidationError("User already exists"))

		req, err := http.NewRequest("POST", "/users", reqBodyOf(validParams))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MiddlewareChain(handleCreateUser(store), LoggerMiddleware, ErrorHandler))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("should return 400 if params are invalid", func(t *testing.T) {
		store := &db.MockStore{}

		req, err := http.NewRequest("POST", "/users", reqBodyOf(invalidParams))
		if err != nil {
			t.Fatal(err)
		}

		defer req.Body.Close()
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MiddlewareChain(handleCreateUser(store), LoggerMiddleware, ErrorHandler))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})

	t.Run("should return 400 if body is empty", func(t *testing.T) {
		store := &db.MockStore{}

		req, err := http.NewRequest("POST", "/users", reqBodyOf(db.CreateUserParams{}))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MiddlewareChain(handleCreateUser(store), LoggerMiddleware, ErrorHandler))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusBadRequest)
		}
	})
}

func reqBodyOf(content interface{}) io.Reader {
	jsonBytes, _ := json.Marshal(content)

	return bytes.NewBuffer(jsonBytes)
}

func TestHandleGetUser(t *testing.T) {

	t.Run("should return 200 and user if user exists", func(t *testing.T) {
		userRow := db.GetUserRow{
			ID:        1,
			Email:     "john@doe.com",
			FirstName: "John",
			LastName:  "Doe",
			CreatedAt: time.Time{},
		}
		store := &db.MockStore{}
		store.SetGetUserRow(userRow)

		// Handler can`t read URL and path vars... WHY?!
		req, err := http.NewRequest("GET", "/users/1", bytes.NewBuffer([]byte{}))
		if err != nil {
			t.Fatal(err)
		}
		defer req.Body.Close()

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(MiddlewareChain(handleGetUser(store), LoggerMiddleware, ErrorHandler))
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Log("Error message: ", rr.Body.String())
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var respUser db.User
		err = json.Unmarshal(rr.Body.Bytes(), &respUser)
		if err != nil {
			t.Errorf("failed to unmarshal response: %v", err)
		}

		if !reflect.DeepEqual(respUser, userRow) {
			t.Errorf("handler returned unexpected body: got %v want %v", respUser, userRow)
		}
	})

}
