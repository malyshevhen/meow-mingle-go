package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errEmailRequired = errors.New("email is required")
var errFirstNameRequired = errors.New("first name is required")
var errLastNameRequired = errors.New("last name is required")
var errPasswordRequired = errors.New("password is required")

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (ts *UserService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/users/register", ts.handleCreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", WithJWTAuth(ts.handleGetUser, ts.store)).Methods("GET")
}

func (ts *UserService) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	defer r.Body.Close()

	var user *User
	err = json.Unmarshal(body, &user)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	if err := validateUserPayload(user); err != nil {
		WriteJson(w, http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	hashedPwd, err := HashPwd(user.Password)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}
	user.Password = hashedPwd

	u, err := ts.store.CreateUser(user)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating task"))
		return
	}

	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating task"))
		return
	}

	WriteJson(w, http.StatusCreated, token)
}

func (ts *UserService) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	u, err := ts.store.GetUserById(id)
	if err != nil {
		WriteJson(w, http.StatusNotFound, NewErrorResponse("user is not found"))
		return
	}

	WriteJson(w, http.StatusOK, u)
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	secret := Envs.JWTSecret
	token, err := CreateJwt([]byte(secret), id)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

func validateUserPayload(user *User) error {
	if user.Email == "" {
		return errEmailRequired
	}
	if user.FirstName == "" {
		return errFirstNameRequired
	}
	if user.LastName == "" {
		return errLastNameRequired
	}
	if user.Password == "" {
		return errPasswordRequired
	}

	return nil
}
