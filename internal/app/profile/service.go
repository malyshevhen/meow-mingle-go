package profile

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	profileRepo app.ProfileRepository
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

func NewService(profileRepo app.ProfileRepository) app.ProfileService {
	return &service{
		profileRepo: profileRepo,
	}
}
