package main

import (
	"flag"
	"os"
	"slices"

	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/logger"
)

const (
	Create string = "create"
	Up     string = "up"
	Down   string = "down"
)

func main() {
	logger.InitLogger()

	var command string
	var name string

	flag.StringVar(&command, "command", "up", "Migration command: up, down or create")
	flag.StringVar(&name, "name", "", "Migration name (required for create command")
	flag.Parse()

	commands := []string{Create, Up, Down}

	if !slices.Contains(commands, command) {
		logger.L().Error("Invalid command. Use 'up', 'down' or 'create'", "command", command)
		os.Exit(1)
	}

	if command == Create {
		if name == "" {
			logger.L().Error("Migration name is required. Use: -command=create -name=your_migration_name")
			os.Exit(1)
		}

		upFile, downFile, err := postgresql.MigrateCreate(name)
		if err != nil {
			logger.L().Error("Failed to create migration", "err", err)
			os.Exit(1)
		}

		logger.L().Info("Created migration files", "up_file", upFile, "down_file", downFile)

		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.L().Error("Failed to load config", "err", err)
		os.Exit(1)
	}

	switch command {
	case Up:
		logger.L().Info("Applying migrations...")
		if err := postgresql.MigrateUp(cfg.DB.DatabaseURL); err != nil {
			logger.L().Error("Failed to run migrations", "err", err)
			os.Exit(1)
		}

		logger.L().Info("Migrations applied successfully")
	case Down:
		logger.L().Info("Rolling back migrations...")
		if err := postgresql.MigrateDown(cfg.DB.DatabaseURL); err != nil {
			logger.L().Error("Failed to rollback migrations", "err", err)
			os.Exit(1)
		}

		logger.L().Info("Migrations rolled back successfully")
	}
}
