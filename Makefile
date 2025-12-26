# ============================
# Makefile (root, kubectl-native + helm)
# ============================

# ---- Config (override via env/CLI) ----
REGISTRY          ?= docker.io/masterbogdan0
TAG               ?= latest
ENV               ?= dev              # dev | prod
APP               ?=                  # single app name, e.g. personal-website-ui
# ---- UI static site publish (dev defaults) ----
RCLONE                ?= rclone
UI_NAME_PATTERN       ?= ui
UI_BUILD_OUTPUT_DIRS  ?= build dist out
SEAWEEDFS_S3_REGION   ?= us-east-1
SEAWEEDFS_S3_BUCKET     ?= $(APP)
RCLONE_REMOTE         ?= SEAWEEDFS
RCLONE_REMOTE_ENV     ?= SEAWEEDFS
RCLONE_SYNC_ARGS      ?= --checksum
RCLONE_FORCE_PATH_STYLE ?= true
RCLONE_NO_CHECK_BUCKET  ?= true
ifeq ($(ENV),dev)
SEAWEEDFS_S3_ENDPOINT ?= http://storage.apps.192-168-58-2.sslip.io:30080
SEAWEEDFS_S3_ACCESS_KEY ?= yZSEuJpiHbYcXHxU9Wci
SEAWEEDFS_S3_SECRET_KEY ?= Bp3vkTLC5BxxeqYFADlYvU807I4ryGUFzExB5k4N
endif
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
AUTHENTIK_DIR         := $(KUBERNETES_DIR)/auth/authentik
AUTH_DIR              := $(KUBERNETES_DIR)/auth
AUTH_REFERENCE_GRANT  := $(AUTH_DIR)/reference-grant.yaml
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

# ---- kubectl (cluster-agnostic) ----
KUBECTL               := kubectl

# ============================
# Kustomize helpers (plain apply/delete)
# ============================

# Apply kustomization: prefer overlays/$(ENV), fallback to base (no helmCharts)
define k8s_apply_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🚀 Applying: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) apply -k "$${ROOT}/overlays/$(ENV)"; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
	  echo "🚀 Applying: $${ROOT}"; \
	  $(KUBECTL) apply -k "$${ROOT}"; \
	else \
	  echo "🚀 Applying: $${ROOT}/base"; \
	  $(KUBECTL) apply -k "$${ROOT}/base" 2>/dev/null || \
	  $(KUBECTL) apply -f "$${ROOT}/base"; \
	fi
endef

# Delete kustomization: prefer overlays/$(ENV), fallback to base (no helmCharts)
define k8s_delete_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🔥 Deleting: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) delete -k "$${ROOT}/overlays/$(ENV)" --ignore-not-found; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
	  echo "🔥 Deleting: $${ROOT}"; \
	  $(KUBECTL) delete -k "$${ROOT}" --ignore-not-found; \
	else \
	  echo "🔥 Deleting: $${ROOT}/base"; \
	  $(KUBECTL) delete -k "$${ROOT}/base" --ignore-not-found 2>/dev/null || \
	  $(KUBECTL) delete -f "$${ROOT}/base" --ignore-not-found; \
	fi
endef

# ============================
# HELM+KUSTOMIZE helpers (for helmCharts)
# ============================

