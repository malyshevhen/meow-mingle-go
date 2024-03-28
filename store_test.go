package main

// Mocks

type MockStore struct{}

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
