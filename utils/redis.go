package utils

import (
	"strconv"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient 回傳一個連線到指定 config 的 Redis client
func NewRedisClient(cfg AppConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort),
		Username: cfg.RedisUser,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}
