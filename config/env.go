package config

import (
	"os"
	"github.com/joho/godotenv"
)

// Config holds PostgreSQL configuration details
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// singleton
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		Host: getEnv("DB_HOST", "http://localhost"),
		Port: getEnv("DB_PORT", "8080"),
		User: getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName: getEnv("DB_NAME", "bookstore"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}