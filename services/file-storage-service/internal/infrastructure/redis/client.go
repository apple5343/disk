package redis

import (
	"context"
	"storage/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewClient(cfg *config.RedisConfig) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	if err := c.Ping(context.TODO()).Err(); err != nil {
		return nil, err
	}
	return c, nil
}
