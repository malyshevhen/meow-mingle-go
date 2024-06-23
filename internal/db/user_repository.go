package db

import (
	"context"
	_ "embed"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	//go:embed cypher/create_user.cypher
	createUserCypher string

	//go:embed cypher/match_user_by_id.cypher
	getUserCypher string

	//go:embed cypher/create_subscription.cypher
	createSubscriptionCypher string

	//go:embed cypher/delete_subscription.cypher
	deleteSubscriptionCypher string
)

type IUserReposytory interface {
	CreateUser(ctx context.Context, params CreateUserParams) (user User, err error)
	CreateSubscription(ctx context.Context, params CreateSubscriptionParams) error
	GetUser(ctx context.Context, id int64) (user User, err error)
	DeleteSubscription(ctx context.Context, params DeleteSubscriptionParams) error
}

type UserRepository struct {
	Reposytory[User]
}

func NewUserReposiory(driver neo4j.DriverWithContext) *UserRepository {
	return &UserRepository{
		Reposytory: Reposytory[User]{
			driver: driver,
		},
	}
}

func (ur *UserRepository) CreateUser(ctx context.Context, params CreateUserParams) (user User, err error) {
	return ur.Create(ctx, params, createUserCypher)
}

func (ur *UserRepository) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) error {
	return ur.Write(ctx, createSubscriptionCypher, params)
}

func (ur *UserRepository) GetUser(ctx context.Context, id int64) (user User, err error) {
	return ur.GetById(ctx, getUserCypher, id)
}

func (ur *UserRepository) DeleteSubscription(ctx context.Context, params DeleteSubscriptionParams) error {
	return ur.Delete(ctx, deleteSubscriptionCypher, params)
}
