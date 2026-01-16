package postgresql

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	migrationsPath = "file://migrations"
	migrationsDir  = "migrations"
)

func migrateInstance(databaseURL string) (*migrate.Migrate, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func MigrateUp(databaseURL string) error {
	m, err := migrateInstance(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func MigrateDown(databaseURL string) error {
	m, err := migrateInstance(databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Down()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func MigrateCreate(name string) (string, string, error) {
	if name == "" {
		return "", "", fmt.Errorf("migration name can't be empty")
	}

	err := os.MkdirAll(migrationsDir, 0755)
	if err != nil {
		return "", "", fmt.Errorf("failed to create migrations directory: %w", err)
	}

	timestamp := time.Now().Unix()

	upFile := filepath.Join(migrationsDir, fmt.Sprintf("%d_%s.up.sql", timestamp, name))
	downFile := filepath.Join(migrationsDir, fmt.Sprintf("%d_%s.down.sql", timestamp, name))

	err = os.WriteFile(upFile, []byte("-- Write your UP migration here\n"), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to create up migration file: %w", err)
	}

	err = os.WriteFile(downFile, []byte("-- Write your Down migration here\n"), 0644)
	if err != nil {
		return "", "", fmt.Errorf("failed to create down migration file: %w", err)
	}

	return upFile, downFile, nil
}
