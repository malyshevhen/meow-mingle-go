package db

import (
	"context"
)

type IStore interface {
	CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error)
	CreatePostTx(ctx context.Context, params CreatePostParams) (post Post, err error)
	CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error)
	CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error
	CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error)
	CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error
	GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error)
	GetPostTx(ctx context.Context, id int64) (post PostInfo, err error)
	GetFeed(ctx context.Context, userId int64) (feed []PostInfo, err error)
	ListUserPostsTx(ctx context.Context, userId int64) (posts []PostInfo, err error)
	ListPostCommentsTx(ctx context.Context, id int64) (posts []CommentInfo, err error)
	UpdatePostTx(ctx context.Context, params UpdatePostParams) (post Post, err error)
	UpdateCommentTx(ctx context.Context, params UpdateCommentParams) (comment Comment, err error)
	DeletePostTx(ctx context.Context, userId, postId int64) error
	DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error
	DeleteCommentTx(ctx context.Context, userId, commentId int64) (err error)
	DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error
	DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error
}
