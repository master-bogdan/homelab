.PHONY: migrate-up migrate-down migrate-create

# Apply all pending migrations
migrate-up:
	@echo "Applying migrations..."
	@go run cmd/migrations/main.go -command=up

# Rollback all migrations
migrate-down:
	@echo "Rolling back all migrations..."
	@go run cmd/migrations/main.go -command=down

# Create a new migration file pair
# Usage: make migrate-create name=create_users_table
migrate-create:
ifndef name
	@echo "Error: Migration name is required"
	@echo "Usage: make migrate-create name=your_migration_name"
	@exit 1
endif
	@go run cmd/migrations/main.go -command=create -name=$(name)
