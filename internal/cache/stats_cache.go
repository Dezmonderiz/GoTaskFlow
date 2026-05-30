package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"GoTaskFlow/internal/model"

	"github.com/redis/go-redis/v9"
)

var ErrCacheMiss = errors.New("cache miss")

type StatsCache interface {
	Get(ctx context.Context) (model.TaskStats, error)
	Set(ctx context.Context, stats model.TaskStats) error
	Delete(ctx context.Context) error
}

type RedisStatsCache struct {
	client *redis.Client
	key    string
	ttl    time.Duration
}

func NewRedisStatsCache(client *redis.Client, ttl time.Duration) *RedisStatsCache {
	return &RedisStatsCache{
		client: client,
		key:    "task_stats",
		ttl:    ttl,
	}
}

func (c *RedisStatsCache) Get(ctx context.Context) (model.TaskStats, error) {
	value, err := c.client.Get(ctx, c.key).Result()
	if errors.Is(err, redis.Nil) {
		return model.TaskStats{}, ErrCacheMiss
	}
	if err != nil {
		return model.TaskStats{}, err
	}

	var stats model.TaskStats
	if err := json.Unmarshal([]byte(value), &stats); err != nil {
		return model.TaskStats{}, err
	}

	return stats, nil
}

func (c *RedisStatsCache) Set(ctx context.Context, stats model.TaskStats) error {
	value, err := json.Marshal(stats)
	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.key, value, c.ttl).Err()
}

func (c *RedisStatsCache) Delete(ctx context.Context) error {
	return c.client.Del(ctx, c.key).Err()
}
