package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type ICache interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Exists(key string) bool
	Del(key string) error
	HSet(key string, field string, value interface{}) error
	HGet(key, field string) (string, error)
	HGetAll(key string) (map[string]string, error)
	Expire(key string, expiration time.Duration) error
}

type RedisCache struct {
	client *redis.Client
}

func NewRedisCache(addr, password string, db int) ICache {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisCache{
		client: client,
	}
}

func (c *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(context.Background(), key, value, expiration).Err()
}

func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(context.Background(), key).Result()
}

func (c *RedisCache) Exists(key string) bool {
	result, err := c.client.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return result > 0
}

func (c *RedisCache) Del(key string) error {
	return c.client.Del(context.Background(), key).Err()
}

func (c *RedisCache) HSet(key string, field string, value interface{}) error {
	return c.client.HSet(context.Background(), key, field, value).Err()
}

func (c *RedisCache) HGet(key, field string) (string, error) {
	return c.client.HGet(context.Background(), key, field).Result()
}

func (c *RedisCache) HGetAll(key string) (map[string]string, error) {
	return c.client.HGetAll(context.Background(), key).Result()
}

func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(context.Background(), key, expiration).Err()
}
