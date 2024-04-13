package main

import "time"

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Error     string    `json:"message"`
}

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error:     message,
		Timestamp: time.Now(),
	}
}

type User struct {
	CreatedAt time.Time `json:"createdAt"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Password  string    `json:"password"`
	ID        int64     `json:"id"`
}

type CommentRequest struct {
	Content string `json:"content"`
}

type CommentResponse struct {
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Content  string    `json:"content"`
	Id       int64     `json:"id"`
	AuthorId int64     `json:"authorId"`
	PostId   int64     `json:"postId"`
	Likes    int       `json:"likes"`
}

type PostRequest struct {
	Content string `json:"content"`
}

type PostResponse struct {
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Content  string    `json:"content"`
	Id       int64     `json:"id"`
	AuthorId int64     `json:"authorId"`
	Likes    int       `json:"likes"`
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

type PostLike struct {
	UserId int64 `json:"userId"`
	PostId int64 `json:"postId"`
}

type CommentLike struct {
	UserId    int64 `json:"userId"`
	CommentId int64 `json:"commentId"`
}
