package db

import "context"

type MockStore struct {
	err error
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
	return Comment{}, m.err
}

// CreatePostLikeTx implements IStore.
func (m *MockStore) CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error {
	return m.err
}

// CreatePostTx implements IStore.
func (m *MockStore) CreatePostTx(ctx context.Context, authorId int64, content string) (post Post, err error) {
	return Post{}, m.err
}

// CreateSubscriptionTx implements IStore.
func (m *MockStore) CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error {
	return m.err
}

// CreateUserTx implements IStore.
func (m *MockStore) CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error) {
	return User{}, m.err
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
	return []ListUserPostsRow{}, m.err
}

// GetPostTx implements IStore.
func (m *MockStore) GetPostTx(ctx context.Context, id int64) (post GetPostRow, err error) {
	return GetPostRow{}, m.err
}

// GetUserTx implements IStore.
func (m *MockStore) GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error) {
	return GetUserRow{}, m.err
}

// ListPostCommentsTx implements IStore.
func (m *MockStore) ListPostCommentsTx(ctx context.Context, id int64) (posts []ListPostCommentsRow, err error) {
	return []ListPostCommentsRow{}, m.err
}

// ListUserPostsTx implements IStore.
func (m *MockStore) ListUserPostsTx(ctx context.Context, userId int64) (posts []ListUserPostsRow, err error) {
	return []ListUserPostsRow{}, m.err
}

// UpdateCommentTx implements IStore.
func (m *MockStore) UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error) {
	return Comment{}, m.err
}

// UpdatePostTx implements IStore.
func (m *MockStore) UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error) {
	return Post{}, m.err
}
