package types

import (
	"encoding/json"
	"fmt"
	"time"

	db "github.com/malyshEvhen/meow_mingle/db/sqlc"
)

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
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"firstName" validate:"required"`
	LastName  string    `json:"lastName" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	ID        int64     `json:"id"`
}

func UserFromParams(up db.CreateUserParams) User {
	return User{
		Email:     up.Email,
		FirstName: up.FirstName,
		LastName:  up.LastName,
		Password:  up.Password,
	}
}

func (u *User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CreatedAt time.Time `json:"createdAt"`
		Email     string    `json:"email"`
		FirstName string    `json:"firstName"`
		LastName  string    `json:"lastName"`
		Password  string    `json:"password"`
		ID        int64     `json:"id"`
	}{
		CreatedAt: u.CreatedAt,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  "******",
		ID:        u.ID,
	})
}

func (u *User) String() string {
	return fmt.Sprintf(
		"User { %d, %s, %s, %s, %s, %s }",
		u.ID,
		u.Email,
		u.FirstName,
		u.LastName,
		"******",
		u.CreatedAt,
	)
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
	Content string `json:"content" validate:"required"`
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
	PostId int64 `json:"postId" validate:"required"`
}

type CommentLike struct {
	UserId    int64 `json:"userId"`
	CommentId int64 `json:"commentId" validate:"required"`
}
