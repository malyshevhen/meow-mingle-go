package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

var failError = &BasicError{
	Code:    http.StatusInternalServerError,
	Message: "Internal server error",
}

type TxProvider struct {
	db *sql.DB
}

func (s *TxProvider) Begin() (*sql.Tx, error) {
	return s.db.Begin()
}

type UserService struct {
	*TxProvider
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		TxProvider: &TxProvider{
			db: db,
		},
	}
}

func (s *UserService) GetUserById(id int64) (*UserRequest, error) {
	log.Printf("%-15s ==> 🧐 Looking for user with I %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &UserRequest{}, failError
	}
	defer tx.Rollback()

	u, err := fetchUserById(tx, id)
	if err != nil {
		return &UserRequest{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: fmt.Sprintf("User with ID: %d was not found", id),
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &UserRequest{}, failError
	}

	return u, nil
}

func (s *UserService) CreateUser(u *UserRequest) (*UserRequest, error) {
	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &UserRequest{}, failError
	}
	defer tx.Rollback()

	id, err := saveNewUser(tx, u)
	if err != nil {
		return &UserRequest{}, &BasicError{
			Code:    http.StatusBadRequest,
			Message: "User already exists",
		}
	}

	u, err = fetchUserById(tx, id)
	if err != nil {
		return &UserRequest{}, err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &UserRequest{}, failError
	}

	return u, nil
}

type PostService struct {
	TxProvider
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		TxProvider: TxProvider{
			db: db,
		},
	}
}

func (s *PostService) CreatePost(userId int64, p *PostRequest) (*PostResponse, error) {
	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}
	defer tx.Rollback()

	id, err := savePost(tx, p, userId)
	if err != nil {
		return &PostResponse{}, &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Post is not created",
		}
	}

	pr, err := fetchPostById(tx, id)
	if err != nil {
		log.Printf("%-15s ==> Post is not found: %v\n", "Store", err)
		return &PostResponse{}, failError
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}

	log.Printf("%-15s ==> 🙌 Successfully created post\n", "Store")

	return pr, nil
}

func (s *PostService) DeletePostById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Attempting to delete post with Id %d\n", "Store ", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return failError
	}
	defer tx.Rollback()

	relatedCom, err := fetchPostComments(tx, id)
	if err != nil {
		log.Printf("%-15s ==> Error getting comments by post Id: %v\n", "Store", err)
		return failError
	}

	for _, c := range relatedCom {
		log.Printf("%-15s ==> 🗑️ Attempting to delete comment with Id %d\n", "Store ", c.Id)

		if err := deleteComment(tx, c.Id); err != nil {
			log.Printf("%-15s ==> � Error deleting comment by post Id: %v\n", "Store", err)
			return failError
		}
	}

	if err := deletePost(tx, id); err != nil {
		// TODO: log
		return &BasicError{
			Code:    http.StatusNotFound,
			Message: "Post is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return failError
	}

	log.Printf("%-15s ==> 🎉 Successfully deleted post with Id %v\n", "Store", id)

	return nil
}

func (s *PostService) GetPostById(id int64) (*PostResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving post with Id %v\n", "Store ", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}
	defer tx.Rollback()

	pr, err := fetchPostById(tx, id)
	if err != nil {
		return &PostResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Post is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}

	return pr, nil
}

func (s *PostService) GetUserPosts(id int64) (*[]PostResponse, error) {
	log.Printf("%-15s ==> 📚 Retrieving posts for user with Id:%v\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &[]PostResponse{}, failError
	}
	defer tx.Rollback()

	posts, err := fetchUserPosts(tx, id)
	if err != nil {
		return &[]PostResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "User is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &[]PostResponse{}, failError
	}

	return posts, nil
}

