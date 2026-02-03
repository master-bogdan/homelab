# Staging on Minikube (homelab-staging)

This guide brings up the **staging** environment on Minikube with **prod-like configs** and **HTTPS**.
Staging uses Vault + External Secrets (no hardcoded secrets) and HA Vault in raft mode.

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
If your IP is not `192.168.58.2`, update the staging overlays + Makefile.

```bash
IP=$(minikube ip --profile homelab-staging)
IP_DASH=${IP//./-}

# Replace all staging sslip.io hostnames + Makefile base URL
rg -l "192-168-58-2" k8s/**/overlays/staging Makefile | xargs sed -i "s/192-168-58-2/${IP_DASH}/g"
```

If you're on macOS, use `sed -i ''` instead of `sed -i`.

## 3) Deploy namespaces + networking

```bash
make deploy-namespaces ENV=staging
make deploy-networking ENV=staging
```

This installs Gateway API, Traefik, cert-manager, and the gateway routes.

## 4) Deploy secrets stack (Vault + External Secrets)

```bash
make deploy-secrets ENV=staging
```

### 4a) Initialize and unseal Vault (HA)

Staging uses HA raft. Initialize **once** on `vault-0`, then unseal **each** pod with
the required number of unseal keys (default: 3 of 5).

```bash
NS=staging-secrets
kubectl -n $NS get pods

kubectl -n $NS exec -it vault-0 -- vault operator init

# Unseal vault-0 with 3 different keys:
kubectl -n $NS exec -it vault-0 -- vault operator unseal <unseal-key-1>
kubectl -n $NS exec -it vault-0 -- vault operator unseal <unseal-key-2>
kubectl -n $NS exec -it vault-0 -- vault operator unseal <unseal-key-3>

# Repeat for vault-1 and vault-2:
kubectl -n $NS exec -it vault-1 -- vault operator unseal <unseal-key-1>
kubectl -n $NS exec -it vault-1 -- vault operator unseal <unseal-key-2>
kubectl -n $NS exec -it vault-1 -- vault operator unseal <unseal-key-3>

kubectl -n $NS exec -it vault-2 -- vault operator unseal <unseal-key-1>
kubectl -n $NS exec -it vault-2 -- vault operator unseal <unseal-key-2>
kubectl -n $NS exec -it vault-2 -- vault operator unseal <unseal-key-3>

kubectl -n $NS exec -it vault-0 -- vault login <root-token>
```

### 4b) Enable KV v2 + Kubernetes auth

```bash
kubectl -n $NS exec -it vault-0 -- vault secrets enable -path=kv kv-v2
kubectl -n $NS exec -it vault-0 -- vault auth enable kubernetes
```

### 4c) Configure Kubernetes auth for External Secrets

```bash
K8S_HOST=$(kubectl config view --raw --minify --output 'jsonpath={.clusters[0].cluster.server}')
kubectl -n $NS create token external-secrets > /tmp/es.jwt
kubectl config view --raw --minify --output 'jsonpath={.clusters[0].cluster.certificate-authority-data}' \
  | base64 -d > /tmp/ca.crt

kubectl -n $NS exec -i vault-0 -- vault write auth/kubernetes/config \
  token_reviewer_jwt=@/tmp/es.jwt \
  kubernetes_host="$K8S_HOST" \
  kubernetes_ca_cert=@/tmp/ca.crt
```

### 4d) Create policy + role for staging

```bash
ENV=staging
kubectl -n $NS exec -i vault-0 -- vault policy write homelab-$ENV - <<EOF
path "kv/data/$ENV/*" {
  capabilities = ["read"]
}
EOF

kubectl -n $NS exec -i vault-0 -- vault write auth/kubernetes/role/homelab-$ENV \
  bound_service_account_names=external-secrets \
  bound_service_account_namespaces=$NS \
  policies=homelab-$ENV \
  ttl=1h
```

### 4e) Add secrets in the Vault UI (recommended)

Open the Vault UI and create KV v2 secrets under `kv/staging/*`.

1) Open the UI: `https://vault.apps.<IP_DASH>.sslip.io`
2) Login with the **root token** from `vault operator init`.
3) Go to **Secrets Engines** → ensure **KV v2** is enabled at path `kv`.
4) Create the following secrets (plain strings; no base64).

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

## 5) Deploy core stacks

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

## 6) Verify

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
