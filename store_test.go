package main

// Mocks

type MockStore struct {
	Err error
}

func (m *MockStore) SetError(err error) {
	m.Err = err
}

// CreateComment implements Store.
func (m *MockStore) CreateComment(postId int64, userId int64, c *CommentRequest) (*CommentResponse, error) {
	return &CommentResponse{}, m.Err
}

// DeleteCommentById implements Store.
func (m *MockStore) DeleteCommentById(id int64) error {
	return m.Err
}

// GetCommentById implements Store.
func (m *MockStore) GetCommentById(id int64) (*CommentResponse, error) {
	return &CommentResponse{}, m.Err
}

// GetCommentsByPostId implements Store.
func (m *MockStore) GetCommentsByPostId(id int64) ([]*CommentResponse, error) {
	return []*CommentResponse{}, m.Err
}

func (m *MockStore) UpdateCommentById(id int64, c *CommentRequest) (*CommentResponse, error) {
	return &CommentResponse{}, m.Err
}

// CreatePost implements Store.
func (m *MockStore) CreatePost(userId int64, p *PostRequest) (*PostResponse, error) {
	return &PostResponse{}, m.Err
}

// DeletePostById implements Store.
func (m *MockStore) DeletePostById(id int64) error {
	return m.Err
}

// GetPostById implements Store.
func (m *MockStore) GetPostById(id int64) (*PostResponse, error) {
	return &PostResponse{}, m.Err
}

// GetUserPosts implements Store.
func (m *MockStore) GetUserPosts(id int64) (*[]PostResponse, error) {
	return &[]PostResponse{}, m.Err
}

// UpdatePostById implements Store.
func (m *MockStore) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	return &PostResponse{}, m.Err
}

// GetUserById implements Store.
func (m *MockStore) GetUserById(id string) (*User, error) {
	return &User{}, m.Err
}

// CreateUser implements Store.
func (m *MockStore) CreateUser(user *User) (*User, error) {
	return &User{}, m.Err
}

// GetTask implements Store.
func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, m.Err
}

// CreateTask implements Store.
func (m *MockStore) CreateTask(task *Task) (*Task, error) {
	return &Task{}, m.Err
}