func (s *PostService) UpdatePostById(id int64, p *PostRequest) (*PostResponse, error) {
	log.Printf("%-15s ==> 📝 Updating post with Id %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}
	defer tx.Rollback()

	if err := updatePost(tx, p, id); err != nil {
		return &PostResponse{}, &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Post is not updated",
		}
	}

	pr, err := fetchPostById(tx, id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting updated post with Id %d:%v\n", "Store", id, err)
		return &PostResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Post is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &PostResponse{}, failError
	}

	log.Printf("%-15s ==> 🙌 Successfully updated and retrieved post with Id:%d\n", "Store", id)

	return pr, nil
}

type CommentService struct {
	TxProvider
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{
		TxProvider: TxProvider{
			db: db,
		},
	}
}

func (s *CommentService) CreateComment(postId int64, userId int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Creating new comment for post Id %d, by user Id %d\n", "Store", postId, userId)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> Transaction opening is fail: %v\n", "Store", err)
		return &CommentResponse{}, failError
	}
	defer tx.Rollback()

	id, err := saveComment(tx, c, userId, postId)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting Id for new comment for post Id %d, by user Id %d %v\n", "Store", postId, userId, err)
		return &CommentResponse{}, &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Comment is not created",
		}
	}

	log.Printf("%-15s ==> 📚 Retrieving new comment with Id %d,for post Id %d by user Id %d\n", "Store", id, postId, userId)
	cr, err := fetchCommentById(tx, id)
	if err != nil {
		log.Printf("%-15s ==> 😞 Error getting new comment with Id %d for post Id %d, by user Id %d: %v\n", "Store", id, postId, userId, err)
		return &CommentResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Comment is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> Transaction committing is fail: %v\n", "Store", err)
		return &CommentResponse{}, failError
	}

	log.Printf("%-15s ==> 🙌 Successfully created and retrieved new comment with Id %d for post Id %d, by user Id %d\n", "Store", id, postId, userId)

	return cr, nil
}

func (s *CommentService) DeleteCommentById(id int64) error {
	log.Printf("%-15s ==> 🗑️ Deleting comment with Id: %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return failError
	}
	defer tx.Rollback()

	if err := deleteComment(tx, id); err != nil {
		return &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Comment is not deleted",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return failError
	}

	log.Printf("%-15s ==> 🙌 Successfully deleted comment with Id: %d\n", "Store", id)

	return nil
}

func (s *CommentService) GetCommentById(id int64) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comment with Id: %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return &CommentResponse{}, failError
	}
	defer tx.Rollback()

	cr, err := fetchCommentById(tx, id)
	if err != nil {
		return &CommentResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Comment is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return &CommentResponse{}, failError
	}

	return cr, nil
}

func (s *CommentService) GetCommentsByPostId(id int64) ([]*CommentResponse, error) {
	log.Printf("%-15s ==> 📖 Retrieving comments for post with Id: %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return nil, failError
	}
	defer tx.Rollback()

	cs, err := fetchPostComments(tx, id)
	if err != nil || cs == nil {
		return []*CommentResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Post is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return nil, failError
	}

	return cs, nil
}

func (s *CommentService) UpdateCommentById(id int64, c *CommentRequest) (*CommentResponse, error) {
	log.Printf("%-15s ==> 📝 Updating comment with Id: %d\n", "Store", id)

	tx, err := s.Begin()
	if err != nil {
		log.Printf("%-15s ==> ☹️ Transaction is not open: %v\n", "Store", err)
		return &CommentResponse{}, failError
	}
	defer tx.Rollback()

	if err := updateComment(tx, c, id); err != nil {
		return &CommentResponse{}, &BasicError{
			Code:    http.StatusBadRequest,
			Message: "Comment is not updated",
		}
	}

	cr, err := fetchCommentById(tx, id)
	if err != nil {
		return &CommentResponse{}, &BasicError{
			Code:    http.StatusNotFound,
			Message: "Comment is not found",
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("%-15s ==> ☹️ Transaction commit is fail: %v\n", "Store", err)
		return nil, failError
	}

	return cr, nil
}
