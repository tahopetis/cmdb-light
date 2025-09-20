# CMDB Lite Helm Chart

This Helm chart deploys CMDB Lite on a Kubernetes cluster.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.0+

## Installing the Chart

To install the chart with the release name `cmdb-lite`:

```bash
helm repo add cmdb-lite https://your-repo-url
helm install cmdb-lite cmdb-lite/cmdb-lite
```

## Uninstalling the Chart

To uninstall/delete the `cmdb-lite` deployment:

```bash
helm delete cmdb-lite
```

## Parameters

The following table lists the configurable parameters of the CMDB Lite chart and their default values.

| Parameter | Description | Default |
|-----------|-------------|---------|
| `replicaCount` | Number of replicas | `1` |
| `image.repository` | Image repository | `cmdb-lite` |
| `image.pullPolicy` | Image pull policy | `IfNotPresent` |
| `image.tag` | Image tag | `""` (defaults to Chart appVersion) |
| `imagePullSecrets` | Image pull secrets | `[]` |
| `nameOverride` | String to partially override release name | `""` |
| `fullnameOverride` | String to fully override release name | `""` |
| `serviceAccount.create` | Create service account | `true` |
| `serviceAccount.annotations` | Service account annotations | `{}` |
| `serviceAccount.name` | Service account name | `""` |
| `podAnnotations` | Pod annotations | `{}` |
| `podSecurityContext` | Pod security context | `{}` |
| `securityContext` | Container security context | `{}` |
| `service.type` | Service type | `ClusterIP` |
| `service.port` | Service port | `80` |
| `ingress.enabled` | Enable ingress | `false` |
| `ingress.className` | Ingress class name | `""` |
| `ingress.annotations` | Ingress annotations | `{}` |
| `ingress.hosts` | Ingress hosts | See values.yaml |
| `ingress.tls` | Ingress TLS configuration | `[]` |
| `resources` | Resource limits and requests | `{}` |
| `autoscaling.enabled` | Enable autoscaling | `false` |
| `autoscaling.minReplicas` | Minimum replicas for autoscaling | `1` |
| `autoscaling.maxReplicas` | Maximum replicas for autoscaling | `100` |
| `autoscaling.targetCPUUtilizationPercentage` | Target CPU utilization percentage | `80` |
| `nodeSelector` | Node selector | `{}` |
| `tolerations` | Tolerations | `[]` |
| `affinity` | Affinity | `{}` |
| `env.environment` | Environment (development/production) | `production` |
| `env.database.host` | Database host | `""` |
| `env.database.port` | Database port | `5432` |
| `env.database.name` | Database name | `cmdb_lite` |
| `env.database.user` | Database user | `cmdb_user` |
| `env.database.password` | Database password | `""` |
| `env.server.port` | Server port | `8080` |
| `env.server.jwtSecret` | JWT secret | `""` |
| `env.frontend.port` | Frontend port | `3000` |
| `postgresql.enabled` | Enable PostgreSQL | `true` |
| `postgresql.global.postgresql.auth.postgresPassword` | PostgreSQL password | `""` |
| `postgresql.global.postgresql.auth.database` | PostgreSQL database | `cmdb_lite` |
| `postgresql.global.postgresql.auth.username` | PostgreSQL username | `cmdb_user` |
| `postgresql.global.postgresql.auth.password` | PostgreSQL password | `""` |
| `backend.replicaCount` | Backend replica count | `1` |
| `backend.image.repository` | Backend image repository | `cmdb-lite-backend` |
| `backend.image.pullPolicy` | Backend image pull policy | `IfNotPresent` |
| `backend.image.tag` | Backend image tag | `latest` |
| `backend.service.type` | Backend service type | `ClusterIP` |
| `backend.service.port` | Backend service port | `8080` |
| `backend.resources` | Backend resource limits and requests | See values.yaml |
| `backend.livenessProbe` | Backend liveness probe | See values.yaml |
| `backend.readinessProbe` | Backend readiness probe | See values.yaml |
| `frontend.replicaCount` | Frontend replica count | `1` |
| `frontend.image.repository` | Frontend image repository | `cmdb-lite-frontend` |
| `frontend.image.pullPolicy` | Frontend image pull policy | `IfNotPresent` |
| `frontend.image.tag` | Frontend image tag | `latest` |
| `frontend.service.type` | Frontend service type | `ClusterIP` |
| `frontend.service.port` | Frontend service port | `80` |
| `frontend.resources` | Frontend resource limits and requests | See values.yaml |
| `frontend.livenessProbe` | Frontend liveness probe | See values.yaml |
| `frontend.readinessProbe` | Frontend readiness probe | See values.yaml |
| `adminer.enabled` | Enable Adminer | `false` |
| `adminer.image.repository` | Adminer image repository | `adminer` |
| `adminer.image.pullPolicy` | Adminer image pull policy | `IfNotPresent` |
| `adminer.image.tag` | Adminer image tag | `latest` |
| `adminer.service.type` | Adminer service type | `ClusterIP` |
| `adminer.service.port` | Adminer service port | `8080` |
| `adminer.ingress.enabled` | Enable Adminer ingress | `false` |
| `adminer.ingress.className` | Adminer ingress class name | `""` |
| `adminer.ingress.annotations` | Adminer ingress annotations | `{}` |
| `adminer.ingress.hosts` | Adminer ingress hosts | See values.yaml |
| `adminer.ingress.tls` | Adminer ingress TLS configuration | `[]` |
| `persistence.enabled` | Enable persistence | `true` |
| `persistence.storageClass` | Storage class | `""` |
| `persistence.accessMode` | Access mode | `ReadWriteOnce` |
| `persistence.size` | Storage size | `8Gi` |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```bash
helm install cmdb-lite \
  --set env.database.password=secretpassword \
  cmdb-lite/cmdb-lite
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
helm install cmdb-lite -f values.yaml cmdb-lite/cmdb-lite
```

## Persistence

The chart mounts a Persistent Volume for the PostgreSQL database. The volume is created using dynamic volume provisioning. If you want to disable persistence, set `persistence.enabled` to `false`.

## Ingress

This chart provides support for Ingress resources. If you want to enable Ingress, set `ingress.enabled` to `true` and configure the `ingress.hosts` and `ingress.tls` parameters.

## Adminer

Adminer is a database management tool that can be enabled by setting `adminer.enabled` to `true`. If enabled, you can access it at the specified port or through the Ingress if configured.

## Database Migrations

The chart does not automatically run database migrations. You need to run them manually after installing the chart. You can do this by running the following command:

```bash
kubectl exec -it deployment/cmdb-lite-backend -- ./scripts/db-migrate.sh all
```

## Backup and Restore

The chart does not include automatic backup and restore functionality. You need to implement this separately. You can use the provided scripts in the `scripts` directory to create backups and restore them.

## Security

The chart uses Kubernetes Secrets to store sensitive information like database passwords and JWT secrets. Make sure to set strong passwords and secrets.

## Monitoring

The chart does not include monitoring by default. You can add monitoring by configuring Prometheus and Grafana separately.