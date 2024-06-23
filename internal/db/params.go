package db

type CreateCommentLikeParams struct {
	UserID    int64 `json:"user_id" validate:"required"`
	CommentID int64 `json:"comment_id" validate:"required"`
}

type DeleteCommentLikeParams struct {
	CommentID int64 `json:"comment_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
}

type GetCommentLikeParams struct {
	CommentID int64 `json:"comment_id" validate:"required"`
	UserID    int64 `json:"user_id" validate:"required"`
}

type CreateCommentParams struct {
	Content  string `json:"content" validate:"required"`
	AuthorID int64  `json:"author_id" validate:"required"`
	PostID   int64  `json:"post_id" validate:"required"`
}

type UpdateCommentParams struct {
	ID       int64  `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId int64  `json:"author_id"`
}

type CreatePostLikeParams struct {
	UserID int64 `json:"user_id" validate:"required"`
	PostID int64 `json:"post_id" validate:"required"`
}

type DeletePostLikeParams struct {
	PostID int64 `json:"post_id" validate:"required"`
	UserID int64 `json:"user_id" validate:"required"`
}

type GetPostLikeParams struct {
	PostID int64 `json:"post_id" validate:"required"`
	UserID int64 `json:"user_id" validate:"required"`
}

type CreatePostParams struct {
	Content  string `json:"content" validate:"required"`
	AuthorID int64  `json:"author_id" validate:"required"`
}

type UpdatePostParams struct {
	ID       int64  `json:"id"`
	Content  string `json:"content" validate:"required"`
	AuthorId int64  `json:"author_id"`
}

type CreateSubscriptionParams struct {
	UserID         int64 `json:"user_id"`
	SubscriptionID int64 `json:"subscription_id"`
}

type DeleteSubscriptionParams struct {
	UserID         int64 `json:"user_id"`
	SubscriptionID int64 `json:"subscription_id"`
}

type GetSubscriptionParams struct {
	UserID         int64 `json:"user_id"`
	SubscriptionID int64 `json:"subscription_id"`
}

type CreateUserParams struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

type UpdateUserParams struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
}
