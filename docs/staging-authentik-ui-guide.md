# Staging Authentik UI Guide

This guide is only for the `staging` environment.

It is not a general Authentik guide. It tells you what to click in the Authentik UI for the apps that already exist in this repo's staging manifests.

## Scope

This guide covers:

- Authentik in staging
- Grafana in staging
- Headlamp in staging
- OpenSearch and OpenSearch Dashboards in staging
- n8n in staging
- Prometheus in staging
- Vault UI in staging
- SeaweedFS dashboard in staging

This guide does not cover:

- the public notes API
- the public SeaweedFS bucket route
- dev or prod
- Terraform or blueprints

## Staging URLs used by the current repo

These are the hostnames currently committed in the staging overlays:

- Authentik: `https://auth.apps.10-96-11-221.sslip.io`
- Grafana: `https://grafana.apps.10-96-11-221.sslip.io`
- Headlamp: `https://headlamp.apps.10-96-11-221.sslip.io`
- OpenSearch Dashboards: `https://opensearch.apps.10-96-11-221.sslip.io`
- n8n: `https://n8n.apps.10-96-11-221.sslip.io`
- Prometheus: `https://prometheus.apps.10-96-11-221.sslip.io`
- Vault UI: `https://vault.apps.10-96-11-221.sslip.io`
- SeaweedFS dashboard: `https://seaweedfs.apps.10-96-11-221.sslip.io`

If you run `make update-minikube-ip ENV=staging`, the `10-96-11-221` part can change. When that happens, use the new hostname everywhere in Authentik.

## What you should create in Authentik for staging

You need three kinds of objects in Authentik:

1. groups and policies
2. OIDC applications/providers for apps that support OIDC
3. proxy applications/providers for apps that do not support OIDC well

For this repo, staging should look like this:

| App | Authentik pattern | Needs outpost? |
| --- | --- | --- |
| Grafana | OIDC | No |
| Headlamp | Proxy | Yes |
| OpenSearch + OpenSearch Dashboards | OIDC | No |
| n8n | Proxy | Yes |
| Prometheus | Proxy | Yes |
| Vault UI | Proxy | Yes |
| SeaweedFS dashboard | Proxy | Yes |

## Before you start in the Authentik UI

Make sure these are true first:

- Authentik is reachable at `https://auth.apps.10-96-11-221.sslip.io`
- you can log in to the Authentik admin UI
- Vault has the staging paths from `docs/staging-minikube.md`

You should also understand one repo-specific rule:

- the staging `HTTPRoute` manifests already expect the shared proxy outpost service to be called `ak-outpost-shared-outpost`

That means the outpost name inside Authentik should be:

- `shared-outpost`

If you create the outpost with a different name, the service name will differ and your existing staging routes will not match.

## Step 1: Create groups first

In Authentik UI:

1. Open `Directory -> Groups`.
2. Create these groups if they do not already exist:
   - `staging-internal-users`
   - `authentik Admins`
   - `opensearch-admins`
   - `opensearch-users`

Recommended meaning:

- `staging-internal-users`: baseline access to internal staging UIs
- `authentik Admins`: full admin access where the app maps admins
- `opensearch-admins`: full OpenSearch admin users
- `opensearch-users`: normal OpenSearch dashboard users

You can add your users to these groups now or later.

## Step 2: Read this naming rule before you create proxy apps

Do not create the outpost yet if it does not already exist.

Create the proxy applications first, then create the outpost and attach those applications to it.

### Important repo note

Your Git-managed staging routes already point to:

- `ak-outpost-shared-outpost`

So do not invent another outpost name unless you also plan to update these files:

- `k8s/observability/headlamp/overlays/staging/httproute-patch.yaml`
- `k8s/platform/n8n/overlays/staging/httproute-patch.yaml`
- `k8s/observability/prometheus/overlays/staging/httproute-patch.yaml`
- `k8s/secrets/vault/overlays/staging/httproute-ui-patch.yaml`
- `k8s/platform/seaweedfs/overlays/staging/httproute-dashboard-patch.yaml`
- `k8s/auth/authentik/overlays/staging/reference-grant-patch.yaml`

