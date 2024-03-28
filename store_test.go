package main

// Mocks

type MockStore struct{}

// CreatePost implements Store.
func (m *MockStore) CreatePost(userId int64, p *PostRequest) (*PostResponse, error) {
	panic("unimplemented")
}

// DeletePostById implements Store.
func (m *MockStore) DeletePostById(id int64) error {
	panic("unimplemented")
}

// GetPostById implements Store.
func (m *MockStore) GetPostById(id int64) (*PostResponse, error) {
	panic("unimplemented")
}

// GetUserPosts implements Store.
func (m *MockStore) GetUserPosts(id int64) (*[]PostResponse, error) {
	panic("unimplemented")
}

// UpdatePostById implements Store.
func (m *MockStore) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	panic("unimplemented")
}

// GetUserById implements Store.
func (m *MockStore) GetUserById(id string) (*User, error) {
	return &User{}, nil
}

// CreateUser implements Store.
func (m *MockStore) CreateUser(user *User) (*User, error) {
	return &User{}, nil
}

// GetTask implements Store.
func (m *MockStore) GetTask(id string) (*Task, error) {
	return &Task{}, nil
}

// CreateTask implements Store.
func (m *MockStore) CreateTask(task *Task) (*Task, error) {
	return &Task{}, nil
}
