# Staging on Minikube (homelab-staging)

This guide brings up the **staging** environment on Minikube with **prod-like configs** and **HTTPS**.
Staging uses Vault + External Secrets (no hardcoded secrets) with a single Vault replica.

## Prereqs

- minikube
- kubectl
- helm
- Docker driver for minikube

## 1) Start Minikube profile

```bash
minikube start --profile homelab-staging --nodes=3 --driver=docker --cpus=4 --memory=8192
minikube addons enable storage-provisioner --profile homelab-staging
minikube addons enable default-storageclass --profile homelab-staging
```

Note: the Minikube `ingress` addon is not required here (we use Traefik + Gateway API).

Set kubectl context:

```bash
kubectl config use-context homelab-staging
kubectl config current-context
```

## 2) Update sslip.io hostnames for your Minikube IP

Staging TLS certs and routes use `*.sslip.io` with the Minikube IP encoded as dashes.
Run the helper target to rewrite all staging overlays + `Makefile` URLs in one step:

```bash
make update-minikube-ip ENV=staging MINIKUBE_PROFILE=homelab-staging
```

If you cannot query Minikube from your current shell/session, pass the IP directly:

```bash
ENV=staging PROFILE=homelab-staging MINIKUBE_IP=<minikube-ip> scripts/update-minikube-ip.sh
```

Examples:
- `MINIKUBE_IP=192.168.76.2`
- `MINIKUBE_IP=192.168.58.2`

## 3) Deploy namespaces + networking

```bash
make deploy-namespaces ENV=staging
make deploy-networking ENV=staging
```

This installs Gateway API, Traefik, cert-manager, and the gateway routes.

## 4) Start reverse proxy exposure in Minikube

Staging uses Traefik Service type `LoadBalancer`. On Minikube, run tunnel so Gateway traffic is reachable on ports `80/443`.

In a separate terminal (keep it running):

```bash
make minikube-tunnel ENV=staging MINIKUBE_PROFILE=homelab-staging
```

Quick check:

```bash
kubectl -n staging-networking get svc traefik-staging
```

The `EXTERNAL-IP` should be populated (not `<pending>`).

## 5) Deploy secrets stack (Vault + External Secrets)

```bash
make deploy-secrets ENV=staging
```

### Vault setup (UI, click-by-click)

1. Open Vault setup UI:
   - URL: `https://vault-setup.apps.<IP_DASH>.sslip.io/ui/`

2. Initialize Vault:
   - Screen: `Let's set up the initial set of root keys`
   - `Key shares`: `5`
   - `Key threshold`: `3`
   - Click `Initialize`
   - Save all generated:
     - `Unseal Key 1..5`
     - `Initial Root Token`

3. Unseal Vault:
   - Screen: `Unseal Vault`
   - Paste 3 different unseal keys (one by one), click `Unseal` after each.
   - When unsealed, login screen appears.
   - Login method: `Token`
   - Paste `Initial Root Token` and sign in.

4. Enable KV secrets engine:
   - Left menu: `Secrets Engines`
   - Click `Enable new engine`
   - Engine: `KV`
   - `Path`: `kv`
   - `Version`: `2`
   - Click `Enable Engine`

5. Enable Kubernetes auth method:
   - Left menu: `Access`
   - Tab: `Auth Methods`
   - Click `Enable new method`
   - Method: `Kubernetes`
   - `Path`: `kubernetes`
   - Click `Enable Method`

6. Configure Kubernetes auth method:
   - Open `Access` -> `Auth Methods` -> `kubernetes` -> `Configure`
   - Fill fields:
     - `Kubernetes host`:
       ```bash
       kubectl config view --raw --minify -o jsonpath='{.clusters[0].cluster.server}'; echo
       ```
     - `Kubernetes CA Certificate`:
       ```bash
       kubectl -n staging-secrets get configmap kube-root-ca.crt -o jsonpath='{.data.ca\.crt}'; echo
       ```
     - `Token Reviewer JWT`:
       ```bash
       kubectl -n staging-secrets create token external-secrets
       ```
   - Click `Save`

