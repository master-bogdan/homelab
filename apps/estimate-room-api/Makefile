.PHONY: build run dev test lint compose-up compose-down run-compose dev-compose migrate-up migrate-down migrate-create swagger-generate

# Build the API binary
build:
	@mkdir -p bin
	@go build -o bin/estimate-room-api ./cmd/server

# Run the API locally
run:
	@go run ./cmd/server

# Run the API locally with live reload via air
dev:
	@AIR_BIN="$$(command -v air || true)"; \
	if [ -z "$$AIR_BIN" ]; then \
		GOBIN="$$(go env GOBIN)"; \
		GOPATH="$$(go env GOPATH)"; \
		if [ -n "$$GOBIN" ] && [ -x "$$GOBIN/air" ]; then \
			AIR_BIN="$$GOBIN/air"; \
		elif [ -n "$$GOPATH" ] && [ -x "$$GOPATH/bin/air" ]; then \
			AIR_BIN="$$GOPATH/bin/air"; \
		else \
			echo "air is not installed. Install it with: go install github.com/air-verse/air@latest"; \
			exit 1; \
		fi; \
	fi; \
	"$$AIR_BIN" -c .air.toml

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

# Run required services via docker compose, then start the API with live reload
dev-compose: compose-up
	@$(MAKE) dev

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
