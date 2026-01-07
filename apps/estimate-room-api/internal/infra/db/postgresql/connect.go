// Package postgresql provides postgresql connection and utils
package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*pgxpool.Pool, error) {
	dbCfg, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, err
	}

	dbpool, err := pgxpool.NewWithConfig(context.Background(), dbCfg)
	if err != nil {
		return nil, err
	}

	err = dbpool.Ping(context.Background())
	if err != nil {
		dbpool.Close()
		return nil, err
	}

	return dbpool, nil
}
