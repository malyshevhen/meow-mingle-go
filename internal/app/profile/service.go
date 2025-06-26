package profile

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type repository interface {
	Save(ctx context.Context, userId, email, firstName, lastName string) (user app.Profile, err error)
	GetById(ctx context.Context, id string) (user app.Profile, err error)
	GetByEmail(ctx context.Context, email string) (user app.Profile, err error)
}

type service struct {
	profileRepo repository
}

// Create implements app.ProfileService.
func (s *service) Create(ctx context.Context, profile *app.Profile) error {
	panic("unimplemented")
}

// GetByEmail implements app.ProfileService.
func (s *service) GetByEmail(ctx context.Context, profileEmail string) (user *app.Profile, err error) {
	panic("unimplemented")
}

// GetById implements app.ProfileService.
func (s *service) GetById(ctx context.Context, profileId string) (user *app.Profile, err error) {
	panic("unimplemented")
}

func NewService(profileRepo repository) app.ProfileService {
	return &service{
		profileRepo: profileRepo,
	}
}