# Apply with kubectl kustomize --enable-helm (Traefik, Authentik, n8n, obs stack, etc.)
define k8s_apply_helm_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🚀 Applying (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" | $(KUBECTL) apply -f -; \
	else \
	  echo "🚀 Applying (helm+kustomize): $${ROOT}/base"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) apply -f -; \
	fi
endef

# Delete with kubectl kustomize --enable-helm (mirror of apply_helm)
define k8s_delete_helm_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🔥 Deleting (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" | $(KUBECTL) delete --ignore-not-found -f -; \
	else \
	  echo "🔥 Deleting (helm+kustomize): $${ROOT}/base"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) delete --ignore-not-found -f -; \
	fi
endef

# ============================
# Validation helpers (dry-run/build-only)
# ============================

# Validate normal kustomization with client-side dry-run, no schema/API required
define k8s_dry_run_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🧪 Dry-run: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) apply -k "$${ROOT}/overlays/$(ENV)" --dry-run=client --validate=false -o yaml >/dev/null; \
	elif [ -f "$${ROOT}/kustomization.yaml" ]; then \
	  echo "🧪 Dry-run: $${ROOT}"; \
	  $(KUBECTL) apply -k "$${ROOT}" --dry-run=client --validate=false -o yaml >/dev/null; \
	else \
	  echo "🧪 Dry-run: $${ROOT}/base"; \
	  $(KUBECTL) apply -k "$${ROOT}/base" --dry-run=client --validate=false -o yaml >/dev/null 2>/dev/null || \
	  $(KUBECTL) apply -f "$${ROOT}/base" --dry-run=client --validate=false -o yaml >/dev/null; \
	fi
endef

# Validate helm+kustomize kustomization (build-only via kubectl kustomize --enable-helm)
define k8s_dry_run_helm_or_base
	set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🧪 Build-only (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" >/dev/null; \
	else \
	  echo "🧪 Build-only (helm+kustomize): $${ROOT}/base"; \
	  $(KUBECTL) kustomize --enable-helm "$${ROOT}/base" >/dev/null; \
	fi
endef

# ============================
# Dockerfile resolver
# ============================

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
	deploy-ui-list deploy-ui-build deploy-ui-build-all deploy-ui-upload deploy-ui-upload-all deploy-ui deploy-ui-all \
	deploy-namespaces delete-namespaces validate-namespaces \
	deploy-app deploy-apps delete-app delete-apps validate-app validate-apps \
	deploy-networking delete-networking validate-networking \
	deploy-auth delete-auth validate-auth \
	deploy-platform delete-platform validate-platform \
	deploy-observability delete-observability validate-observability \
	deploy-databases delete-databases validate-databases \
	validate-all \
	deploy-all deploy-all-dev deploy-all-prod \
	delete-all delete-all-dev delete-all-prod \
	clean-charts \

# ============================
# Help
# ============================
help:
	@echo "🎛  Root targets (ENV=$(ENV))"
	@echo "  apps-list                          # List discovered apps"
	@echo ""
	@echo "🐳 Docker"
	@echo "  docker-build        APP=<name>     # Build one image"
	@echo "  docker-push         APP=<name>     # Push one image"
	@echo "  docker-build-push   APP=<name>     # Build + push one image"
	@echo "  docker-build-all                   # Build all images"
	@echo "  docker-push-all                    # Push all images"
	@echo "  docker-build-push-all              # Build + push all images"
	@echo "🎨 UI (static -> SeaweedFS)"
	@echo "  deploy-ui-list                     # List UI apps (name contains UI, no Dockerfile)"
	@echo "  deploy-ui-build    APP=<name>      # Build UI app"
	@echo "  deploy-ui-upload   APP=<name>      # rclone build output to SeaweedFS"
	@echo "  deploy-ui          APP=<name>      # Build + upload UI app"
	@echo "  deploy-ui-build-all               # Build all UI apps"
	@echo "  deploy-ui-upload-all              # Upload all UI apps"
	@echo "  deploy-ui-all                     # Build + upload all UI apps"
	@echo "☸  Kubernetes (per layer)"
	@echo "  deploy-namespaces / delete-namespaces / validate-namespaces"
	@echo "  deploy-networking  / delete-networking  / validate-networking"
	@echo "  deploy-auth        / delete-auth        / validate-auth"
	@echo "  deploy-platform    / delete-platform    / validate-platform"
	@echo "  deploy-observability / delete-observability / validate-observability"
	@echo "  deploy-databases   / delete-databases   / validate-databases  (dev only)"
	@echo ""
	@echo "  deploy-app         APP=<name>      # Build+push APP, then apply k8s"
	@echo "  delete-app         APP=<name>      # Delete k8s for APP"
	@echo "  validate-app       APP=<name>      # Dry-run validate k8s for APP"
	@echo "  deploy-apps / delete-apps / validate-apps (all apps)"
	@echo ""
	@echo "🧪 Validation"
	@echo "  clean-charts                       # Delete all k8s/**/charts before regen"
	@echo "  validate-all                       # Clean charts + dry-run/build-only EVERYTHING"
	@echo ""
	@echo "🎯 High-level"
	@echo "  ENV=dev  make deploy-all           # dev: validate + deploy full stack (with DBs)"
	@echo "  ENV=prod make deploy-all           # prod: validate + deploy full stack (no DBs)"
	@echo "  ENV=dev|prod make delete-all       # Full teardown"
	@echo ""
	@echo "🔧 Vars: REGISTRY=$(REGISTRY) TAG=$(TAG) ENV=$(ENV)"
	@echo "   kubectl must be configured to point at your cluster (minikube/k3s/k8s)."

apps-list:
	@echo "$(APPS)"

# ============================
# UI static site build + publish (SeaweedFS rclone)
# ============================
deploy-ui-list:
	@for app in $(APPS); do \
	  if printf "%s" "$$app" | grep -qi "$(UI_NAME_PATTERN)"; then \
	    if [ ! -f "$(APPS_DIR)/$$app/Dockerfile" ] && [ ! -f "$(APPS_DIR)/$$app/deployments/docker/Dockerfile" ]; then \
	      echo "$$app"; \
	    fi; \
	  fi; \
	done

deploy-ui-build:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@APP_DIR="$(APPS_DIR)/$(APP)"; \
	if [ ! -d "$$APP_DIR" ]; then echo "❌ $$APP_DIR not found"; exit 1; fi; \
	if ! printf "%s" "$(APP)" | grep -qi "$(UI_NAME_PATTERN)"; then \
	  echo "❌ $(APP) is not a UI app (name must contain UI)"; exit 1; \
	fi; \
	if [ -f "$$APP_DIR/Dockerfile" ] || [ -f "$$APP_DIR/deployments/docker/Dockerfile" ]; then \
	  echo "❌ $(APP) has a Dockerfile; static UI build is skipped"; exit 1; \
	fi; \
	if [ ! -f "$$APP_DIR/package.json" ]; then \
	  echo "❌ $$APP_DIR/package.json not found"; exit 1; \
	fi; \
	if [ -f "$$APP_DIR/pnpm-lock.yaml" ]; then PM="pnpm"; \
	elif [ -f "$$APP_DIR/yarn.lock" ]; then PM="yarn"; \
	else PM="npm"; fi; \
	echo "🎨 Building UI $$APP_DIR (using $$PM)"; \
	cd "$$APP_DIR" && $$PM run build

deploy-ui-build-all:
	@for app in $(APPS); do \
	  if printf "%s" "$$app" | grep -qi "$(UI_NAME_PATTERN)"; then \
	    if [ ! -f "$(APPS_DIR)/$$app/Dockerfile" ] && [ ! -f "$(APPS_DIR)/$$app/deployments/docker/Dockerfile" ]; then \
	      APP="$$app" $(MAKE) --no-print-directory deploy-ui-build; \
	    fi; \
	  fi; \
	done
	@echo "✅ Built all UI apps"

deploy-ui-upload:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@APP_DIR="$(APPS_DIR)/$(APP)"; \
	if [ ! -d "$$APP_DIR" ]; then echo "❌ $$APP_DIR not found"; exit 1; fi; \
	if ! printf "%s" "$(APP)" | grep -qi "$(UI_NAME_PATTERN)"; then \
	  echo "❌ $(APP) is not a UI app (name must contain UI)"; exit 1; \
	fi; \
	if [ -f "$$APP_DIR/Dockerfile" ] || [ -f "$$APP_DIR/deployments/docker/Dockerfile" ]; then \
	  echo "❌ $(APP) has a Dockerfile; static UI upload is skipped"; exit 1; \
	fi; \
	OUT_DIR=""; \
	for candidate in $(UI_BUILD_OUTPUT_DIRS); do \
	  if [ -d "$$APP_DIR/$$candidate" ]; then OUT_DIR="$$APP_DIR/$$candidate"; break; fi; \
	done; \
	if [ -z "$$OUT_DIR" ]; then \
	  echo "❌ No build output dir found for $(APP). Expected one of: $(UI_BUILD_OUTPUT_DIRS)"; \
	  exit 1; \
	fi; \
	if ! command -v "$(RCLONE)" >/dev/null 2>&1; then \
	  echo "❌ rclone is required but not installed"; exit 1; \
	fi; \
	ENDPOINT="$(SEAWEEDFS_S3_ENDPOINT)"; \
	ACCESS_KEY="$(SEAWEEDFS_S3_ACCESS_KEY)"; \
	SECRET_KEY="$(SEAWEEDFS_S3_SECRET_KEY)"; \
	REGION="$(SEAWEEDFS_S3_REGION)"; \
	BUCKET="$(SEAWEEDFS_S3_BUCKET)"; \
	REMOTE="$(RCLONE_REMOTE)"; \
	REMOTE_ENV="$(RCLONE_REMOTE_ENV)"; \
	if [ -z "$$ENDPOINT" ]; then echo "❌ SEAWEEDFS_S3_ENDPOINT is required"; exit 1; fi; \
	if [ -z "$$ACCESS_KEY" ] || [ -z "$$SECRET_KEY" ]; then \
	  echo "❌ SEAWEEDFS_S3_ACCESS_KEY and SEAWEEDFS_S3_SECRET_KEY are required"; exit 1; \
	fi; \
	if [ -z "$$BUCKET" ]; then echo "❌ SEAWEEDFS_S3_BUCKET is required"; exit 1; fi; \
	echo "📦 rclone $$OUT_DIR to $$REMOTE:$$BUCKET"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_TYPE=s3; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_PROVIDER=Other; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_ENDPOINT="$$ENDPOINT"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_ACCESS_KEY_ID="$$ACCESS_KEY"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_SECRET_ACCESS_KEY="$$SECRET_KEY"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_REGION="$$REGION"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_FORCE_PATH_STYLE="$(RCLONE_FORCE_PATH_STYLE)"; \
	export RCLONE_CONFIG_$${REMOTE_ENV}_NO_CHECK_BUCKET="$(RCLONE_NO_CHECK_BUCKET)"; \
	$(RCLONE) sync "$$OUT_DIR" "$$REMOTE:$$BUCKET" $(RCLONE_SYNC_ARGS)

deploy-ui-upload-all:
	@for app in $(APPS); do \
	  if printf "%s" "$$app" | grep -qi "$(UI_NAME_PATTERN)"; then \
	    if [ ! -f "$(APPS_DIR)/$$app/Dockerfile" ] && [ ! -f "$(APPS_DIR)/$$app/deployments/docker/Dockerfile" ]; then \
	      APP="$$app" $(MAKE) --no-print-directory deploy-ui-upload; \
	    fi; \
	  fi; \
	done
	@echo "✅ Uploaded all UI apps"

deploy-ui:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@$(MAKE) --no-print-directory deploy-ui-build APP=$(APP)
	@$(MAKE) --no-print-directory deploy-ui-upload APP=$(APP)

deploy-ui-all:
	@for app in $(APPS); do \
	  if printf "%s" "$$app" | grep -qi "$(UI_NAME_PATTERN)"; then \
	    if [ ! -f "$(APPS_DIR)/$$app/Dockerfile" ] && [ ! -f "$(APPS_DIR)/$$app/deployments/docker/Dockerfile" ]; then \
	      APP="$$app" $(MAKE) --no-print-directory deploy-ui; \
	    fi; \
	  fi; \
	done
	@echo "✅ Deployed all UI apps"

# ============================
# Docker (no pattern wrappers)
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
# Namespaces
# ============================
deploy-namespaces:
	$(call k8s_apply_or_base,$(NAMESPACES_DIR))
	@echo "✅ Namespaces applied for $(ENV)"

delete-namespaces:
	$(call k8s_delete_or_base,$(NAMESPACES_DIR))
	@echo "🗑  Namespaces deleted for $(ENV)"

validate-namespaces:
	$(call k8s_dry_run_or_base,$(NAMESPACES_DIR))
	@echo "✅ Namespaces manifest validation passed for $(ENV)"

# ============================
# Apps
# ============================
deploy-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then echo "⚠  Skipping $(APP): no k8s manifests"; exit 0; fi
	@$(MAKE) --no-print-directory docker-build-push APP=$(APP)
	$(call k8s_apply_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ Deployed app: $(APP) (env=$(ENV))"

deploy-apps:
	@for app in $(APPS); do \
	  echo "🚀 Deploying app $$app (env=$(ENV))"; \
	  if [ ! -d "$(KUBERNETES_DIR)/apps/$$app" ]; then echo "⚠  Skipping $$app: no k8s manifests"; continue; fi; \
	  APP="$$app" $(MAKE) --no-print-directory deploy-app; \
	done
	@echo "🎉 All apps deployed (env=$(ENV))"

delete-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then echo "⚠  Skipping $(APP): no k8s manifests"; exit 0; fi
	$(call k8s_delete_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "🗑  Deleted app: $(APP) (env=$(ENV))"

delete-apps:
	@for app in $(APPS); do \
	  echo "🧹 Deleting app $$app (env=$(ENV))"; \
	  if [ ! -d "$(KUBERNETES_DIR)/apps/$$app" ]; then echo "⚠  Skipping $$app: no k8s manifests"; continue; fi; \
	  APP="$$app" $(MAKE) --no-print-directory delete-app; \
	done
	@echo "🧹 All apps deleted (env=$(ENV))"

validate-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then echo "⚠  Skipping $(APP): no k8s manifests"; exit 0; fi
	$(call k8s_dry_run_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ App $(APP) manifests are valid (env=$(ENV))"

validate-apps:
	@for app in $(APPS); do \
	  echo "🧪 Validating app $$app (env=$(ENV))"; \
	  if [ ! -d "$(KUBERNETES_DIR)/apps/$$app" ]; then echo "⚠  Skipping $$app: no k8s manifests"; continue; fi; \
	  APP="$$app" $(MAKE) --no-print-directory validate-app; \
	done
	@echo "✅ All apps manifest validation passed (env=$(ENV))"

# ============================
# Networking (Traefik + Gateway)
# ============================
deploy-networking: deploy-namespaces
	# Traefik uses helmCharts -> kubectl kustomize --enable-helm
	$(call k8s_apply_helm_or_base,$(TRAEFIK_DIR))
	# Gateway is plain kustomize
	$(call k8s_apply_or_base,$(GATEWAY_DIR))
	@echo "✅ Networking (Traefik + Gateway) deployed (env=$(ENV))"

delete-networking:
	$(call k8s_delete_or_base,$(GATEWAY_DIR))
	$(call k8s_delete_helm_or_base,$(TRAEFIK_DIR))
	@echo "🗑  Networking (Traefik + Gateway) deleted (env=$(ENV))"

validate-networking:
	# Traefik: helm+kustomize build-only (no apply)
	$(call k8s_dry_run_helm_or_base,$(TRAEFIK_DIR))
	# Gateway: build-only (no apply, avoids CRD discovery issues)
	set -e; ROOT="$(GATEWAY_DIR)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🧪 Build-only: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) kustomize "$${ROOT}/overlays/$(ENV)" >/dev/null; \
	else \
	  echo "🧪 Build-only: $${ROOT}/base"; \
	  $(KUBECTL) kustomize "$${ROOT}/base" >/dev/null; \
	fi
	@echo "✅ Networking manifests are valid (env=$(ENV))"

# ============================
# Auth (Authentik + supporting resources)
# ============================
deploy-auth: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_apply_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	@echo "🚀 Applying ReferenceGrant: $(AUTH_REFERENCE_GRANT)"
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)"
	@echo "🔐 Auth deployed (env=$(ENV))"

delete-auth:
	@echo "🔥 Deleting ReferenceGrant: $(AUTH_REFERENCE_GRANT)"
	@$(KUBECTL) delete --ignore-not-found -f "$(AUTH_REFERENCE_GRANT)"
	$(call k8s_delete_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	$(call k8s_delete_helm_or_base,$(AUTHENTIK_DIR))
	@echo "🧹 Auth deleted (env=$(ENV))"

validate-auth:
	$(call k8s_dry_run_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_dry_run_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	@echo "🧪 Dry-run: $(AUTH_REFERENCE_GRANT)"
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)" --dry-run=client --validate=false -o yaml >/dev/null
	@echo "✅ Auth manifests are valid (env=$(ENV))"

# ============================
# Platform (authentik, n8n, seaweedfs)
# ============================
deploy-platform: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_apply_helm_or_base,$(N8N_DIR))
	$(call k8s_apply_helm_or_base,$(SEAWEEDFS_DIR))
	@echo "🎉 Platform services deployed (env=$(ENV))"

delete-platform:
	$(call k8s_delete_helm_or_base,$(SEAWEEDFS_DIR))
	$(call k8s_delete_helm_or_base,$(N8N_DIR))
	$(call k8s_delete_helm_or_base,$(AUTHENTIK_DIR))
	@echo "🧹 Platform services deleted (env=$(ENV))"

validate-platform:
	$(call k8s_dry_run_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_dry_run_helm_or_base,$(N8N_DIR))
	$(call k8s_dry_run_helm_or_base,$(SEAWEEDFS_DIR))
	@echo "✅ Platform manifests are valid (env=$(ENV))"

# ============================
# Observability (all via helmCharts)
# ============================
deploy-observability: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(FLUENTBIT_DIR))
	$(call k8s_apply_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_apply_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_apply_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_apply_helm_or_base,$(OPENSEARCH_DASH_DIR))
	@echo "📊 Observability stack deployed (env=$(ENV))"

delete-observability:
	$(call k8s_delete_helm_or_base,$(OPENSEARCH_DASH_DIR))
	$(call k8s_delete_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_delete_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_delete_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_delete_helm_or_base,$(FLUENTBIT_DIR))
	@echo "🧻 Observability stack deleted (env=$(ENV))"

validate-observability:
	$(call k8s_dry_run_helm_or_base,$(FLUENTBIT_DIR))
	$(call k8s_dry_run_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_dry_run_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_dry_run_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_dry_run_helm_or_base,$(OPENSEARCH_DASH_DIR))
	@echo "✅ Observability manifests are valid (env=$(ENV))"

# ============================
# Databases (dev only)
# ============================
deploy-databases: deploy-namespaces
ifeq ($(ENV),dev)
	$(call k8s_apply_or_base,$(REDIS_DIR))
	$(call k8s_apply_or_base,$(POSTGRESQL_DIR))
	@echo "🧰 Databases deployed (env=$(ENV))"
else
	@echo "⚠  deploy-databases skipped: ENV=$(ENV) (no DBs in prod)"
endif

delete-databases:
ifeq ($(ENV),dev)
	$(call k8s_delete_or_base,$(POSTGRESQL_DIR))
	$(call k8s_delete_or_base,$(REDIS_DIR))
	@echo "🧨 Databases deleted (env=$(ENV))"
else
	@echo "⚠  delete-databases skipped: ENV=$(ENV) (no DBs in prod)"
endif

validate-databases:
ifeq ($(ENV),dev)
	$(call k8s_dry_run_or_base,$(REDIS_DIR))
	$(call k8s_dry_run_or_base,$(POSTGRESQL_DIR))
	@echo "✅ Database manifests are valid (env=$(ENV))"
else
	@echo "⚠  validate-databases skipped: ENV=$(ENV) (no DBs in prod)"
endif

# ============================
# Charts cleanup
# ============================
clean-charts:
	@echo "🧹 Cleaning local Helm chart cache (k8s/**/charts)..."
	@find $(KUBERNETES_DIR) -type d -name charts -prune -exec rm -rf {} +
	@echo "✅ charts/ directories removed"

# ============================
# Global validation (dry-run/build-only)
# ============================
validate-all: clean-charts
	@echo "🧪 Validating ALL manifests (ENV=$(ENV))..."
	@$(MAKE) --no-print-directory validate-namespaces
	@$(MAKE) --no-print-directory validate-networking
	@$(MAKE) --no-print-directory validate-platform
	@$(MAKE) --no-print-directory validate-observability
	@$(MAKE) --no-print-directory validate-databases
	@$(MAKE) --no-print-directory validate-apps
	@echo "✅ validate-all finished (ENV=$(ENV))"
	@$(MAKE) --no-print-directory clean-charts

# ============================
# Everything (dev/prod)
# ============================
deploy-all-dev: validate-all deploy-namespaces deploy-networking deploy-databases deploy-platform deploy-observability deploy-apps deploy-ui-all
	@echo "✅ Full stack applied for dev (with databases)"

deploy-all-prod: validate-all deploy-namespaces deploy-networking deploy-platform deploy-observability deploy-apps deploy-ui-all
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
