# CMDB Lite Deployment Guide

This guide explains how to deploy CMDB Lite using Docker Compose and Kubernetes.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Quick Start](#quick-start)
- [Docker Compose Deployment](#docker-compose-deployment)
  - [Development Environment](#development-environment)
  - [Production Environment](#production-environment)
- [Kubernetes Deployment](#kubernetes-deployment)
  - [Using Helm Chart](#using-helm-chart)
- [Database Management](#database-management)
  - [Migrations](#migrations)
  - [Backups and Restoration](#backups-and-restoration)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)

## Prerequisites

Before deploying CMDB Lite, ensure you have the following installed:

- [Docker](https://docs.docker.com/get-docker/) (version 20.10 or higher)
- [Docker Compose](https://docs.docker.com/compose/install/) (version 1.29 or higher)
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [OpenSSL](https://www.openssl.org/) (for generating secrets)

For Kubernetes deployment:
- [Kubernetes cluster](https://kubernetes.io/docs/setup/) (version 1.19 or higher)
- [Helm](https://helm.sh/docs/intro/install/) (version 3.0 or higher)
- [kubectl](https://kubernetes.io/docs/tasks/tools/) (configured to access your cluster)

## Quick Start

For a quick start in development mode:

```bash
# Clone the repository
git clone https://github.com/yourorg/cmdb-lite.git
cd cmdb-lite

# Run the initialization script
chmod +x scripts/init.sh
./scripts/init.sh
```

This will set up the environment, build the Docker images, start the services, run database migrations, and create an admin user.

## Docker Compose Deployment

### Development Environment

The development environment uses `docker-compose.dev.yml` for configuration, which includes hot-reload capabilities for both frontend and backend.

#### Starting the Development Environment

```bash
# Set environment variable
export ENVIRONMENT=development

# Start services
./scripts/deploy.sh start
```

#### Accessing the Services

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Adminer (Database UI): http://localhost:8081

#### Stopping the Development Environment

```bash
./scripts/deploy.sh stop
```

#### Viewing Logs

```bash
# View all logs
./scripts/deploy.sh logs

# View logs for a specific service
./scripts/deploy.sh logs backend
```

### Production Environment

The production environment uses `docker-compose.yml` for configuration, which includes optimized settings and security best practices.

#### Starting the Production Environment

```bash
# Set environment variable
export ENVIRONMENT=production

# Start services
./scripts/deploy.sh start
```

#### Accessing the Services

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

#### Stopping the Production Environment

```bash
./scripts/deploy.sh stop
```

#### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Environment
ENVIRONMENT=production

# Database
DB_HOST=db
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=secure-password

# Backend
SERVER_PORT=8080
JWT_SECRET=very-secure-jwt-secret

# Frontend
FRONTEND_PORT=3000
```

## Kubernetes Deployment

### Using Helm Chart

CMDB Lite includes a Helm chart for deployment on Kubernetes clusters.

#### Installing the Chart

1. Add the Helm repository:

```bash
helm repo add cmdb-lite https://your-repo-url
helm repo update
```

2. Create a values file for your environment:

```bash
cp k8s/helm/values.yaml my-values.yaml
```

3. Edit the `my-values.yaml` file to configure your deployment:

```yaml
# Set secure passwords and secrets
env:
  database:
    password: "secure-password"
  server:
    jwtSecret: "very-secure-jwt-secret"

# Configure resource limits
backend:
  resources:
    requests:
      memory: "256Mi"
      cpu: "250m"
    limits:
      memory: "512Mi"
      cpu: "500m"

frontend:
  resources:
    requests:
      memory: "128Mi"
      cpu: "125m"
    limits:
      memory: "256Mi"
      cpu: "250m"

# Enable ingress if needed
ingress:
  enabled: true
  className: "nginx"
  hosts:
    - host: cmdb-lite.example.com
      paths:
        - path: /
          pathType: Prefix
  tls:
    - secretName: cmdb-lite-tls
      hosts:
        - cmdb-lite.example.com
```

4. Install the chart:

```bash
helm install cmdb-lite cmdb-lite/cmdb-lite -f my-values.yaml
```

#### Upgrading the Chart

```bash
helm upgrade cmdb-lite cmdb-lite/cmdb-lite -f my-values.yaml
```

#### Uninstalling the Chart

```bash
helm uninstall cmdb-lite
```

#### Running Database Migrations

After installing the chart, you need to run database migrations:

```bash
# Get the backend pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=cmdb-lite-backend -o jsonpath='{.items[0].metadata.name}')

# Run migrations
kubectl exec -it $POD_NAME -- ./scripts/db-migrate.sh all
```

## Database Management

### Migrations

#### Creating a New Migration

```bash
# For Docker Compose
./scripts/db-migrate.sh create add_new_table

# For Kubernetes
# First, get the backend pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=cmdb-lite-backend -o jsonpath='{.items[0].metadata.name}')

# Then create the migration
kubectl exec -it $POD_NAME -- ./scripts/db-migrate.sh create add_new_table
```

#### Running Migrations

```bash
# For Docker Compose
./scripts/db-migrate.sh all

# For Kubernetes
# First, get the backend pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=cmdb-lite-backend -o jsonpath='{.items[0].metadata.name}')

# Then run migrations
kubectl exec -it $POD_NAME -- ./scripts/db-migrate.sh all
```

### Backups and Restoration

#### Creating a Backup

```bash
# For Docker Compose
./scripts/db-backup.sh create

# For Kubernetes
# First, get the database pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/instance=cmdb-lite-postgresql -o jsonpath='{.items[0].metadata.name}')

# Create a directory for backups on your local machine
mkdir -p ./backups

# Then create the backup
kubectl exec $POD_NAME -- pg_dump -U cmdb_user cmdb_lite > ./backups/cmdb_lite_backup_$(date +%Y%m%d_%H%M%S).sql
```

#### Restoring from a Backup

```bash
# For Docker Compose
./scripts/db-backup.sh restore ./backups/cmdb_lite_backup_20230101_120000.sql

# For Kubernetes
# First, get the database pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/instance=cmdb-lite-postgresql -o jsonpath='{.items[0].metadata.name}')

# Then restore the backup
kubectl exec -i $POD_NAME -- psql -U cmdb_user cmdb_lite < ./backups/cmdb_lite_backup_20230101_120000.sql
```

#### Listing Backups

```bash
# For Docker Compose
./scripts/db-backup.sh list

# For Kubernetes
ls -la ./backups
```

## Environment Variables

### Common Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `DB_HOST` | Database host | `db` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `cmdb_lite` |
| `DB_USER` | Database user | `cmdb_user` |
| `DB_PASSWORD` | Database password | `cmdb_password` |
| `SERVER_PORT` | Backend server port | `8080` |
| `JWT_SECRET` | JWT secret key | Randomly generated |
| `FRONTEND_PORT` | Frontend port | `3000` |
| `ADMINER_PORT` | Adminer port | `8081` |

### Backend Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GIN_MODE` | Gin mode (debug/release) | `debug` for development, `release` for production |

### Frontend Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VITE_API_URL` | Backend API URL | `http://localhost:8080` |
| `VITE_ENVIRONMENT` | Frontend environment | `development` |

## Troubleshooting

### Common Issues

#### Service Fails to Start

1. Check the logs:

```bash
# For Docker Compose
./scripts/deploy.sh logs [service-name]

# For Kubernetes
kubectl logs [pod-name]
```

2. Check if all required environment variables are set:

```bash
# For Docker Compose
docker-compose config

# For Kubernetes
kubectl get configmap [configmap-name] -o yaml
kubectl get secret [secret-name] -o yaml
```

#### Database Connection Issues

1. Check if the database is running:

```bash
# For Docker Compose
docker-compose ps db

# For Kubernetes
kubectl get pods -l app=postgresql
```

2. Check database logs:

```bash
# For Docker Compose
docker-compose logs db

# For Kubernetes
kubectl logs [postgresql-pod-name]
```

3. Check database connection:

```bash
# For Docker Compose
docker-compose exec db pg_isready -U cmdb_user -d cmdb_lite

# For Kubernetes
kubectl exec [postgresql-pod-name] -- pg_isready -U cmdb_user -d cmdb_lite
```

#### Migration Issues

1. Check if migrations are applied correctly:

```bash
# For Docker Compose
./scripts/db-migrate.sh status

# For Kubernetes
# First, get the database pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/instance=cmdb-lite-postgresql -o jsonpath='{.items[0].metadata.name}')

# Then check migration status
kubectl exec -it $POD_NAME -- ./scripts/db-migrate.sh status
```

2. Manually run migrations:

```bash
# For Docker Compose
./scripts/db-migrate.sh all

# For Kubernetes
# First, get the backend pod name
POD_NAME=$(kubectl get pods -l app.kubernetes.io/name=cmdb-lite-backend -o jsonpath='{.items[0].metadata.name}')

# Then run migrations
kubectl exec -it $POD_NAME -- ./scripts/db-migrate.sh all
```

### Getting Help

If you encounter any issues not covered in this guide, please:

1. Check the [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues) page for similar problems.
2. Create a new issue with detailed information about your problem.
3. Include relevant logs, environment details, and steps to reproduce the issue.