# Homelab — compact

Lightweight Kubernetes manifests and local app projects for a home cluster. Uses Kustomize (`base/` + `overlays/{dev,prod}`) for environment-specific deployments.

## Quick start

### Prerequisites

- **Docker** (build & push images)
- **kubectl** with Kustomize support (`kubectl apply -k ...`)
- **make**
- **minikube** (for local cluster)
- Container registry access (default: `REGISTRY=docker.io/masterbogdan0`)

### Start Minikube (3-node cluster)

```bash
minikube start --profile homelab-dev --nodes=3 --driver=docker --cpus=4 --memory=8192
```

### Example workflow

```bash
# 1) Start minikube
minikube start --profile homelab-dev --nodes=3 --driver=docker --cpus=4 --memory=8192

# 2) Build & push images for one app
make docker-build-push APP=ephermal-notes-api REGISTRY=docker.io/yourname TAG=latest

# 3) Deploy namespaces + app (MINIKUBE=1 uses minikube's kubectl)
make deploy-namespaces ENV=dev MINIKUBE=1
make deploy-app APP=ephermal-notes-api ENV=dev REGISTRY=docker.io/yourname TAG=latest MINIKUBE=1

# 4) Or deploy the full stack (dev — includes databases)
make deploy-all ENV=dev REGISTRY=docker.io/yourname TAG=latest MINIKUBE=1

# Check deployments
kubectl get pods -A
kubectl logs -f -n <namespace> <pod>
```

## Makefile targets

### Docker (build & push)

- `make apps-list` — List discovered apps
- `make docker-build APP=<name> [REGISTRY=<r>] [TAG=<t>]` — Build one image
- `make docker-push APP=<name> [REGISTRY=<r>] [TAG=<t>]` — Push one image
- `make docker-build-push APP=<name> [REGISTRY=<r>] [TAG=<t>]` — Build & push one image
- `make docker-build-all [REGISTRY=<r>] [TAG=<t>]` — Build all images
- `make docker-push-all [REGISTRY=<r>] [TAG=<t>]` — Push all images
- `make docker-build-push-all [REGISTRY=<r>] [TAG=<t>]` — Build & push all images

### Kubernetes (deploy by layer)

Use `ENV=dev` or `ENV=prod` and set `MINIKUBE=1` to use minikube's kubectl context.

**Namespaces:**
- `make deploy-namespaces ENV=<env> [MINIKUBE=1]`
- `make delete-namespaces ENV=<env> [MINIKUBE=1]`
- `make validate-namespaces ENV=<env>`

**Networking (Traefik + Gateway):**
- `make deploy-networking ENV=<env> [MINIKUBE=1]`
- `make delete-networking ENV=<env> [MINIKUBE=1]`
- `make validate-networking ENV=<env>`

**Platform (Authentik, n8n, SeaweedFS):**
- `make deploy-platform ENV=<env> [MINIKUBE=1]`
- `make delete-platform ENV=<env> [MINIKUBE=1]`
- `make validate-platform ENV=<env>`

**Observability (Prometheus, Grafana, OpenSearch, Fluent Bit):**
- `make deploy-observability ENV=<env> [MINIKUBE=1]`
- `make delete-observability ENV=<env> [MINIKUBE=1]`
- `make validate-observability ENV=<env>`

**Databases (PostgreSQL, Redis — dev only):**
- `make deploy-databases ENV=<env> [MINIKUBE=1]`
- `make delete-databases ENV=<env> [MINIKUBE=1]`
- `make validate-databases ENV=<env>`

**Applications (personal-website-ui, ephermal-notes-api):**
- `make deploy-app APP=<name> ENV=<env> [REGISTRY=<r>] [TAG=<t>] [MINIKUBE=1]` — Build, push, then deploy
- `make delete-app APP=<name> ENV=<env> [MINIKUBE=1]`
- `make validate-app APP=<name> ENV=<env>`
- `make deploy-apps ENV=<env> [REGISTRY=<r>] [TAG=<t>] [MINIKUBE=1]` — Deploy all apps
- `make delete-apps ENV=<env> [MINIKUBE=1]`
- `make validate-apps ENV=<env>`

### Full-stack deployments

- `make deploy-all ENV=dev [REGISTRY=<r>] [TAG=<t>] [MINIKUBE=1]` — Deploy everything (with DBs, only dev)
- `make deploy-all ENV=prod [REGISTRY=<r>] [TAG=<t>] [MINIKUBE=1]` — Deploy everything (no DBs, prod-ready)
- `make delete-all ENV=<env> [MINIKUBE=1]` — Delete full stack
- `make validate-all ENV=<env>` — Dry-run validate all manifests

