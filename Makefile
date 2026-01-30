# ============================
# Makefile (Optimized)
# ============================

# ---- Config ----
REGISTRY          ?= docker.io/masterbogdan0
TAG               ?= latest
ENV               ?= dev
APP               ?=
# Optional fallback overlay for mixed environments (e.g., staging -> prod)
ENV_FALLBACK      ?= $(if $(filter staging,$(ENV)),prod,)
ROOT_DIR          := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
SCRIPTS_DIR       := $(ROOT_DIR)/scripts

# ---- UI Static Site Config ----
UI_NAME_PATTERN         ?= ui
UI_BUILD_OUTPUT_DIRS    ?= out build dist
UI_PUBLIC_BUCKET        ?= public
SEAWEEDFS_NAMESPACE     ?= platform
SEAWEED_PUBLIC_BASE_DEV      ?= http://storage.apps.192-168-58-2.sslip.io:30080
SEAWEED_PUBLIC_BASE_PROD     ?= https://storage.example.invalid
SEAWEED_PUBLIC_BASE_STAGING ?= https://storage.apps.192-168-58-2.sslip.io
ifeq ($(ENV),prod)
SEAWEED_PUBLIC_BASE     ?= $(SEAWEED_PUBLIC_BASE_PROD)
else ifeq ($(ENV),staging)
SEAWEED_PUBLIC_BASE     ?= $(SEAWEED_PUBLIC_BASE_STAGING)
else
SEAWEED_PUBLIC_BASE     ?= $(SEAWEED_PUBLIC_BASE_DEV)
endif
SEAWEED_INTERNAL_BASE   ?= http://seaweedfs-s3.$(SEAWEEDFS_NAMESPACE).svc.cluster.local:8333
SEAWEEDFS_MASTER_ADDR   ?= seaweedfs-master.$(SEAWEEDFS_NAMESPACE):9333
SEAWEEDFS_FILER_ADDR    ?= seaweedfs-filer-client.$(SEAWEEDFS_NAMESPACE):8888
SEAWEED_ANON_ACTIONS    ?= Read,Write
SEAWEED_UPLOAD_IMAGE    ?= curlimages/curl:8.5.0
SEAWEED_UPLOAD_SETUP_CMD ?= true
SEAWEED_UPLOAD_POD_PREFIX ?= seaweedfs-uploader
SEAWEED_UPLOAD_WAIT     ?= 90s
SEAWEED_UPLOAD_CONNECT_TIMEOUT ?= 5
SEAWEED_UPLOAD_MAX_TIME ?= 60

# UI base path passed into build
UI_BASE_PATH            ?= /$(UI_PUBLIC_BUCKET)/$(APP)

# ---- Paths ----
APPS_DIR              := $(ROOT_DIR)/apps
KUBERNETES_DIR        := $(ROOT_DIR)/k8s
NAMESPACES_DIR        := $(KUBERNETES_DIR)/namespaces
NETWORKING_DIR        := $(KUBERNETES_DIR)/networking
TRAEFIK_DIR           := $(NETWORKING_DIR)/traefik
GATEWAY_DIR           := $(NETWORKING_DIR)/gateway
DATABASES_DIR         := $(KUBERNETES_DIR)/databases
REDIS_DIR             := $(DATABASES_DIR)/redis
POSTGRESQL_DIR        := $(DATABASES_DIR)/postgresql
PLATFORM_DIR          := $(KUBERNETES_DIR)/platform
AUTH_DIR              := $(KUBERNETES_DIR)/auth
SECRETS_DIR           := $(KUBERNETES_DIR)/secrets
AUTHENTIK_DIR         := $(AUTH_DIR)/authentik
AUTHENTIK_FWD_AUTH_DIR := $(AUTH_DIR)/forward-auth
AUTH_REFERENCE_GRANT  := $(AUTH_DIR)/reference-grant.yaml
EXTERNAL_SECRETS_DIR  := $(SECRETS_DIR)/external-secrets
VAULT_DIR             := $(SECRETS_DIR)/vault
N8N_DIR               := $(PLATFORM_DIR)/n8n
SEAWEEDFS_DIR          := $(PLATFORM_DIR)/seaweedfs
OBSERVABILITY_DIR     := $(KUBERNETES_DIR)/observability
FLUENTBIT_DIR         := $(OBSERVABILITY_DIR)/fluent-bit
GRAFANA_DIR           := $(OBSERVABILITY_DIR)/grafana
PROMETHEUS_DIR        := $(OBSERVABILITY_DIR)/prometheus
OPENSEARCH_DIR        := $(OBSERVABILITY_DIR)/opensearch
OPENSEARCH_DASH_DIR   := $(OBSERVABILITY_DIR)/opensearch-dashboards

