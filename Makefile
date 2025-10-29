# ============================
# Monorepo Makefile (root)
# Build & Push images + Apply K8s
# ============================

# ---- Config (override via env/CLI) ----
REGISTRY       ?= docker.io/masterbogdan0
TAG            ?= latest
ENV            ?= dev              # dev | prod
APP            ?=                  # target app, e.g. personal-website
PARALLEL_JOBS  ?= 1                # >1 to build/push in parallel
MINIKUBE       ?=                  # set 1 to use "minikube kubectl --"

# Namespaces
NS_APP         ?= $(ENV)
NS_DB          ?= db
NS_OBS         ?= observability

APPS_DIR       := apps
K8S_DIR        := k8s
APPS           := $(notdir $(wildcard $(APPS_DIR)/*))

# kubectl (optionally via minikube)
KUBECTL        := $(if $(MINIKUBE),minikube kubectl --,kubectl)

# ---- Helpers: prefer overlays/$(ENV), fallback to base ----
define k8s_apply_or_base
	@set -e; ROOT="$(1)"; NS="$(2)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🚀 Applying: $${ROOT}/overlays/$(ENV) (ns=$${NS})"; \
	  $(KUBECTL) apply -k "$${ROOT}/overlays/$(ENV)" -n "$${NS}"; \
	else \
	  echo "🚀 Applying: $${ROOT}/base (ns=$${NS})"; \
	  $(KUBECTL) apply -k "$${ROOT}/base" -n "$${NS}" 2>/dev/null || \
	  $(KUBECTL) apply -f "$${ROOT}/base" -n "$${NS}"; \
	fi
endef

define k8s_delete_or_base
	@set -e; ROOT="$(1)"; NS="$(2)"; \
	if [ -d "$${ROOT}/overlays/$(ENV)" ]; then \
	  echo "🔥 Deleting: $${ROOT}/overlays/$(ENV) (ns=$${NS})"; \
	  $(KUBECTL) delete -k "$${ROOT}/overlays/$(ENV)" -n "$${NS}" --ignore-not-found; \
	else \
	  echo "🔥 Deleting: $${ROOT}/base (ns=$${NS})"; \
	  $(KUBECTL) delete -k "$${ROOT}/base" -n "$${NS}" --ignore-not-found 2>/dev/null || \
	  $(KUBECTL) delete -f "$${ROOT}/base" -n "$${NS}" --ignore-not-found; \
	fi
endef

# Resolve Dockerfile per app (supports different structures)
# 1) apps/<app>/deployments/docker/Dockerfile
# 2) apps/<app>/Dockerfile
define resolve_dockerfile
	@set -e; APP="$(1)"; \
	if [ -f "$(APPS_DIR)/$$APP/deployments/docker/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP/deployments/docker/Dockerfile"; \
	elif [ -f "$(APPS_DIR)/$$APP/Dockerfile" ]; then \
	  echo "$(APPS_DIR)/$$APP/Dockerfile"; \
	else \
	  echo "❌ No Dockerfile for $$APP" >&2; exit 1; \
	fi
endef

.PHONY: help apps-list \
        docker-build docker-push docker-build-all docker-push-all \
        create-namespaces \
        deploy-app deploy-apps delete-app delete-apps \
        deploy-redis delete-redis \
        deploy-observability delete-observability \
        deploy-all delete-all restart-app logs

help:
	@echo "🎛️  Root targets (ENV=$(ENV))"
	@echo "  apps-list                         # List discovered apps"
	@echo "  docker-build  APP=<name>          # Build one image"
	@echo "  docker-push   APP=<name>          # Push one image"
	@echo "  docker-build-all                  # Build all images"
	@echo "  docker-push-all                   # Push all images"
	@echo "  create-namespaces                 # Ensure namespaces"
	@echo "  deploy-app    APP=<name>          # Apply k8s for one app"
	@echo "  deploy-apps                       # Apply k8s for all apps"
	@echo "  delete-app    APP=<name>          # Delete k8s for one app"
	@echo "  delete-apps                       # Delete k8s for all apps"
	@echo "  deploy-redis                      # Apply Redis"
	@echo "  delete-redis                      # Delete Redis"
	@echo "  deploy-observability              # Apply FB/Prom/Graf/OS"
	@echo "  delete-observability              # Delete observability"
	@echo "  deploy-all                        # Redis + Observability + All Apps"
	@echo "  delete-all                        # Tear everything down"
	@echo "  restart-app  APP=<name>           # Rollout restart by label app=<name>"
	@echo "  logs         APP=<name>           # Follow logs by label app=<name>"
	@echo ""
	@echo "🔧 Vars: REGISTRY=$(REGISTRY) TAG=$(TAG) MINIKUBE=$(MINIKUBE) PARALLEL_JOBS=$(PARALLEL_JOBS)"

apps-list:
	@echo "$(APPS)"

# ============================
# Docker: Build / Push
# ============================
docker-build:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@DFILE="$$( $(call resolve_dockerfile,$(APP)) )"; \
	echo "🐳 Building $(REGISTRY)/$(APP):$(TAG) using $$DFILE"; \
	docker build -t "$(REGISTRY)/$(APP):$(TAG)" -f "$$DFILE" "$(APPS_DIR)/$(APP)"

docker-push:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@echo "📤 Pushing $(REGISTRY)/$(APP):$(TAG)"
	@docker push "$(REGISTRY)/$(APP):$(TAG)"

docker-build-all:
	@echo "🐳 Building all: $(APPS)"
	@echo "$(APPS)" | xargs -n1 -P"$(PARALLEL_JOBS)" -I{} \
		sh -c 'DFILE="$$( \
		  if [ -f "$(APPS_DIR)"/{}/deployments/docker/Dockerfile ]; then echo "$(APPS_DIR)"/{}/deployments/docker/Dockerfile; \
		  elif [ -f "$(APPS_DIR)"/{}/Dockerfile ]; then echo "$(APPS_DIR)"/{}/Dockerfile; \
		  else echo "❌ No Dockerfile for {}" >&2; exit 1; fi \
		)"; \
		echo "— Building $(REGISTRY)/{}:$(TAG) from $$DFILE"; \
		docker build -t "$(REGISTRY)/{}:$(TAG)" -f "$$DFILE" "$(APPS_DIR)/{}"'

docker-push-all:
	@echo "📤 Pushing all: $(APPS)"
	@echo "$(APPS)" | xargs -n1 -P"$(PARALLEL_JOBS)" -I{} \
		sh -c 'echo "— Pushing $(REGISTRY)/{}:$(TAG)"; docker push "$(REGISTRY)/{}:$(TAG)"'

# ============================
# Kubernetes: Namespaces
# ============================
create-namespaces:
	-@$(KUBECTL) create namespace $(NS_APP) || true
	-@$(KUBECTL) create namespace $(NS_DB)  || true
	-@$(KUBECTL) create namespace $(NS_OBS) || true
	@echo "✅ Namespaces ready: $(NS_APP), $(NS_DB), $(NS_OBS)"

# ============================
# Kubernetes: Apps
# ============================
deploy-app: create-namespaces
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	$(call k8s_apply_or_base,$(K8S_DIR)/apps/$(APP),$(NS_APP))
	@echo "✅ Deployed app: $(APP) (ns=$(NS_APP), env=$(ENV))"

deploy-apps: create-namespaces
	@set -e; for a in $(APPS); do \
	  $(MAKE) --no-print-directory deploy-app APP=$${a}; \
	done
	@echo "🎉 All apps deployed"

delete-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	$(call k8s_delete_or_base,$(K8S_DIR)/apps/$(APP),$(NS_APP))
	@echo "🗑️  Deleted app: $(APP) (ns=$(NS_APP), env=$(ENV))"

delete-apps:
	@set -e; for a in $(APPS); do \
	  $(MAKE) --no-print-directory delete-app APP=$${a}; \
	done
	@echo "🧹 All apps deleted"

restart-app:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@echo "🔁 Restarting deployments with label app=$(APP) in ns $(NS_APP)"
	@$(KUBECTL) rollout restart deployment -n $(NS_APP) -l app=$(APP)

logs:
	@test -n "$(APP)" || (echo "❌ APP is required"; exit 1)
	@$(KUBECTL) logs -n $(NS_APP) -l app=$(APP) -f --max-log-requests=5

# ============================
# Kubernetes: Redis
# ============================
deploy-redis: create-namespaces
	$(call k8s_apply_or_base,$(K8S_DIR)/redis,$(NS_DB))
	@echo "🧰 Redis deployed (ns=$(NS_DB))"

delete-redis:
	$(call k8s_delete_or_base,$(K8S_DIR)/redis,$(NS_DB))
	@echo "🧨 Redis deleted (ns=$(NS_DB))"

# ============================
# Kubernetes: Observability
# ============================
define apply_obs
	@[ -d "$(K8S_DIR)/observability/$(1)" ] && \
	  $(call k8s_apply_or_base,$(K8S_DIR)/observability/$(1),$(NS_OBS)) || true
endef

define delete_obs
	@[ -d "$(K8S_DIR)/observability/$(1)" ] && \
	  $(call k8s_delete_or_base,$(K8S_DIR)/observability/$(1),$(NS_OBS)) || true
endef

deploy-observability: create-namespaces
	$(call apply_obs,fluent-bit)
	$(call apply_obs,grafana)
	$(call apply_obs,prometheus)
	$(call apply_obs,opensearch)
	@echo "📈 Observability deployed (ns=$(NS_OBS))"

delete-observability:
	$(call delete_obs,opensearch)
	$(call delete_obs,prometheus)
	$(call delete_obs,grafana)
	$(call delete_obs,fluent-bit)
	@echo "🧻 Observability deleted (ns=$(NS_OBS))"

# ============================
# Everything
# ============================
deploy-all: deploy-redis deploy-observability deploy-apps
	@echo "✅ Full stack applied (env=$(ENV))"

delete-all: delete-apps delete-observability delete-redis
	@echo "✅ Full stack deleted (env=$(ENV))"
