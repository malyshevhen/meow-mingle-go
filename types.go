package main

import "time"

type ErrorResponse struct {
	Error     string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now(),
	}
}

type Task struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	ProjectID    int64     `json:"projectID"`
	AssignedToID int64     `json:"assignedTo"`
	CreatedAt    time.Time `json:"createdAt"`
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	Id       int64     `json:"id"`
	Content  string    `json:"content"`
	AuthorId int64     `json:"authorId"`
	PostId   int64     `json:"postId"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

type PostRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Id       int64             `json:"id"`
	Content  string            `json:"content"`
	AuthorId int64             `json:"authorId"`
	Comments []CommentResponse `json:"comments"`
	Created  time.Time         `json:"created"`
	Updated  time.Time         `json:"updated"`
}

type Page[T any] struct {
	Content          []T   `json:"content"`
	Number           int64 `json:"number"`
	NumberOfElements int64 `json:"numberOfElements"`
	TotalPages       int64 `json:"totalPages"`
	TotalElements    int64 `json:"totalElements"`
	Size             int64 `json:"size"`
	First            bool  `json:"first"`
	Last             bool  `json:"last"`
	Empty            bool  `json:"empty"`
}
