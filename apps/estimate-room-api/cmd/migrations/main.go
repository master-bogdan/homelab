package main

import (
	"flag"
	"fmt"
	"log"
	"slices"

	"github.com/master-bogdan/estimate-room-api/config"
	"github.com/master-bogdan/estimate-room-api/internal/infra/db/postgresql"
)

const (
	Create string = "create"
	Up     string = "up"
	Down   string = "down"
)

func main() {
	var command string
	var name string

	flag.StringVar(&command, "command", "up", "Migration command: up, down or create")
	flag.StringVar(&name, "name", "", "Migration name (required for create command")
	flag.Parse()

	commands := []string{Create, Up, Down}

	if !slices.Contains(commands, command) {
		log.Fatalf("Invalid command: %s. Use 'up', 'down' or 'create'", command)
	}

	if command == Create {
		if name == "" {
			log.Fatal("Migration name is required. Use: -command=create -name=your_migration_name")
		}

		upFile, downFile, err := postgresql.MigrateCreate(name)
		if err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}

		fmt.Printf("✓ Created migration files:\n")
		fmt.Printf("  - %s\n", upFile)
		fmt.Printf("  - %s\n", downFile)

		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	switch command {
	case Up:
		fmt.Println("Applying migrations...")
		if err := postgresql.MigrateUp(cfg.DB.DatabaseURL); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		fmt.Println("✓ Migrations applied successfully")
	case Down:
		fmt.Println("Rolling back migrations...")
		if err := postgresql.MigrateDown(cfg.DB.DatabaseURL); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}

		fmt.Println("✓ Migrations rolled back successfully")
	}
}
