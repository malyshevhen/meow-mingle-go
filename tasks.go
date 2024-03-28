package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

var errUserIDRequired = errors.New("user id is required")
var errNameRequired = errors.New("name is required")
var errProjectIDRequired = errors.New("user id is required")

type TaskService struct {
	store Store
}

func NewTaskService(s Store) *TaskService {
	return &TaskService{store: s}
}

func (ts *TaskService) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/tasks", WithJWTAuth(ts.handleCreateTask, ts.store)).Methods("POST")
	r.HandleFunc("/tasks/{id}", WithJWTAuth(ts.handleGetTask, ts.store)).Methods("GET")
}

func (ts *TaskService) handleCreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid payload"})
		return
	}

	defer r.Body.Close()

	var task *Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid payload"})
		return
	}

	if err := validateTaskPayload(task); err != nil {
		WriteJson(w, http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	t, err := ts.store.CreateTask(task)
	if err != nil {
		WriteJson(w, http.StatusInternalServerError, ErrorResponse{Error: "Error creating task"})
		return
	}

	WriteJson(w, http.StatusCreated, t)
}

func validateTaskPayload(task *Task) error {
	if task.Name == "" {
		return errNameRequired
	}

	if task.ProjectID == 0 {
		return errProjectIDRequired
	}

	if task.ProjectID == 0 {
		return errUserIDRequired
	}

	return nil
}

func (ts *TaskService) handleGetTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	t, err := ts.store.GetTask(id)
	if err != nil {
		WriteJson(w, http.StatusNotFound, ErrorResponse{Error: "task is not found"})
		return
	}

	WriteJson(w, http.StatusOK, t)
}
