package db

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/malyshEvhen/meow_mingle/pkg/errors"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/db"
)

type Neo4jRepository[T any] struct {
	driver neo4j.DriverWithContext
}

func (s *Neo4jRepository[T]) Create(ctx context.Context, params any, cypher string) (entity T, execErr error) {
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

func (s *Neo4jRepository[T]) Retrieve(ctx context.Context, cypher string, params map[string]any) (entity T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

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

func (s *Neo4jRepository[T]) Write(ctx context.Context, cypher string, params any) error {
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

func (s *Neo4jRepository[T]) Update(ctx context.Context, cypher string, params any) (entity T, execErr error) {
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

func (s *Neo4jRepository[T]) List(ctx context.Context, cypher string, id string) (list []T, execErr error) {
	session := s.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)

	params := map[string]any{
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

func (s *Neo4jRepository[T]) Delete(ctx context.Context, cypher string, params any) error {
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

	return parseRecord[T](record)
}

func retrieveSingle[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	params map[string]any,
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

	return parseRecord[T](record)
}

func retrieveMany[T any](
	ctx context.Context,
	tx neo4j.ManagedTransaction,
	cypher string,
	params map[string]any,
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
		res, err := parseRecord[T](record)
		if err != nil {
			return nil, err
		}
		obj = append(obj, res)
	}

	return
}

func parseRecord[T any](r *db.Record) (result T, err error) {
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

	return Unmarshal[T](respJson)
}

func readJsonTags(typ reflect.Type) []string {
	rows := []string{}
	for i := range typ.NumField() {
		field := typ.Field(i)
		jsonTag := field.Tag.Get("json")
		rows = append(rows, jsonTag)
	}
	return rows
}

func getRecordValues(r *db.Record, rows []string) map[string]any {
	vals := make(map[string]any)
	for _, key := range rows {
		val, _ := r.Get(key)
		vals[key] = val
	}
	return vals
}

func mapToProperties[T any](params T) (map[string]any, error) {
	marshaled, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	props, err := Unmarshal[map[string]any](marshaled)
	if err != nil {
		return nil, err
	}

	return props, nil
}

func Unmarshal[T any](v []byte) (value T, err error) {
	if err := json.Unmarshal(v, &value); err != nil {
		return value, errors.NewValidationError("error parse JSON payload")
	}
	return
}
