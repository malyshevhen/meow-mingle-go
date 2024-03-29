package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
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
		log.Printf("%-15s ==> 😞 Error reading request body: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	defer r.Body.Close()

	var user *User
	err = json.Unmarshal(body, &user)
	if err != nil {
		log.Printf("%-15s ==> 😕 Error unmarshal JSON: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}

	log.Printf("%-15s ==> 👀 Validating user payload...", "UserService")
	if err := validateUserPayload(user); err != nil {
		log.Printf("%-15s ==> ❌ Validation failed: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse(err.Error()))
		return
	}

	log.Printf("%-15s ==> 🔑 Hashing password...", "UserService")
	hashedPwd, err := HashPwd(user.Password)
	if err != nil {
		log.Printf("%-15s ==> 🔒 Error hashing password: %v\n", "UserService", err)
		WriteJson(w, http.StatusBadRequest, NewErrorResponse("Invalid payload"))
		return
	}
	user.Password = hashedPwd

	log.Printf("%-15s ==> 📝 Creating user in database...\n", "UserService")
	u, err := ts.store.CreateUser(user)
	if err != nil {
		log.Printf("%-15s ==> 🛑 Error creating user: %v\n", "UserService", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating user"))
		return
	}

	log.Printf("%-15s ==> 🔐 Creating auth token...\n", "UserService")
	token, err := createAndSetAuthCookie(u.ID, w)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error creating auth token: %v\n", "UserService", err)
		WriteJson(w, http.StatusInternalServerError, NewErrorResponse("Error creating auth token"))
		return
	}

	log.Printf("%-15s ==> ✅ User created successfully!\n", "UserService")
	WriteJson(w, http.StatusCreated, token)
}

func (ts *UserService) handleGetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	log.Printf("%-15s ==> 🕵️ Searching for user with ID:%s\n", "UserService", id)

	u, err := ts.store.GetUserById(id)
	if err != nil {
		log.Printf("%-15s ==> 😕 User not found for ID:%s\n", "UserService", id)
		WriteJson(w, http.StatusNotFound, NewErrorResponse("user is not found"))
		return
	}

	log.Printf("%-15s ==> 👍 Found user: %d\n", "UserService", u.ID)
	WriteJson(w, http.StatusOK, u)
}

func createAndSetAuthCookie(id int64, w http.ResponseWriter) (string, error) {
	log.Printf("%-15s ==> 🔑 Generating JWT token..\n", "UserService.")
	secret := Envs.JWTSecret
	token, err := CreateJwt([]byte(secret), id)
	if err != nil {
		log.Printf("%-15s ==> ❌ Error generating JWT token: %s\n", "UserService", err)
		return "", err
	}

	log.Printf("%-15s ==> 🍪 Setting auth cookie..\n", "UserService.")
	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	log.Printf("%-15s ==> ✅ Auth cookie set successfully!\n", "UserService")
	return token, nil
}

func validateUserPayload(user *User) error {
	log.Printf("%-15s ==> 📧 Checking if email is provided..", "UserService.")
	if user.Email == "" {
		log.Printf("%-15s ==> ❌ Email is required but not provided", "UserService")
		return errEmailRequired
	}

	log.Printf("%-15s ==> 📛 Checking if first name is provided..", "UserService.")
	if user.FirstName == "" {
		log.Printf("%-15s ==> ❌ First name is required but not provided", "UserService")
		return errFirstNameRequired
	}

	log.Printf("%-15s ==> 📛 Checking if last name is provided..", "UserService.")
	if user.LastName == "" {
		log.Printf("%-15s ==> ❌ Last name is required but not provided", "UserService")
		return errLastNameRequired
	}

	log.Printf("%-15s ==> 🔑 Checking if password is provided..", "UserService.")
	if user.Password == "" {
		log.Printf("%-15s ==> ❌ Password is required but not provided", "UserService")
		return errPasswordRequired
	}

	log.Printf("%-15s ==> ✅ User payload validation passed!", "UserService")
	return nil
}
