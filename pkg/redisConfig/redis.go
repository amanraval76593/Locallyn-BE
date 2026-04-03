package redisConfig

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitRedis(redisAdd string) {
	RedisClient = redis.NewClient(
		&redis.Options{
			Addr: redisAdd,
		},
	)

	_, err := RedisClient.Ping(context.Background()).Result()

	if err != nil {
		log.Fatalf("Error initialization redis : %v", err)
	}

	log.Println("Redis Client initialized")
}
