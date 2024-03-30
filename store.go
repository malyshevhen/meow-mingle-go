package main

import (
	"database/sql"
	"log"
)

type Store interface {
	// Users
	CreateUser(user *UserRequest) (*UserRequest, error)
	GetUserById(id int64) (*UserRequest, error)
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

// GetUserById implements Store
func (s *Storage) GetUserById(id int64) (*UserRequest, error) {
	log.Printf("%-15s ==> 🧐 Looking for user with I %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	u, err := fetchUserById(tx, id)

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return u, err
}

// CreateUser implements Store
func (s *Storage) CreateUser(u *UserRequest) (*UserRequest, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	id, err := saveNewUser(tx, u)
	if err != nil {
		return &UserRequest{}, err
	}

	u, err = fetchUserById(tx, id)
	if err != nil {
		return &UserRequest{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return u, nil
}

// CreatePost implements Store.
func (s *Storage) CreatePost(userId int64, p *PostRequest) (*PostResponse, error) {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	id, err := savePost(tx, p, userId)
	if err != nil {
		return &PostResponse{}, err
	}

	pr, err := fetchPostById(tx, id)
	if err != nil {
		return &PostResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	log.Printf("%-15s ==> 🙌 Successfully created post\n", "Store")

	return pr, nil
}

// DeletePostById implements Store.
//
// TODO: replace Methods to helper functions
func (s *Storage) DeletePostById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Attempting to delete post with Id %d\n", "Store ", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	relatedCom, err := fetchPostComments(tx, id)
	if err != nil {
		log.Printf("%-15s ==> Error getting comments by post Id: %v\n", "Store", err)
		return err
	}

	for _, c := range relatedCom {
		log.Printf("%-15s ==> 🗑️ Attempting to delete comment with Id %d\n", "Store ", c.Id)

		if err := deleteComment(tx, c.Id); err != nil {
			log.Printf("%-15s ==> � Error deleting comment by post Id: %v\n", "Store", err)
			return err
		}
	}

	if err := deletePost(tx, id); err != nil {
		// TODO: log
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted post with Id %v\n", "Store", id)

	return nil
}

// GetPostById implements Store.
func (s *Storage) GetPostById(id int64) (*PostResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving post with Id %v\n", "Store ", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	pr, err := fetchPostById(tx, id)
	if err != nil {
		return &PostResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return pr, nil
}

// GetUserPosts implements Store.
func (s *Storage) GetUserPosts(id int64) (*[]PostResponse, error) {
	log.Printf("%-15s ==> 📚 Retrieving posts for user with Id:%v\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	posts, err := fetchUserPosts(tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	return posts, nil
}

// UpdatePostById implements Store.
func (s *Storage) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	log.Printf("%-15s ==> 📝 Updating post with Id %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
	}
	defer tx.Rollback()

	if err := updatePost(tx, p, id); err != nil {
		return &PostResponse{}, err
	}

	pr, err := fetchPostById(tx, id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting updated post with Id %d:%v\n", "Store", id, err)
		return &PostResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
	}

	log.Printf("%-15s ==> 🙌 Successfully updated and retrieved post with Id:%d\n", "Store", id)

	return pr, nil
}

// CreateComment implements Store.
func (s *Storage) CreateComment(postId int64, userId int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Creating new comment for post Id %d, by user Id %d\n", "Store", postId, userId)

	rows, err := s.db.Exec(`
	INSERT INTO comments (content, author_id, post_id)
	VALUES (?, ?, ?)`,
		c.Content,
		userId,
		postId,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error creating comment for post Id %d by user Id %d %v\n", "Store", postId, userId, err)
		return &CommentResponse{}, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting Id for new comment for post Id %d, by user Id %d %v\n", "Store", postId, userId, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 📚 Retrieving new comment with Id %d,for post Id %d by user Id %d\n", "Store", id, postId, userId)
	cr, err := s.GetCommentById(id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting new comment with Id %d for post Id %d, by user Id %d: %v\n", "Store", id, postId, userId, err)
		return &CommentResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully created and retrieved new comment with Id %d for post Id %d, by user Id %d\n", "Store", id, postId, userId)

	return cr, nil
}

// DeleteCommentById implements Store.
func (s *Storage) DeleteCommentById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Deleting comment with Id: %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return err
	}
	defer tx.Rollback()

	if err := deleteComment(tx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return err
	}

	log.Printf("%-15s ==> 🙌 Successfully deleted comment with Id: %d\n", "Store", id)

	return nil
}

// GetCommentById implements Store.
func (s *Storage) GetCommentById(id int64) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comment with Id: %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return &CommentResponse{}, err
	}
	defer tx.Rollback()

	cr, err := fetchCommentById(tx, id)
	if err != nil {
		return &CommentResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return &CommentResponse{}, err
	}

	return cr, nil
}

func (s *Storage) GetCommentsByPostId(id int64) ([]*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comments for post with Id: %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return nil, err
	}
	defer tx.Rollback()

	cs, err := fetchPostComments(tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return nil, err
	}

	return cs, nil
}

// UpdateCommentById implements Store.
func (s *Storage) UpdateCommentById(id int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Updating comment with Id: %d\n", "Store", id)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return nil, err
	}
	defer tx.Rollback()

	if err := updateComment(tx, c, id); err != nil {
		return &CommentResponse{}, err
	}

	cr, err := fetchCommentById(tx, id)
	if err != nil {
		return &CommentResponse{}, err
	}

	return cr, nil
}

func saveNewUser(tx *sql.Tx, u *UserRequest) (int64, error) {
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
		return 0, err
	}

	log.Printf("%-15s ==> 🎉 Successfully inserted user\n", "Store")

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> Error getting last insert Id: %v\n", "Store", err)
		return 0, err
	}

	log.Printf("%-15s ==> 🆔 Got user Id %v\n", "Store", id)

	return id, nil
}

