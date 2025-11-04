# ============================
# Monorepo Makefile (root)
# DRY: pattern rules + fan-out deps
# ============================

# ---- Config (override via env/CLI) ----
REGISTRY       ?= docker.io/masterbogdan0
TAG            ?= latest
ENV            ?= dev              # dev | prod
APP            ?=                  # single app name, e.g. personal-website-ui
PARALLEL_JOBS  ?= 1                # pass -jN to make for true parallel; kept for docs
MINIKUBE       ?=                  # set 1 to use "minikube kubectl --"

# ---- Paths ----
APPS_DIR       := apps
K8S_DIR        := k8s
NAMESPACES_DIR := $(K8S_DIR)/namespaces
NET_DIR        := $(K8S_DIR)/networking
TRAEFIK_DIR    := $(NET_DIR)/traefik
GATEWAY_DIR    := $(NET_DIR)/gateway

# ---- Discover apps ----
APPS           := $(notdir $(wildcard $(APPS_DIR)/*))

# ---- kubectl (optionally via minikube) ----
KUBECTL        := $(if $(MINIKUBE),minikube kubectl --,kubectl)

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
	APP="$(1)"; \
	if [ -f "$(APPS_DIR)/$$APP/deployments/docker/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP/deployments/docker/Dockerfile"; \
	elif [ -f "$(APPS_DIR)/$$APP/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP/Dockerfile"; \
	else \
	  echo "❌ No Dockerfile for $$APP" >&2; exit 1; \
	fi
endef

# ----------------------------
# PHONY
# ----------------------------
.PHONY: help apps-list \
        docker-build-% docker-push-% docker-build-push-% \
        docker-build docker-push docker-build-push \
        docker-build-all docker-push-all docker-build-push-all \
        deploy-namespaces delete-namespaces \
        deploy-app-% deploy-app deploy-apps delete-app-% delete-app delete-apps \
        deploy-traefik delete-traefik deploy-gateway delete-gateway \
        deploy-networking delete-networking \
        deploy-redis delete-redis \
        deploy-observability delete-observability \
        deploy-all delete-all restart-app logs

# ============================
# Help
# ============================
help:
	@echo "🎛️  Root targets (ENV=$(ENV))"
	@echo "  apps-list                          # List discovered apps"
	@echo "  docker-build        APP=<name>     # Build one image"
	@echo "  docker-push         APP=<name>     # Push one image"
	@echo "  docker-build-push   APP=<name>     # Build + push one image"
	@echo "  docker-build-all                   # Build all images"
	@echo "  docker-push-all                    # Push all images"
	@echo "  docker-build-push-all              # Build + push all images"
	@echo "  deploy-namespaces                  # Apply namespaces (env-aware)"
	@echo "  delete-namespaces                  # Delete namespaces (env-aware)"
	@echo "  deploy-app         APP=<name>      # Build+push APP, then apply k8s for APP"
	@echo "  deploy-apps                        # Build+push ALL, then apply k8s for ALL"
	@echo "  delete-app         APP=<name>      # Delete k8s for APP"
	@echo "  delete-apps                        # Delete k8s for ALL"
	@echo "  deploy-networking                  # Traefik + Gateway (env-aware)"
	@echo "  delete-networking                  # Delete networking"
	@echo "  deploy-redis                       # Apply Redis (env-aware)"
	@echo "  delete-redis                       # Delete Redis"
	@echo "  deploy-observability               # Apply FB/Prom/Graf/OS (env-aware)"
	@echo "  delete-observability               # Delete observability"
	@echo "  deploy-all                         # Build+push ALL, namespaces, networking, infra, apps"
	@echo "  delete-all                         # Full teardown"
	@echo "  restart-app       APP=<name>       # Rollout restart by label app=<name>"
	@echo "  logs             APP=<name>        # Follow logs by label app=<name>"
	@echo ""
	@echo "🔧 Vars: REGISTRY=$(REGISTRY) TAG=$(TAG) MINIKUBE=$(MINIKUBE)"

apps-list:
	@echo "$(APPS)"

# ============================
# Docker — Pattern rules (DRY)
# ============================
# $* is the stem from the pattern rule match. Use it only in pattern/static rules. :contentReference[oaicite:3]{index=3}

docker-build-%:
	@DFILE="$$( $(call resolve_dockerfile,$*) )"; \
	echo "🐳 Building $(REGISTRY)/$*:$(TAG) using $$DFILE"; \
	docker build -t "$(REGISTRY)/$*:$(TAG)" -f "$$DFILE" "$(APPS_DIR)/$*"

docker-push-%:
	@echo "📤 Pushing $(REGISTRY)/$*:$(TAG)"; \
	docker push "$(REGISTRY)/$*:$(TAG)"

# IMPORTANT: give this its own recipe (don’t chain to docker-build-% directly),
# otherwise the stem for docker-build-% becomes 'push-<app>'.
docker-build-push-%:
	@$(MAKE) --no-print-directory docker-build APP=$*
	@$(MAKE) --no-print-directory docker-push  APP=$*

# Single-app shims reuse pattern rules
docker-build:      docker-build-$(APP)
docker-push:       docker-push-$(APP)
docker-build-push: docker-build-push-$(APP)

# Bulk targets depend on per-app goals (parallel with -j)
docker-build-all:      $(APPS:%=docker-build-%)
docker-push-all:       $(APPS:%=docker-push-%)
docker-build-push-all: $(APPS:%=docker-build-push-%)

# ============================
# Kubernetes — Namespaces (env-aware; overlays/<ENV>)
# ============================
deploy-namespaces:
	@set -e; ROOT="$(NAMESPACES_DIR)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "📦 Applying namespaces overlay: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) apply -k "$${ROOT}/overlays/$(ENV)"; \
	elif [ -d "$${ROOT}/$(ENV)" ]; then \
	  echo "📦 Applying namespaces overlay: $${ROOT}/$(ENV)"; \
	  $(KUBECTL) apply -k "$${ROOT}/$(ENV)"; \
	else \
	  echo "📦 Applying namespaces base: $${ROOT}/base"; \
	  $(KUBECTL) apply -k "$${ROOT}/base"; \
	fi
	@echo "✅ Namespaces applied for $(ENV)"

delete-namespaces:
	@set -e; ROOT="$(NAMESPACES_DIR)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🧹 Deleting namespaces overlay: $${ROOT}/overlays/$(ENV)"; \
	  $(KUBECTL) delete -k "$${ROOT}/overlays/$(ENV)" --ignore-not-found; \
	elif [ -d "$${ROOT}/$(ENV)" ]; then \
	  echo "🧹 Deleting namespaces overlay: $${ROOT}/$(ENV)"; \
	  $(KUBECTL) delete -k "$${ROOT}/$(ENV)" --ignore-not-found; \
	else \
	  echo "🧹 Deleting namespaces base: $${ROOT}/base"; \
	  $(KUBECTL) delete -k "$${ROOT}/base" --ignore-not-found; \
	fi
	@echo "🗑️  Namespaces deleted for $(ENV)"

# ============================
# Kubernetes — Apps (env-aware; overlays/<ENV>)
# ============================
deploy-app-%: docker-build-push-% deploy-namespaces
	$(call k8s_apply_or_base,$(K8S_DIR)/apps/$*)
	@echo "✅ Deployed app: $* (env=$(ENV))"

deploy-app: deploy-app-$(APP)
deploy-apps: $(APPS:%=deploy-app-%)
	@echo "🎉 All apps deployed (env=$(ENV))"

delete-app-%:
	$(call k8s_delete_or_base,$(K8S_DIR)/apps/$*)
	@echo "🗑️  Deleted app: $* (env=$(ENV))"
delete-app:  delete-app-$(APP)
delete-apps: $(APPS:%=delete-app-%)
	@echo "🧹 All apps deleted (env=$(ENV))"

# ============================
# Kubernetes — Networking (Traefik + Gateway)
# ============================
deploy-traefik: deploy-namespaces
ifeq ($(ENV),dev)
	$(call k8s_apply_or_base,$(TRAEFIK_DIR))
	@echo "✅ Traefik (dev) applied"
else ifeq ($(ENV),prod)
	@echo "⚙️  Applying HelmChartConfig to built-in k3s Traefik"
	$(KUBECTL) apply -k $(TRAEFIK_DIR)/overlays/prod
	@echo "✅ Traefik (prod) configured"
else
	@echo "❌ Unknown ENV=$(ENV). Use dev|prod"; exit 1
endif

delete-traefik:
ifeq ($(ENV),dev)
	$(call k8s_delete_or_base,$(TRAEFIK_DIR))
	@echo "🗑️  Traefik (dev) deleted"
else ifeq ($(ENV),prod)
	$(KUBECTL) delete -k $(TRAEFIK_DIR)/overlays/prod --ignore-not-found
	@echo "🗑️  Traefik (prod) HelmChartConfig removed"
else
	@echo "❌ Unknown ENV=$(ENV). Use dev|prod"; exit 1
endif

deploy-gateway: deploy-namespaces
	$(call k8s_apply_or_base,$(GATEWAY_DIR))
	@echo "✅ Gateway applied (ENV=$(ENV))"

delete-gateway:
	$(call k8s_delete_or_base,$(GATEWAY_DIR))
	@echo "🗑️  Gateway deleted (ENV=$(ENV))"

deploy-networking: deploy-traefik deploy-gateway
delete-networking: delete-gateway delete-traefik

# ============================
# Kubernetes — Redis & Observability
# ============================
deploy-redis: deploy-namespaces
	$(call k8s_apply_or_base,$(K8S_DIR)/redis)
	@echo "🧰 Redis deployed"

delete-redis:
	$(call k8s_delete_or_base,$(K8S_DIR)/redis)
	@echo "🧨 Redis deleted"

define apply_obs
	@[ -d "$(K8S_DIR)/observability/$(1)" ] && \
	  $(call k8s_apply_or_base,$(K8S_DIR)/observability/$(1)) || true
endef
define delete_obs
	@[ -d "$(K8S_DIR)/observability/$(1)" ] && \
	  $(call k8s_delete_or_base,$(K8S_DIR)/observability/$(1)) || true
endef

deploy-observability: deploy-namespaces
	$(call apply_obs,fluent-bit)
	$(call apply_obs,grafana)
	$(call apply_obs,prometheus)
	$(call apply_obs,opensearch)
	@echo "📈 Observability deployed"

delete-observability:
	$(call delete_obs,opensearch)
	$(call delete_obs,prometheus)
	$(call delete_obs,grafana)
	$(call delete_obs,fluent-bit)
	@echo "🧻 Observability deleted"

# ============================
# Everything (DRY)
# ============================
deploy-all: docker-build-push-all deploy-namespaces deploy-networking deploy-redis deploy-observability deploy-apps
	@echo "✅ Full stack applied (env=$(ENV))"

delete-all: delete-apps delete-observability delete-redis delete-networking delete-namespaces
	@echo "✅ Full stack deleted (env=$(ENV))"

# ============================
# Ops utilities
# ============================
restart-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@echo "🔁 Restarting deployments with label app=$(APP)"
	@$(KUBECTL) rollout restart deployment -l app=$(APP)

logs:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@$(KUBECTL) logs -l app=$(APP) -f --max-log-requests=5
