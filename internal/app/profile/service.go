package profile

import (
	"context"

	"github.com/malyshEvhen/meow_mingle/internal/app"
)

type service struct {
	profileRepo app.ProfileRepository
}

// CreateProfile implements app.ProfileService.
func (s *service) CreateProfile(ctx context.Context, profile *app.Profile) error {
	panic("unimplemented")
}

// GetProfileByEmail implements app.ProfileService.
func (s *service) GetProfileByEmail(ctx context.Context, profileEmail string) (user *app.Profile, err error) {
	panic("unimplemented")
}

// GetProfileById implements app.ProfileService.
func (s *service) GetProfileById(ctx context.Context, profileId string) (user *app.Profile, err error) {
	panic("unimplemented")
}

func NewService(profileRepo app.ProfileRepository) app.ProfileService {
	return &service{
		profileRepo: profileRepo,
	}
}
