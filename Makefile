# ============================
# Makefile (Optimized)
# ============================

# ---- Config ----
REGISTRY          ?= docker.io/masterbogdan0
TAG               ?= latest
ENV               ?= dev
APP               ?=
ROOT_DIR          := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))

# ---- UI Static Site Config ----
UI_NAME_PATTERN         ?= ui
UI_BUILD_OUTPUT_DIRS    ?= out build dist
UI_PUBLIC_BUCKET        ?= public
SEAWEEDFS_NAMESPACE     ?= platform
SEAWEED_PUBLIC_BASE     ?= http://storage.apps.192-168-58-2.sslip.io:30080
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
AUTHENTIK_DIR         := $(AUTH_DIR)/authentik
AUTHENTIK_FWD_AUTH_DIR := $(AUTH_DIR)/forward-auth
AUTH_REFERENCE_GRANT  := $(AUTH_DIR)/reference-grant.yaml
N8N_DIR               := $(PLATFORM_DIR)/n8n
SEAWEEDFS_DIR          := $(PLATFORM_DIR)/seaweedfs
OBSERVABILITY_DIR     := $(KUBERNETES_DIR)/observability
FLUENTBIT_DIR         := $(OBSERVABILITY_DIR)/fluent-bit
GRAFANA_DIR           := $(OBSERVABILITY_DIR)/grafana
PROMETHEUS_DIR        := $(OBSERVABILITY_DIR)/prometheus
OPENSEARCH_DIR        := $(OBSERVABILITY_DIR)/opensearch
OPENSEARCH_DASH_DIR   := $(OBSERVABILITY_DIR)/opensearch-dashboards

