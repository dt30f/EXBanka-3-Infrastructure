# EXBanka PostgreSQL replication foundation

This folder contains the CloudNativePG foundation for the shared EXBanka PostgreSQL database.

It creates one logical database cluster:

- cluster name: `exbanka-postgres`
- database: `bankdb`
- namespace: `exbanka-db`
- instances: `2`
- topology: one primary plus one hot standby
- synchronization target: one standby

All EXBanka services still use one database, `bankdb`. This PR does not split the database per service and does not add application-level read/write routing.

## Prerequisites

- Kubernetes cluster
- `kubectl`
- CloudNativePG operator installed
- `exbanka-db` namespace from `deploy/k8s/base/namespaces.yaml`
- StorageClass that supports `ReadWriteOnce`

Install the CloudNativePG operator before applying this folder. Use the official installation method for the version selected by the team.

## Secrets

Do not commit real database credentials.

Create the database secret manually:

```powershell
kubectl create secret generic bankdb-credentials `
  --namespace exbanka-db `
  --type=kubernetes.io/basic-auth `
  --from-literal=username=exbanka_app `
  --from-literal=password='<replace-me>'
```

`bankdb-secret.example.yaml` documents the required shape only.

Because Kubernetes secrets are namespace-scoped, the application namespace still needs
`exbanka-app-secrets` with matching `DB_USER` and `DB_PASSWORD` values. Keep the
database credentials in `exbanka-db/bankdb-credentials` and
`exbanka/exbanka-app-secrets` synchronized until an external secret sync mechanism is
introduced.

## Apply

Render manifests:

```powershell
kubectl kustomize .\deploy\postgres
```

Apply after the CloudNativePG operator is installed:

```powershell
kubectl apply -k .\deploy\postgres
```

## Verify readiness

```powershell
kubectl get cluster -n exbanka-db
kubectl get pods -n exbanka-db
kubectl get svc -n exbanka-db
```

CloudNativePG creates stable services similar to:

- `exbanka-postgres-rw` for the current primary
- `exbanka-postgres-ro` for read-only replicas
- `exbanka-postgres-r` for all Postgres instances

The application base config points `DB_HOST` at:

```text
exbanka-postgres-rw.exbanka-db.svc.cluster.local
```

That is intentional. Until repository-level read routing is reviewed, all services should keep using the primary endpoint.

## Primary/standby check

Run on each Postgres pod:

```powershell
kubectl exec -n exbanka-db <postgres-pod> -- `
  psql -U exbanka_app -d bankdb -c "select pg_is_in_recovery();"
```

Expected:

- one pod returns `false` and is the primary
- one pod returns `true` and is the standby

## Basic failover smoke

Delete the current primary pod:

```powershell
kubectl delete pod -n exbanka-db <primary-pod>
kubectl get pods -n exbanka-db -w
```

Expected:

- standby is promoted automatically
- `exbanka-postgres-rw` follows the new primary
- cluster returns to a healthy state

## Out of scope

- Redis
- GORM `dbresolver`
- `DB_READ_HOST`
- read/write split in application code
- WAL archival / PITR object store
- multi-region disaster recovery

WAL archival and PITR should be added as a follow-up once the team chooses MinIO/S3-compatible storage.
