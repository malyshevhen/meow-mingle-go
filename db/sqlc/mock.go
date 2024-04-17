package db

import (
	"context"
)

type MockStore struct {
	err                 error
	user                User
	post                Post
	comment             Comment
	listPostsRows       []ListUserPostsRow
	getPostRow          GetPostRow
	getUserRow          GetUserRow
	listPostCommentRows []ListPostCommentsRow
	listPostUserRows    []ListUserPostsRow
}

func (m *MockStore) SetUser(user User) {
	m.user = user
}

func (m *MockStore) SetPost(post Post) {
	m.post = post
}

func (m *MockStore) SetComment(comment Comment) {
	m.comment = comment
}

func (m *MockStore) SetListPostsRows(rows []ListUserPostsRow) {
	m.listPostsRows = rows
}

func (m *MockStore) SetGetPostRow(row GetPostRow) {
	m.getPostRow = row
}

func (m *MockStore) SetGetUserRow(row GetUserRow) {
	m.getUserRow = row
}

func (m *MockStore) SetListPostCommentRows(rows []ListPostCommentsRow) {
	m.listPostCommentRows = rows
}

func (m *MockStore) SetListUserPostsRows(rows []ListUserPostsRow) {
	m.listPostUserRows = rows
}

func (m *MockStore) SetError(err error) {
	m.err = err
}

// CreateCommentLikeTx implements IStore.
func (m *MockStore) CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error) {
	return m.err
}

// CreateCommentTx implements IStore.
func (m *MockStore) CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error) {
	return m.comment, m.err
}

// CreatePostLikeTx implements IStore.
func (m *MockStore) CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error {
	return m.err
}

// CreatePostTx implements IStore.
func (m *MockStore) CreatePostTx(ctx context.Context, params CreatePostParams) (post Post, err error) {
	return m.post, m.err
}

// CreateSubscriptionTx implements IStore.
func (m *MockStore) CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error {
	return m.err
}

// CreateUserTx implements IStore.
func (m *MockStore) CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error) {
	return m.user, m.err
}

// DeleteCommentLikeTx implements IStore.
func (m *MockStore) DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error {
	return m.err
}

// DeleteCommentTx implements IStore.
func (m *MockStore) DeleteCommentTx(ctx context.Context, userId int64, commentId int64) (err error) {
	return m.err
}

// DeletePostLikeTx implements IStore.
func (m *MockStore) DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error {
	return m.err
}

// DeletePostTx implements IStore.
func (m *MockStore) DeletePostTx(ctx context.Context, userId int64, postId int64) error {
	return m.err
}

// DeleteSubscriptionTx implements IStore.
func (m *MockStore) DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error {
	return m.err
}

// GetFeed implements IStore.
func (m *MockStore) GetFeed(ctx context.Context, userId int64) (feed []ListUserPostsRow, err error) {
	return m.listPostsRows, m.err
}

// GetPostTx implements IStore.
func (m *MockStore) GetPostTx(ctx context.Context, id int64) (post GetPostRow, err error) {
	return m.getPostRow, m.err
}

// GetUserTx implements IStore.
func (m *MockStore) GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error) {
	return m.getUserRow, m.err
}

// ListPostCommentsTx implements IStore.
func (m *MockStore) ListPostCommentsTx(ctx context.Context, id int64) (posts []ListPostCommentsRow, err error) {
	return m.listPostCommentRows, m.err
}

// ListUserPostsTx implements IStore.
func (m *MockStore) ListUserPostsTx(ctx context.Context, userId int64) (posts []ListUserPostsRow, err error) {
	return m.listPostUserRows, m.err
}

// UpdateCommentTx implements IStore.
func (m *MockStore) UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error) {
	return m.comment, m.err
}

// UpdatePostTx implements IStore.
func (m *MockStore) UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error) {
	return m.post, m.err
}
