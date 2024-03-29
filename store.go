package main

import (
	"database/sql"
	"log"
)

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
	CreateComment(postId int64, userId int64, c *CommentRequest) (*CommentResponse, error)
	GetCommentById(id int64) (*CommentResponse, error)
	GetCommentsByPostId(id int64) ([]*CommentResponse, error)
	UpdateCommentById(id int64, c *CommentRequest) (*CommentResponse, error)
	DeleteCommentById(id int64) error
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// GetUserId implements Store
func (s *Storage) GetUserById(id string) (*User, error) {
	log.Printf("%-15s ==> 🧐 Looking for user with I %s\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	var u User
	err = tx.QueryRow(`
	SELECT id, email, first_name, last_name, password, created_at
	FROM users WHERE id = ?`, id).Scan(
		&u.ID,
		&u.Email,
		&u.FirstName,
		&u.LastName,
		&u.Password,
		&u.CreatedAt,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Failed to find user with I %s\n", "Store", id)
	} else {
		log.Printf("%-15s ==> 🎉 Found user with I %s\n", "Store", id)
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return &u, err
}

// CreateUser implements Store
func (s *Storage) CreateUser(u *User) (*User, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	rows, err := tx.Exec(`
	INSERT INTO users (email, first_name, last_name, password)
	VALUES (?, ?, ?, ?)`,
		u.Email,
		u.FirstName,
		u.LastName,
		u.Password,
	)
	if err != nil {
		log.Printf("%-15s ==> Error inserting user: %v\n", "Store", err)
		return &User{}, err
	}

	log.Printf("%-15s ==> 🎉 Successfully inserted user\n", "Store")

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> Error getting last insert ID: %v\n", "Store", err)
		return &User{}, err
	}

	log.Printf("%-15s ==> 🆔 Got user ID %v\n", "Store", id)

	u.ID = id

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return u, nil
}

// CreateTask implements Store
func (s *Storage) CreateTask(t *Task) (*Task, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	rows, err := tx.Exec(`
	INSERT INTO tasks (name, status, project_id, assigned_to)
	VALUES (?, ?, ?, ?)`,
		t.Name,
		t.Status,
		t.ProjectID,
		t.AssignedToID,
	)
	if err != nil {
		log.Printf("%-15s ==> Error inserting task: %v\n", "Store", err)
		return nil, err
	}

	log.Printf("%-15s ==> 🎉 Successfully inserted task\n", "Store")

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> Error getting last insert ID: %v\n", "Store", err)
		return nil, err
	}

	log.Printf("%-15s ==> 🆔 Got task ID %v\n", "Store", id)

	t.ID = id

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return t, nil
}

