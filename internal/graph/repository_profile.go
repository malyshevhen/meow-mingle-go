package graph

import (
	"context"
	_ "embed"

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
	query Neo4jQuerier[app.Profile]
}

func NewProfileRepository(driver neo4j.DriverWithContext) *profileNeo4jRepository {
	return &profileNeo4jRepository{
		query: Neo4jQuerier[app.Profile]{
			driver: driver,
		},
	}
}

// Save implements app.ProfileRepository.
func (ur *profileNeo4jRepository) Save(ctx context.Context, userId, email, firstName, lastName string) (user app.Profile, err error) {
	params := struct {
		ID        string `json:"id" validate:"required"`
		Email     string `json:"email" validate:"required,email"`
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
	}{
		ID:        userId,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
	}

	return ur.query.Create(ctx, params, createProfileCypher)
}

// GetById implements app.ProfileRepository.
func (ur *profileNeo4jRepository) GetById(ctx context.Context, id string) (user app.Profile, err error) {
	return ur.query.Retrieve(ctx, getProfileByIdCypher, map[string]any{
		"id": id,
	})
}
