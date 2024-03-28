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
	CreatePost(userId int64, p *PostRequest) (*PostResponse, error)
	GetUserPosts(id int64) (*[]PostResponse, error)
	GetPostById(id int64) (*PostResponse, error)
	UpdatePostById(id int64, p *PostRequest) (*PostResponse, error)
	DeletePostById(id int64) error
}

type Storage struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// GetUserId implements Store
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

// CreateUser implements Store
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

// CreateTask implements Store
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

// GetTask implements Store
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
func (s *Storage) CreatePost(userId int64, p *PostRequest) (*PostResponse, error) {
	rows, err := s.db.Exec(`
	INSERT INTO posts (content, author_id)
	VALUES (?, ?)`,
		p.Content,
		userId,
	)
	if err != nil {
		return &PostResponse{}, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		return &PostResponse{}, err
	}

	pr, err := s.GetPostById(id)
	if err != nil {
		return &PostResponse{}, err
	}

	return pr, nil
}

// DeletePostById implements Store.
func (s *Storage) DeletePostById(id int64) error {
	if _, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id); err != nil {
		return err
	}

	return nil
}

// GetPostById implements Store.
func (s *Storage) GetPostById(id int64) (*PostResponse, error) {
	var pr PostResponse
	err := s.db.QueryRow(`
	SELECT id, content, author_id, created_at, updated_at
	FROM posts WHERE id = ?`, id).Scan(
		&pr.Id,
		&pr.Content,
		&pr.AuthorId,
		&pr.Created,
		&pr.Updated,
	)
	return &pr, err
}

// GetUserPosts implements Store.
func (s *Storage) GetUserPosts(id int64) (*[]PostResponse, error) {
	var (
		record   = PostResponse{}
		pubsResp = []PostResponse{}
	)
	rows, err := s.db.Query(`
	SELECT id, content, author_id, created_at, updated_at
	FROM posts WHERE author_id = ?;`, id)

	for rows.Next() {
		rows.Scan(
			&record.Id,
			&record.Content,
			&record.AuthorId,
			&record.Created,
			&record.Updated,
		)

		pubsResp = append(pubsResp, record)
	}

	return &pubsResp, err
}

// UpdatePostById implements Store.
func (s *Storage) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	_, err := s.db.Exec(`
	UPDATE posts p SET p.content = ?
	WHERE p.id = ?;`,
		p.Content,
		id,
	)
	if err != nil {
		return &PostResponse{}, err
	}

	pr, err := s.GetPostById(id)
	if err != nil {
		return &PostResponse{}, err
	}

	return pr, nil
}
