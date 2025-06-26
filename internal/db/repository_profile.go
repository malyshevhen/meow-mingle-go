package db

import (
	"context"
	_ "embed"

	"github.com/google/uuid"
	"github.com/malyshEvhen/meow_mingle/internal/app"
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

type profileNeo4jRepository struct {
	Neo4jRepository[app.Profile]
}

func NewProfileRepository(driver neo4j.DriverWithContext) *profileNeo4jRepository {
	return &profileNeo4jRepository{
		Neo4jRepository: Neo4jRepository[app.Profile]{
			driver: driver,
		},
	}
}

func (ur *profileNeo4jRepository) CreateProfile(ctx context.Context, params app.CreateProfileParams) (user app.Profile, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return
	}
	params.ID = id.String()

	return ur.Create(ctx, params, createProfileCypher)
}

func (ur *profileNeo4jRepository) GetProfileById(ctx context.Context, id string) (user app.Profile, err error) {
	return ur.Retrieve(ctx, getProfileByIdCypher, map[string]any{
		"id": id,
	})
}

func (ur *profileNeo4jRepository) GetProfileByEmail(ctx context.Context, email string) (user app.Profile, err error) {
	return ur.Retrieve(ctx, getProfileByEmailCypher, map[string]any{
		"email": email,
	})
}
