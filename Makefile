# ============================
# Makefile (root)
# ============================

# ---- Config (override via env/CLI) ----
REGISTRY          ?= docker.io/masterbogdan0
TAG               ?= latest
ENV               ?= dev              # dev | prod
APP               ?=                  # single app name, e.g. personal-website-ui
MINIKUBE          ?=                  # set 1 to use "minikube kubectl --"
MINIKUBE_PROFILE  ?= homelab-dev      # profile name for minikube

# ---- Paths ----
APPS_DIR              := apps
KUBERNETES_DIR        := k8s

NAMESPACES_DIR        := $(KUBERNETES_DIR)/namespaces

NETWORKING_DIR        := $(KUBERNETES_DIR)/networking
TRAEFIK_DIR           := $(NETWORKING_DIR)/traefik
GATEWAY_DIR           := $(NETWORKING_DIR)/gateway

DATABASES_DIR         := $(KUBERNETES_DIR)/databases
REDIS_DIR             := $(DATABASES_DIR)/redis
POSTGRESQL_DIR        := $(DATABASES_DIR)/postgresql

PLATFORM_DIR          := $(KUBERNETES_DIR)/platform
AUTHENTIK_DIR         := $(PLATFORM_DIR)/authentik
N8N_DIR               := $(PLATFORM_DIR)/n8n
SEAWEEDFS_DIR         := $(PLATFORM_DIR)/seaweedfs

OBSERVABILITY_DIR     := $(KUBERNETES_DIR)/observability
FLUENTBIT_DIR         := $(OBSERVABILITY_DIR)/fluent-bit
GRAFANA_DIR           := $(OBSERVABILITY_DIR)/grafana
PROMETHEUS_DIR        := $(OBSERVABILITY_DIR)/prometheus
OPENSEARCH_DIR        := $(OBSERVABILITY_DIR)/opensearch
OPENSEARCH_DASH_DIR   := $(OBSERVABILITY_DIR)/opensearch-dashboards

