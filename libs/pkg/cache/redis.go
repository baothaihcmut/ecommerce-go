package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	client *redis.Client
}

// Remove implements CacheService.
func (r *RedisService) Remove(ctx context.Context, key string) error {
	_, err := r.client.Del(ctx, key).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func (r *RedisService) SetValue(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	val, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetValue(ctx context.Context, key string, ouput any) error {
	jsonData, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}
	err = json.Unmarshal([]byte(jsonData), ouput)
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) SetString(ctx context.Context, key string, val string, ttl time.Duration) error {
	err := r.client.Set(ctx, key, val, ttl).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisService) GetString(ctx context.Context, key string) (string, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return res, nil
}

func NewRedisService(client *redis.Client) CacheService {
	return &RedisService{
		client: client,
	}
}
