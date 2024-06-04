package db

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

//go:embed cypher/create_user.cypher
var createUserCypher string

type VStore struct {
	driver neo4j.DriverWithContext
}

func (s *VStore) CreateUserTx(ctx context.Context, userForm CreateUserParams) (user User, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, err := s.persistUser(ctx, tx, userForm)
		if err != nil {
			return nil, err
		}
		user = record.(User)

		return record, err
	}); err != nil {
		execErr = err
		return
	}

	return
}

func (s *VStore) persistUser(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	userForm CreateUserParams,
) (user any, err error) {
	params := map[string]interface{}{
		"email":    userForm.Email,
		"password": userForm.Password,
	}

	fail := func(msg string) (any, error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occured while %s: %s", msg, userForm.Email))
		return nil, err
	}
	result, err := tx.Run(ctx, createUserCypher, params)
	if err != nil {
		return fail("executing transaction")
	}

	record, err := result.Single(ctx)
	if err != nil {
		return fail("saving new user wiht email")
	}

	id, ok := record.Values[0].(int64)
	if !ok {
		return fail("extracting the ID")
	}

	user = User{
		ID:    id,
		Email: userForm.Email,
	}

	return user, nil
}

func (s *VStore) CreatePostTx(ctx context.Context, params CreatePostParams) (post Post, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) CreateCommentTx(ctx context.Context, params CreateCommentParams) (comment Comment, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) CreatePostLikeTx(ctx context.Context, params CreatePostLikeParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) CreateCommentLikeTx(ctx context.Context, params CreateCommentLikeParams) (err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) CreateSubscriptionTx(ctx context.Context, params CreateSubscriptionParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) GetUserTx(ctx context.Context, id int64) (user GetUserRow, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) GetPostTx(ctx context.Context, id int64) (post PostInfo, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) GetFeed(ctx context.Context, userId int64) (feed []PostInfo, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) ListUserPostsTx(ctx context.Context, userId int64) (posts []PostInfo, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) ListPostCommentsTx(ctx context.Context, id int64) (posts []CommentInfo, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) UpdatePostTx(ctx context.Context, userId int64, params UpdatePostParams) (post Post, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) UpdateCommentTx(ctx context.Context, userId int64, params UpdateCommentParams) (comment Comment, err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) DeletePostTx(ctx context.Context, userId int64, postId int64) error {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) DeletePostLikeTx(ctx context.Context, params DeletePostLikeParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) DeleteCommentTx(ctx context.Context, userId int64, commentId int64) (err error) {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) DeleteCommentLikeTx(ctx context.Context, params DeleteCommentLikeParams) error {
	panic("not implemented") // TODO: Implement
}

func (s *VStore) DeleteSubscriptionTx(ctx context.Context, params DeleteSubscriptionParams) error {
	panic("not implemented") // TODO: Implement
}
