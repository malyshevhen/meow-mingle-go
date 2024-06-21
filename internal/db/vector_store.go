package db

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
)

var (
	//go:embed cypher/create_user.cypher
	createUserCypher string

	//go:embed cypher/match_user_by_id.cypher
	getUserCypher string

	//go:embed cypher/create_post.cypher
	createPostCypher string
)

type VStore struct {
	driver neo4j.DriverWithContext
}

func NewVstore(d neo4j.DriverWithContext) IStore {
	return &VStore{
		driver: d,
	}
}

type Neo4jResponse[T any] struct {
	ID        int64    `json:"Id"`
	ElementID string   `json:"ElementId"`
	Labels    []string `json:"Labels"`
	Props     T
}

func (s *VStore) CreateUserTx(ctx context.Context, userForm CreateUserParams) (user User, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, id, err := persist[User](ctx, tx, createUserCypher, "u", userForm)
		if err != nil {
			return nil, err
		}
		user = record
		user.ID = id

		return record, err
	}); err != nil {
		execErr = err
		return
	}

	return
}

func (s *VStore) GetUserTx(ctx context.Context, id int64) (user GetUserRow, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"id": id,
	}

	if _, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, err := retrieve[GetUserRow](ctx, tx, "u", params)
		if err != nil {
			return nil, err
		}
		user = record
		user.ID = id

		return nil, err
	}); err != nil {
		execErr = err
		return
	}
	return
}

func (s *VStore) CreatePostTx(ctx context.Context, params CreatePostParams) (post Post, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		record, id, err := persist[Post](ctx, tx, createPostCypher, "p", params)
		if err != nil {
			return nil, err
		}
		post = record
		post.AuthorID = params.AuthorID
		post.ID = id

		return record, err
	}); err != nil {
		execErr = err
		return
	}

	return
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

func persist[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	alias string,
	form any,
) (obj T, id int64, err error) {
	fail := func(msg string) (res T, id int64, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	params, err := mapToProperties(form)
	if err != nil {
		return fail(fmt.Sprintf("convert to properties: %s", err.Error()))
	}

	result, err := tx.Run(ctx, cypher, params)
	if err != nil {
		return fail(fmt.Sprintf("executing transaction: %s", err.Error()))
	}

	record, err := result.Single(ctx)
	if err != nil {
		return fail(fmt.Sprintf("saving new record: %s", err.Error()))
	}

	return parceRecord[T](record, alias)
}

func retrieve[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	alias string,
	params map[string]interface{},
) (obj T, err error) {
	fail := func(msg string) (res T, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	result, err := tx.Run(ctx, getUserCypher, params)
	if err != nil {
		return fail(fmt.Sprintf("executing transaction: %s", err.Error()))
	}

	record, err := result.Single(ctx)
	if err != nil {
		return fail(fmt.Sprintf("saving new record: %s", err.Error()))
	}

	obj, _, err = parceRecord[T](record, alias)

	return
}

func parceRecord[T any](r *db.Record, alias string) (result T, id int64, err error) {
	fail := func(msg string) (res T, id int64, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}
	propsMap := r.AsMap()[alias]

	respJson, err := json.MarshalIndent(propsMap, "", "  ")
	if err != nil {
		return fail(fmt.Sprintf("marshal JSON from DB: %s", err.Error()))
	}

	fmt.Printf("DB RECORD: %s", string(respJson))

	resp, err := utils.Unmarshal[Neo4jResponse[T]](respJson)
	if err != nil {
		return fail(fmt.Sprintf("unmarshal response from DB: %s", err.Error()))
	}
	result = resp.Props
	id = resp.ID

	return
}

func mapToProperties[T any](params T) (map[string]interface{}, error) {
	marshaled, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	props, err := utils.Unmarshal[map[string]interface{}](marshaled)
	if err != nil {
		return nil, err
	}

	return props, nil
}