# ---- Layer Groups ----
NETWORKING_DIRS       := $(GATEWAY_DIR)
NETWORKING_HELM_DIRS  := $(TRAEFIK_DIR)
AUTH_DIRS             := $(AUTHENTIK_FWD_AUTH_DIR)
AUTH_HELM_DIRS        := $(AUTHENTIK_DIR)
PLATFORM_HELM_DIRS    := $(AUTHENTIK_DIR) $(N8N_DIR) $(SEAWEEDFS_DIR)
SECRETS_HELM_DIRS     := $(VAULT_DIR) $(EXTERNAL_SECRETS_DIR)
OBS_HELM_DIRS         := $(FLUENTBIT_DIR) $(PROMETHEUS_DIR) $(GRAFANA_DIR) $(OPENSEARCH_DIR) $(OPENSEARCH_DASH_DIR)
OBS_HELM_DELETE_DIRS  := $(OPENSEARCH_DASH_DIR) $(OPENSEARCH_DIR) $(GRAFANA_DIR) $(PROMETHEUS_DIR) $(FLUENTBIT_DIR)
DATABASE_DIRS         := $(REDIS_DIR) $(POSTGRESQL_DIR)
DATABASE_DELETE_DIRS  := $(POSTGRESQL_DIR) $(REDIS_DIR)

