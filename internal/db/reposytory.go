package db

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/malyshEvhen/meow_mingle/internal/errors"
	"github.com/malyshEvhen/meow_mingle/internal/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
)

type Reposytory[T any] struct {
	driver neo4j.DriverWithContext
}

func (s *Reposytory[T]) Create(ctx context.Context, params interface{}, cypher string) (entity T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := persist[T](ctx, tx, cypher, params)
		if err != nil {
			return nil, err
		}
		entity = result
		return result, nil
	}); err != nil {
		execErr = err
		return
	}
	return
}

func (s *Reposytory[T]) GetById(ctx context.Context, cypher string, id string) (entity T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"id": id,
	}

	if _, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := retrieveSingle[T](ctx, tx, cypher, params)
		if err != nil {
			return nil, err
		}
		entity = result

		return result, nil
	}); err != nil {
		execErr = err
		return
	}
	return
}

func (s *Reposytory[T]) Write(ctx context.Context, cypher string, params interface{}) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err := exec(ctx, tx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func (s *Reposytory[T]) Update(ctx context.Context, cypher string, params interface{}) (entity T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := persist[T](ctx, tx, updatePostCypher, params)
		if err != nil {
			return nil, err
		}
		entity = result

		return result, nil
	}); err != nil {
		execErr = err
		return
	}
	return
}

func (s *Reposytory[T]) List(ctx context.Context, cypher string, id string) (list []T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]interface{}{
		"id": id,
	}

	if _, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := retrieveMany[T](ctx, tx, cypher, params)
		if err != nil {
			return nil, err
		}
		list = result

		return result, nil
	}); err != nil {
		execErr = err
		return
	}
	return
}

func (s *Reposytory[T]) Delete(ctx context.Context, cypher string, params any) error {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	if _, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		err := exec(ctx, tx, cypher, params)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}); err != nil {
		return err
	}
	return nil
}

func exec(
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	form any,
) (err error) {
	fail := func(msg string) (err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	params, err := mapToProperties(form)
	if err != nil {
		return fail(fmt.Sprintf("convert to properties: %s", err.Error()))
	}

	_, err = tx.Run(ctx, cypher, params)
	if err != nil {
		return fail(fmt.Sprintf("executing transaction: %s", err.Error()))
	}
	return nil
}

func persist[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	form any,
) (obj T, err error) {
	fail := func(msg string) (res T, err error) {
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

	return parceRecord[T](record)
}

func retrieveSingle[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	params map[string]interface{},
) (obj T, err error) {
	fail := func(msg string) (res T, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	result, err := tx.Run(ctx, cypher, params)
	if err != nil {
		return fail(fmt.Sprintf("executing transaction: %s", err.Error()))
	}

	record, err := result.Single(ctx)
	if err != nil {
		return fail(fmt.Sprintf("retrieving the record: %s", err.Error()))
	}

	return parceRecord[T](record)
}

func retrieveMany[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	params map[string]interface{},
) (obj []T, err error) {
	fail := func(msg string) (res []T, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	result, err := tx.Run(ctx, cypher, params)
	if err != nil {
		return fail(fmt.Sprintf("executing transaction: %s", err.Error()))
	}

	records, err := result.Collect(ctx)
	if err != nil {
		return fail(fmt.Sprintf("retrieving the record: %s", err.Error()))
	}

	for _, record := range records {
		res, err := parceRecord[T](record)
		if err != nil {
			return nil, err
		}
		obj = append(obj, res)
	}

	return
}

func parceRecord[T any](r *db.Record) (result T, err error) {
	fail := func(msg string) (res T, err error) {
		err = errors.NewDatabaseError(fmt.Errorf("error occurred while %s", msg))
		return
	}

	typ := reflect.TypeFor[T]()

	rows := readJsonTags(typ)
	vals := getRecordValues(r, rows)

	respJson, err := json.MarshalIndent(vals, "", "  ")
	if err != nil {
		return fail(fmt.Sprintf("marshal JSON from DB: %s", err.Error()))
	}

	fmt.Printf("DB RECORD: \n %s", string(respJson))

	return utils.Unmarshal[T](respJson)
}

func readJsonTags(typ reflect.Type) []string {
	rows := []string{}
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		rows = append(rows, jsonTag)
	}
	return rows
}

func getRecordValues(r *db.Record, rows []string) map[string]interface{} {
	vals := make(map[string]interface{})
	for _, key := range rows {
		val, _ := r.Get(key)
		vals[key] = val
	}
	return vals
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
