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
		DBName:    getEnv("DB_NAME", "mew_mingle_db"),
		JWTSecret: getEnv("JWT_SECRET", "rsa256randomstring"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
