# EXBanka deployment runbook

This runbook describes the current Kubernetes rollout path for EXBanka.

It is intentionally operational: follow it from top to bottom when preparing a
local or team Kubernetes smoke deployment.

## 1. Scope

This runbook covers:

- Kubernetes base manifests in `deploy/k8s/base`
- CloudNativePG PostgreSQL cluster in `deploy/postgres`
- Bitnami Redis deployment in `deploy/redis`
- application readiness/liveness checks
- basic smoke tests after deployment
- common troubleshooting paths

It does not cover:

- production-grade external secret management
- image build/publish automation
- HPA/autoscaling
- cron leader election
- full CI/CD release promotion

## 2. Prerequisites

Required local tools:

- Docker Desktop or another Kubernetes-capable runtime
- Kubernetes cluster enabled and reachable via `kubectl`
- Helm
- nginx Ingress controller, if testing ingress routes
- CloudNativePG operator
- access to built EXBanka container images

Quick checks:

```powershell
kubectl cluster-info
kubectl get nodes
helm version
```

Expected:

- cluster API is reachable
- at least one node is `Ready`
- Helm prints client version information

## 3. Deployment order

Use this order. Do not apply app workloads before PostgreSQL secrets and the
CloudNativePG cluster are ready.

1. Render and apply base namespaces.
2. Create database and application secrets.
3. Install CloudNativePG operator.
4. Apply PostgreSQL cluster.
5. Install Redis.
6. Build and publish application images.
7. Replace placeholder image tags.
8. Apply application manifests.
9. Run smoke checks.

## 4. Base namespaces

From repository root:

```powershell
kubectl apply -f .\deploy\k8s\base\namespaces.yaml
```

Verify:

```powershell
kubectl get ns exbanka
kubectl get ns exbanka-db
```

## 5. Secrets

Never commit real secrets.

Create PostgreSQL credentials:

```powershell
kubectl create secret generic bankdb-credentials `
  --namespace exbanka-db `
  --type=kubernetes.io/basic-auth `
  --from-literal=username=exbanka_app `
  --from-literal=password='<replace-me>'
```

Create Redis credentials:

```powershell
kubectl create secret generic redis-credentials `
  --namespace exbanka `
  --from-literal=redis-password='<replace-me>'
```

Create application credentials. Keep `DB_PASSWORD` equal to
`bankdb-credentials.password`, and `REDIS_PASSWORD` equal to
`redis-credentials.redis-password`.

```powershell
kubectl create secret generic exbanka-app-secrets `
  --namespace exbanka `
  --from-literal=DB_USER=exbanka_app `
  --from-literal=DB_PASSWORD='<replace-me>' `
  --from-literal=JWT_SECRET='<replace-me>' `
  --from-literal=ALPHA_VANTAGE_KEY='<replace-me>' `
  --from-literal=REDIS_PASSWORD='<replace-me>'
```

Verify secrets exist:

```powershell
kubectl get secret -n exbanka-db bankdb-credentials
kubectl get secret -n exbanka redis-credentials
kubectl get secret -n exbanka exbanka-app-secrets
```

## 6. PostgreSQL

Install the CloudNativePG operator using the team-approved operator version.

After the operator is installed, render and apply the PostgreSQL manifests:

```powershell
kubectl kustomize .\deploy\postgres
kubectl apply -k .\deploy\postgres
```

Wait for the cluster:

```powershell
kubectl get cluster -n exbanka-db
kubectl get pods -n exbanka-db -w
```

Expected:

- CloudNativePG cluster `exbanka-postgres` exists
- two PostgreSQL pods become healthy
- one primary and one standby exist

Check services:

```powershell
kubectl get svc -n exbanka-db
```

Expected service used by applications:

```text
exbanka-postgres-rw.exbanka-db.svc.cluster.local
```

Primary/standby check:

```powershell
kubectl exec -n exbanka-db <postgres-pod> -- `
  psql -U exbanka_app -d bankdb -c "select pg_is_in_recovery();"
```

Expected:

- primary returns `f` or `false`
- standby returns `t` or `true`

## 7. Redis

Add and update the Bitnami chart repository:

```powershell
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update
```

Install Redis:

```powershell
helm upgrade --install exbanka-redis bitnami/redis `
  --namespace exbanka `
  --values .\deploy\redis\values.yaml
