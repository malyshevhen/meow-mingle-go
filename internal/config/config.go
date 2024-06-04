package config

import (
	"os"
)

const (
	SERVER_PORT_ENV_KEY string = "SERVER_PORT"
	DB_URL_ENV_KEY      string = "DB_URL"
	DB_USER_ENV_KEY     string = "DB_USER"
	DB_PASSWORD_ENV_KEY string = "DB_PASS"
	JWT_EVN_KEY         string = "JWT_SECRET"
	DEFAULT_SERVER_PORT string = ":8080"
	DEFAULT_DB_URL      string = "bolt://localhost:7687"
	DEFAULT_DB_USER     string = "neo4j"
	DEFAULT_DB_PASSWORD string = "example"
	DEFAULT_JWT_KEY     string = "13bb62a3f8a44d0523918228c3ea7643547495c7ba74c893f9546d6de37ad996"
)

type Config struct {
	ServerPort string
	DBConnURL  string
	DBUser     string
	DBPassword string
	JWTSecret  string
}

func InitConfig() Config {
	return Config{
		ServerPort: getEnv(SERVER_PORT_ENV_KEY, DEFAULT_SERVER_PORT),
		DBConnURL:  getEnv(DB_URL_ENV_KEY, DEFAULT_DB_URL),
		DBUser:     getEnv(DB_USER_ENV_KEY, DEFAULT_DB_USER),
		DBPassword: getEnv(DB_PASSWORD_ENV_KEY, DEFAULT_DB_PASSWORD),
		JWTSecret:  getEnv(JWT_EVN_KEY, DEFAULT_JWT_KEY),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
