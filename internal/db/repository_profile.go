package db

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	//go:embed cypher/create_user.cypher
	createProfileCypher string

	//go:embed cypher/match_user_by_id.cypher
	getProfileByIdCypher string

	//go:embed cypher/match_user_by_email.cypher
	getProfileByEmailCypher string

	//go:embed cypher/create_subscription.cypher
	subscribeCypher string

	//go:embed cypher/delete_subscription.cypher
	unsubscribeCypher string
)

type IProfileRepository interface {
	CreateProfile(ctx context.Context, params CreateProfileParams) (user Profile, err error)
	CreateSubscription(ctx context.Context, params CreateSubscriptionParams) error
	GetProfileById(ctx context.Context, id string) (user Profile, err error)
	GetProfileByEmail(ctx context.Context, email string) (user Profile, err error)
	DeleteSubscription(ctx context.Context, params DeleteSubscriptionParams) error
}

type ProfileNeo4jRepository struct {
	Neo4jRepository[Profile]
}

func NewUserRepository(driver neo4j.DriverWithContext) *ProfileNeo4jRepository {
	return &ProfileNeo4jRepository{
		Neo4jRepository: Neo4jRepository[Profile]{
			driver: driver,
		},
	}
}

func (ur *ProfileNeo4jRepository) CreateProfile(ctx context.Context, params CreateProfileParams) (user Profile, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return ur.Create(ctx, params, createProfileCypher)
}

func (ur *ProfileNeo4jRepository) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	params.ID = id.String()

	return ur.Write(ctx, subscribeCypher, params)
}

func (ur *ProfileNeo4jRepository) GetProfileById(ctx context.Context, id string) (user Profile, err error) {
	return ur.Retrieve(ctx, getProfileByIdCypher, map[string]any{
		"id": id,
	})
}

func (ur *ProfileNeo4jRepository) GetProfileByEmail(ctx context.Context, email string) (user Profile, err error) {
	return ur.Retrieve(ctx, getProfileByEmailCypher, map[string]any{
		"email": email,
	})
}

func (ur *ProfileNeo4jRepository) DeleteSubscription(ctx context.Context, params DeleteSubscriptionParams) error {
	return ur.Delete(ctx, unsubscribeCypher, params)
}
