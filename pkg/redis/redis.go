package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func InitRedis(redisAdd string) {
	redisClient = redis.NewClient(
		&redis.Options{
			Addr: redisAdd,
		},
	)

	_, err := redisClient.Ping(context.Background()).Result()

	if err != nil {
		log.Fatalf("Error initialization redis : %v", err)
	}

	log.Println("Redis Client initialized")
}
