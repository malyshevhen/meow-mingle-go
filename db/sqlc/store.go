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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *Store) CreateUserTx(ctx context.Context, userParams CreateUserParams) (User, error) {
	var (
		user User = User{}
		err  error
		fail = func(message string) error { return errors.NewValidationError(message) }
	)

	log.Printf("%-15s ==> 📝 Creating user in database...\n", "UserService")

	err = s.execTx(ctx, func(q *Queries) error {
		i, err := s.IsUserExists(ctx, userParams.Email)
		if err != nil {
			return err
		}
		if i > 0 {
			return fail(fmt.Sprintf("user with email: %s already exists", userParams.Email))
		}

		if user, err = s.CreateUser(ctx, userParams); err != nil {
			return fail(fmt.Sprintf("error create user with email: %s", userParams.Email))
		}

		return nil
	})

	return user, err
}
