package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sort"

	"github.com/malyshEvhen/meow_mingle/errors"
)

type IStore interface {
	CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error)
	CreatePostTx(ctx context.Context, authorId int64, content string) (post Post, err error)
	CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error)
	CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error
	CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error)
	CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error
	GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error)
	GetPostTx(ctx context.Context, id int64) (post GetPostRow, err error)
	GetFeed(ctx context.Context, userId int64) (feed []ListUserPostsRow, err error)
	ListUserPostsTx(ctx context.Context, userId int64) (posts []ListUserPostsRow, err error)
	ListPostCommentsTx(ctx context.Context, id int64) (posts []ListPostCommentsRow, err error)
	UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error)
	UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error)
	DeletePostTx(ctx context.Context, userId, postId int64) error
	DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error
	DeleteCommentTx(ctx context.Context, userId, commentId int64) (err error)
	DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error
	DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
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

func (s *SQLStore) CreateUserTx(ctx context.Context, params CreateUserParams) (user User, err error) {
	log.Printf("%-15s ==> Creating user in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		count, err := s.IsUserExists(ctx, params.Email)
		if err != nil {
			return errors.NewDatabaseError(err)
		}
		if count != 0 {
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

func (s *SQLStore) GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error) {
	log.Printf("%-15s ==> Retrieving User from database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if user, err = s.GetUser(ctx, id); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error {
	log.Printf("%-15s ==> Subscribe User with ID: %d to User with ID: %d from database...\n",
		"Store",
		params.UserID,
		params.SubscriptionID,
	)

	err := s.execTx(ctx, func(q *Queries) error {
		if err := s.CreateSubscription(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error {
	log.Printf("%-15s ==> Unsubscribe User with ID: %d from User with ID: %d from database...\n",
		"Store",
		params.UserID,
		params.SubscriptionID,
	)

	err := s.execTx(ctx, func(q *Queries) error {
		if _, err := s.GetSubscription(ctx, GetSubscriptionParams(params)); err != nil {
			msg := fmt.Sprintf("Subscription of User %d to User %d was not found",
				params.UserID,
				params.SubscriptionID,
			)
			return errors.NewNotFoundError(msg)
		}

		if err := s.DeleteSubscription(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) GetFeed(ctx context.Context, userId int64) (feed []ListUserPostsRow, err error) {
	log.Printf("%-15s ==> Retrieve feed of User with ID: %d from database...\n", "Store", userId)

	err = s.execTx(ctx, func(q *Queries) error {
		if _, err := s.GetUser(ctx, userId); err != nil {
			msg := fmt.Sprintf("User with ID: %d not found", userId)
			return errors.NewNotFoundError(msg)
		}

		subs, err := s.ListSubscriptions(ctx, userId)
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		for _, sub := range subs {
			posts, err := s.ListUserPosts(ctx, sub)
			if err != nil {
				return errors.NewDatabaseError(err)
			}

			feed = append(feed, posts...)
		}
		return nil
	})

	sort.Slice(feed, func(i, j int) bool {
		return feed[i].CreatedAt.Before(feed[j].CreatedAt)
	})

	return
}

func (s *SQLStore) CreatePostTx(ctx context.Context, authorId int64, content string) (post Post, err error) {
	log.Printf("%-15s ==> Creating post in database...\n", "Store")

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

func (s *SQLStore) ListUserPostsTx(ctx context.Context, userId int64) (posts []ListUserPostsRow, err error) {
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

func (s *SQLStore) GetPostTx(ctx context.Context, id int64) (post GetPostRow, err error) {
	log.Printf("%-15s ==> Retrieving post from database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if post, err = s.GetPost(ctx, id); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error) {
	log.Printf("%-15s ==> Updating post in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		err, ok := s.isPostsAuthor(ctx, params.ID, userId)
		if !ok {
			return errors.NewForbiddenError()
		}
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		if post, err = s.UpdatePost(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) DeletePostTx(ctx context.Context, userId, postId int64) error {
	log.Printf("%-15s ==> Deleting post from database...\n", "Store")

	err := s.execTx(ctx, func(q *Queries) error {
		err, ok := s.isPostsAuthor(ctx, postId, userId)
		if !ok {
			return errors.NewForbiddenError()
		}
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		if err := s.DeletePost(ctx, postId); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error {
	log.Printf("%-15s ==> Create comment like in database...\n", "Store")

	err := s.execTx(ctx, func(q *Queries) error {
		if err := s.CreatePostLike(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error {
	log.Printf("%-15s ==> Delete comment like from database...\n", "Store")

	err := s.execTx(ctx, func(q *Queries) error {
		_, err := s.GetPostLike(ctx, GetPostLikeParams(params))
		if err != nil {
			log.Printf("%-15s ==> Error: %s\n", "Store", err.Error())

			msg := fmt.Sprintf("Like from user with ID: %d on post with ID: %d is not found",
				params.UserID,
				params.PostID,
			)
			return errors.NewNotFoundError(msg)
		}

		if err := s.DeletePostLike(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error) {
	log.Printf("%-15s ==> Create comment in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if comment, err = s.CreateComment(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) ListPostCommentsTx(ctx context.Context, id int64) (posts []ListPostCommentsRow, err error) {
	log.Printf("%-15s ==> Retrieving post comments from database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if posts, err = s.ListPostComments(ctx, id); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error) {
	log.Printf("%-15s ==> Updating comment in database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		err, ok := s.isCommentsAuthor(ctx, params.ID, userId)
		if !ok {
			return errors.NewForbiddenError()
		}
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		if comment, err = s.UpdateComment(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error) {
	log.Printf("%-15s ==> Add comment like to database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		if err = s.CreateCommentLike(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return
}

func (s *SQLStore) DeleteCommentTx(ctx context.Context, userId, commentId int64) (err error) {
	log.Printf("%-15s ==> Retrieve Comments author ID from database...\n", "Store")

	err = s.execTx(ctx, func(q *Queries) error {
		err, ok := s.isCommentsAuthor(ctx, commentId, userId)
		if !ok {
			return errors.NewForbiddenError()
		}
		if err != nil {
			return errors.NewDatabaseError(err)
		}

		if err = s.DeleteComment(ctx, commentId); err != nil {
			return errors.NewDatabaseError(err)
		}

		return nil
	})

	return
}

func (s *SQLStore) DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error {
	log.Printf("%-15s ==> Delete comments like from database...\n", "Store")

	err := s.execTx(ctx, func(q *Queries) error {
		_, err := s.GetCommentLike(ctx, GetCommentLikeParams(params))
		if err != nil {
			log.Printf("%-15s ==> Error: %s\n", "Store", err.Error())

			msg := fmt.Sprintf("Like from user with ID: %d on comment with ID: %d is not found",
				params.UserID,
				params.CommentID,
			)
			return errors.NewNotFoundError(msg)
		}

		if err := s.DeleteCommentLike(ctx, params); err != nil {
			return errors.NewDatabaseError(err)
		}
		return nil
	})
	return err
}

func (s *SQLStore) isCommentsAuthor(ctx context.Context, commentId int64, userId int64) (error, bool) {
	authorId, err := s.GetCommentsAuthorID(ctx, commentId)
	if err != nil {
		return errors.NewDatabaseError(err), false
	}
	if authorId != userId {
		return nil, false
	}
	return nil, true
}

func (s *SQLStore) isPostsAuthor(ctx context.Context, postId int64, userId int64) (error, bool) {
	authorId, err := s.GetPostsAuthorID(ctx, postId)
	if err != nil {
		return errors.NewDatabaseError(err), false
	}
	if authorId != userId {
		return nil, false
	}
	return nil, true
}
