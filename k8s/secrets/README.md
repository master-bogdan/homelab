# Secrets (Vault + External Secrets Operator)

This stack lives under `k8s/secrets` and provides:
- **Vault** (in-cluster) for secret storage
- **External Secrets Operator (ESO)** to sync Vault values into Kubernetes Secrets

## Deploy

```bash
make deploy-namespaces ENV=staging   # or ENV=prod/dev
make deploy-secrets ENV=staging
```

## After deploy (Vault bootstrap)

### 1) Initialize + unseal Vault

```bash
NS=staging-secrets # or dev-secrets / prod-secrets
kubectl -n $NS exec -it deploy/vault -- vault operator init
kubectl -n $NS exec -it deploy/vault -- vault operator unseal
kubectl -n $NS exec -it deploy/vault -- vault login <root-token>
```

> In `staging` we use dev mode (no unseal, root token in logs).

### 2) Enable KV v2 and Kubernetes auth

```bash
kubectl -n $NS exec -it deploy/vault -- vault secrets enable -path=kv kv-v2
kubectl -n $NS exec -it deploy/vault -- vault auth enable kubernetes
```

### 3) Configure Kubernetes auth (ESO)

```bash
K8S_HOST=$(kubectl config view --raw --minify --output 'jsonpath={.clusters[0].cluster.server}')
kubectl -n $NS create token external-secrets > /tmp/es.jwt
kubectl config view --raw --minify --output 'jsonpath={.clusters[0].cluster.certificate-authority-data}' \
  | base64 -d > /tmp/ca.crt

kubectl -n $NS exec -i deploy/vault -- vault write auth/kubernetes/config \
  token_reviewer_jwt=@/tmp/es.jwt \
  kubernetes_host="$K8S_HOST" \
  kubernetes_ca_cert=@/tmp/ca.crt
```

### 4) Create policy + role

```bash
kubectl -n $NS exec -i deploy/vault -- vault policy write homelab-prod - <<'EOF'
path "kv/data/prod/*" {
  capabilities = ["read"]
}
EOF

kubectl -n $NS exec -i deploy/vault -- vault write auth/kubernetes/role/homelab-prod \
  bound_service_account_names=external-secrets \
  bound_service_account_namespaces=$NS \
  policies=homelab-prod \
  ttl=1h
```

For `ENV=staging`, use role name `homelab-staging` and the same policy.

### 5) Put secrets into Vault

Fill in `.env` at repo root, then write values into Vault under `kv/prod/*`:

```bash
# Example (grafana)
kubectl -n $NS exec -i deploy/vault -- vault kv put kv/prod/grafana \
  admin_user="$VAULT_GRAFANA_ADMIN_USER" \
  admin_password="$VAULT_GRAFANA_ADMIN_PASSWORD" \
  oidc_client_id="$VAULT_GRAFANA_OIDC_CLIENT_ID" \
  oidc_client_secret="$VAULT_GRAFANA_OIDC_CLIENT_SECRET"
```

Repeat for each app based on `.env` and the ExternalSecret manifests.

## Vault UI routes

- **Setup (direct, no Authentik):** `vault-setup.*`
- **UI (via Authentik):** `vault.*`

Use the setup route only during bootstrap, then lock it down.

## Vault paths & keys expected

These are the exact keys ESO reads from Vault (KV v2):

- `kv/prod/authentik`
  - `secret_key`, `postgresql_host`, `postgresql_port`, `postgresql_name`, `postgresql_user`, `postgresql_password`
  - `redis_host`, `redis_port`, `redis_password`
  - `email_host`, `email_port`, `email_username`, `email_password`, `email_from`, `email_use_tls`, `email_use_ssl`
- `kv/prod/ephermal-notes-api`
  - `server_host`, `server_port`, `redis_addr`, `redis_password`
- `kv/prod/grafana`
  - `admin_user`, `admin_password`, `oidc_client_id`, `oidc_client_secret`
- `kv/prod/opensearch`
  - `admin_password`, `oidc_client_id`, `oidc_client_secret`
- `kv/prod/opensearch-dashboards`
  - `dashboards_username`, `dashboards_password`, `oidc_client_id`, `oidc_client_secret`
- `kv/prod/fluent-bit`
  - `opensearch_username`, `opensearch_password`
- `kv/prod/n8n`
  - `postgres_password`, `redis_password`, `encryption_key`
