package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/malyshEvhen/meow_mingle/cmd/mingle"
	"github.com/malyshEvhen/meow_mingle/pkg/logger"
)

const (
	CONFIG_PATH_ENV_KEY string = "CONFIG_PATH"
	DEFAULT_CONFIG_PATH string = "/opt/mingle/config.yaml"
)

type Config struct {
	Logger logger.LoggerConfig `yaml:"logger"`
	Mingle mingle.Config       `yaml:"mingle"`
}

func (cfg *Config) SetEnv() {
	cfg.Mingle.SetEnv()
}

func (cfg *Config) Validate() error {
	if err := cfg.Mingle.Validate(); err != nil {
		return err
	}

	return nil
}

func InitConfig() (Config, error) {
	configLogger := logger.GetLogger()
	configLogger.WithComponent("config").Info("Initializing config")

	cfg, err := readConfigFromFile()
	if err != nil {
		return Config{}, err
	}

	cfg.SetEnv()

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func readConfigFromFile() (Config, error) {
	configLogger := logger.GetLogger().WithComponent("config")
	configLogger.Info("Reading config file")

	filePath := os.Getenv(CONFIG_PATH_ENV_KEY)
	if filePath == "" {
		filePath = DEFAULT_CONFIG_PATH
	}

	configLogger.Info("Config file path", "path", filePath)

	contentBytes, err := os.ReadFile(filePath)
	if err != nil {
		configLogger.Warn("Config file not found, using environment variables only",
			"attempted_path", filePath,
		)
		return Config{}, nil
	}

	configLogger.Info("Config file content", "content", string(contentBytes))

	var cfg Config
	if err := yaml.Unmarshal(contentBytes, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config file: %w", err)
	}
	configLogger.Info("Application config read successfully", "config", cfg)

	return cfg, nil
}
