package main

import (
	"context"
	"log"

	"github.com/martin0b101/book-rental-api/api"
	"github.com/martin0b101/book-rental-api/config"
	"github.com/martin0b101/book-rental-api/database"
	"github.com/redis/go-redis/v9"
)

func main() {

	db, err := database.NewPostgresSQLStorage(config.Config{
		Host:     config.Envs.Host,
		Port:     config.Envs.Port,
		DBName:   config.Envs.DBName,
		User:     config.Envs.User,
		Password: config.Envs.Password,
	})

	if err != nil {
		log.Fatal(err)
	}

	var ctx = context.Background()

	redisCache := redis.NewClient(&redis.Options{
		Addr: config.Envs.RedisAddress,
		Password: "",
		DB: 0,
	})

	errCache := redisCache.Ping(ctx).Err()

	if errCache != nil {
		log.Fatal("Error connecting to redis.")
	} else {
		log.Printf("Connected to redis!")
	}

	server := api.NewApiServer(":8080", db, redisCache, ctx)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
