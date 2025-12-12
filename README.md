# Homelab

Lightweight Kubernetes manifests and local app projects for a home cluster. Uses Kustomize (`base/` + `overlays/{dev,prod}`) for environment-specific deployments.

## Quick start

### Prerequisites

- **Docker** (build & push images)
- **kubectl** with Kustomize support (`kubectl apply -k ...`)
- **make**
- **Kubernetes cluster** (Minikube, k3s, k0s, bare metal, etc.) accessible via your active kubecontext
- Container registry access (default: `REGISTRY=docker.io/masterbogdan0`)

> Using Minikube? Start your cluster (example 3-node profile), then enable the storage addons once all nodes are ready (enabling them inline prevents secondary nodes from booting and kube-proxy from starting):
> ```bash
> minikube start --profile homelab-dev --nodes=3 --driver=docker --cpus=4 --memory=8192
> minikube addons enable storage-provisioner --profile homelab-dev
> minikube addons enable default-storageclass --profile homelab-dev
> ```
>
> If the profile gets wedged (e.g., kube-proxy or nodes never become Ready), delete and recreate it:
> ```bash
> minikube delete --profile homelab-dev
> minikube start --profile homelab-dev --nodes=3 --driver=docker --cpus=4 --memory=8192
> minikube addons enable storage-provisioner --profile homelab-dev
> minikube addons enable default-storageclass --profile homelab-dev
> ```

### Example workflow

```bash
# 1) Confirm kubectl points at the cluster (minikube, k3s, etc.)
kubectl config get-contexts
kubectl config use-context homelab

# 2) Build & push images for one app
make docker-build-push APP=ephermal-notes-api REGISTRY=docker.io/yourname TAG=latest

# 3) Deploy namespaces + app
make deploy-namespaces ENV=dev
make deploy-app APP=ephermal-notes-api ENV=dev REGISTRY=docker.io/yourname TAG=latest

# 4) Or deploy the full stack (dev — includes databases)
make deploy-all ENV=dev REGISTRY=docker.io/yourname TAG=latest

# Check deployments
kubectl get pods -A
kubectl logs -f -n <namespace> <pod>
```

## Networking & access control

- **Gateway API + Traefik:** `k8s/networking/gateway` defines the shared `homelab-gw`. Dev overlays ride the HTTP listener for speed, while prod HTTPRoutes bind to the HTTPS listener with TLS terminated by the `homelab-gw-tls` secret.
- **Authentik forward auth:** The `authentik-forward-auth` middleware (Traefik CR) forwards requests to `authentik-server.platform.svc.cluster.local/outpost.goauthentik.io/auth/traefik` and injects identity headers. Every sensitive platform/observability route (Grafana, n8n, OpenSearch Dashboards, SeaweedFS filer, etc.) references this middleware so users must authenticate; only the personal site and the ephemeral notes API remain public on purpose.
- **SeaweedFS routing:** `/personal-website` and `/internal-artifacts` continue to flow through the S3 route, while the filer UI now uses `http-seaweedfs-filer` and is SSO-protected. This avoids exposing the filer/console without Authentik.
- **Namespace scoping:** Both HTTP and HTTPS listeners restrict `allowedRoutes` to the `networking` namespace so misconfigured routes elsewhere cannot bind to the shared gateway.

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

Use `ENV=dev` or `ENV=prod` and make sure `kubectl` points at the intended cluster (Minikube profile, kubeadm, etc.).

- **Namespaces:**
- `make deploy-namespaces ENV=<env>`
- `make delete-namespaces ENV=<env>`
- `make validate-namespaces ENV=<env>`

**Networking (Traefik + Gateway):**
- `make deploy-networking ENV=<env>`
- `make delete-networking ENV=<env>`
- `make validate-networking ENV=<env>`

**Platform (Authentik, n8n, SeaweedFS):**
- `make deploy-platform ENV=<env>`
- `make delete-platform ENV=<env>`
- `make validate-platform ENV=<env>`

**Observability (Prometheus, Grafana, OpenSearch, Fluent Bit):**
- `make deploy-observability ENV=<env>`
- `make delete-observability ENV=<env>`
- `make validate-observability ENV=<env>`

**Databases (PostgreSQL, Redis — dev only):**
- `make deploy-databases ENV=<env>`
- `make delete-databases ENV=<env>`
- `make validate-databases ENV=<env>`

- **Applications (personal-website-ui, ephermal-notes-api):**
- `make deploy-app APP=<name> ENV=<env> [REGISTRY=<r>] [TAG=<t>]` — Build, push, then deploy
- `make delete-app APP=<name> ENV=<env>`
- `make validate-app APP=<name> ENV=<env>`
- `make deploy-apps ENV=<env> [REGISTRY=<r>] [TAG=<t>]` — Deploy all apps
- `make delete-apps ENV=<env>`
- `make validate-apps ENV=<env>`

### Full-stack deployments

- `make deploy-all ENV=dev [REGISTRY=<r>] [TAG=<t>]` — Deploy everything (with DBs, only dev)
- `make deploy-all ENV=prod [REGISTRY=<r>] [TAG=<t>]` — Deploy everything (no DBs, prod-ready)
- `make delete-all ENV=<env>` — Delete full stack
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
- `APP` — Application name (required for single-app targets like `deploy-app`, `docker-build`)

### Kubernetes overlays

- `k8s/*/base/` — Base manifests (shared)
- `k8s/*/overlays/dev/` — Dev-specific (with databases, resource limits, etc.)
- `k8s/*/overlays/prod/` — Prod-specific (no databases, higher resource limits)

The Makefile automatically prefers `overlays/$(ENV)` when present; falls back to `base/` if overlay not found.

## Notes

- **Cluster context**: Make sure `kubectl config current-context` points at the cluster you intend to manage (minikube, k3s, etc.) before running make targets.
- **Dev vs Prod**: `ENV=dev` includes PostgreSQL and Redis; `ENV=prod` omits them.
- **Validate before deploy**: Use `make validate-all ENV=dev` (dry-run) before deploying.
- **Kustomize**: All k8s manifests use Kustomize. Overlays can override base patches, replicas, resource limits, etc.

## License

MIT