# ---- Discovery ----
APPS                  := $(notdir $(wildcard $(APPS_DIR)/*))
KUBECTL               := kubectl

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

# Generic function to iterate over apps with a target
define for_each_app
	@for app in $(APPS); do \
		APP="$$app" $(MAKE) --no-print-directory $(1); \
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

# Kustomize delete
define k8s_delete_or_base
	@set -e; ROOT="$(1)"; \
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

# Helm+Kustomize apply
define k8s_apply_helm_or_base
	@set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
		echo "🚀 Applying (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" | $(KUBECTL) apply -f -; \
	else \
		echo "🚀 Applying (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) apply -f -; \
	fi
endef

# Helm+Kustomize delete
define k8s_delete_helm_or_base
	@set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
		echo "🔥 Deleting (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" | $(KUBECTL) delete --ignore-not-found -f -; \
	else \
		echo "🔥 Deleting (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" | $(KUBECTL) delete --ignore-not-found -f -; \
	fi
endef

# Validate kustomize (dry-run)
define k8s_dry_run_or_base
	@set -e; ROOT="$(1)"; \
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

# Validate helm+kustomize (build-only)
define k8s_dry_run_helm_or_base
	@set -e; ROOT="$(1)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
		echo "🧪 Build-only (helm+kustomize): $${ROOT}/overlays/$(ENV)"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/overlays/$(ENV)" >/dev/null; \
	else \
		echo "🧪 Build-only (helm+kustomize): $${ROOT}/base"; \
		$(KUBECTL) kustomize --enable-helm "$${ROOT}/base" >/dev/null; \
	fi
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
	@echo "  deploy-namespaces   / delete-namespaces   / validate-namespaces"
	@echo "  deploy-networking   / delete-networking   / validate-networking"
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
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@DFILE="$(call resolve_dockerfile,$(APP))"; \
	if [ -z "$$DFILE" ]; then echo "❌ No Dockerfile for $(APP)"; exit 1; fi; \
	echo "🐳 Building $(REGISTRY)/$(APP):$(TAG)"; \
	docker build -t "$(REGISTRY)/$(APP):$(TAG)" -f "$$DFILE" "$(APPS_DIR)/$(APP)"

docker-push:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
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
	@set -e; \
	$(KUBECTL) -n "$(SEAWEEDFS_NAMESPACE)" exec -i statefulset/seaweedfs-master -- sh -c \
		"echo \"s3.bucket.create --name $(UI_PUBLIC_BUCKET)\" | weed shell -master=\"$(SEAWEEDFS_MASTER_ADDR)\" -filer=\"$(SEAWEEDFS_FILER_ADDR)\"" || true; \
	$(KUBECTL) -n "$(SEAWEEDFS_NAMESPACE)" exec -i statefulset/seaweedfs-master -- sh -c \
		"echo \"s3.configure --user anonymous --buckets $(UI_PUBLIC_BUCKET) --actions $(SEAWEED_ANON_ACTIONS) --apply true\" | weed shell -master=\"$(SEAWEEDFS_MASTER_ADDR)\" -filer=\"$(SEAWEEDFS_FILER_ADDR)\"" || true

deploy-ui-build:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
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
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@set -e; \
	APP_DIR="$(APPS_DIR)/$(APP)"; \
	if [ -z "$(call is_ui_app,$(APP))" ]; then \
		echo "❌ $(APP) is not a UI app"; exit 1; \
	fi; \
	OUT_DIR=""; \
	for dir in $(UI_BUILD_OUTPUT_DIRS); do \
		if [ -d "$$APP_DIR/$$dir" ]; then OUT_DIR="$$APP_DIR/$$dir"; break; fi; \
	done; \
	if [ -z "$$OUT_DIR" ]; then \
		echo "❌ No build output dir found for $(APP)"; exit 1; \
	fi; \
	$(MAKE) --no-print-directory deploy-ui-bucket; \
	NS="$(SEAWEEDFS_NAMESPACE)"; \
	POD="$(SEAWEED_UPLOAD_POD_PREFIX)-$(APP)"; \
	cleanup() { $(KUBECTL) -n "$$NS" delete pod "$$POD" --ignore-not-found >/dev/null; }; \
	trap cleanup EXIT; \
	cleanup; \
	$(KUBECTL) -n "$$NS" run "$$POD" --image="$(SEAWEED_UPLOAD_IMAGE)" --restart=Never --command -- \
		sh -c "$(SEAWEED_UPLOAD_SETUP_CMD) && sleep 3600"; \
	$(KUBECTL) -n "$$NS" wait --for=condition=Ready pod/"$$POD" --timeout="$(SEAWEED_UPLOAD_WAIT)"; \
	echo "📦 Uploading $$OUT_DIR -> $(SEAWEED_INTERNAL_BASE)/$(UI_PUBLIC_BUCKET)/$(APP) (in-cluster)"; \
	OUT_BASE="$$(basename "$$OUT_DIR")"; \
	$(KUBECTL) -n "$$NS" exec "$$POD" -- sh -c "mkdir -p /tmp/upload"; \
	$(KUBECTL) -n "$$NS" cp "$$OUT_DIR" "$$POD:/tmp/upload"; \
	$(KUBECTL) -n "$$NS" exec "$$POD" -- env \
		SEAWEED_BASE="$(SEAWEED_INTERNAL_BASE)" \
		BUCKET="$(UI_PUBLIC_BUCKET)" \
		APP="$(APP)" \
		UPLOAD_DIR="/tmp/upload/$$OUT_BASE" \
		sh -c 'set -e; cd "$$UPLOAD_DIR" && find . -type f | while IFS= read -r f; do \
			rel="$${f#./}"; \
			url="$$SEAWEED_BASE/$$BUCKET/$$APP/$$rel"; \
			content_type="application/octet-stream"; \
			case "$$rel" in \
				*.html) content_type="text/html; charset=utf-8" ;; \
				*.css) content_type="text/css; charset=utf-8" ;; \
				*.js|*.mjs) content_type="application/javascript; charset=utf-8" ;; \
				*.json|*.map) content_type="application/json; charset=utf-8" ;; \
				*.txt) content_type="text/plain; charset=utf-8" ;; \
				*.svg) content_type="image/svg+xml" ;; \
				*.png) content_type="image/png" ;; \
				*.jpg|*.jpeg) content_type="image/jpeg" ;; \
				*.gif) content_type="image/gif" ;; \
				*.webp) content_type="image/webp" ;; \
				*.ico) content_type="image/x-icon" ;; \
				*.woff2) content_type="font/woff2" ;; \
				*.woff) content_type="font/woff" ;; \
				*.ttf) content_type="font/ttf" ;; \
				*.otf) content_type="font/otf" ;; \
				*.eot) content_type="application/vnd.ms-fontobject" ;; \
				*.pdf) content_type="application/pdf" ;; \
			esac; \
			echo "  ↥ $$rel"; \
			curl -sS --fail -X PUT -H "Expect:" -H "Content-Type: $$content_type" \
				--connect-timeout "$(SEAWEED_UPLOAD_CONNECT_TIMEOUT)" \
				--max-time "$(SEAWEED_UPLOAD_MAX_TIME)" \
				--upload-file "$$f" "$$url" -o /dev/null; \
		done'; \
	cleanup; \
	echo "✅ Uploaded $(APP)"

deploy-ui-upload-all:
	$(call for_each_ui_app,deploy-ui-upload)
	@echo "✅ Uploaded all UI apps"

deploy-ui-url:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
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
	$(call k8s_apply_helm_or_base,$(TRAEFIK_DIR))
	$(call k8s_apply_or_base,$(GATEWAY_DIR))
	@echo "✅ Networking deployed"

delete-networking:
	$(call k8s_delete_or_base,$(GATEWAY_DIR))
	$(call k8s_delete_helm_or_base,$(TRAEFIK_DIR))

validate-networking:
	$(call k8s_dry_run_helm_or_base,$(TRAEFIK_DIR))
	@$(KUBECTL) kustomize "$(GATEWAY_DIR)/$(if $(wildcard $(GATEWAY_DIR)/overlays/$(ENV)),overlays/$(ENV),base)" >/dev/null
	@echo "✅ Networking validated"

# ---- Auth ----
deploy-auth: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_apply_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)"
	@echo "✅ Auth deployed"

delete-auth:
	@$(KUBECTL) delete --ignore-not-found -f "$(AUTH_REFERENCE_GRANT)"
	$(call k8s_delete_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	$(call k8s_delete_helm_or_base,$(AUTHENTIK_DIR))

validate-auth:
	$(call k8s_dry_run_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_dry_run_or_base,$(AUTHENTIK_FWD_AUTH_DIR))
	@$(KUBECTL) apply -f "$(AUTH_REFERENCE_GRANT)" --dry-run=client -o yaml >/dev/null
	@echo "✅ Auth validated"

# ---- Platform ----
deploy-platform: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_apply_helm_or_base,$(N8N_DIR))
	$(call k8s_apply_helm_or_base,$(SEAWEEDFS_DIR))
	@echo "✅ Platform deployed"

delete-platform:
	$(call k8s_delete_helm_or_base,$(SEAWEEDFS_DIR))
	$(call k8s_delete_helm_or_base,$(N8N_DIR))
	$(call k8s_delete_helm_or_base,$(AUTHENTIK_DIR))

validate-platform:
	$(call k8s_dry_run_helm_or_base,$(AUTHENTIK_DIR))
	$(call k8s_dry_run_helm_or_base,$(N8N_DIR))
	$(call k8s_dry_run_helm_or_base,$(SEAWEEDFS_DIR))
	@echo "✅ Platform validated"

# ---- Observability ----
deploy-observability: deploy-namespaces
	$(call k8s_apply_helm_or_base,$(FLUENTBIT_DIR))
	$(call k8s_apply_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_apply_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_apply_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_apply_helm_or_base,$(OPENSEARCH_DASH_DIR))
	@echo "✅ Observability deployed"

delete-observability:
	$(call k8s_delete_helm_or_base,$(OPENSEARCH_DASH_DIR))
	$(call k8s_delete_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_delete_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_delete_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_delete_helm_or_base,$(FLUENTBIT_DIR))

validate-observability:
	$(call k8s_dry_run_helm_or_base,$(FLUENTBIT_DIR))
	$(call k8s_dry_run_helm_or_base,$(PROMETHEUS_DIR))
	$(call k8s_dry_run_helm_or_base,$(GRAFANA_DIR))
	$(call k8s_dry_run_helm_or_base,$(OPENSEARCH_DIR))
	$(call k8s_dry_run_helm_or_base,$(OPENSEARCH_DASH_DIR))
	@echo "✅ Observability validated"

# ---- Databases (dev only) ----
deploy-databases: deploy-namespaces
ifeq ($(ENV),dev)
	$(call k8s_apply_or_base,$(REDIS_DIR))
	$(call k8s_apply_or_base,$(POSTGRESQL_DIR))
	@echo "✅ Databases deployed"
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

delete-databases:
ifeq ($(ENV),dev)
	$(call k8s_delete_or_base,$(POSTGRESQL_DIR))
	$(call k8s_delete_or_base,$(REDIS_DIR))
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

validate-databases:
ifeq ($(ENV),dev)
	$(call k8s_dry_run_or_base,$(REDIS_DIR))
	$(call k8s_dry_run_or_base,$(POSTGRESQL_DIR))
	@echo "✅ Databases validated"
else
	@echo "⚠️  Databases skipped in $(ENV) environment"
endif

# ---- Apps ----
deploy-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then \
		echo "⚠️  No k8s manifests for $(APP)"; exit 0; \
	fi
	@$(MAKE) --no-print-directory docker-build-push APP=$(APP)
	$(call k8s_apply_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ Deployed $(APP)"

deploy-apps:
	@for app in $(APPS); do \
		if [ -d "$(KUBERNETES_DIR)/apps/$$app" ]; then \
			APP="$$app" $(MAKE) --no-print-directory deploy-app; \
		fi; \
	done
	@echo "✅ All apps deployed"

delete-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then \
		echo "⚠️  No k8s manifests for $(APP)"; exit 0; \
	fi
	$(call k8s_delete_or_base,$(KUBERNETES_DIR)/apps/$(APP))

delete-apps:
	@for app in $(APPS); do \
		if [ -d "$(KUBERNETES_DIR)/apps/$$app" ]; then \
			APP="$$app" $(MAKE) --no-print-directory delete-app; \
		fi; \
	done

validate-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@if [ ! -d "$(KUBERNETES_DIR)/apps/$(APP)" ]; then \
		echo "⚠️  No k8s manifests for $(APP)"; exit 0; \
	fi
	$(call k8s_dry_run_or_base,$(KUBERNETES_DIR)/apps/$(APP))
	@echo "✅ $(APP) validated"

validate-apps:
	@for app in $(APPS); do \
		if [ -d "$(KUBERNETES_DIR)/apps/$$app" ]; then \
			APP="$$app" $(MAKE) --no-print-directory validate-app; \
		fi; \
	done
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
deploy-all: validate-all deploy-namespaces deploy-networking deploy-auth deploy-databases deploy-platform deploy-observability deploy-apps deploy-ui-all
	@echo "🎉 Full stack deployed (ENV=$(ENV))"

delete-all: delete-apps delete-observability delete-platform delete-databases delete-auth delete-networking delete-namespaces
	@echo "✅ Full stack deleted (ENV=$(ENV))"
