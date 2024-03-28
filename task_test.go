package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateTask(t *testing.T) {
	ms := &MockStore{}
	service := NewTaskService(ms)

	t.Run("should return an error if the name is empty", func(t *testing.T) {
		payload := &Task{
			Name: "",
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.handleCreateTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Fatal("invalid status code")
		}
	})

	t.Run("should create a task", func(t *testing.T) {
		payload := &Task{
			Name:         "Valid name",
			ProjectID:    1,
			AssignedToID: 4,
		}

		b, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}

		req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(b))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks", service.handleCreateTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

func TestGetTask(t *testing.T) {
	ms := &MockStore{}
	service := NewTaskService(ms)

	t.Run("should return a task", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/tasks/42", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := mux.NewRouter()

		router.HandleFunc("/tasks/{id}", service.handleGetTask)

		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatal("invalid status code")
		}
	})

}
