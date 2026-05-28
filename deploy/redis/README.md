# EXBanka Redis deployment foundation

This folder contains the Kubernetes Redis deployment foundation for EXBanka.

Redis is introduced as shared cross-instance state for cache and coordination use
cases. It is not the system of record. Money, account balances, SAGA state, and
durable interbank state remain in PostgreSQL.

## Scope

This folder provides the Redis deployment foundation:

- Bitnami Redis Helm values
- Redis authentication via a Kubernetes Secret
- master plus one replica topology
- application ConfigMap placeholders for `REDIS_ADDR` and `REDIS_DB`
- rollout and verification documentation

Current backend usage:

- auth token revocation for employee and client logout

It does not make Redis a system of record and it does not yet add FX cache
replacement, rate limiting, or cron leadership.

## Prerequisites

- Kubernetes cluster
- `kubectl`
- Helm
- `exbanka` namespace from `deploy/k8s/base/namespaces.yaml`
- Bitnami Helm repository

Add the Bitnami repository:

```powershell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

## Secret

Do not commit real Redis credentials.

Create the Redis secret manually:

```powershell
kubectl create secret generic redis-credentials `
  --namespace exbanka `
  --from-literal=redis-password='<replace-me>'
```

`redis-secret.example.yaml` documents the expected shape only.

The app namespace secret `exbanka-app-secrets` also contains `REDIS_PASSWORD`.
Keep that value synchronized with `redis-credentials` until an external secret sync
mechanism is introduced.

## Install

Install Redis with a stable release name:

```powershell
helm upgrade --install exbanka-redis bitnami/redis `
  --namespace exbanka `
  --values .\deploy\redis\values.yaml
```

The current public Bitnami chart may emit a rolling-image warning depending on the
available free image catalog. Treat that as a production hardening follow-up: pin a
team-approved Redis image tag or digest before a real production rollout.

The application ConfigMap points `REDIS_ADDR` at:

```text
exbanka-redis-master.exbanka.svc.cluster.local:6379
```

This keeps all initial Redis traffic on the master. Read-from-replica can be added
later if metrics justify it.

## Verify

Check Redis pods and services:

```powershell
kubectl get pods -n exbanka -l app.kubernetes.io/instance=exbanka-redis
kubectl get svc -n exbanka -l app.kubernetes.io/instance=exbanka-redis
```

Expected services include:

- `exbanka-redis-master`
- `exbanka-redis-replicas`
- `exbanka-redis-headless`

Smoke test with the stored password:

```powershell
$redisPassword = kubectl get secret redis-credentials `
  -n exbanka `
  -o jsonpath="{.data.redis-password}"
$redisPassword = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($redisPassword))

kubectl run redis-client `
  --namespace exbanka `
  --rm `
  --tty `
  -i `
  --restart='Never' `
  --image docker.io/bitnami/redis:7.4 `
  --env REDISCLI_AUTH=$redisPassword `
  -- redis-cli -h exbanka-redis-master ping
```

Expected output:

```text
PONG
```

## Auth revocation smoke

After application pods are rolled out with `REDIS_ADDR`, `REDIS_DB`, and
`REDIS_PASSWORD`, verify that logout writes revocation entries and blocks token
reuse:

1. Login through `/api/v1/auth/login`.
2. Call a protected endpoint with the returned access token and expect success.
3. Call `/api/v1/auth/logout` with the access token and refresh token.
4. Reuse the old access token and expect `401`.
5. Reuse the old refresh token on `/api/v1/auth/refresh` and expect `401`.

Optional Redis key check:

```powershell
$redisPassword = kubectl get secret redis-credentials `
  -n exbanka `
  -o jsonpath="{.data.redis-password}"
$redisPassword = [Text.Encoding]::UTF8.GetString([Convert]::FromBase64String($redisPassword))

kubectl run redis-client `
  --namespace exbanka `
  --rm `
  --tty `
  -i `
  --restart='Never' `
  --image docker.io/bitnami/redis:7.4 `
  --env REDISCLI_AUTH=$redisPassword `
  -- redis-cli -h exbanka-redis-master keys "auth:revoked:jti:*"
```

## Failure model

Consumers must treat Redis as cache/coordination only:

- Redis outage must not crash services.
- Auth revocation checks should fail open and log loudly.
- FX cache should fall back to the current provider path.
- AlphaVantage rate limiting should fall back to the local limiter.
- No money or durable workflow state may be stored in Redis.

## Follow-up work

Next PRs should add backend integration gradually:

1. `exchange-service`: shared FX-rate cache.
2. `exchange-service`: distributed AlphaVantage rate limiter.
3. Cron leadership only after the team chooses Redis lock vs Kubernetes Lease.
