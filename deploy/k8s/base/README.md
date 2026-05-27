# EXBanka Kubernetes base skeleton

This folder contains the first Kubernetes deployment skeleton for EXBanka.

It is intentionally limited to the application layer shape:

- namespaces
- shared app ConfigMap
- Secret example
- one Deployment and Service per existing application service
- dev Mailhog Deployment and Service
- Ingress route map translated from the current frontend `nginx.conf`

It does not deploy PostgreSQL, Redis, CloudNativePG, real secrets, or production observability.

## Layout

```text
deploy/k8s/base/
  kustomization.yaml
  namespaces.yaml
  app-config.yaml
  app-secret.example.yaml
  deployments.yaml
  services.yaml
  ingress.yaml
  mailhog.yaml
```

## Prerequisites

- Kubernetes cluster
- `kubectl`
- nginx Ingress controller if applying `ingress.yaml`
- Image registry with built EXBanka service images
- PostgreSQL layer from the next PR before the application pods can become fully ready

## Secrets

Do not commit real secret values.

Create the app secret manually from secure values:

```powershell
kubectl create secret generic exbanka-app-secrets `
  --namespace exbanka `
  --from-literal=DB_USER=exbanka_app `
  --from-literal=DB_PASSWORD='<replace-me>' `
  --from-literal=JWT_SECRET='<replace-me>' `
  --from-literal=ALPHA_VANTAGE_KEY='<replace-me>' `
  --from-literal=REDIS_PASSWORD='<replace-me>'
```

`app-secret.example.yaml` exists only as a shape reference.

The `DB_USER` and `DB_PASSWORD` values must match the Postgres credentials created
for CloudNativePG in the `exbanka-db` namespace.

`REDIS_PASSWORD` must match the Redis credentials created from
`deploy/redis/redis-secret.example.yaml`.

## Validate locally

Render the base:

```powershell
kubectl kustomize .\deploy\k8s\base
```

Dry-run the rendered manifests:

```powershell
kubectl apply --dry-run=client -k .\deploy\k8s\base
```

## Health probes

Go services use separate probe endpoints:

- readiness probes call `/ready`, which checks PostgreSQL connectivity before a pod receives traffic.
- liveness probes call `/health`, which stays lightweight and does not depend on PostgreSQL or Redis.

The frontend keeps `/` for both probes because it only serves static assets/nginx routing.

## Apply order

For the full project rollout, use this order:

1. Apply this Kubernetes base skeleton.
2. Apply PostgreSQL/CloudNativePG from the Postgres replication PR.
3. Apply Redis from the Redis PR.
4. Build and push application images with immutable tags.
5. Replace `replace-with-git-sha` image tags with actual image tags or Kustomize overlays.
6. Apply app manifests and verify Ingress routes.

## Scaling rule

Keep these services at one replica until cron leadership is implemented:

- `exchange-service`
- `loan-service`

Both services own scheduled jobs today. Scaling them before leader election can run the same job multiple times.

## Current limitations

- The base points `DB_HOST` to the CloudNativePG read/write endpoint `exbanka-postgres-rw.exbanka-db.svc.cluster.local`.
- Application pods are not expected to become fully ready until the Postgres PR is applied.
- Redis values are reserved in the Secret example and consumed by Redis-aware services, but Redis is not a hard readiness dependency.
- This skeleton uses raw Kubernetes plus Kustomize because it is easy to validate with the currently available local tooling.
