package initialize

import (
	"context"

	"github.com/baothaihcmut/Ecommerce-go/users/internal/config"
	"github.com/redis/go-redis/v9"
)

func InitializeRedis(cfg *config.RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Endpoint,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}
