package core_cache_redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Ping(ctx context.Context) error
	Close() error
}

type client struct {
	rdb *redis.Client
}

func New(cfg Config) (Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}
	return &client{rdb: rdb}, nil
}

func (c *client) Set(
	ctx context.Context,
	key string,
	value any,
	expiration time.Duration,
) error {
	return c.rdb.Set(ctx, key, value, expiration).Err()
}

func (c *client) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return c.rdb.Get(ctx, key).Result()
}

func (c *client) Del(
	ctx context.Context,
	keys ...string,
) error {
	return c.rdb.Del(ctx, keys...).Err()
}

func (c *client) Ping(ctx context.Context) error {
	return c.rdb.Ping(ctx).Err()
}

func (c *client) Close() error {
	return c.rdb.Close()
}