## Step 3: Create the OIDC app for Grafana

Grafana in staging is already configured to expect Authentik OIDC credentials from:

- Vault path: `kv/staging/grafana`
- keys:
  - `oidc_client_id`
  - `oidc_client_secret`

Relevant repo files:

- `k8s/observability/grafana/overlays/staging/values-staging.yaml`
- `k8s/observability/grafana/overlays/staging/external-secret.yaml`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `Grafana Staging`
   - Slug: `grafana-staging`
   - Launch URL: `https://grafana.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `OAuth2/OpenID Connect`
5. In the provider section, set:
   - Name: `Grafana Staging Provider`
   - Redirect URI: `https://grafana.apps.10-96-11-221.sslip.io/login/generic_oauth`
   - Client type: `Confidential`
   - Signing key: leave default or choose your default signing key
6. Save the application/provider.
7. Copy the generated client ID and client secret.

### Bind access

On the bindings step, or after creation:

- allow `staging-internal-users`
- keep `authentik Admins` as members of that group or bind them directly

### Make sure Grafana gets groups

Grafana staging expects:

- `groups`
- `email`
- `profile`

So in Authentik, make sure the provider emits the `groups` claim. If your provider only emits basic OpenID claims, add the mapping that includes groups.

### Store the secret in Vault

Put the Authentik values here:

- `kv/staging/grafana`
  - `oidc_client_id`
  - `oidc_client_secret`

Do not create a Kubernetes secret by hand. The repo already syncs this through External Secrets.

## Step 4: Create the OIDC app for OpenSearch and OpenSearch Dashboards

This repo uses one Authentik OIDC provider slug for both:

- OpenSearch
- OpenSearch Dashboards

That provider slug is:

- `opensearch`

Relevant repo files:

- `k8s/observability/opensearch/overlays/staging/values-staging.yaml`
- `k8s/observability/opensearch/overlays/staging/external-secret.yaml`
- `k8s/observability/opensearch-dashboards/overlays/staging/values-staging.yaml`
- `k8s/observability/opensearch-dashboards/overlays/staging/external-secret.yaml`

### Important rule

Do not create two unrelated Authentik OIDC providers unless you also change the manifests.

Right now both components point at:

- `https://auth.apps.10-96-11-221.sslip.io/application/o/opensearch/.well-known/openid-configuration`

So the simplest correct setup is:

- create one Authentik OIDC app/provider for `opensearch`
- store that client ID and client secret in `kv/staging/opensearch`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `OpenSearch Staging`
   - Slug: `opensearch`
   - Launch URL: `https://opensearch.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `OAuth2/OpenID Connect`
5. In the provider section, set:
   - Name: `OpenSearch Staging Provider`
   - Redirect URI: `https://opensearch.apps.10-96-11-221.sslip.io/auth/openid/login`
   - Client type: `Confidential`
6. Save and copy the client ID and client secret.

### Bind access

Allow these groups:

- `opensearch-users`
- `opensearch-admins`
- `authentik Admins`

### Make sure OpenSearch gets groups

This is required because the staging OpenSearch config maps roles from:

- `roles_key: groups`

So the Authentik provider must emit the `groups` claim.

### Store the OpenSearch OIDC client in one staging Vault path

Put these values here:

- `kv/staging/opensearch`
  - `oidc_client_id`
  - `oidc_client_secret`

The current staging manifests make both OpenSearch and OpenSearch Dashboards read the same Kubernetes secret sourced from `kv/staging/opensearch`.

### Important TLS note for OpenSearch and Dashboards

Authentik in staging is served with a certificate chained to the internal `homelab-ca`.

That means both components need a namespace-local CA bundle:

- OpenSearch Dashboards trusts the IdP with `opensearch_security.openid.root_ca`
- OpenSearch itself does not use that Dashboards setting
- OpenSearch must trust the IdP through the security plugin OIDC TLS block:
  - `openid_connect_idp.enable_ssl: true`
  - `openid_connect_idp.verify_hostnames: true`
  - `openid_connect_idp.pemtrustedcas_filepath: /usr/share/opensearch/config/oidc-ca/ca.crt`

If Dashboards redirects back from Authentik and then returns `401`, check OpenSearch logs first. The common failure is OpenSearch being unable to fetch IdP metadata or JWKS because the IdP CA is not configured correctly.

## Step 4b: Create the proxy app for Headlamp

Relevant repo files:

- `k8s/observability/headlamp/overlays/staging/clusterrolebinding.yaml`
- `k8s/observability/headlamp/overlays/staging/values-staging.yaml`
- `k8s/observability/headlamp/overlays/staging/httproute-patch.yaml`

Public URL:

- `https://headlamp.apps.10-96-11-221.sslip.io`

Expected internal service URL in this repo:

- `http://headlamp.staging-observability.svc.cluster.local`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `Headlamp Staging`
   - Slug: `headlamp-staging`
   - Launch URL: `https://headlamp.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `Proxy`
5. In the provider section, set:
   - Name: `Headlamp Staging Proxy`
   - External host: `https://headlamp.apps.10-96-11-221.sslip.io`
   - Internal host: `http://headlamp.staging-observability.svc.cluster.local`
6. Save.
7. If `shared-outpost` already exists, attach `Headlamp Staging` to it now.

### Bind access

Bind access to:

- `staging-internal-users`

### Important Headlamp note

This repo no longer uses Authentik OIDC for Headlamp.

The shared outpost protects the public route, and this repo runs Headlamp in-cluster with its own `headlamp` service account bound to the built-in Kubernetes `view` ClusterRole.

If Headlamp still prompts for a token after the Authentik redirect, generate one with:

```bash
kubectl -n staging-observability create token headlamp
```

Then use `Use A Token` in the Headlamp UI.

## Step 5: Create the proxy app for n8n

Relevant repo files:

- `k8s/platform/n8n/overlays/staging/values-staging.yaml`
- `k8s/platform/n8n/overlays/staging/httproute-patch.yaml`

Public URL:

- `https://n8n.apps.10-96-11-221.sslip.io`

Expected internal service URL in this repo:

- `http://n8n.staging-platform.svc.cluster.local:5678`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `n8n Staging`
   - Slug: `n8n-staging`
   - Launch URL: `https://n8n.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `Proxy`
5. In the provider section, set:
   - Name: `n8n Staging Proxy`
   - External host: `https://n8n.apps.10-96-11-221.sslip.io`
   - Internal host: `http://n8n.staging-platform.svc.cluster.local:5678`
6. Save the application/provider.

### Outpost note

If `shared-outpost` already exists, you can attach `n8n Staging` to it now.

If it does not exist yet, create the app first and attach it later in Step 9.

### Important warning for n8n

If you protect the whole hostname with Authentik, webhooks can break.

Typical webhook paths:

- `/webhook`
- `/webhook-test`

So for n8n you must choose one of these:

1. create a second public hostname for webhooks
2. allow webhook paths without auth in the proxy provider
3. keep all webhook traffic off the protected hostname

Do not ignore this. It is the most common mistake with n8n behind SSO.

## Step 6: Create the proxy app for Prometheus

Relevant repo files:

- `k8s/observability/prometheus/overlays/staging/httproute-patch.yaml`
- `k8s/observability/grafana/overlays/staging/values-staging.yaml`

Public URL:

- `https://prometheus.apps.10-96-11-221.sslip.io`

Expected internal service URL in this repo:

- `http://prometheus-server.staging-observability.svc.cluster.local`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `Prometheus Staging`
   - Slug: `prometheus-staging`
   - Launch URL: `https://prometheus.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `Proxy`
5. In the provider section, set:
   - Name: `Prometheus Staging Proxy`
   - External host: `https://prometheus.apps.10-96-11-221.sslip.io`
   - Internal host: `http://prometheus-server.staging-observability.svc.cluster.local`
6. Save.
7. If `shared-outpost` already exists, attach `Prometheus Staging` to it now.

Bind access to:

- `staging-internal-users`

## Step 7: Create the proxy app for Vault UI

Relevant repo files:

- `k8s/secrets/vault/overlays/staging/httproute-ui-patch.yaml`
- `k8s/secrets/vault/base/http-vault-ui.yaml`

Public URL:

- `https://vault.apps.10-96-11-221.sslip.io`

Confirmed internal service URL from the cluster:

- `http://vault-ui.staging-secrets.svc.cluster.local:8200`

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `Vault UI Staging`
   - Slug: `vault-ui-staging`
   - Launch URL: `https://vault.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `Proxy`
5. In the provider section, set:
   - Name: `Vault UI Staging Proxy`
   - External host: `https://vault.apps.10-96-11-221.sslip.io`
   - Internal host: `http://vault-ui.staging-secrets.svc.cluster.local:8200`
6. Save.
7. If `shared-outpost` already exists, attach `Vault UI Staging` to it now.

Bind access to:

- `staging-internal-users`

## Step 8: Create the proxy app for SeaweedFS dashboard

Relevant repo files:

- `k8s/platform/seaweedfs/overlays/staging/httproute-dashboard-patch.yaml`
- `k8s/platform/seaweedfs/overlays/staging/values-staging.yaml`

Public URL:

- `https://seaweedfs.apps.10-96-11-221.sslip.io`

Expected internal service URL in this repo:

- likely `http://seaweedfs-filer-client.staging-platform.svc.cluster.local:8888`

Because the staging values reference:

- `WEED_CLUSTER_SW_FILER: "seaweedfs-filer-client.staging-platform:8888"`

### Before creating it

If the staging platform stack is deployed, verify the exact service name:

```bash
kubectl -n staging-platform get svc | rg seaweedfs
```

Then use that verified internal URL in Authentik.

### Create it in Authentik UI

1. Open `Applications -> Applications`.
2. Click `Create with provider`.
3. In the application section, set:
   - Name: `SeaweedFS Staging`
   - Slug: `seaweedfs-staging`
   - Launch URL: `https://seaweedfs.apps.10-96-11-221.sslip.io`
4. Choose provider type:
   - `Proxy`
5. In the provider section, set:
   - Name: `SeaweedFS Staging Proxy`
   - External host: `https://seaweedfs.apps.10-96-11-221.sslip.io`
   - Internal host: use the verified filer UI service URL
6. Save.
7. If `shared-outpost` already exists, attach `SeaweedFS Staging` to it now.

Bind access to:

- `staging-internal-users`

## Step 9: Create the shared proxy outpost and attach the proxy apps

Only do this once for staging.

### Applications that should be attached to this outpost

When you create the outpost, add these applications:

- `Headlamp Staging`
- `n8n Staging`
- `Prometheus Staging`
- `Vault UI Staging`
- `SeaweedFS Staging`

### Create it in the UI

In Authentik UI:

1. Open `Applications -> Outposts`.
2. Click `Create`.
3. Set:
   - Name: `shared-outpost`
   - Type: `Proxy`
   - Integration: Kubernetes
4. Under `Applications`, select:
   - `Headlamp Staging`
   - `n8n Staging`
   - `Prometheus Staging`
   - `Vault UI Staging`
   - `SeaweedFS Staging`
5. Set the namespace to `staging-auth` if that setting is available.
6. Save.

### Verify the outpost exists in Kubernetes

After Authentik deploys it, verify:

```bash
kubectl -n staging-auth get deploy,svc | rg ak-outpost
```

You should see a service named:

- `ak-outpost-shared-outpost`