// GetTask implements Store
func (s *Storage) GetTask(id string) (*Task, error) {
	var t Task

	log.Printf("%-15s ==> 🕵️ Retrieving for task with ID %vsn", "Store", id)

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

	if err != nil {
		log.Printf("%-15s ==> Error querying for task: %v\n", "Store", err)
		return nil, err
	}

	log.Printf("%-15s ==> 🎉 Successfully queried for task\n", "Store")

	return &t, nil
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
		log.Printf("%-15s ==> Error inserting post: %v\n", "Store", err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 🎉 Successfully inserted post\n", "Store!")

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> Error getting last insert ID: %v\n", "Store", err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 🆔 Got post ID %vdn", "Store ", id)

	pr, err := s.GetPostById(id)
	if err != nil {
		log.Printf("%-15s ==> Error getting post by ID: %v\n", "Store", err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully created post\n", "Store")

	return pr, nil
}

// DeletePostById implements Store.
func (s *Storage) DeletePostById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Attempting to delete post with ID %d\n", "Store ", id)

	relatedCom, err := s.GetCommentsByPostId(id)
	if err != nil {
		log.Printf("%-15s ==> Error getting comments by post ID: %v\n", "Store", err)
		return err
	}

	for _, c := range relatedCom {
		log.Printf("%-15s ==> 🗑️ Attempting to delete comment with ID %d\n", "Store ", c.Id)
		if err := s.DeleteCommentById(c.Id); err != nil {
			log.Printf("%-15s ==> � Error deleting comment by post ID: %v\n", "Store", err)
			return err
		}
	}

	if _, err := s.db.Exec("DELETE FROM posts WHERE id = ?", id); err != nil {
		log.Printf("%-15s ==> Error deleting post: %v\n", "Store", err)
		return err
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted post with ID %v\n", "Store", id)

	return nil
}

// GetPostById implements Store.
func (s *Storage) GetPostById(id int64) (*PostResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving post with ID %v\n", "Store ", id)

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

	if err != nil {
		log.Printf("%-15s ==> Error getting post by ID: %v\n", "Store", err)
		return nil, err
	}

	log.Printf("%-15s ==> 🙌 Successfully retrieved post with ID %v\n", "Store", id)

	return &pr, nil
}

// GetUserPosts implements Store.
func (s *Storage) GetUserPosts(id int64) (*[]PostResponse, error) {
	log.Printf("%-15s ==> 📚 Retrieving posts for user with ID:%v\n", "Store", id)

	var (
		record   = PostResponse{}
		pubsResp = []PostResponse{}
	)
	rows, err := s.db.Query(`
	SELECT id, content, author_id, created_at, updated_at
	FROM posts WHERE author_id = ?`, id)

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

	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting posts for user with ID %d: %v\n", "Store", id, err)
		return nil, err
	}

	log.Printf("%-15s ==> 🙌 Successfully retrieved posts for user with ID: %d\n", "Store", id)

	return &pubsResp, nil
}

// UpdatePostById implements Store.
func (s *Storage) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	log.Printf("%-15s ==> 📝 Updating post with ID %d\n", "Store", id)

	_, err := s.db.Exec(`
	UPDATE posts p SET p.content = ? 
	WHERE p.id = ?`,
		p.Content,
		id,
	)

	if err != nil {
		log.Printf("%-15s ==> 😞 Error updating post with ID %d: %v\n", "Store", id, err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 📚 Retrieving updated post with ID %d\n", "Store", id)
	pr, err := s.GetPostById(id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting updated post with ID %d:%v\n", "Store", id, err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully updated and retrieved post with ID:%d\n", "Store", id)

	return pr, nil
}

// CreateComment implements Store.
func (s *Storage) CreateComment(postId int64, userId int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Creating new comment for post ID %d, by user ID %d\n", "Store", postId, userId)

	rows, err := s.db.Exec(`
	INSERT INTO comments (content, author_id, post_id)
	VALUES (?, ?, ?)`,
		c.Content,
		userId,
		postId,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error creating comment for post ID %d by user ID %d %v\n", "Store", postId, userId, err)
		return &CommentResponse{}, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting ID for new comment for post ID %d, by user ID %d %v\n", "Store", postId, userId, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 📚 Retrieving new comment with ID %d,for post ID %d by user ID %d\n", "Store", id, postId, userId)
	cr, err := s.GetCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting new comment with ID %d for post ID %d, by user ID %d: %v\n", "Store", id, postId, userId, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully created and retrieved new comment with ID %d for post ID %d, by user ID %d\n", "Store", id, postId, userId)

	return cr, nil
}

// DeleteCommentById implements Store.
func (s *Storage) DeleteCommentById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Deleting comment with ID: %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM comments WHERE id = ?", id); err != nil {
		log.Printf("%-15s ==> 😞 Error deleting comment with ID: %d %v\n", "Store", id, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return err
	}

	log.Printf("%-15s ==> 🙌 Successfully deleted comment with ID: %d\n", "Store", id)

	return nil
}

// GetCommentById implements Store.
func (s *Storage) GetCommentById(id int64) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comment with ID: %d\n", "Store", id)

	var cr CommentResponse
	err := s.db.QueryRow(`
	SELECT id, content, author_id, post_id, created_at, updated_at 
	FROM comments WHERE id = ?`, id).Scan(
		&cr.Id,
		&cr.Content,
		&cr.AuthorId,
		&cr.PostId,
		&cr.Created,
		&cr.Updated,
	)

	if err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving comment with ID: %d %v\n", "Store", id, err)
	} else {
		log.Printf("%-15s ==> 🙌 Successfully retrieved comment with ID: %d\n", "Store", id)
	}

	return &cr, err
}

func (s *Storage) GetCommentsByPostId(id int64) ([]*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comments for post with ID: %d\n", "Store", id)
	var cr []*CommentResponse
	rows, err := s.db.Query(`
	SELECT id, content, author_id, post_id, created_at, updated_at
	FROM comments WHERE post_id = ?`, id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving comments for post with ID: %d %v\n", "Store", id, err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		record := &CommentResponse{}
		if err := rows.Scan(
			&record.Id,
			&record.Content,
			&record.AuthorId,
			&record.PostId,
			&record.Created,
			&record.Updated,
		); err != nil {
			log.Printf("%-15s ==> 😞 Error retrieving comments for post with ID: %d %v\n", "Store", id, err)
		} else {
			cr = append(cr, record)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving comments for post with ID: %d %v\n", "Store", id, err)
	} else {
		log.Printf("%-15s ==> 🙌 Successfully retrieved comments for post with ID: %d\n", "Store", id)
	}
	return cr, nil
}

// UpdateCommentById implements Store.
func (s *Storage) UpdateCommentById(id int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Updating comment with ID: %d\n", "Store", id)

	_, err := s.db.Exec(`
	UPDATE comments c SET c.content = ? 
	WHERE c.id = ?`,
		c.Content,
		id,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error updating comment with ID: %d %v\n", "Store", id, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 🔎 Retrieving updated comment with ID: %d\n", "Store", id)
	cr, err := s.GetCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving updated comment with ID: %d %v\n", "Store", id, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully updated and retrieved comment with ID: %d\n", "Store", id)

	return cr, nil
}
