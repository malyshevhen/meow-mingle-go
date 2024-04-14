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

func (s *Store) CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error) {
	log.Printf("%-15s ==> 📝 Creating user in database...\n", "UserService")

	err = s.execTx(ctx, func(q *Queries) error {
		count, err := s.IsUserExists(ctx, params.Email)
		if err != nil {
			return errors.NewDatabaseError(err)
		} else if count > 0 {
			message := fmt.Sprintf("user with email: %s already exists", params.Email)
			return errors.NewValidationError(message)
		}

		if user, err = s.CreateUser(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}

		return nil
	})
	return
}
