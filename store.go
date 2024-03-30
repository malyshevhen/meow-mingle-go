package main

import (
	"database/sql"
	"log"
)

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

func saveComment(tx *sql.Tx, c *CommentRequest, userId, postId int64) (int64, error) {
	rows, err := tx.Exec(`
	INSERT INTO comments (content, author_id, post_id)
	VALUES (?, ?, ?)`,
		c.Content,
		userId,
		postId,
	)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error creating comment for post Id %d by user Id %d %v\n", "Store", postId, userId, err)
		return 0, err
	}

	id, err := rows.LastInsertId()
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting Id for new comment for post Id %d, by user Id %d %v\n", "Store", postId, userId, err)
		return 0, err
	}

	return id, nil
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
