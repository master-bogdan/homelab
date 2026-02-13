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

## 4) Deploy secrets stack (Vault + External Secrets)

```bash
make deploy-secrets ENV=staging
```

### Vault setup (UI + join nodes)

1. Open `https://vault-setup.apps.<IP_DASH>.sslip.io`.
2. In UI choose `Create a new Raft cluster`, set shares/threshold (recommended `5` / `3`), save unseal keys + root token, and unseal `vault-0`.
3. Join follower nodes:

```bash
NS=staging-secrets
kubectl -n $NS exec -it vault-1 -- sh -c 'VAULT_ADDR=http://127.0.0.1:8200 vault operator raft join http://vault-0.vault-internal:8200'
kubectl -n $NS exec -it vault-2 -- sh -c 'VAULT_ADDR=http://127.0.0.1:8200 vault operator raft join http://vault-0.vault-internal:8200'
```

4. Unseal `vault-1` and `vault-2` with the same unseal keys.
5. Verify cluster:

```bash
kubectl -n $NS exec -it vault-0 -- sh -c 'VAULT_ADDR=http://127.0.0.1:8200 vault operator raft list-peers'
kubectl -n $NS get pods -l app.kubernetes.io/name=vault
```

6. In Vault UI:
   - Enable `KV` at path `kv` (version `v2`).
   - Enable auth method `Kubernetes` at path `kubernetes`.
   - Configure `kubernetes` auth with:
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
7. In UI create ACL policy `homelab-staging`:

```hcl
path "kv/data/staging/*" {
  capabilities = ["read"]
}
path "kv/metadata/staging/*" {
  capabilities = ["read", "list"]
}
```

8. In UI create Kubernetes role `homelab-staging`:
   - `bound_service_account_names`: `external-secrets`
   - `bound_service_account_namespaces`: `staging-secrets`
   - `token_policies`: `homelab-staging`
   - `ttl`: `1h`
9. In UI create KV v2 secrets under `kv/staging/*` (plain strings, no base64).

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
