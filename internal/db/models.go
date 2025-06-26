package db

import (
	"time"
)

type Comment struct {
	ID        string    `json:"id"`
	Content   string    `json:"content" validate:"required"`
	AuthorID  string    `json:"author_id" validate:"required"`
	PostID    string    `json:"post_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CommentInfo struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	PostID    string    `json:"post_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int64     `json:"likes"`
}

type CommentLike struct {
	UserID    string `json:"user_id" validate:"required"`
	CommentID string `json:"comment_id" validate:"required"`
}

type Post struct {
	ID        string    `json:"id"`
	Content   string    `json:"content" validate:"required"`
	AuthorID  string    `json:"author_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostInfo struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Likes     int64     `json:"likes"`
}

type PostLike struct {
	UserID string `json:"user_id" validate:"required"`
	PostID string `json:"post_id" validate:"required"`
}

type Profile struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type UserInfo struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	FirstName string    `json:"first_name" validate:"required"`
	LastName  string    `json:"last_name" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type UsersSubscription struct {
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
}
