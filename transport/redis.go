package transport

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisConfig(addr, password string, db int) *RedisConfig {
	return &RedisConfig{
		Addr:     addr,
		Password: password,
		DB:       db,
	}
}

type RedisTransport struct {
	config *RedisConfig
	client *redis.Client
}

func NewRedisTransport(config *RedisConfig) *RedisTransport {
	cl := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	return &RedisTransport{
		config: config,
		client: cl,
	}
}

func (r *RedisTransport) Get(ctx context.Context, key string) (string, error) {
	cmd := r.client.Get(ctx, key)
	if cmd.Err() != nil {
		return "", cmd.Err()
	}
	return cmd.Result()
}

func (r *RedisTransport) Set(ctx context.Context, key string, value string) error {
	cmd := r.client.Set(ctx, key, value, 0)
	return cmd.Err()
}
