package config

import (
	"os"
)

type Config struct {
	ServerPort string
	DBSource   string
	JWTSecret  string
}

func InitConfig() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", ":8080"),
		DBSource: getEnv(
			"DB_SOURCE",
			"postgresql://postgres:example@localhost:5432/mingle_db?sslmode=disable",
		),
		JWTSecret: getEnv(
			"JWT_SECRET",
			"13bb62a3f8a44d0523918228c3ea7643547495c7ba74c893f9546d6de37ad996",
		),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
