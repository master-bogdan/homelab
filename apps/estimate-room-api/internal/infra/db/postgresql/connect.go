// Package postgresql provides postgresql connection and utils
package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/master-bogdan/estimate-room-api/config"
)

func Connect(cfg *config.Config) (*pgxpool.Pool, error) {
	dbCfg, err := pgx.ParseConfig(cfg.DB.DatabaseURL)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DB.DatabaseURL)
	if err != nil {
		return nil, err
	}
	defer dbpool.Close()

	return dbpool, nil
}
