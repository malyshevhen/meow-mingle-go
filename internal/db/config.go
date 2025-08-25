package db

import (
	"errors"
	"os"
)

const (
	SCYLLA_URL_ENV_KEY      string = "SCYLLA_URL"
	SCYLLA_USER_ENV_KEY     string = "SCYLLA_USER"
	SCYLLA_PASSWORD_ENV_KEY string = "SCYLLA_PASS"
)

var (
	ErrMissingDBConnectionURL error = errors.New("missing database connection URL")
	ErrMissingDBUser          error = errors.New("missing database user")
	ErrMissingDBPassword      error = errors.New("missing database password")
)

// Config is the scylla DB configuration
type Config struct {
	URL      string `yaml:"url" json:"url"`
	User     string `yaml:"user" json:"user"`
	Password string `yaml:"password" json:"password"`
}

// SetEnv Updates config with values from environment if available
func (c *Config) SetEnv() {
	if url := os.Getenv(SCYLLA_URL_ENV_KEY); url != "" {
		c.URL = url
	}
	if user := os.Getenv(SCYLLA_USER_ENV_KEY); user != "" {
		c.User = user
	}
	if password := os.Getenv(SCYLLA_PASSWORD_ENV_KEY); password != "" {
		c.Password = password
	}
}

func (c *Config) Validate() error {
	_errors := make([]error, 0)
	if len(c.URL) == 0 {
		_errors = append(_errors, ErrMissingDBConnectionURL)
	}

	if len(c.User) == 0 {
		_errors = append(_errors, ErrMissingDBUser)
	}

	if len(c.Password) == 0 {
		_errors = append(_errors, ErrMissingDBPassword)
	}

	if len(_errors) > 0 {
		return errors.Join(_errors...)
	}

	return nil
}
