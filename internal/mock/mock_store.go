package mock

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/db"
)

type MockStore struct {
	err                 error
	user                db.User
	post                db.Post
	comment             db.Comment
	getPostRow          db.PostInfo
	getUserRow          db.GetUserRow
	listPostCommentRows []db.CommentInfo
	listPostRows        []db.PostInfo
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

func (m *MockStore) SetUser(user db.User) {
	m.user = user
}

func (m *MockStore) SetPost(post db.Post) {
	m.post = post
}

func (m *MockStore) SetComment(comment db.Comment) {
	m.comment = comment
}

func (m *MockStore) SetGetPostRow(row db.PostInfo) {
	m.getPostRow = row
}

func (m *MockStore) SetGetUserRow(row db.GetUserRow) {
	m.getUserRow = row
}

func (m *MockStore) SetListPostCommentRows(rows []db.CommentInfo) {
	m.listPostCommentRows = rows
}

func (m *MockStore) AddListUserPostsRows(row db.PostInfo) {
	m.listPostRows = append(m.listPostRows, row)
}

func (m *MockStore) SetListUserPostRows(rows []db.PostInfo) {
	m.listPostRows = rows
}

func (m *MockStore) AddComments(comment db.CommentInfo) {
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
func (m *MockStore) CreateCommentLikeTx(
	ctx context.Context,
	params db.CreateCommentLikeParams,
) (err error) {
	m.likeCommentCalled = true
	return m.err
}

// CreateCommentTx implements IStore.
func (m *MockStore) CreateCommentTx(
	ctx context.Context,
	params db.CreateCommentParams,
) (comment db.Comment, err error) {
	return m.comment, m.err
}

// CreatePostLikeTx implements IStore.
func (m *MockStore) CreatePostLikeTx(ctx context.Context, params db.CreatePostLikeParams) error {
	m.likePostCalled = true
	return m.err
}

// CreatePostTx implements IStore.
func (m *MockStore) CreatePostTx(
	ctx context.Context,
	params db.CreatePostParams,
) (post db.Post, err error) {
	return m.post, m.err
}

// CreateSubscriptionTx implements IStore.
func (m *MockStore) CreateSubscriptionTx(
	ctx context.Context,
	params db.CreateSubscriptionParams,
) error {
	m.createSubCalled = true
	return m.err
}

// CreateUserTx implements IStore.
func (m *MockStore) CreateUserTx(
	ctx context.Context,
	params db.CreateUserParams,
) (user db.User, err error) {
	return m.user, m.err
}

// DeleteCommentLikeTx implements IStore.
func (m *MockStore) DeleteCommentLikeTx(
	ctx context.Context,
	params db.DeleteCommentLikeParams,
) error {
	m.unlikeCommentCalled = true
	return m.err
}

// DeleteCommentTx implements IStore.
func (m *MockStore) DeleteCommentTx(
	ctx context.Context,
	userId int64,
	commentId int64,
) (err error) {
	m.deleteCommentCalled = true
	return m.err
}

// DeletePostLikeTx implements IStore.
func (m *MockStore) DeletePostLikeTx(ctx context.Context, params db.DeletePostLikeParams) error {
	m.unlikePostCalled = true
	return m.err
}

// DeletePostTx implements IStore.
func (m *MockStore) DeletePostTx(ctx context.Context, userId int64, postId int64) error {
	m.deletePostCalled = true
	return m.err
}

// DeleteSubscriptionTx implements IStore.
func (m *MockStore) DeleteSubscriptionTx(
	ctx context.Context,
	params db.DeleteSubscriptionParams,
) error {
	m.deleteSubCalled = true
	return m.err
}

// GetFeed implements IStore.
func (m *MockStore) GetFeed(
	ctx context.Context,
	userId int64,
) (feed []db.PostInfo, err error) {
	return m.listPostRows, m.err
}

// GetPostTx implements IStore.
func (m *MockStore) GetPostTx(ctx context.Context, id int64) (post db.PostInfo, err error) {
	return m.getPostRow, m.err
}

// GetUserTx implements IStore.
func (m *MockStore) GetUserTx(ctx context.Context, id int64) (user db.GetUserRow, err error) {
	return m.getUserRow, m.err
}

// ListPostCommentsTx implements IStore.
func (m *MockStore) ListPostCommentsTx(
	ctx context.Context,
	id int64,
) (posts []db.CommentInfo, err error) {
	return m.listPostCommentRows, m.err
}

// ListUserPostsTx implements IStore.
func (m *MockStore) ListUserPostsTx(
	ctx context.Context,
	userId int64,
) (posts []db.PostInfo, err error) {
	return m.listPostRows, m.err
}

// UpdateCommentTx implements IStore.
func (m *MockStore) UpdateCommentTx(
	ctx context.Context,
	params db.UpdateCommentParams,
) (comment db.Comment, err error) {
	return m.comment, m.err
}

// UpdatePostTx implements IStore.
func (m *MockStore) UpdatePostTx(
	ctx context.Context,
	params db.UpdatePostParams,
) (post db.Post, err error) {
	return m.post, m.err
}
