.PHONY: build run test lint compose-up compose-down run-compose migrate-up migrate-down migrate-create docs-gen swagger-gen

# Build the API binary
build:
	@mkdir -p bin
	@go build -o bin/estimate-room-api ./cmd/server

# Run the API locally
run:
	@go run ./cmd/server

# Run all tests
test:
	@go test ./...

# Lint (uses golangci-lint if available, otherwise go vet)
lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || (echo "golangci-lint not installed; running go vet" && go vet ./...)

# Start only the required services via docker compose
compose-up:
	@docker compose -f docker-compose.dev.yaml up -d

# Stop services started by docker compose
compose-down:
	@docker compose -f docker-compose.dev.yaml down

# Run required services via docker compose, then start the API locally
run-compose: compose-up
	@go run ./cmd/server

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

# Generate Swagger/OpenAPI docs
swagger-generate:
	@go run cmd/docs/main.go