```

Verify pods and services:

```powershell
kubectl get pods -n exbanka -l app.kubernetes.io/instance=exbanka-redis
kubectl get svc -n exbanka -l app.kubernetes.io/instance=exbanka-redis
```

Expected services:

- `exbanka-redis-master`
- `exbanka-redis-replicas`
- `exbanka-redis-headless`

Redis smoke:

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

Expected:

```text
PONG
```

Redis is cache/coordination only. Application readiness must not depend on Redis.

## 8. Application manifests

Render base manifests:

```powershell
kubectl kustomize .\deploy\k8s\base
```

Dry-run:

```powershell
kubectl apply --dry-run=client -k .\deploy\k8s\base
```

Before applying workloads, replace every image tag:

```text
replace-with-git-sha
```

with the actual immutable image tag for the build being deployed.

Apply:

```powershell
kubectl apply -k .\deploy\k8s\base
```

Watch pods:

```powershell
kubectl get pods -n exbanka -w
```

Expected:

- app pods start
- readiness probes eventually pass
- frontend pod is ready
- exchange-service and loan-service stay at one replica

## 9. Health and readiness smoke

Go services expose:

- `/health` for liveness
- `/ready` for readiness

`/ready` checks PostgreSQL connectivity. Redis is not a hard readiness dependency.

Port-forward one service:

```powershell
kubectl port-forward -n exbanka svc/auth-service 8081:8081
```

In another terminal:

```powershell
Invoke-RestMethod http://localhost:8081/health
Invoke-RestMethod http://localhost:8081/ready
```

Expected:

- `/health` returns status `ok`
- `/ready` returns status `ready`

Repeat for key services if needed:

```powershell
kubectl port-forward -n exbanka svc/account-service 8084:8084
kubectl port-forward -n exbanka svc/exchange-service 8088:8088
kubectl port-forward -n exbanka svc/loan-service 8089:8089
```

Sample checks:

```powershell
Invoke-RestMethod http://localhost:8084/ready
Invoke-RestMethod http://localhost:8088/ready
Invoke-RestMethod http://localhost:8089/ready
```

## 10. Ingress smoke

If nginx Ingress is installed and `exbanka.local` points to the cluster:

```powershell
Invoke-RestMethod http://exbanka.local/
Invoke-RestMethod http://exbanka.local/api/v1/exchange/rates
```

If DNS is not configured, port-forward the frontend service:

```powershell
kubectl port-forward -n exbanka svc/frontend 8080:80
```

Then open:

```text
http://localhost:8080
```

## 11. Basic application smoke checklist

Run these after pods are ready:

- Frontend loads through gateway.
- Auth service `/ready` returns ready.
- Account service `/ready` returns ready.
- Exchange service `/ready` returns ready.
- Loan service `/ready` returns ready.
- Redis ping returns `PONG`.
- PostgreSQL cluster has one primary and one standby.
- No pod is in `CrashLoopBackOff`.
- No app pod is stuck in `0/1 Ready`.

Useful commands:

```powershell
kubectl get pods -n exbanka
kubectl get pods -n exbanka-db
kubectl describe pod -n exbanka <pod-name>
kubectl logs -n exbanka <pod-name>
kubectl get events -n exbanka --sort-by=.lastTimestamp
```

## 12. Troubleshooting

### Pod is running but not ready

Check readiness details:

```powershell
kubectl describe pod -n exbanka <pod-name>
kubectl logs -n exbanka <pod-name>
```

Most likely causes:

- PostgreSQL cluster is not ready.
- `DB_HOST` does not match the CloudNativePG read/write service.
- `DB_USER` or `DB_PASSWORD` in `exbanka-app-secrets` does not match `bankdb-credentials`.
- migrations failed during startup.

### Redis ping fails

Check:

```powershell
kubectl get pods -n exbanka -l app.kubernetes.io/instance=exbanka-redis
kubectl get secret -n exbanka redis-credentials
kubectl logs -n exbanka statefulset/exbanka-redis-master
```

Most likely causes:

- `redis-credentials` was not created before Helm install.
- `REDIS_PASSWORD` in `exbanka-app-secrets` does not match Redis credentials.
- Redis chart did not finish creating master/replica pods.

### App pod crashes on startup

Check logs:

```powershell
kubectl logs -n exbanka <pod-name> --previous
kubectl logs -n exbanka <pod-name>
```

Most likely causes:

- image tag still says `replace-with-git-sha`
- database connection failed
- migration failed
- missing required secret value

### Ingress route returns 404

Check:

```powershell
kubectl get ingress -n exbanka
kubectl describe ingress -n exbanka exbanka
kubectl get svc -n exbanka
```

Most likely causes:

- nginx Ingress controller is not installed
- `exbanka.local` does not resolve to the cluster
- wrong path or service mapping in `deploy/k8s/base/ingress.yaml`

## 13. Current limitations

Keep these limitations explicit:

- `exchange-service` stays at one replica until cron leadership is implemented.
- `loan-service` stays at one replica until cron leadership is implemented.
- Redis-backed auth token revocation is not fully implemented yet.
- Redis is not a system of record.
- CloudNativePG WAL archival/PITR is not configured yet.
- Image build/publish automation is still manual.
- Secrets are created manually, not through Sealed Secrets or External Secrets.

## 14. Definition of deploy smoke success

The deployment is smoke-ready when:

- all namespaces exist
- all required secrets exist
- PostgreSQL primary/standby is healthy
- Redis responds with `PONG`
- app manifests render and dry-run successfully
- all app pods are running
- backend service `/health` returns ok
- backend service `/ready` returns ready
- frontend loads through service or ingress
- no app pods are crash-looping
