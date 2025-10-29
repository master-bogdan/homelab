# ============================
# ephermal notes api
# ============================

APP_NAME    ?= ephermal-notes-api
BINARY_NAME ?= $(APP_NAME)
GOBASE      := $(shell pwd)
GOBIN       := $(GOBASE)/bin

.PHONY: help build clean run lint test test-local

help:
	@echo "🧰 $(APP_NAME) — local dev"
	@echo "  build        🔨  Build ./bin/$(BINARY_NAME)"
	@echo "  clean        🧹  Remove ./bin"
	@echo "  run          ▶️   Start Redis (docker-compose.dev.yaml) then go run"
	@echo "  lint         🔍  golangci-lint run"
	@echo "  test         ✅  go test -v ./..."
	@echo "  test-local   🧪  Start Redis (test), run tests, stop Redis"

build:
	@echo "🔨 Building Go binary → $(GOBIN)/$(BINARY_NAME)"
	@mkdir -p $(GOBIN)
	@go build -o $(GOBIN)/$(BINARY_NAME) ./...

clean:
	@echo "🧹 Cleaning build artifacts"
	@rm -rf $(GOBIN)

run:
	@echo "🐘 Starting Redis for local development…"
	@docker-compose -f docker-compose.dev.yaml up -d
	@sleep 3
	@echo "▶️  Running Go application…"
	@go run ./...
	@echo "🛑 Stopping Redis…"
	@docker-compose -f docker-compose.dev.yaml down

lint:
	@echo "🔍 Running linter…"
	@golangci-lint run

test:
	@echo "✅ Running tests…"
	@go test -v ./...

test-local:
	@echo "🧪 Spinning up Redis for tests…"
	@docker-compose -f docker-compose.test.yaml up -d
	@sleep 3
	@echo "🧪 Running tests…"
	@go test -v ./...
	@echo "🧹 Tearing down test Redis…"
	@docker-compose -f docker-compose.test.yaml down