# ---- Discovery ----
APPS                  := $(notdir $(wildcard $(APPS_DIR)/*))
KUBECTL               := kubectl
export APPS_DIR KUBECTL UI_PUBLIC_BUCKET UI_BUILD_OUTPUT_DIRS SEAWEEDFS_NAMESPACE \
	SEAWEED_PUBLIC_BASE SEAWEED_INTERNAL_BASE SEAWEEDFS_MASTER_ADDR SEAWEEDFS_FILER_ADDR \
	SEAWEED_ANON_ACTIONS SEAWEED_UPLOAD_IMAGE SEAWEED_UPLOAD_SETUP_CMD SEAWEED_UPLOAD_POD_PREFIX \
	SEAWEED_UPLOAD_WAIT SEAWEED_UPLOAD_CONNECT_TIMEOUT SEAWEED_UPLOAD_MAX_TIME APP

# ============================
# Helper Functions
# ============================

# Check if app is a UI app (contains pattern and no Dockerfile)
define is_ui_app
$(shell \
	if printf "%s" "$(1)" | grep -qi "$(UI_NAME_PATTERN)"; then \
		if [ ! -f "$(APPS_DIR)/$(1)/Dockerfile" ] && \
		   [ ! -f "$(APPS_DIR)/$(1)/deployments/docker/Dockerfile" ]; then \
			echo "yes"; \
		fi; \
	fi \
)
endef

# Get all UI apps
UI_APPS := $(foreach app,$(APPS),$(if $(call is_ui_app,$(app)),$(app)))

# Resolve Dockerfile path
define resolve_dockerfile
$(shell \
	if [ -f "$(APPS_DIR)/$(1)/deployments/docker/Dockerfile" ]; then \
		echo "$(APPS_DIR)/$(1)/deployments/docker/Dockerfile"; \
	elif [ -f "$(APPS_DIR)/$(1)/Dockerfile" ]; then \
		echo "$(APPS_DIR)/$(1)/Dockerfile"; \
	fi \
)
endef

# Require APP to be set
define require_app
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
endef

# Guard for k8s manifests on APP
define require_app_k8s
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then \
		echo "⚠️  No k8s manifests for $(APP)"; exit 0; \
	fi
endef

# Generic function to iterate over apps with a target
define for_each_app
	@for app in $(APPS); do \
		APP="$$app" $(MAKE) --no-print-directory $(1); \
	done
endef

# Generic function to iterate over apps with k8s manifests
define for_each_k8s_app
	@for app in $(APPS); do \
		if [ -d "$(KUBERNETES_DIR)/apps/$$app" ]; then \
			APP="$$app" $(MAKE) --no-print-directory $(1); \
		fi; \
	done
endef

# Generic function to iterate over UI apps with a target
define for_each_ui_app
	@for app in $(UI_APPS); do \
		APP="$$app" $(MAKE) --no-print-directory $(1); \
	done
endef

# Kustomize apply (prefer overlay, fallback to base)
define k8s_apply_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🚀 Applying: $$OVERLAY"; \
		$(KUBECTL) apply -k "$$OVERLAY"; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🚀 Applying: $$FALLBACK (fallback)"; \
		$(KUBECTL) apply -k "$$FALLBACK"; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
		echo "🚀 Applying: $${ROOT}"; \
		$(KUBECTL) apply -k "$${ROOT}"; \
	else \
		echo "🚀 Applying: $${ROOT}/base"; \
		$(KUBECTL) apply -k "$${ROOT}/base" 2>/dev/null || \
		$(KUBECTL) apply -f "$${ROOT}/base"; \
	fi
endef

# Kustomize delete
define k8s_delete_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🔥 Deleting: $$OVERLAY"; \
		$(KUBECTL) delete -k "$$OVERLAY" --ignore-not-found; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🔥 Deleting: $$FALLBACK (fallback)"; \
		$(KUBECTL) delete -k "$$FALLBACK" --ignore-not-found; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
		echo "🔥 Deleting: $${ROOT}"; \
		$(KUBECTL) delete -k "$${ROOT}" --ignore-not-found; \
	else \
		echo "🔥 Deleting: $${ROOT}/base"; \
		$(KUBECTL) delete -k "$${ROOT}/base" --ignore-not-found 2>/dev/null || \
		$(KUBECTL) delete -f "$${ROOT}/base" --ignore-not-found; \
	fi
endef

# Helm+Kustomize apply
define k8s_apply_helm_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🚀 Applying (helm+kustomize): $$OVERLAY"; \
		$(KUBECTL) kustomize --enable-helm "$$OVERLAY" | $(KUBECTL) apply -f -; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🚀 Applying (helm+kustomize): $$FALLBACK (fallback)"; \
		$(KUBECTL) kustomize --enable-helm "$$FALLBACK" | $(KUBECTL) apply -f -; \
	else \
		echo "🚀 Applying (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) apply -f -; \
	fi
endef

# Helm+Kustomize delete
define k8s_delete_helm_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🔥 Deleting (helm+kustomize): $$OVERLAY"; \
		$(KUBECTL) kustomize --enable-helm "$$OVERLAY" | $(KUBECTL) delete --ignore-not-found -f -; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🔥 Deleting (helm+kustomize): $$FALLBACK (fallback)"; \
		$(KUBECTL) kustomize --enable-helm "$$FALLBACK" | $(KUBECTL) delete --ignore-not-found -f -; \
	else \
		echo "🔥 Deleting (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) delete --ignore-not-found -f -; \
	fi
endef

# Validate kustomize (dry-run)
define k8s_dry_run_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🧪 Dry-run: $$OVERLAY"; \
		$(KUBECTL) apply -k "$$OVERLAY" --dry-run=client --validate=false -o yaml >/dev/null; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🧪 Dry-run: $$FALLBACK (fallback)"; \
		$(KUBECTL) apply -k "$$FALLBACK" --dry-run=client --validate=false -o yaml >/dev/null; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
		echo "🧪 Dry-run: $${ROOT}"; \
		$(KUBECTL) apply -k "$${ROOT}" --dry-run=client --validate=false -o yaml >/dev/null; \
	else \
		echo "🧪 Dry-run: $${ROOT}/base"; \
		$(KUBECTL) apply -k "$${ROOT}/base" --dry-run=client --validate=false -o yaml >/dev/null 2>/dev/null || \
		$(KUBECTL) apply -f "$${ROOT}/base" --dry-run=client --validate=false -o yaml >/dev/null; \
	fi
endef

# Validate helm+kustomize (build-only)
define k8s_dry_run_helm_or_base
	@set -e; ROOT="$(1)"; \
	OVERLAY="$${ROOT}/overlays/$(ENV)"; \
	FALLBACK="$${ROOT}/overlays/$(ENV_FALLBACK)"; \
	if [ -d "$$OVERLAY" ]; then \
		echo "🧪 Build-only (helm+kustomize): $$OVERLAY"; \
		$(KUBECTL) kustomize --enable-helm "$$OVERLAY" >/dev/null; \
	elif [ -n "$(ENV_FALLBACK)" ] && [ -d "$$FALLBACK" ]; then \
		echo "🧪 Build-only (helm+kustomize): $$FALLBACK (fallback)"; \
		$(KUBECTL) kustomize --enable-helm "$$FALLBACK" >/dev/null; \
	else \
		echo "🧪 Build-only (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" >/dev/null; \
	fi
endef

define k8s_apply_list
	@for dir in $(1); do \
		$(call k8s_apply_or_base,$$dir); \
	done
endef

define k8s_apply_helm_list
	@for dir in $(1); do \
		$(call k8s_apply_helm_or_base,$$dir); \
	done
endef

define k8s_delete_list
	@for dir in $(1); do \
		$(call k8s_delete_or_base,$$dir); \
	done
endef

define k8s_delete_helm_list
	@for dir in $(1); do \
		$(call k8s_delete_helm_or_base,$$dir); \
	done
endef

define k8s_dry_run_list
	@for dir in $(1); do \
		$(call k8s_dry_run_or_base,$$dir); \
	done
endef

define k8s_dry_run_helm_list
	@for dir in $(1); do \
		$(call k8s_dry_run_helm_or_base,$$dir); \
	done
endef

# ============================
# PHONY Declarations
# ============================
.PHONY: help apps-list ui-apps-list \
	docker-build docker-push docker-build-push \
	docker-build-all docker-push-all docker-build-push-all \
	deploy-ui-build deploy-ui-build-all deploy-ui-upload deploy-ui-upload-all \
	deploy-ui deploy-ui-all deploy-ui-bucket deploy-ui-url \
	deploy-namespaces delete-namespaces validate-namespaces \
	deploy-app deploy-apps delete-app delete-apps validate-app validate-apps \
	deploy-networking delete-networking validate-networking \
	deploy-secrets delete-secrets validate-secrets \
	deploy-auth delete-auth validate-auth \
	deploy-platform delete-platform validate-platform \
	deploy-observability delete-observability validate-observability \
	deploy-databases delete-databases validate-databases \
	validate-all clean-charts \
	deploy-all delete-all

# ============================
# Help
# ============================
help:
	@echo "🎛  Root targets (ENV=$(ENV), TAG=$(TAG))"
	@echo ""
	@echo "📋 Discovery"
	@echo "  apps-list                          # List all apps"
	@echo "  ui-apps-list                       # List UI apps (static sites)"
	@echo ""
	@echo "🐳 Docker"
	@echo "  docker-build        APP=<name>     # Build one image"
	@echo "  docker-push         APP=<name>     # Push one image"
	@echo "  docker-build-push   APP=<name>     # Build + push one"
	@echo "  docker-build-all                   # Build all images"
	@echo "  docker-push-all                    # Push all images"
	@echo "  docker-build-push-all              # Build + push all"
	@echo ""
	@echo "🎨 UI Static Sites (SeaweedFS)"
	@echo "  deploy-ui-build     APP=<name>     # Build UI app"
	@echo "  deploy-ui-bucket                  # Ensure public bucket exists"
	@echo "  deploy-ui-upload    APP=<name>     # Upload to SeaweedFS (in-cluster)"
	@echo "  deploy-ui-url       APP=<name>     # Print public URL"
	@echo "  deploy-ui           APP=<name>     # Build + upload"
	@echo "  deploy-ui-build-all                # Build all UI apps"
	@echo "  deploy-ui-upload-all               # Upload all UI apps"
	@echo "  deploy-ui-all                      # Build + upload all"
	@echo ""
	@echo "☸  Kubernetes (Layer-based)"
	@echo "  ENV=staging uses overlays/staging when present, otherwise falls back to prod"
	@echo "  deploy-namespaces   / delete-namespaces   / validate-namespaces"
	@echo "  deploy-networking   / delete-networking   / validate-networking"
	@echo "  deploy-secrets      / delete-secrets      / validate-secrets"
	@echo "  deploy-auth         / delete-auth         / validate-auth"
	@echo "  deploy-platform     / delete-platform     / validate-platform"
	@echo "  deploy-observability / delete-observability / validate-observability"
	@echo "  deploy-databases    / delete-databases    / validate-databases (dev only)"
	@echo ""
	@echo "  deploy-app          APP=<name>     # Build+push+deploy single app"
	@echo "  delete-app          APP=<name>     # Delete single app"
	@echo "  validate-app        APP=<name>     # Validate single app"
	@echo "  deploy-apps / delete-apps / validate-apps (all apps)"
	@echo ""
	@echo "🧪 Validation"
	@echo "  clean-charts                       # Remove k8s/**/charts"
	@echo "  validate-all                       # Validate entire stack"
	@echo ""
	@echo "🎯 High-level"
	@echo "  deploy-all          ENV=dev|prod   # Deploy full stack"
	@echo "  delete-all          ENV=dev|prod   # Delete full stack"

apps-list:
	@echo "$(APPS)"

ui-apps-list:
	@echo "$(UI_APPS)"

# ============================
# Docker Operations
# ============================
docker-build:
	$(call require_app)
	@DFILE="$(call resolve_dockerfile,$(APP))"; \
	if [ -z "$$DFILE" ]; then echo "❌ No Dockerfile for $(APP)"; exit 1; fi; \
	echo "🐳 Building $(REGISTRY)/$(APP):$(TAG)"; \
	docker build -t "$(REGISTRY)/$(APP):$(TAG)" -f "$$DFILE" "$(APPS_DIR)/$(APP)"

docker-push:
	$(call require_app)
	@echo "📤 Pushing $(REGISTRY)/$(APP):$(TAG)"; \
	docker push "$(REGISTRY)/$(APP):$(TAG)"

docker-build-push: docker-build docker-push

docker-build-all:
	$(call for_each_app,docker-build)
	@echo "✅ Built all images"

docker-push-all:
	$(call for_each_app,docker-push)
	@echo "✅ Pushed all images"

docker-build-push-all:
	$(call for_each_app,docker-build-push)
	@echo "✅ Built and pushed all images"

# ============================
# UI Static Site Operations (SeaweedFS in-cluster upload)
# ============================

deploy-ui-bucket:
	@echo "🪣 Ensuring bucket exists: $(UI_PUBLIC_BUCKET)"
	@$(SCRIPTS_DIR)/seaweedfs-ui-bucket.sh

deploy-ui-build:
	$(call require_app)
	@APP_DIR="$(APPS_DIR)/$(APP)"; \
	if [ ! -d "$$APP_DIR" ]; then echo "❌ $$APP_DIR not found"; exit 1; fi; \
	if [ -z "$(call is_ui_app,$(APP))" ]; then \
		echo "❌ $(APP) is not a UI app"; exit 1; \
	fi; \
	if [ ! -f "$$APP_DIR/package.json" ]; then \
		echo "❌ package.json not found"; exit 1; \
	fi; \
	PM="npm"; \
	if [ -f "$$APP_DIR/pnpm-lock.yaml" ]; then PM="pnpm"; \
	elif [ -f "$$APP_DIR/yarn.lock" ]; then PM="yarn"; fi; \
	echo "🎨 Building $$APP_DIR (using $$PM)"; \
	cd "$$APP_DIR" && NEXT_PUBLIC_BASE_PATH="$(UI_BASE_PATH)" $$PM run build

deploy-ui-build-all:
	$(call for_each_ui_app,deploy-ui-build)
	@echo "✅ Built all UI apps"

deploy-ui-upload:
	$(call require_app)
	@if [ -z "$(call is_ui_app,$(APP))" ]; then \
		echo "❌ $(APP) is not a UI app"; exit 1; \
	fi
	@$(MAKE) --no-print-directory deploy-ui-bucket
	@$(SCRIPTS_DIR)/seaweedfs-ui-upload.sh
	@echo "✅ Uploaded $(APP)"

deploy-ui-upload-all:
	$(call for_each_ui_app,deploy-ui-upload)
	@echo "✅ Uploaded all UI apps"

deploy-ui-url:
	$(call require_app)
	@echo "$(SEAWEED_PUBLIC_BASE)/$(UI_PUBLIC_BUCKET)/$(APP)/index.html"

deploy-ui: deploy-ui-build deploy-ui-upload

deploy-ui-all:
	$(call for_each_ui_app,deploy-ui)
	@echo "✅ Deployed all UI apps"

# ============================
# Kubernetes Layers
# ============================
deploy-namespaces:
	$(call k8s_apply_or_base,$(NAMESPACES_DIR))
	@echo "✅ Namespaces deployed"

delete-namespaces:
	$(call k8s_delete_or_base,$(NAMESPACES_DIR))

validate-namespaces:
	$(call k8s_dry_run_or_base,$(NAMESPACES_DIR))
	@echo "✅ Namespaces validated"

# ---- Networking ----
deploy-networking: deploy-namespaces
	$(call k8s_apply_helm_list,$(NETWORKING_HELM_DIRS))
	$(call k8s_apply_list,$(NETWORKING_DIRS))
	@echo "✅ Networking deployed"

delete-networking:
	$(call k8s_delete_list,$(NETWORKING_DIRS))
	$(call k8s_delete_helm_list,$(NETWORKING_HELM_DIRS))

validate-networking:
	$(call k8s_dry_run_helm_list,$(NETWORKING_HELM_DIRS))
	@$(KUBECTL) kustomize "$(GATEWAY_DIR)/$(if $(wildcard $(GATEWAY_DIR)/overlays/$(ENV)),overlays/$(ENV),$(if $(ENV_FALLBACK),$(if $(wildcard $(GATEWAY_DIR)/overlays/$(ENV_FALLBACK)),overlays/$(ENV_FALLBACK),base),base))" >/dev/null
	@echo "✅ Networking validated"

# ---- Secrets (External Secrets Operator) ----
deploy-secrets: deploy-namespaces
	$(call k8s_apply_helm_list,$(SECRETS_HELM_DIRS))
	@echo "✅ Secrets deployed"

delete-secrets:
	$(call k8s_delete_helm_list,$(SECRETS_HELM_DIRS))

validate-secrets:
	$(call k8s_dry_run_helm_list,$(SECRETS_HELM_DIRS))
	@echo "✅ Secrets validated"

# ---- Auth ----
deploy-auth: deploy-namespaces deploy-secrets
	$(call k8s_apply_helm_list,$(AUTH_HELM_DIRS))
	$(call k8s_apply_list,$(AUTH_DIRS))
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)"
	@echo "✅ Auth deployed"

delete-auth:
	@$(KUBECTL) delete --ignore-not-found -f "$(AUTH_REFERENCE_GRANT)"
	$(call k8s_delete_list,$(AUTH_DIRS))
	$(call k8s_delete_helm_list,$(AUTH_HELM_DIRS))

validate-auth:
	$(call k8s_dry_run_helm_list,$(AUTH_HELM_DIRS))
	$(call k8s_dry_run_list,$(AUTH_DIRS))
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)" --dry-run=client -o yaml >/dev/null
	@echo "✅ Auth validated"

