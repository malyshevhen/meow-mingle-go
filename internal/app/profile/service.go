package profile

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	Save(ctx context.Context, userID, email, firstName, lastName string) (user app.Profile, err error)
	GetByID(ctx context.Context, id string) (user app.Profile, err error)
}

type service struct {
	profileRepo repository
}

// Create implements app.ProfileService.
func (s *service) Create(ctx context.Context, profile *app.Profile) error {
	panic("unimplemented")
}

// GetById implements app.ProfileService.
func (s *service) GetByID(ctx context.Context, profileID string) (user *app.Profile, err error) {
	panic("unimplemented")
}

func NewService(profileRepo repository) app.ProfileService {
	return &service{
		profileRepo: profileRepo,
	}
}