func fetchUserById(tx *sql.Tx, id int64) (*UserRequest, error) {
	var u UserRequest
	err := tx.QueryRow(`
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
		log.Printf("%-15s ==> 😞 Failed to find user with I %d\n", "Store", id)
	} else {
		log.Printf("%-15s ==> 🎉 Found user with I %d\n", "Store", id)
	}
	return &u, err
}

func savePost(tx *sql.Tx, p *PostRequest, userId int64) (int64, error) {
	rows, err := tx.Exec(`
	INSERT INTO posts (content, author_id)
	VALUES (?, ?)`,
		p.Content,
		userId,
	)
	if err != nil {
		log.Printf("%-15s ==> Error inserting post: %v\n", "Store", err)
		return 0, err
	}

	log.Printf("%-15s ==> 🎉 Successfully inserted post\n", "Store!")

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> Error getting last insert Id: %v\n", "Store", err)
		return 0, err
	}

	log.Printf("%-15s ==> 🆔 Got post Id %vdn", "Store ", id)

	return id, nil
}

func fetchPostById(tx *sql.Tx, id int64) (*PostResponse, error) {
	var pr PostResponse
	err := tx.QueryRow(`
	SELECT id, content, author_id, created_at, updated_at
	FROM posts WHERE id = ?`, id).Scan(
		&pr.Id,
		&pr.Content,
		&pr.AuthorId,
		&pr.Created,
		&pr.Updated,
	)

	if err != nil {
		log.Printf("%-15s ==> Error getting post by Id: %v\n", "Store", err)
		return &PostResponse{}, err
	}

	log.Printf("%-15s ==> 🙌 Successfully retrieved post with Id %v\n", "Store", id)
	return &pr, nil
}

func fetchUserPosts(tx *sql.Tx, id int64) (*[]PostResponse, error) {
	var (
		record   = PostResponse{}
		pubsResp = []PostResponse{}
	)
	rows, err := tx.Query(`
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
		log.Printf("%-15s ==> 😞 Error getting posts for user with Id %d: %v\n", "Store", id, err)
		return nil, err
	}

	log.Printf("%-15s ==> 🙌 Successfully retrieved posts for user with Id: %d\n", "Store", id)

	return &pubsResp, nil
}

func updatePost(tx *sql.Tx, p *PostRequest, id int64) error {
	_, err := tx.Exec(`
	UPDATE posts p SET p.content = ?
	WHERE p.id = ?`,
		p.Content,
		id,
	)

	if err != nil {
		log.Printf("%-15s ==> 😞 Error updating post with Id %d: %v\n", "Store", id, err)
		return err
	}

	log.Printf("%-15s ==> 📚 Retrieving updated post with Id %d\n", "Store", id)

	return nil
}

func deletePost(tx *sql.Tx, id int64) error {
	if _, err := tx.Exec("DELETE FROM posts WHERE id = ?", id); err != nil {
		log.Printf("%-15s ==> Error deleting post: %v\n", "Store", err)
		return err
	}

	return nil
}

func fetchCommentById(tx *sql.Tx, id int64) (*CommentResponse, error) {
	var cr CommentResponse
	err := tx.QueryRow(`
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
		log.Printf("%-15s ==> 😞 Error retrieving comment with Id: %d %v\n", "Store", id, err)
	} else {
		log.Printf("%-15s ==> 🙌 Successfully retrieved comment with Id: %d\n", "Store", id)
	}

	return &cr, err
}

func fetchPostComments(tx *sql.Tx, id int64) ([]*CommentResponse, error) {
	var cr []*CommentResponse
	rows, err := tx.Query(`
	SELECT id, content, author_id, post_id, created_at, updated_at
	FROM comments WHERE post_id = ?`, id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving comments for post with Id: %d %v\n", "Store", id, err)
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
			log.Printf("%-15s ==> 😞 Error retrieving comments for post with Id: %d %v\n", "Store", id, err)
		} else {
			cr = append(cr, record)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("%-15s ==> 😞 Error retrieving comments for post with Id: %d %v\n", "Store", id, err)
	} else {
		log.Printf("%-15s ==> 🙌 Successfully retrieved comments for post with Id: %d\n", "Store", id)
	}
	return cr, nil
}

func updateComment(tx *sql.Tx, c *CommentRequest, id int64) error {
	_, err := tx.Exec(`
	UPDATE comments c SET c.content = ?
	WHERE c.id = ?`,
		c.Content,
		id,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error updating comment with Id: %d %v\n", "Store", id, err)
		return err
	}

	return nil
}

func deleteComment(tx *sql.Tx, id int64) error {
	if _, err := tx.Exec("DELETE FROM comments WHERE id = ?", id); err != nil {
		log.Printf("%-15s ==> 😞 Error deleting comment with Id: %d %v\n", "Store", id, err)
		return err
	}
	return nil
}
