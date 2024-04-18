package db

import (
	"context"
)

type MockStore struct {
	err                 error
	user                User
	post                Post
	comment             Comment
	getPostRow          GetPostRow
	getUserRow          GetUserRow
	listPostCommentRows []ListPostCommentsRow
	listPostRows        []ListUserPostsRow
	createSubCalled     bool
	deleteSubCalled     bool
	likeCommentCalled   bool
	unlikeCommentCalled bool
	deleteCommentCalled bool
	createPostCalled    bool
	deletePostCalled    bool
	likePostCalled      bool
	unlikePostCalled    bool
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

func (m *MockStore) SetGetPostRow(row GetPostRow) {
	m.getPostRow = row
}

func (m *MockStore) SetGetUserRow(row GetUserRow) {
	m.getUserRow = row
}

func (m *MockStore) SetListPostCommentRows(rows []ListPostCommentsRow) {
	m.listPostCommentRows = rows
}

func (m *MockStore) AddListUserPostsRows(row ListUserPostsRow) {
	m.listPostRows = append(m.listPostRows, row)
}

func (m *MockStore) AddComments(comment ListPostCommentsRow) {
	m.listPostCommentRows = append(m.listPostCommentRows, comment)
}

func (m *MockStore) SetError(err error) {
	m.err = err
}

func (m *MockStore) CreateSubscriptionCalled() bool {
	result := m.createSubCalled
	m.createSubCalled = false
	return result
}

func (m *MockStore) DeleteSubscriptionCalled() bool {
	result := m.deleteSubCalled
	m.deleteSubCalled = false
	return result
}

func (m *MockStore) LikeCommentCalled() bool {
	result := m.likeCommentCalled
	m.likeCommentCalled = false
	return result
}

func (m *MockStore) UnlikeCommentCalled() bool {
	result := m.unlikeCommentCalled
	m.unlikeCommentCalled = false
	return result
}

func (m *MockStore) CreatePostCalled() bool {
	result := m.createPostCalled
	m.createPostCalled = false
	return result
}

func (m *MockStore) DeletePostCalled() bool {
	result := m.deletePostCalled
	m.deletePostCalled = false
	return result
}

func (m *MockStore) LikePostCalled() bool {
	result := m.likePostCalled
	m.likePostCalled = false
	return result
}

func (m *MockStore) UnlikePostCalled() bool {
	result := m.unlikePostCalled
	m.unlikePostCalled = false
	return result
}

func (m *MockStore) DeleteCommentCalled() bool {
	result := m.deleteCommentCalled
	m.deleteCommentCalled = false
	return result
}

// CreateCommentLikeTx implements IStore.
func (m *MockStore) CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error) {
	m.likeCommentCalled = true
	return m.err
}

// CreateCommentTx implements IStore.
func (m *MockStore) CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error) {
	return m.comment, m.err
}

// CreatePostLikeTx implements IStore.
func (m *MockStore) CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error {
	m.likePostCalled = true
	return m.err
}

// CreatePostTx implements IStore.
func (m *MockStore) CreatePostTx(ctx context.Context, params CreatePostParams) (post Post, err error) {
	return m.post, m.err
}

// CreateSubscriptionTx implements IStore.
func (m *MockStore) CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error {
	m.createSubCalled = true
	return m.err
}

// CreateUserTx implements IStore.
func (m *MockStore) CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error) {
	return m.user, m.err
}

// DeleteCommentLikeTx implements IStore.
func (m *MockStore) DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error {
	m.unlikeCommentCalled = true
	return m.err
}

// DeleteCommentTx implements IStore.
func (m *MockStore) DeleteCommentTx(ctx context.Context, userId int64, commentId int64) (err error) {
	m.deleteCommentCalled = true
	return m.err
}

// DeletePostLikeTx implements IStore.
func (m *MockStore) DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error {
	m.unlikePostCalled = true
	return m.err
}

// DeletePostTx implements IStore.
func (m *MockStore) DeletePostTx(ctx context.Context, userId int64, postId int64) error {
	m.deletePostCalled = true
	return m.err
}

// DeleteSubscriptionTx implements IStore.
func (m *MockStore) DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error {
	m.deleteSubCalled = true
	return m.err
}

// GetFeed implements IStore.
func (m *MockStore) GetFeed(ctx context.Context, userId int64) (feed []ListUserPostsRow, err error) {
	return m.listPostRows, m.err
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
	return m.listPostRows, m.err
}

// UpdateCommentTx implements IStore.
func (m *MockStore) UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error) {
	return m.comment, m.err
}

// UpdatePostTx implements IStore.
func (m *MockStore) UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error) {
	return m.post, m.err
}
