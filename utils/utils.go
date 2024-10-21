package utils

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func DeleteFromRedis(context context.Context, cache *redis.Client, key string) {
	err := cache.Del(context, key).Err()
	if err != nil{
		log.Printf("Error: deleting redis %s", err.Error())
	}
}