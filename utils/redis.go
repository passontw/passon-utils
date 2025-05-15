package utils

import (
	"github.com/redis/go-redis/v9"
)

// NewRedisClient 回傳一個連線到 localhost:6379 的 Redis client
func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
}
