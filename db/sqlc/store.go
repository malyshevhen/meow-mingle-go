package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/malyshEvhen/meow_mingle/errors"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Queries: New(db),
		db:      db,
	}
}

func (s *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	fail := func(err error) error { return errors.NewDatabaseError(err) }

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fail(err)
	}

	query := New(tx)

	if err := fn(query); err != nil {
		if rErr := tx.Rollback(); rErr != nil {
			return fmt.Errorf("%v %v", err, rErr)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fail(err)
	}

	return nil
}

func (s *Store) CreateUserTx(ctx context.Context, ur CreateUserParams) (user User, err error) {
	log.Printf("%-15s ==> 📝 Creating user in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		count, err := s.IsUserExists(ctx, ur.Email)
		if err != nil {
			return errors.NewDatabaseError(err)
		} else if count > 0 {
			message := fmt.Sprintf("user with email: %s already exists", ur.Email)
			return errors.NewValidationError(message)
		}

		dbUser, err := s.CreateUser(ctx, CreateUserParams{
			Email:     ur.Email,
			FirstName: ur.FirstName,
			LastName:  ur.LastName,
			Password:  ur.Password,
		})
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		user.ID = dbUser.ID
		user.Email = dbUser.Email
		user.FirstName = dbUser.FirstName
		user.LastName = dbUser.LastName

		return nil
	})
	return
}

func (s *Store) CreatePostTx(ctx context.Context, authorId int64, content string) (post Post, err error) {
	log.Printf("%-15s ==> 📝 Creating post in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if post, err = s.CreatePost(ctx, CreatePostParams{
			Content:  content,
			AuthorID: authorId,
		}); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *Store) ListUserPostsTx(ctx context.Context, userId int64) (posts []ListUserPostsRow, err error) {
	log.Printf("%-15s ==> Retrieving users post from database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		posts, err = s.ListUserPosts(ctx, userId)
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		return nil
	})
	return
}
