package main

import "database/sql"

type Store interface {
	// Users
	CreateUser(user *User) (*User, error)
	GetUserById(id string) (*User, error)
	// Tasks
	CreateTask(task *Task) (*Task, error)
	GetTask(id string) (*Task, error)
	// Posts
	CreatePost(p *PostRequest) (*PostResponse, error)
	GetUserPosts(id string) (*Page[PostResponse], error)
	GetPostsById(id string) (*PostResponse, error)
	UpdatePostsById(id string, p *PostRequest) (*PostResponse, error)
	DeletePostsById(id string) error
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetUserById(id string) (*User, error) {
	var u User
	err := s.db.QueryRow(`
	SELECT id, email, first_name, last_name, password, created_at
	FROM users WHERE id = ?`, id).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Password,
		&u.CreatedAt,
	)
	return &u, err
}

func (s *Storage) CreateUser(u *User) (*User, error) {
	rows, err := s.db.Exec(`
	INSERT INTO users (email, first_name, last_name, password)
	VALUES (?, ?, ?, ?)`,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
	)
	if err != nil {
		return &User{}, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return &User{}, err
	}

	u.ID = id

	return u, nil
}

func (s *Storage) CreateTask(t *Task) (*Task, error) {
	rows, err := s.db.Exec(`
	INSERT INTO tasks (name, status, project_id, assigned_to)
	VALUES (?, ?, ?, ?)`,
		t.Name,
		t.Status,
		t.ProjectID,
		t.AssignedToID,
	)
	if err != nil {
		return nil, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return nil, err
	}

	t.ID = id

	return t, nil
}

func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task
	err := s.db.QueryRow(`
	SELECT id, name, status, project_id, assigned_to, created_at
	FROM tasks WHERE id = ?`, id).Scan(
		&t.ID,
		&t.Name,
		&t.Status,
		&t.ProjectID,
		&t.AssignedToID,
		&t.CreatedAt,
	)
	return &t, err
}

// CreatePost implements Store.
func (s *Storage) CreatePost(p *PostRequest) (*PostResponse, error) {
	panic("unimplemented")
}

// DeletePostsById implements Store.
func (s *Storage) DeletePostsById(id string) error {
	panic("unimplemented")
}

// GetPostsById implements Store.
func (s *Storage) GetPostsById(id string) (*PostResponse, error) {
	panic("unimplemented")
}

// GetUserPosts implements Store.
func (s *Storage) GetUserPosts(id string) (*Page[PostResponse], error) {
	panic("unimplemented")
}

// UpdatePostsById implements Store.
func (s *Storage) UpdatePostsById(id string, p *PostRequest) (*PostResponse, error) {
	panic("unimplemented")
}