# ---- Platform ----
deploy-platform: deploy-namespaces deploy-secrets
	$(call k8s_apply_helm_list,$(PLATFORM_HELM_DIRS))
	@echo "✅ Platform deployed"

delete-platform:
	$(call k8s_delete_helm_list,$(PLATFORM_HELM_DIRS))

validate-platform:
	$(call k8s_dry_run_helm_list,$(PLATFORM_HELM_DIRS))
	@echo "✅ Platform validated"

# ---- Observability ----
deploy-observability: deploy-namespaces deploy-secrets
	$(call k8s_apply_helm_list,$(OBS_HELM_DIRS))
	@echo "✅ Observability deployed"

delete-observability:
	$(call k8s_delete_helm_list,$(OBS_HELM_DELETE_DIRS))

validate-observability:
	$(call k8s_dry_run_helm_list,$(OBS_HELM_DIRS))
	@echo "✅ Observability validated"

# ---- Databases (dev only) ----
deploy-databases: deploy-namespaces
ifeq ($(ENV),dev)
	$(call k8s_apply_list,$(DATABASE_DIRS))
	@echo "✅ Databases deployed"
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

delete-databases:
ifeq ($(ENV),dev)
	$(call k8s_delete_list,$(DATABASE_DELETE_DIRS))
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

