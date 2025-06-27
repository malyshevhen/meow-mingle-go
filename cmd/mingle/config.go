package mingle

import (
	"errors"

	"github.com/malyshEvhen/meow_mingle/internal/api"
	"github.com/malyshEvhen/meow_mingle/internal/db"
)

type Config struct {
	Server   api.Config `yaml:"server"`
	Database db.Config  `yaml:"database"`
}

func (c Config) Validate() error {
	_errors := make([]error, 0)

	if err := c.Database.Validate(); err != nil {
		_errors = append(_errors, err)
	}

	if err := c.Server.Validate(); err != nil {
		_errors = append(_errors, err)
	}

	if len(_errors) > 0 {
		return errors.Join(_errors...)
	}

	return nil
}

func (cfg *Config) SetEnv() {
	cfg.Server.SetEnv()
	cfg.Database.SetEnv()
}
