package db

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type profileRepository struct {
	session *gocql.Session
}

// Save implements profile.repository.
func (ur *profileRepository) Save(ctx context.Context, userId, email, firstName, lastName string) (user app.Profile, err error) {
	panic("unimplemented")
}

// GetById implements profile.repository.
func (ur *profileRepository) GetById(ctx context.Context, id string) (user app.Profile, err error) {
	panic("unimplemented")
}

func NewProfileRepository(session *gocql.Session) *profileRepository {
	return &profileRepository{
		session: session,
	}
}
