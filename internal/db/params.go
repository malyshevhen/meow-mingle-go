package db

type CreateCommentLikeParams struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id" validate:"required"`
	CommentID string `json:"comment_id" validate:"required"`
}

type DeleteCommentLikeParams struct {
	CommentID string `json:"comment_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

type GetCommentLikeParams struct {
	CommentID string `json:"comment_id" validate:"required"`
	UserID    string `json:"user_id" validate:"required"`
}

type CreateCommentParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
	PostID   string `json:"post_id" validate:"required"`
}

type UpdateCommentParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId string `json:"author_id"`
}

type CreatePostLikeParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id" validate:"required"`
	PostID string `json:"post_id" validate:"required"`
}

type DeletePostLikeParams struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}

type GetPostLikeParams struct {
	PostID string `json:"post_id" validate:"required"`
	UserID string `json:"user_id" validate:"required"`
}

type CreatePostParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorID string `json:"author_id" validate:"required"`
}

type UpdatePostParams struct {
	ID       string `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId string `json:"author_id"`
}

type CreateSubscriptionParams struct {
	ID             string `json:"id"`
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
}

type DeleteSubscriptionParams struct {
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
}

type GetSubscriptionParams struct {
	UserID         string `json:"user_id"`
	SubscriptionID string `json:"subscription_id"`
}

type CreateUserParams struct {
	ID        string `json:"id"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUserParams struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
