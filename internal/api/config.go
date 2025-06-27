package api

import (
	"errors"
	"os"
)

const (
	SERVER_PORT_ENV_KEY string = "SERVER_PORT"
	DEFAULT_SERVER_PORT string = "3000"
)

type Config struct {
	Port string `yaml:"port"`
}

func (cfg *Config) SetEnv() {
	if serverPort := os.Getenv(SERVER_PORT_ENV_KEY); serverPort != "" {
		cfg.Port = ":" + serverPort
	} else if cfg.Port == "" {
		cfg.Port = ":" + DEFAULT_SERVER_PORT
	} else if cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
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
