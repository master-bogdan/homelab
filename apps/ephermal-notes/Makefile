# ------------------------
# Project variables
# ------------------------
APP_NAME      := ephermal-notes
BINARY_NAME   := $(APP_NAME)
IMAGE_NAME    := $(APP_NAME):latest
DOCKER_IMAGE  := masterbogdan0/$(APP_NAME):latest

K8S_NAMESPACE_APP := dev
K8S_NAMESPACE_DB  := db
K8S_NAMESPACE_OBS := observability

GOBASE        := $(shell pwd)
GOBIN         := $(GOBASE)/bin

.PHONY: all build clean run test lint docker-build docker-push \
        k8s-deploy-app k8s-deploy-db k8s-deploy-observability \
        k8s-deploy k8s-delete minikube-deploy

# ------------------------
# Go targets
# ------------------------

all: clean build

build:
	@echo "➡️  Building Go binary..."
	@go build -o $(GOBIN)/$(BINARY_NAME) ./...

clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(GOBIN)

run:
	@echo "Starting Redis for local development..."
	@docker-compose -f docker-compose.dev.yaml up -d
	@sleep 3 # wait for Redis to start
	@echo "Running Go application..."
	@go run ./...
	@echo "Stopping Redis..."
	@docker-compose -f docker-compose.dev.yaml down

lint:
	@echo "🔍 Running linter..."
	@golangci-lint run

# -----------------------------
# Tests
# -----------------------------
test:
	@echo "Running tests..."
	@go test -v ./...

# Local tests with Redis container
test-local:
	@echo "Starting Redis for tests..."
	@docker-compose -f docker-compose.test.yaml up -d
	@sleep 3 # wait for Redis to start
	@go test -v ./...
	@docker-compose -f docker-compose.test.yaml down

# ------------------------
# Docker targets
# ------------------------
docker-build:
	@echo "🐳 Building Docker image..."
	@docker build -t $(IMAGE_NAME) -f ./deployments/docker/Dockerfile .

docker-push:
	@echo "📤 Pushing Docker image..."
	@docker tag $(IMAGE_NAME) $(DOCKER_IMAGE)
	@docker push $(DOCKER_IMAGE)

# ------------------------
# Kubernetes targets
# ------------------------
k8s-deploy-app:
	@echo "📦 Creating namespace $(K8S_NAMESPACE_APP)..."
	minikube kubectl -- create namespace $(K8S_NAMESPACE_APP) || echo "Namespace already exists"
	@echo "🚀 Deploying APP to Kubernetes..."
	minikube kubectl -- apply -R -f ./deployments/k8s/app -n $(K8S_NAMESPACE_APP)
	@echo "📤 Restarting pods..."
	minikube kubectl -- rollout restart deployment notes-api-deploy -n $(K8S_NAMESPACE_APP)

k8s-deploy-db:
	@echo "📦 Creating namespace $(K8S_NAMESPACE_DB)..."
	minikube kubectl -- create namespace $(K8S_NAMESPACE_DB) || echo "Namespace already exists"
	@echo "🚀 Deploying DB to Kubernetes..."
	minikube kubectl -- apply -R -f ./deployments/k8s/db -n $(K8S_NAMESPACE_DB)
	@echo "📤 Restarting pods..."
	minikube kubectl -- rollout restart deployment redis-deploy -n $(K8S_NAMESPACE_DB)

k8s-deploy-observability:
	@echo "📦 Creating namespace $(K8S_NAMESPACE_OBS)..."
	minikube kubectl -- create namespace $(K8S_NAMESPACE_OBS) || echo "Namespace already exists"
	@echo "🚀 Deploying Observability to Kubernetes..."
	minikube kubectl -- apply -R -f ./deployments/k8s/observability -n $(K8S_NAMESPACE_OBS)

k8s-deploy: k8s-deploy-db k8s-deploy-app k8s-deploy-observability

k8s-delete:
	@echo "🔥 Deleting from Kubernetes..."
	minikube kubectl -- delete -R -f ./deployments/k8s/ -n $(K8S_NAMESPACE_APP)

# ------------------------
# Minikube helpers (simulate prod)
# ------------------------
minikube-deploy: docker-build docker-push
	@echo "🚀 Deploying to local Minikube"
	$(MAKE) k8s-deploy
