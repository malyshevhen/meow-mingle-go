// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CountPostLikes(ctx context.Context, postID int64) (int64, error)
	CreateComment(ctx context.Context, arg CreateCommentParams) (Comment, error)
	CreateCommentLike(ctx context.Context, arg CreateCommentLikeParams) error
	CreatePost(ctx context.Context, arg CreatePostParams) (Post, error)
	CreatePostLike(ctx context.Context, arg CreatePostLikeParams) error
	CreateSubscription(ctx context.Context, arg CreateSubscriptionParams) error
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteComment(ctx context.Context, id int64) error
	DeleteCommentLike(ctx context.Context, arg DeleteCommentLikeParams) error
	DeletePost(ctx context.Context, id int64) error
	DeletePostLike(ctx context.Context, arg DeletePostLikeParams) error
	DeleteSubscription(ctx context.Context, arg DeleteSubscriptionParams) error
	DeleteUser(ctx context.Context, id int64) error
	GetComment(ctx context.Context, id int64) (CommentInfo, error)
	GetCommentLike(ctx context.Context, arg GetCommentLikeParams) (CommentLike, error)
	GetCommentsAuthorID(ctx context.Context, id int64) (int64, error)
	GetPost(ctx context.Context, id int64) (PostInfo, error)
	GetPostLike(ctx context.Context, arg GetPostLikeParams) (PostLike, error)
	GetPostsAuthorID(ctx context.Context, id int64) (int64, error)
	GetSubscription(ctx context.Context, arg GetSubscriptionParams) (UsersSubscription, error)
	GetUser(ctx context.Context, id int64) (GetUserRow, error)
	IsUserExists(ctx context.Context, email string) (int64, error)
	ListCommentLikes(ctx context.Context, commentID int64) ([]CommentLike, error)
	ListPostComments(ctx context.Context, postID int64) ([]CommentInfo, error)
	ListPostLikes(ctx context.Context, postID int64) ([]PostLike, error)
	ListSubscriptions(ctx context.Context, userID int64) ([]int64, error)
	ListUserPosts(ctx context.Context, authorID int64) ([]PostInfo, error)
	ListUsers(ctx context.Context) ([]ListUsersRow, error)
	UpdateComment(ctx context.Context, arg UpdateCommentParams) (Comment, error)
	UpdatePost(ctx context.Context, arg UpdatePostParams) (Post, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)