# ---- Discover apps ----
APPS                  := $(notdir $(wildcard $(APPS_DIR)/*))

# ---- kubectl (optionally via minikube) ----
KUBECTL               := $(if $(MINIKUBE),minikube kubectl --,kubectl)

# ---- Helpers: prefer overlays/$(ENV), fallback to base ----
define k8s_apply_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🚀 Applying: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) apply -k "$${ROOT}/overlays/$(ENV)"; \
	else \
	  echo "🚀 Applying: $${ROOT}/base"; \
	  $(KUBECTL) apply -k "$${ROOT}/base" 2>/dev/null || \
	  $(KUBECTL) apply -f "$${ROOT}/base"; \
	fi
endef

define k8s_delete_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🔥 Deleting: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) delete -k "$${ROOT}/overlays/$(ENV)" --ignore-not-found; \
	else \
	  echo "🔥 Deleting: $${ROOT}/base"; \
	  $(KUBECTL) delete -k "$${ROOT}/base" --ignore-not-found 2>/dev/null || \
	  $(KUBECTL) delete -f "$${ROOT}/base" --ignore-not-found; \
	fi
endef

# Resolve Dockerfile per app:
# 1) apps/<app>/deployments/docker/Dockerfile
# 2) apps/<app>/Dockerfile
define resolve_dockerfile
	APP_NAME="$(1)"; \
	if [ -f "$(APPS_DIR)/$$APP_NAME/deployments/docker/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP_NAME/deployments/docker/Dockerfile"; \
	elif [ -f "$(APPS_DIR)/$$APP_NAME/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP_NAME/Dockerfile"; \
	else \
	  echo "❌ No Dockerfile for $$APP_NAME" >&2; exit 1; \
	fi
endef

# ----------------------------
# PHONY
# ----------------------------
.PHONY: help apps-list \
        docker-build docker-push docker-build-push \
        docker-build-all docker-push-all docker-build-push-all \
        minikube-start minikube-stop minikube-delete \
        deploy-namespaces delete-namespaces \
        deploy-app deploy-apps delete-app delete-apps \
        deploy-networking delete-networking \
        deploy-platform delete-platform \
        deploy-observability delete-observability \
        deploy-databases delete-databases \
        deploy-all deploy-all-dev deploy-all-prod \
        delete-all delete-all-dev delete-all-prod

# ============================
# Help
# ============================
help:
	@echo "🎛️  Root targets (ENV=$(ENV))"
	@echo "  apps-list                          # List discovered apps"
	@echo ""
	@echo "🐳 Docker"
	@echo "  docker-build        APP=<name>     # Build one image"
	@echo "  docker-push         APP=<name>     # Push one image"
	@echo "  docker-build-push   APP=<name>     # Build + push one image"
	@echo "  docker-build-all                   # Build all images"
	@echo "  docker-push-all                    # Push all images"
	@echo "  docker-build-push-all              # Build + push all images"
	@echo ""
	@echo "☸️  Minikube"
	@echo "  MINIKUBE=1 minikube-start          # Start 3-node minikube cluster"
	@echo "  MINIKUBE=1 minikube-stop           # Stop minikube cluster"
	@echo "  MINIKUBE=1 minikube-delete         # Delete minikube cluster"
	@echo ""
	@echo "☸️  Kubernetes (building blocks)"
	@echo "  deploy-namespaces                  # Apply namespaces"
	@echo "  delete-namespaces                  # Delete namespaces"
	@echo "  deploy-networking                  # Traefik + Gateway"
	@echo "  delete-networking                  # Delete networking"
	@echo "  deploy-platform                    # Authentik, n8n, SeaweedFS"
	@echo "  delete-platform                    # Delete platform services"
	@echo "  deploy-observability               # Fluent Bit, Prometheus, Grafana, OpenSearch, Dashboards"
	@echo "  delete-observability               # Delete observability"
	@echo "  deploy-databases                   # Redis + PostgreSQL (dev only)"
	@echo "  delete-databases                   # Delete databases (dev only)"
	@echo "  deploy-app         APP=<name>      # Build+push APP, then apply k8s for APP"
	@echo "  deploy-apps                        # Build+push ALL, then apply k8s for ALL"
	@echo "  delete-app         APP=<name>      # Delete k8s for APP"
	@echo "  delete-apps                        # Delete k8s for ALL"
	@echo ""
	@echo "🎯 High-level"
	@echo "  deploy-all ENV=dev                 # dev: namespaces, net, DBs, platform, obs, apps"
	@echo "  deploy-all ENV=prod                # prod: namespaces, net, platform, obs, apps (no DBs)"
	@echo "  delete-all ENV=dev                 # Full teardown dev (with DBs)"
	@echo "  delete-all ENV=prod                # Full teardown prod (without DBs)"
	@echo ""
	@echo "🔧 Vars: REGISTRY=$(REGISTRY) TAG=$(TAG) ENV=$(ENV) MINIKUBE=$(MINIKUBE) MINIKUBE_PROFILE=$(MINIKUBE_PROFILE)"

apps-list:
	@echo "$(APPS)"

# ============================
# Docker — NO pattern rules
# ============================
docker-build:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@DFILE="$$( $(call resolve_dockerfile,$(APP)) )"; \
	echo "🐳 Building $(REGISTRY)/$(APP):$(TAG) using $$DFILE"; \
	docker build -t "$(REGISTRY)/$(APP):$(TAG)" -f "$$DFILE" "$(APPS_DIR)/$(APP)"

docker-push:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@echo "📤 Pushing $(REGISTRY)/$(APP):$(TAG)"; \
	docker push "$(REGISTRY)/$(APP):$(TAG)"

docker-build-push:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@$(MAKE) --no-print-directory docker-build APP=$(APP)
	@$(MAKE) --no-print-directory docker-push  APP=$(APP)

docker-build-all:
	@for app in $(APPS); do \
	  echo "🐳 Building $(REGISTRY)/$$app:$(TAG)"; \
	  APP="$$app" $(MAKE) --no-print-directory docker-build; \
	done
	@echo "✅ Built images for all apps"

docker-push-all:
	@for app in $(APPS); do \
	  echo "📤 Pushing $(REGISTRY)/$$app:$(TAG)"; \
	  APP="$$app" $(MAKE) --no-print-directory docker-push; \
	done
	@echo "✅ Pushed images for all apps"

docker-build-push-all:
	@for app in $(APPS); do \
	  echo "🚀 Build+push image for $$app"; \
	  APP="$$app" $(MAKE) --no-print-directory docker-build-push; \
	done
	@echo "✅ Build+push for all apps"

# ============================
# Minikube — 3-node cluster (dev)
# ============================
minikube-start:
	@echo "🚀 Starting 3-node minikube cluster (profile=$(MINIKUBE_PROFILE))..."
	minikube start --profile $(MINIKUBE_PROFILE) \
	  --nodes=3 \
	  --driver=docker \
	  --cpus=4 \
	  --memory=8192
	@echo "🔌 Enabling volume-related addons..."
	minikube --profile $(MINIKUBE_PROFILE) addons enable volumesnapshots
	minikube --profile $(MINIKUBE_PROFILE) addons enable storage-provisioner
	minikube --profile $(MINIKUBE_PROFILE) addons enable csi-hostpath-driver
	@echo "✅ Minikube up. Use: kubectl config use-context $(MINIKUBE_PROFILE)"

minikube-stop:
	@echo "⏸️  Stopping minikube (profile=$(MINIKUBE_PROFILE))..."
	minikube stop --profile $(MINIKUBE_PROFILE) || true

minikube-delete:
	@echo "🗑️  Deleting minikube cluster (profile=$(MINIKUBE_PROFILE))..."
	minikube delete --profile $(MINIKUBE_PROFILE) || true

# ============================
# Kubernetes — Namespaces
# ============================
deploy-namespaces:
	$(call k8s_apply_or_base,$(NAMESPACES_DIR))
	@echo "✅ Namespaces applied for $(ENV)"

delete-namespaces:
	$(call k8s_delete_or_base,$(NAMESPACES_DIR))
	@echo "🗑️  Namespaces deleted for $(ENV)"

# ============================
# Kubernetes — Apps (build+deploy)
# ============================
deploy-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@$(MAKE) --no-print-directory docker-build-push APP=$(APP)
	$(call k8s_apply_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ Deployed app: $(APP) (env=$(ENV))"

deploy-apps:
	@for app in $(APPS); do \
	  echo "🚀 Deploying app $$app (env=$(ENV))"; \
	  APP="$$app" $(MAKE) --no-print-directory docker-build-push; \
	  APP="$$app" $(call k8s_apply_or_base,$(KUBERNETES_DIR)/apps/$$app); \
	done
	@echo "🎉 All apps deployed (env=$(ENV))"

delete-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	$(call k8s_delete_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "🗑️  Deleted app: $(APP) (env=$(ENV))"

delete-apps:
	@for app in $(APPS); do \
	  echo "🧹 Deleting app $$app (env=$(ENV))"; \
	  APP="$$app" $(call k8s_delete_or_base,$(KUBERNETES_DIR)/apps/$$app); \
	done
	@echo "🧹 All apps deleted (env=$(ENV))"

# ============================
# Kubernetes — Networking (Traefik + Gateway)
# ============================
deploy-networking: deploy-namespaces
	$(call k8s_apply_or_base,$(TRAEFIK_DIR))
	$(call k8s_apply_or_base,$(GATEWAY_DIR))
	@echo "✅ Networking (Traefik + Gateway) deployed (env=$(ENV))"

delete-networking:
	$(call k8s_delete_or_base,$(GATEWAY_DIR))
	$(call k8s_delete_or_base,$(TRAEFIK_DIR))
	@echo "🗑️  Networking (Traefik + Gateway) deleted (env=$(ENV))"

# ============================
# Kubernetes — Platform
# ============================
deploy-platform: deploy-namespaces
	$(call k8s_apply_or_base,$(AUTHENTIK_DIR))
	$(call k8s_apply_or_base,$(N8N_DIR))
	$(call k8s_apply_or_base,$(SEAWEEDFS_DIR))
	@echo "🎉 Platform services deployed (env=$(ENV))"

delete-platform:
	$(call k8s_delete_or_base,$(SEAWEEDFS_DIR))
	$(call k8s_delete_or_base,$(N8N_DIR))
	$(call k8s_delete_or_base,$(AUTHENTIK_DIR))
	@echo "🧹 Platform services deleted (env=$(ENV))"

# ============================
# Kubernetes — Observability
# ============================
deploy-observability: deploy-namespaces
	$(call k8s_apply_or_base,$(FLUENTBIT_DIR))
	$(call k8s_apply_or_base,$(PROMETHEUS_DIR))
	$(call k8s_apply_or_base,$(GRAFANA_DIR))
	$(call k8s_apply_or_base,$(OPENSEARCH_DIR))
	$(call k8s_apply_or_base,$(OPENSEARCH_DASH_DIR))
	@echo "📊 Observability stack deployed (env=$(ENV))"

delete-observability:
	$(call k8s_delete_or_base,$(OPENSEARCH_DASH_DIR))
	$(call k8s_delete_or_base,$(OPENSEARCH_DIR))
	$(call k8s_delete_or_base,$(GRAFANA_DIR))
	$(call k8s_delete_or_base,$(PROMETHEUS_DIR))
	$(call k8s_delete_or_base,$(FLUENTBIT_DIR))
	@echo "🧻 Observability stack deleted (env=$(ENV))"

# ============================
# Kubernetes — Databases (dev-only)
# ============================
deploy-databases: deploy-namespaces
ifeq ($(ENV),dev)
	$(call k8s_apply_or_base,$(REDIS_DIR))
	$(call k8s_apply_or_base,$(POSTGRESQL_DIR))
	@echo "🧰 Databases deployed (env=$(ENV))"
else
	@echo "⚠️  deploy-databases skipped: ENV=$(ENV) (no DBs in prod)"
endif

delete-databases:
ifeq ($(ENV),dev)
	$(call k8s_delete_or_base,$(POSTGRESQL_DIR))
	$(call k8s_delete_or_base,$(REDIS_DIR))
	@echo "🧨 Databases deleted (env=$(ENV))"
else
	@echo "⚠️  delete-databases skipped: ENV=$(ENV) (no DBs in prod)"
endif

# ============================
# Everything (dev/prod)
# ============================
deploy-all-dev: docker-build-push-all deploy-namespaces deploy-networking deploy-databases deploy-platform deploy-observability deploy-apps
	@echo "✅ Full stack applied for dev (with databases)"

deploy-all-prod: docker-build-push-all deploy-namespaces deploy-networking deploy-platform deploy-observability deploy-apps
	@echo "✅ Full stack applied for prod (no databases)"

deploy-all:
ifeq ($(ENV),dev)
	@$(MAKE) --no-print-directory ENV=dev deploy-all-dev
else ifeq ($(ENV),prod)
	@$(MAKE) --no-print-directory ENV=prod deploy-all-prod
else
	@echo "❌ Unknown ENV=$(ENV). Use dev|prod"; exit 1
endif

delete-all-dev: delete-apps delete-observability delete-platform delete-databases delete-networking delete-namespaces
	@echo "✅ Full stack deleted for dev (with databases)"

delete-all-prod: delete-apps delete-observability delete-platform delete-networking delete-namespaces
	@echo "✅ Full stack deleted for prod (no databases)"

delete-all:
ifeq ($(ENV),dev)
	@$(MAKE) --no-print-directory ENV=dev delete-all-dev
else ifeq ($(ENV),prod)
	@$(MAKE) --no-print-directory ENV=prod delete-all-prod
else
	@echo "❌ Unknown ENV=$(ENV). Use dev|prod"; exit 1
endif
