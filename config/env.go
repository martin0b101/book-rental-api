package config

import (
	"os"
	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	RedisAddress string
}

// singleton
var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()
	return Config{
		Host: getEnv("DB_HOST", "127.0.0.1"),
		Port: getEnv("DB_PORT", "5432"),
		User: getEnv("DB_USER", "user"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName: getEnv("DB_NAME", "postgres"),
		RedisAddress: getEnv("REDIS_ADDRESS", "localhost:6379"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}