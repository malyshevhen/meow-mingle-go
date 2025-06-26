package api

type CreateProfileForm struct {
	UserID    string `json:"user_id" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}

type CreateCommentRequest struct {
	PostID  string `json:"post_id" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type ContentForm struct {
	Content string `json:"content" validate:"required"`
}