If you do not, stop here and fix the outpost before testing the proxy apps.

## Step 10: What not to create in Authentik

Do not create staging Authentik apps for these public routes:

- `notes.10-96-11-221.sslip.io`
- the public SeaweedFS bucket route on `storage.apps.10-96-11-221.sslip.io`

Those are intentionally public in the repo.

## Step 11: Quick validation after each app

After creating an app in Authentik, verify the matching Kubernetes side is already expecting it.

### For OIDC apps

Check the Vault path and ExternalSecret:

```bash
kubectl -n staging-observability get externalsecret
```

Then:

1. confirm the client ID/secret are in Vault
2. restart or redeploy the app if needed
3. test login through the public URL

### For proxy apps

Check:

```bash
kubectl -n staging-auth get svc
kubectl get httproutes -A | rg 'headlamp|n8n|prometheus|vault|seaweedfs'
```

Then:

1. open the public URL
2. make sure Authentik redirects you to login
3. log in
4. make sure you land inside the app, not on an Authentik error page

## Fast troubleshooting

### I created the proxy app but get a bad gateway

Usually one of these is wrong:

- the outpost name is not `shared-outpost`
- the outpost never deployed into `staging-auth`
- the internal host URL is wrong
- the app was not added to `shared-outpost`

### Grafana login returns but role mapping is wrong

Usually:

- the provider is not emitting `groups`
- the user is not in the expected Authentik group

### OpenSearch login works but permissions are wrong

Usually:

- the token does not contain `groups`
- the user is not in `opensearch-users` or `opensearch-admins`
- the `oidc_client_id` or `oidc_client_secret` in `kv/staging/opensearch` does not match the Authentik `opensearch` provider

### OpenSearch login returns from Authentik and then shows `401`

Usually:

- OpenSearch cannot validate the IdP TLS certificate
- the OIDC CA secret is missing in `staging-observability`
- the OpenSearch security config is missing `openid_connect_idp.pemtrustedcas_filepath`

Check:

- `kubectl -n staging-observability logs statefulset/homelab-opensearch-staging-master`
- `kubectl -n staging-observability get secret opensearch-oidc-ca`

The failure in logs usually looks like:

- `PKIX path building failed`
- `unable to find valid certification path to requested target`

### Headlamp passes Authentik but still asks for a token or shows limited data

Usually:

- the shared outpost is working, but Headlamp still needs Kubernetes credentials
- you have not provided the `headlamp` service account token yet
- or the built-in `view` ClusterRole is too limited for what you expect to see

Check:

- `kubectl -n staging-observability create token headlamp`
- `kubectl get clusterrolebinding staging-headlamp-view`
- `k8s/observability/headlamp/overlays/staging/clusterrolebinding.yaml`

### n8n login works but webhooks fail

That is expected if you protected the whole hostname and did not exempt webhook traffic.

Fix it by:

- splitting webhooks to a second hostname
- or allowlisting webhook paths in the proxy provider

## Source links

Official docs used for the UI flow and provider types:

- Authentik applications and provider creation: https://docs.goauthentik.io/providers/
- Create OAuth2/OIDC providers: https://docs.goauthentik.io/add-secure-apps/providers/oauth2/create-oauth2-provider/
- Authentik proxy providers: https://docs.goauthentik.io/add-secure-apps/providers/proxy/
- Authentik outposts: https://docs.goauthentik.io/add-secure-apps/outposts/
- Authentik Kubernetes integration: https://docs.goauthentik.io/add-secure-apps/outposts/integrations/kubernetes/
- Headlamp in-cluster access: https://headlamp.dev/docs/latest/installation/in-cluster/
- n8n OIDC setup and edition limits: https://docs.n8n.io/user-management/oidc/setup/
- n8n community edition features: https://docs.n8n.io/hosting/community-edition-features/
- n8n webhook behavior: https://docs.n8n.io/integrations/builtin/core-nodes/n8n-nodes-base.webhook/workflow-development/
