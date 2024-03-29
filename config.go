package main

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	DBUser    string
	DBPasswd  string
	DBAddress string
	DBName    string
	JWTSecret string
}

var Envs = initConfig()

func initConfig() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		DBUser:    getEnv("DB_USER", "root"),
		DBPasswd:  getEnv("DB_PASSWORD", "example"),
		DBAddress: fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "3306")),
		DBName:    getEnv("DB_NAME", "mingle_db"),
		JWTSecret: getEnv("JWT_SECRET", "13bb62a3f8a44d0523918228c3ea7643547495c7ba74c893f9546d6de37ad996"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