## Repository structure

```
.
├── apps/                          # Application source code with Dockerfiles
│   ├── ephermal-notes-api/        # Go backend API
│   └── personal-website-ui/       # Next.js frontend
├── k8s/                           # Kubernetes manifests (Kustomize)
│   ├── namespaces/                # Namespace definitions (dev/prod overlays)
│   ├── apps/                      # App deployments (base + overlays)
│   ├── networking/                # Traefik & Ingress Gateway
│   ├── platform/                  # Authentik, n8n, SeaweedFS
│   ├── observability/             # Prometheus, Grafana, OpenSearch, Fluent Bit
│   └── databases/                 # PostgreSQL, Redis (dev only)
├── docs/                          # Architecture diagrams
├── .github/workflows/             # CI/CD
├── Makefile                       # Root make targets
└── README.md                      # This file
```

## Dashboard URLs & access

All services use host-based ingress. Choose one access method below:

### Services by namespace

| Service | Host (example) | Namespace | k8s path |
|---------|---|---|---|
| **Authentik** | `auth.<your-domain>` | `platform` | `k8s/platform/authentik` |
| **n8n** | `n8n.<your-domain>` | `platform` | `k8s/platform/n8n` |
| **Grafana** | `grafana.<your-domain>` | `observability` | `k8s/observability/grafana` |
| **Prometheus** | `prometheus.<your-domain>` | `observability` | `k8s/observability/prometheus` |
| **OpenSearch Dashboards** | `opensearch.<your-domain>` | `observability` | `k8s/observability/opensearch-dashboards` |

### Access methods (no `/etc/hosts` edits needed)

**Option 1: Wildcard DNS (nip.io / sslip.io)**
```bash
# Get your Minikube IP
IP=$(minikube ip --profile homelab-dev)

# Access services using nip.io (example with 192.168.49.2)
# https://auth.192.168.49.2.nip.io
# https://n8n.192.168.49.2.nip.io
# http://grafana.192.168.49.2.nip.io
```

**Option 2: minikube tunnel (LoadBalancer services)**
```bash
# In a separate terminal, start tunnel
minikube tunnel --profile homelab-dev

# Then access services on localhost with appropriate ports
# (ports depend on ingress/service configuration)
```

**Option 3: minikube service (direct service URL)**
```bash
# Open service directly in browser (bypasses ingress)
minikube service --url -n observability grafana --profile homelab-dev
minikube service --url -n platform n8n --profile homelab-dev
```

**Option 4: kubectl port-forward (fallback)**
```bash
# Forward local port to service
kubectl -n observability port-forward svc/grafana 3000:3000
kubectl -n observability port-forward svc/prometheus 9090:9090
kubectl -n observability port-forward svc/opensearch-dashboards 5601:5601
kubectl -n platform port-forward svc/n8n 5678:5678
kubectl -n platform port-forward svc/authentik 8000:8000

# Then access on localhost:port
```

### Discover actual ingress hosts

```bash
# List all ingresses across namespaces
kubectl get ingress -A

# Check specific namespace
kubectl -n observability get ingress
kubectl -n platform get ingress
```

## Configuration

### Environment variables (make)

- `ENV` — Deployment environment: `dev` (includes DBs) or `prod` (no DBs). Default: `dev`
- `REGISTRY` — Container registry. Default: `docker.io/masterbogdan0`
- `TAG` — Image tag. Default: `latest`
- `MINIKUBE` — Set to `1` to use `minikube kubectl --` instead of plain `kubectl`
- `APP` — Application name (required for single-app targets like `deploy-app`, `docker-build`)

### Kubernetes overlays

- `k8s/*/base/` — Base manifests (shared)
- `k8s/*/overlays/dev/` — Dev-specific (with databases, resource limits, etc.)
- `k8s/*/overlays/prod/` — Prod-specific (no databases, higher resource limits)

The Makefile automatically prefers `overlays/$(ENV)` when present; falls back to `base/` if overlay not found.

## Notes

- **MINIKUBE=1 required**: After starting minikube, always use `MINIKUBE=1` with make deploy targets so kubectl communicates with the minikube cluster.
- **Dev vs Prod**: `ENV=dev` includes PostgreSQL and Redis; `ENV=prod` omits them.
- **Validate before deploy**: Use `make validate-all ENV=dev` (dry-run) before deploying.
- **Kustomize**: All k8s manifests use Kustomize. Overlays can override base patches, replicas, resource limits, etc.

## License

MIT
