package utils

import (
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient 根據 AppConfig 自動判斷單點或 Cluster，回傳 redis.UniversalClient
func NewRedisClient(cfg AppConfig) redis.UniversalClient {
	addrs := []string{cfg.RedisHost + ":" + strconv.Itoa(cfg.RedisPort)}
	if strings.Contains(cfg.RedisHost, ",") {
		addrs = strings.Split(cfg.RedisHost, ",")
		for i, addr := range addrs {
			if !strings.Contains(addr, ":") {
				addrs[i] = addr + ":" + strconv.Itoa(cfg.RedisPort)
			}
		}
	}
	if len(addrs) > 1 {
		return redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:    addrs,
			Username: cfg.RedisUser,
			Password: cfg.RedisPassword,
		})
	}
	return redis.NewClient(&redis.Options{
		Addr:     addrs[0],
		Username: cfg.RedisUser,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
}
