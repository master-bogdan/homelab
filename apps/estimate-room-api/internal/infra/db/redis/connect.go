// Package redis provides redis connection and utils
package redis

import (
	"context"
	"time"

	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	opts, err := redis.ParseURL(cfg.DB.RedisURL)
	if err != nil {
		return nil, err
	}

	rdb := redis.NewClient(opts)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
