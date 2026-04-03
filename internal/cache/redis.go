package cache

import (
	"context"
	"time"

	"locallyn-be/pkg/redisConfig"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedis() *Redis {
	return &Redis{
		Client: redisConfig.RedisClient,
		ctx:    context.Background(),
	}
}

func (r *Redis) Get(key string) (string, error) {
	return r.Client.Get(r.ctx, key).Result()
}

func (r *Redis) Set(key string, value string, ttl time.Duration) error {
	return r.Client.Set(r.ctx, key, value, ttl).Err()
}

func (r *Redis) Delete(key string) error {
	return r.Client.Del(r.ctx, key).Err()
}
