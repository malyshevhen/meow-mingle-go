package mingle

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/goccy/go-yaml"
)

const (
	CONFIG_PATH_ENV_KEY string = "CONFIG_PATH"
	SERVER_PORT_ENV_KEY string = "SERVER_PORT"
	DB_URL_ENV_KEY      string = "DB_URL"
	DB_USER_ENV_KEY     string = "DB_USER"
	DB_PASSWORD_ENV_KEY string = "DB_PASS"
	JWT_EVN_KEY         string = "JWT_SECRET"
	DEFAULT_CONFIG_PATH string = "/opt/minge/config.yaml"
	DEFAULT_SERVER_PORT string = "3000"
)

var (
	ErrMissingDBConnectionURL error = errors.New("missing database connection URL")
	ErrMissingDBUser          error = errors.New("missing database user")
	ErrMissingDBPassword      error = errors.New("missing database password")
	ErrMissingJWTSecret       error = errors.New("missing JWT secret")
)

type Config struct {
	ServerPort string `yaml:"server.port"`
	DBConnURL  string `yaml:"database.connection_url"`
	DBUser     string `yaml:"database.user"`
	DBPassword string `yaml:"database.password"`
	JWTSecret  string `yaml:"secret"`
}

func (c Config) validate() error {
	_errors := make([]error, 0)
	if len(c.DBConnURL) == 0 {
		_errors = append(_errors, ErrMissingDBConnectionURL)
	}

	if len(c.DBUser) == 0 {
		_errors = append(_errors, ErrMissingDBUser)
	}

	if len(c.DBPassword) == 0 {
		_errors = append(_errors, ErrMissingDBPassword)
	}

	if len(c.JWTSecret) == 0 {
		_errors = append(_errors, ErrMissingJWTSecret)
	}

	if len(_errors) > 0 {
		return errors.Join(_errors...)
	}

	return nil
}

func initConfig() (Config, error) {
	cfg, err := readConfigFromFile()
	if err != nil {
		return Config{}, err
	}

	if serverPort, ok := os.LookupEnv(SERVER_PORT_ENV_KEY); ok {
		cfg.ServerPort = serverPort
	}

	if connURL, ok := os.LookupEnv(DB_URL_ENV_KEY); ok {
		cfg.DBConnURL = connURL
	}

	if dbUser, ok := os.LookupEnv(DB_USER_ENV_KEY); ok {
		cfg.DBUser = dbUser
	}

	if dbPassword, ok := os.LookupEnv(DB_PASSWORD_ENV_KEY); ok {
		cfg.DBPassword = dbPassword
	}

	if jwtSecret, ok := os.LookupEnv(JWT_EVN_KEY); ok {
		cfg.JWTSecret = jwtSecret
	}

	if err := cfg.validate(); err != nil {
		return Config{}, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func readConfigFromFile() (Config, error) {
	filePath, ok := os.LookupEnv(CONFIG_PATH_ENV_KEY)
	if !ok {
		filePath = DEFAULT_CONFIG_PATH
	}

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		slog.Warn("Config file not defined.", "Default config location: ", DEFAULT_CONFIG_PATH)
		return Config{}, nil
	}

	var cfg Config
	if err := yaml.Unmarshal([]byte(contentBytes), &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return cfg, nil
}