7. Create policy for External Secrets:
   - Left menu: `Policies`
   - Click `Create ACL policy`
   - Name: `homelab-staging`
   - Paste policy:

```hcl
path "kv/data/staging/*" {
  capabilities = ["read"]
}
path "kv/metadata/staging/*" {
  capabilities = ["read", "list"]
}
```

   - Click `Create policy`

8. Create Kubernetes auth role:
   - Left menu: `Access` -> `Auth Methods` -> `kubernetes`
   - Tab: `Roles`
   - Click `Create role`
   - `Role Name`: `homelab-staging`
   - `Bound service account names`: `external-secrets`
   - `Bound service account namespaces`: `staging-secrets`
   - `Token policies`: `homelab-staging`
   - `TTL`: `1h`
   - Click `Create`

9. Create secrets in UI (`kv/staging/*`):
   - Left menu: `Secrets Engines` -> `kv`
   - Click `Create secret`
   - For each path below, create a secret and add listed keys as plain string values:
     - `staging/authentik`
     - `staging/n8n`
     - `staging/grafana`
     - `staging/opensearch`
     - `staging/opensearch-dashboards`
     - `staging/fluent-bit`
     - `staging/ephermal-notes-api`
   - Click `Save` for each secret.

Paths and keys required:

**`kv/staging/authentik`**
- `secret_key`
- `postgresql_host`
- `postgresql_port`
- `postgresql_name`
- `postgresql_user`
- `postgresql_password`
- `redis_host`
- `redis_port`
- `redis_password`
- `email_host`
- `email_port`
- `email_username`
- `email_password`
- `email_from`
- `email_use_tls`
- `email_use_ssl`

**`kv/staging/n8n`**
- `postgres_password`
- `redis_password`
- `encryption_key`

**`kv/staging/grafana`**
- `admin_user`
- `admin_password`
- `oidc_client_id`
- `oidc_client_secret`

**`kv/staging/opensearch`**
- `admin_password`
- `oidc_client_id`
- `oidc_client_secret`

**`kv/staging/opensearch-dashboards`**
- `dashboards_username`
- `dashboards_password`
- `oidc_client_id`
- `oidc_client_secret`

**`kv/staging/fluent-bit`**
- `opensearch_username`
- `opensearch_password`

**`kv/staging/ephermal-notes-api`**
- `server_host`
- `server_port`
- `redis_addr`
- `redis_password`

## 6) Deploy core stacks

```bash
make deploy-auth ENV=staging
make deploy-databases ENV=staging
make deploy-platform ENV=staging
make deploy-observability ENV=staging
make deploy-apps ENV=staging
```

Optional: deploy static UI apps

```bash
make deploy-ui-all ENV=staging
```

## 7) Verify

```bash
kubectl get pods -A
kubectl get httproutes -A
kubectl get gateways -A
```

Staging example URLs (HTTPS):

```bash
IP=$(minikube ip --profile homelab-staging)
IP_DASH=${IP//./-}

echo "https://auth.apps.${IP_DASH}.sslip.io"
echo "https://grafana.apps.${IP_DASH}.sslip.io"
echo "https://opensearch.apps.${IP_DASH}.sslip.io"
echo "https://n8n.apps.${IP_DASH}.sslip.io"
echo "https://storage.apps.${IP_DASH}.sslip.io"
```

Note: TLS is issued by the internal `homelab-ca`, so your browser will warn unless you trust that CA.

## Troubleshooting

- If pods are stuck waiting for secrets, check ESO logs: `kubectl -n staging-secrets logs deploy/external-secrets`.
- If HTTPS fails, re-check the `sslip.io` IP in staging overlays and `k8s/networking/cert-manager/overlays/staging/certificate-gateway.yaml`.
- For local access without HTTPS, use port-forward:

```bash
kubectl -n staging-observability port-forward svc/grafana 3000:3000
kubectl -n staging-observability port-forward svc/prometheus 9090:9090
kubectl -n staging-platform port-forward svc/n8n 5678:5678
kubectl -n staging-auth port-forward svc/authentik 8000:8000
```
