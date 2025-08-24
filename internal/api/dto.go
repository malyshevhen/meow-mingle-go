package api

import (
	"errors"

	apperrors "github.com/malyshEvhen/meow_mingle/pkg/errors"
)

type CreateProfileForm struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (f CreateProfileForm) validate() error {
	errs := []error{}

	if f.UserID == "" {
		errs = append(errs, apperrors.NewValidationError("User ID is required"))
	}

	if f.Email == "" {
		errs = append(errs, apperrors.NewValidationError("Email is required"))
	}

	if f.FirstName == "" {
		errs = append(errs, apperrors.NewValidationError("First name is required"))
	}

	if f.LastName == "" {
		errs = append(errs, apperrors.NewValidationError("Last name is required"))
	}

	return apperrors.NewValidationError(errors.Join(errs...).Error())
}

type CreatePostForm struct {
	Title       string   `json:"title"`
	Content     string   `json:"content"`
	Image       string   `json:"image"`
	Tags        []string `json:"tags"`
	PublishedAt string   `json:"published_at"`
}

func (f CreatePostForm) validate() error {
	errs := []error{}

	if f.Title == "" {
		errs = append(errs, apperrors.NewValidationError("Title is required"))
	}

	if f.Content == "" {
		errs = append(errs, apperrors.NewValidationError("Content is required"))
	}

	if f.PublishedAt == "" {
		errs = append(errs, apperrors.NewValidationError("Published at is required"))
	}

	return apperrors.NewValidationError(errors.Join(errs...).Error())
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id"`
	Content string `json:"content"`
}

func (r CreateCommentRequest) validate() error {
	errs := []error{}

	if r.PostID == "" {
		errs = append(errs, apperrors.NewValidationError("Post ID is required"))
	}

	if r.Content == "" {
		errs = append(errs, apperrors.NewValidationError("Content is required"))
	}

	return apperrors.NewValidationError(errors.Join(errs...).Error())
}

type ContentForm struct {
	Content string `json:"content"`
}

func (f ContentForm) validate() error {
	errs := []error{}

	if f.Content == "" {
		errs = append(errs, apperrors.NewValidationError("Content is required"))
	}

	return apperrors.NewValidationError(errors.Join(errs...).Error())
}
