package memory_db

import (
	"context"
	"time"

	"github.com/master-bogdan/ephermal-notes/pkg/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&cfg.Db.Redis)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
