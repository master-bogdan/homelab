// Package postgresql provides postgresql connection and utils
package postgresql

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func Connect(databaseURL string) (*bun.DB, error) {
	sqldb, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	err = db.PingContext(context.Background())
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
