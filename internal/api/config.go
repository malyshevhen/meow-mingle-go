package api

import (
	"errors"
	"os"
)

const (
	ServerPortEnvKey  string = "SERVER_PORT"
	DefaultServerPort string = "3000"
)

type Config struct {
	Port string `yaml:"port"`
}

func (cfg *Config) SetEnv() {
	if serverPort := os.Getenv(ServerPortEnvKey); serverPort != "" {
		cfg.Port = serverPort
	} else if cfg.Port == "" {
		cfg.Port = DefaultServerPort
	}
}

func (c Config) Validate() error {
	_errors := make([]error, 0)

	if c.Port == "" {
		_errors = append(_errors, errors.New("port is required"))
	}

	if len(_errors) > 0 {
		return errors.Join(_errors...)
	}

	return nil
}