validate-databases:
ifeq ($(ENV),dev)
	$(call k8s_dry_run_list,$(DATABASE_DIRS))
	@echo "✅ Databases validated"
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

# ---- Apps ----
deploy-app:
	$(call require_app)
	$(call require_app_k8s)
	@$(MAKE) --no-print-directory docker-build-push APP=$(APP)
	$(call k8s_apply_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ Deployed $(APP)"

deploy-apps:
	$(call for_each_k8s_app,deploy-app)
	@echo "✅ All apps deployed"

delete-app:
	$(call require_app)
	$(call require_app_k8s)
	$(call k8s_delete_or_base,$(KUBERNETES_DIR)/apps/$(APP))

delete-apps:
	$(call for_each_k8s_app,delete-app)

validate-app:
	$(call require_app)
	$(call require_app_k8s)
	$(call k8s_dry_run_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ $(APP) validated"

validate-apps:
	$(call for_each_k8s_app,validate-app)
	@echo "✅ All apps validated"

# ============================
# Validation & Cleanup
# ============================
clean-charts:
	@echo "🧹 Cleaning Helm chart cache..."
	@find $(KUBERNETES_DIR) -type d -name charts -prune -exec rm -rf {} +
	@echo "✅ Charts cleaned"

validate-all: clean-charts
	@echo "🧪 Validating all manifests (ENV=$(ENV))..."
	@$(MAKE) --no-print-directory validate-namespaces
	@$(MAKE) --no-print-directory validate-networking
	@$(MAKE) --no-print-directory validate-secrets
	@$(MAKE) --no-print-directory validate-auth
	@$(MAKE) --no-print-directory validate-platform
	@$(MAKE) --no-print-directory validate-observability
	@$(MAKE) --no-print-directory validate-databases
	@$(MAKE) --no-print-directory validate-apps
	@echo "✅ All validations passed"
	@$(MAKE) --no-print-directory clean-charts

# ============================
# High-Level Operations
# ============================
deploy-all: validate-all deploy-namespaces deploy-networking deploy-secrets deploy-auth deploy-databases deploy-platform deploy-observability deploy-apps deploy-ui-all
	@echo "🎉 Full stack deployed (ENV=$(ENV))"

delete-all: delete-apps delete-observability delete-platform delete-databases delete-auth delete-secrets delete-networking delete-namespaces
	@echo "✅ Full stack deleted (ENV=$(ENV))"
