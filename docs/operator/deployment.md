# CMDB Lite Deployment Guide

This document provides detailed instructions for deploying CMDB Lite in various environments, including development, staging, and production.

## Table of Contents

- [Deployment Options](#deployment-options)
- [Prerequisites](#prerequisites)
- [Environment Setup](#environment-setup)
- [Docker Compose Deployment](#docker-compose-deployment)
  - [Development Environment](#development-environment)
  - [Production Environment](#production-environment)
- [Kubernetes Deployment](#kubernetes-deployment)
  - [Prerequisites](#prerequisites-1)
  - [Deployment Steps](#deployment-steps)
  - [Configuration](#configuration)
  - [Scaling and High Availability](#scaling-and-high-availability)
- [Traditional VM Deployment](#traditional-vm-deployment)
  - [Prerequisites](#prerequisites-2)
  - [Backend Deployment](#backend-deployment)
  - [Frontend Deployment](#frontend-deployment)
  - [Database Setup](#database-setup)
- [Post-Deployment Steps](#post-deployment-steps)
- [Configuration Management](#configuration-management)
- [Monitoring and Logging](#monitoring-and-logging)
- [Backup and Recovery](#backup-and-recovery)
- [Troubleshooting](#troubleshooting)

## Deployment Options

CMDB Lite can be deployed in several ways, depending on your infrastructure and requirements:

1. **Docker Compose**: Simple single-host deployment suitable for development and small production environments
2. **Kubernetes**: Scalable, containerized deployment suitable for medium to large production environments
3. **Traditional VM**: Direct installation on virtual machines suitable for organizations with existing VM infrastructure

### Choosing a Deployment Option

| Deployment Option | Best For | Complexity | Scalability | High Availability |
|-------------------|----------|------------|-------------|-------------------|
| Docker Compose | Development, small production | Low | Limited | No |
| Kubernetes | Medium to large production | High | Excellent | Yes |
| Traditional VM | Organizations with existing VM infrastructure | Medium | Moderate | Requires additional setup |

## Prerequisites

Before deploying CMDB Lite, ensure you have the following prerequisites:

### System Requirements

- **CPU**: 2 cores minimum, 4 cores recommended
- **Memory**: 4 GB RAM minimum, 8 GB recommended
- **Storage**: 20 GB available disk space minimum, 50 GB recommended
- **Database**: PostgreSQL 12 or higher
- **Operating System**: Linux (Ubuntu 20.04+, CentOS 8+, or equivalent)
- **Container Runtime**: Docker 20.10+ (for Docker Compose and Kubernetes deployments)
- **Network**: 1 Gbps network connection

### Software Requirements

#### Docker Compose Deployment

- Docker 20.10+
- Docker Compose 1.29+
- Git (to clone the repository)

#### Kubernetes Deployment

- Kubernetes 1.20+
- kubectl 1.20+
- Helm 3.0+
- Git (to clone the repository)

#### Traditional VM Deployment

- Go 1.18+
- Node.js 16+
- PostgreSQL 12+
- Nginx (for frontend)
- Git (to clone the repository)

### Network Requirements

- **Port 80/443**: For HTTP/HTTPS traffic to the frontend
- **Port 8080**: For the backend API (can be configured)
- **Port 5432**: For PostgreSQL database (can be configured)
- **Outbound Internet Access**: For downloading dependencies and updates

## Environment Setup

Before deploying CMDB Lite, you need to set up your environment:

### 1. Clone the Repository

```bash
git clone https://github.com/yourorg/cmdb-lite.git
cd cmdb-lite
```

### 2. Configure Environment Variables

Copy the example environment file and customize it for your environment:

```bash
cp .env.example .env
```

Edit the `.env` file with your configuration:

```env
# Environment
ENVIRONMENT=production

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=secure_password
DB_SSL_MODE=require

# Backend
SERVER_PORT=8080
JWT_SECRET=your-secure-jwt-secret-key

# Frontend
VITE_API_URL=https://api.your-domain.com
VITE_ENVIRONMENT=production
```

### 3. Generate SSL Certificates (for Production)

For production deployments, you need SSL certificates. You can use Let's Encrypt to generate free certificates:

```bash
# Install Certbot
sudo apt update
sudo apt install certbot

# Generate certificates
sudo certbot certonly --standalone -d your-domain.com -d api.your-domain.com

# Copy certificates to the appropriate location
sudo cp /etc/letsencrypt/live/your-domain.com/fullchain.pem ./certs/
sudo cp /etc/letsencrypt/live/your-domain.com/privkey.pem ./certs/
sudo chown -R $USER:$USER ./certs/
```

## Docker Compose Deployment

Docker Compose is the simplest way to deploy CMDB Lite, suitable for development and small production environments.

### Development Environment

To set up a development environment using Docker Compose:

1. Navigate to the project root directory:

```bash
cd /path/to/cmdb-lite
```

2. Copy the development environment file:

```bash
cp docker-compose.dev.yml docker-compose.override.yml
```

3. Start the development environment:

```bash
docker-compose -f docker-compose.dev.yml up -d
```

4. Wait for the services to start (this may take a few minutes):

```bash
docker-compose ps
```

5. Access the application:

- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Database Admin: http://localhost:8081 (username: cmdb_user, password: cmdb_password)

### Production Environment

To set up a production environment using Docker Compose:

1. Navigate to the project root directory:

```bash
cd /path/to/cmdb-lite
```

2. Configure the production environment:

```bash
cp .env.example .env
```

Edit the `.env` file with your production configuration:

```env
# Environment
ENVIRONMENT=production

# Database
DB_HOST=postgres
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=very-secure-password
DB_SSL_MODE=require

# Backend
SERVER_PORT=8080
JWT_SECRET=very-secure-jwt-secret-key

# Frontend
VITE_API_URL=https://api.your-domain.com
VITE_ENVIRONMENT=production
```

3. Start the production environment:

```bash
docker-compose up -d
```

4. Wait for the services to start:

```bash
docker-compose ps
```

5. Access the application:

- Frontend: https://your-domain.com
- Backend API: https://api.your-domain.com

### Docker Compose Configuration

The `docker-compose.yml` file defines the following services:

- **postgres**: PostgreSQL database
- **backend**: Go backend application
- **frontend**: Vue.js frontend application
- **nginx**: Nginx reverse proxy (for production)

#### Environment Variables

You can override default values using environment variables:

- `POSTGRES_DB`: Database name (default: cmdb_lite)
- `POSTGRES_USER`: Database user (default: cmdb_user)
- `POSTGRES_PASSWORD`: Database password (default: cmdb_password)
- `JWT_SECRET`: JWT secret key (required)
- `VITE_API_URL`: Frontend API URL (required for production)

#### Volume Mounts

The following volumes are mounted:

- `postgres_data`: Persistent storage for PostgreSQL data
- `certs`: SSL certificates (for production)
- `logs`: Application logs (for production)

### Managing the Docker Compose Deployment

#### Starting and Stopping Services

```bash
# Start all services
docker-compose up -d

# Stop all services
docker-compose down

# Stop services and remove volumes
docker-compose down -v

# Restart a specific service
docker-compose restart backend
```

#### Viewing Logs

```bash
# View logs for all services
docker-compose logs

# View logs for a specific service
docker-compose logs backend

# Follow logs for a specific service
docker-compose logs -f backend
```

#### Updating the Application

```bash
# Pull the latest images
docker-compose pull

# Rebuild and restart services
docker-compose up -d --force-recreate
```

#### Running Database Migrations

```bash
# Run database migrations
docker-compose exec backend go run cmd/migrate.go
```

## Kubernetes Deployment

Kubernetes is a container orchestration platform that provides scalability, high availability, and automated management for containerized applications.

### Prerequisites

Before deploying CMDB Lite to Kubernetes, ensure you have the following:

- A Kubernetes cluster (version 1.20+)
- kubectl configured to access your cluster
- Helm 3.0+ installed
- A container registry to store Docker images (e.g., Docker Hub, AWS ECR, Google GCR)

### Deployment Steps

#### 1. Build and Push Docker Images

First, build and push the Docker images to your container registry:

```bash
# Set your container registry
export REGISTRY=your-registry.example.com

# Build and push the backend image
docker build -t $REGISTRY/cmdb-lite-backend:latest ./backend
docker push $REGISTRY/cmdb-lite-backend:latest

# Build and push the frontend image
docker build -t $REGISTRY/cmdb-lite-frontend:latest ./frontend
docker push $REGISTRY/cmdb-lite-frontend:latest
```

#### 2. Install the Helm Chart

CMDB Lite includes a Helm chart for easy deployment to Kubernetes:

```bash
# Navigate to the Helm chart directory
cd k8s/helm

# Install the chart
helm install cmdb-lite . --namespace cmdb-lite --create-namespace
```

#### 3. Verify the Deployment

Check that all pods are running:

```bash
kubectl get pods -n cmdb-lite
```

You should see output similar to this:

```
NAME                                  READY   STATUS    RESTARTS   AGE
cmdb-lite-backend-5f7d8c9f9f-abcde     1/1     Running   0          2m
cmdb-lite-frontend-6f8d9c9f9f-abcde    1/1     Running   0          2m
postgres-0                             1/1     Running   0          5m
```

#### 4. Access the Application

To access the application, you need to set up an Ingress or LoadBalancer service:

```bash
# Get the Ingress IP (if using Ingress)
kubectl get ingress -n cmdb-lite

# Get the LoadBalancer IP (if using LoadBalancer)
kubectl get svc -n cmdb-lite
```

### Configuration

The Helm chart allows you to configure various aspects of the deployment using a `values.yaml` file. Here's an example configuration:

```yaml
# Global configuration
global:
  imageRegistry: your-registry.example.com
  imagePullSecrets:
    - name: registry-secret

# Backend configuration
backend:
  image:
    repository: cmdb-lite-backend
    tag: latest
    pullPolicy: Always
  replicaCount: 3
  resources:
    limits:
      cpu: 500m
      memory: 512Mi
    requests:
      cpu: 200m
      memory: 256Mi
  env:
    - name: ENVIRONMENT
      value: "production"
    - name: DB_HOST
      value: "postgres"
    - name: DB_PORT
      value: "5432"
    - name: DB_NAME
      value: "cmdb_lite"
    - name: DB_USER
      valueFrom:
        secretKeyRef:
          name: postgres-secret
          key: username
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: postgres-secret
          key: password
    - name: JWT_SECRET
      valueFrom:
        secretKeyRef:
          name: jwt-secret
          key: jwt-secret

# Frontend configuration
frontend:
  image:
    repository: cmdb-lite-frontend
    tag: latest
    pullPolicy: Always
  replicaCount: 3
  resources:
    limits:
      cpu: 200m
      memory: 256Mi
    requests:
      cpu: 100m
      memory: 128Mi
  env:
    - name: VITE_API_URL
      value: "https://api.your-domain.com"
    - name: VITE_ENVIRONMENT
      value: "production"

# PostgreSQL configuration
postgres:
  image:
    repository: postgres
    tag: 13
    pullPolicy: IfNotPresent
  persistence:
    enabled: true
    size: 20Gi
    storageClass: fast-ssd
  resources:
    limits:
      cpu: 1000m
      memory: 1Gi
    requests:
      cpu: 500m
      memory: 512Mi
  env:
    - name: POSTGRES_DB
      value: "cmdb_lite"
    - name: POSTGRES_USER
      valueFrom:
        secretKeyRef:
          name: postgres-secret
          key: username
    - name: POSTGRES_PASSWORD
      valueFrom:
        secretKeyRef:
          name: postgres-secret
          key: password

# Ingress configuration
ingress:
  enabled: true
  annotations:
    kubernetes.io/ingress.class: nginx
    cert-manager.io/cluster-issuer: letsencrypt-prod
  hosts:
    - host: your-domain.com
      paths:
        - path: /
          pathType: Prefix
          service: frontend
          port: 80
    - host: api.your-domain.com
      paths:
        - path: /
          pathType: Prefix
          service: backend
          port: 8080
  tls:
    - secretName: cmdb-lite-tls
      hosts:
        - your-domain.com
        - api.your-domain.com
```

### Scaling and High Availability

Kubernetes makes it easy to scale CMDB Lite for high availability:

#### Horizontal Pod Autoscaling

Configure Horizontal Pod Autoscalers (HPAs) to automatically scale the number of pods based on CPU utilization:

```yaml
# Horizontal Pod Autoscaler for the backend
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: cmdb-lite-backend-hpa
  namespace: cmdb-lite
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: cmdb-lite-backend
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
```

#### Pod Disruption Budgets

Configure Pod Disruption Budgets (PDBs) to ensure that a minimum number of pods are always available during voluntary disruptions:

```yaml
# Pod Disruption Budget for the backend
apiVersion: policy/v1
kind: PodDisruptionBudget
metadata:
  name: cmdb-lite-backend-pdb
  namespace: cmdb-lite
spec:
  minAvailable: 2
  selector:
    matchLabels:
      app.kubernetes.io/name: backend
      app.kubernetes.io/instance: cmdb-lite
```

#### Database High Availability

For production deployments, consider using a managed PostgreSQL service with built-in high availability, such as:

- Amazon RDS for PostgreSQL
- Google Cloud SQL for PostgreSQL
- Azure Database for PostgreSQL

Alternatively, you can set up PostgreSQL high availability using tools like Patroni and etcd.

### Managing the Kubernetes Deployment

#### Upgrading the Deployment

To upgrade the deployment to a new version:

```bash
# Update the Helm chart
helm repo update

# Upgrade the deployment
helm upgrade cmdb-lite . --namespace cmdb-lite
```

#### Rolling Back

If an upgrade causes issues, you can roll back to a previous version:

```bash
# List revision history
helm history cmdb-lite -n cmdb-lite

# Roll back to a specific revision
helm rollback cmdb-lite <revision> -n cmdb-lite
```

#### Viewing Logs

```bash
# View logs for all pods in a deployment
kubectl logs -l app.kubernetes.io/name=backend -n cmdb-lite

# View logs for a specific pod
kubectl logs <pod-name> -n cmdb-lite

# Follow logs for a specific pod
kubectl logs -f <pod-name> -n cmdb-lite
```

#### Accessing the Shell

```bash
# Access the shell for a pod
kubectl exec -it <pod-name> -n cmdb-lite -- /bin/bash
```

## Traditional VM Deployment

For organizations that prefer to deploy directly to virtual machines without using containers, CMDB Lite can be deployed as a traditional application.

### Prerequisites

Before deploying CMDB Lite to a virtual machine, ensure you have the following:

- A virtual machine with Ubuntu 20.04+ or CentOS 8+
- Go 1.18+ installed
- Node.js 16+ installed
- PostgreSQL 12+ installed
- Nginx installed
- Git installed

### Backend Deployment

#### 1. Install Dependencies

```bash
# Update package lists
sudo apt update

# Install Go
sudo apt install -y golang-go

# Verify Go installation
go version
```

#### 2. Build the Backend

```bash
# Clone the repository
git clone https://github.com/yourorg/cmdb-lite.git
cd cmdb-lite/backend

# Build the application
go build -o cmdb-lite-backend cmd/main.go

# Create a systemd service file
sudo tee /etc/systemd/system/cmdb-lite-backend.service > /dev/null <<EOF
[Unit]
Description=CMDB Lite Backend
After=network.target

[Service]
Type=simple
User=cmdb
WorkingDirectory=/opt/cmdb-lite/backend
ExecStart=/opt/cmdb-lite/backend/cmdb-lite-backend
Restart=on-failure
RestartSec=5
Environment=ENVIRONMENT=production
Environment=DB_HOST=localhost
Environment=DB_PORT=5432
Environment=DB_NAME=cmdb_lite
Environment=DB_USER=cmdb_user
Environment=DB_PASSWORD=secure_password
Environment=SERVER_PORT=8080
Environment=JWT_SECRET=your-secure-jwt-secret-key

[Install]
WantedBy=multi-user.target
EOF
```

#### 3. Install and Configure the Backend

```bash
# Create a user for the application
sudo useradd -r -s /bin/false cmdb

# Create directories
sudo mkdir -p /opt/cmdb-lite/backend
sudo mkdir -p /var/log/cmdb-lite

# Copy the binary
sudo cp cmdb-lite-backend /opt/cmdb-lite/backend/
sudo chown -R cmdb:cmdb /opt/cmdb-lite/backend

# Set permissions
sudo chmod +x /opt/cmdb-lite/backend/cmdb-lite-backend

# Reload systemd
sudo systemctl daemon-reload

# Enable and start the service
sudo systemctl enable cmdb-lite-backend
sudo systemctl start cmdb-lite-backend

# Check the service status
sudo systemctl status cmdb-lite-backend
```

### Frontend Deployment

#### 1. Install Dependencies

```bash
# Update package lists
sudo apt update

# Install Node.js and npm
sudo apt install -y nodejs npm

# Verify Node.js installation
node --version
npm --version
```

#### 2. Build the Frontend

```bash
# Navigate to the frontend directory
cd ../frontend

# Install dependencies
npm install

# Build the application
npm run build

# Copy the build output to the web root
sudo mkdir -p /var/www/cmdb-lite
sudo cp -r dist/* /var/www/cmdb-lite/
sudo chown -R www-data:www-data /var/www/cmdb-lite
```

#### 3. Configure Nginx

```bash
# Create an Nginx configuration file
sudo tee /etc/nginx/sites-available/cmdb-lite > /dev/null <<EOF
server {
    listen 80;
    server_name your-domain.com;

    # Redirect to HTTPS
    return 301 https://\$host\$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL configuration
    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;

    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Frontend
    location / {
        root /var/www/cmdb-lite;
        index index.html;
        try_files \$uri \$uri/ /index.html;
    }

    # Backend API
    location /api {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    # Logging
    access_log /var/log/nginx/cmdb-lite.access.log;
    error_log /var/log/nginx/cmdb-lite.error.log;
}
EOF

# Enable the site
sudo ln -s /etc/nginx/sites-available/cmdb-lite /etc/nginx/sites-enabled/

# Test Nginx configuration
sudo nginx -t

# Reload Nginx
sudo systemctl reload nginx
```

### Database Setup

#### 1. Install PostgreSQL

```bash
# Install PostgreSQL
sudo apt install -y postgresql postgresql-contrib

# Initialize the database
sudo pg_ctlcluster 13 main start

# Enable PostgreSQL to start on boot
sudo systemctl enable postgresql
```

#### 2. Create the Database and User

```bash
# Switch to the postgres user
sudo -u postgres psql

# Create the database and user
CREATE DATABASE cmdb_lite;
CREATE USER cmdb_user WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE cmdb_lite TO cmdb_user;
\q
```

#### 3. Run Database Migrations

```bash
# Navigate to the backend directory
cd /opt/cmdb-lite/backend

# Run the migrations
sudo -u cmdb ./cmdb-lite-backend migrate
```

### Managing the Traditional VM Deployment

#### Starting and Stopping Services

```bash
# Start the backend service
sudo systemctl start cmdb-lite-backend

# Stop the backend service
sudo systemctl stop cmdb-lite-backend

# Restart the backend service
sudo systemctl restart cmdb-lite-backend

# Check the service status
sudo systemctl status cmdb-lite-backend

# Reload Nginx
sudo systemctl reload nginx
```

#### Viewing Logs

```bash
# View backend service logs
sudo journalctl -u cmdb-lite-backend -f

# View Nginx logs
sudo tail -f /var/log/nginx/cmdb-lite.access.log
sudo tail -f /var/log/nginx/cmdb-lite.error.log
```

#### Updating the Application

```bash
# Update the backend
cd /opt/cmdb-lite/backend
sudo systemctl stop cmdb-lite-backend
sudo -u cmdb git pull
sudo -u cmdb go build -o cmdb-lite-backend cmd/main.go
sudo systemctl start cmdb-lite-backend

# Update the frontend
cd /var/www/cmdb-lite
sudo rm -rf *
sudo -u www-data cp -r /path/to/cmdb-lite/frontend/dist/* .
sudo systemctl reload nginx
```

## Post-Deployment Steps

After deploying CMDB Lite, you need to perform some post-deployment steps:

### 1. Create an Admin User

```bash
# For Docker Compose
docker-compose exec backend go run cmd/create-admin.go

# For Kubernetes
kubectl exec -it <backend-pod-name> -n cmdb-lite -- go run cmd/create-admin.go

# For Traditional VM
sudo -u cmdb /opt/cmdb-lite/backend/cmdb-lite-backend create-admin
```

### 2. Verify the Application

Open a web browser and navigate to your application URL:

- Frontend: https://your-domain.com
- Backend API: https://api.your-domain.com

Verify that you can log in with the admin user you created.

### 3. Set Up Monitoring

Configure monitoring and alerting for your deployment. See the [Monitoring and Logging](#monitoring-and-logging) section for details.

### 4. Set Up Backups

Configure regular backups of your database. See the [Backup and Recovery](#backup-and-recovery) section for details.

## Configuration Management

CMDB Lite uses environment variables for configuration. Here are the key configuration options:

### Backend Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `ENVIRONMENT` | Environment (development, staging, production) | development | No |
| `DB_HOST` | Database host | localhost | No |
| `DB_PORT` | Database port | 5432 | No |
| `DB_NAME` | Database name | cmdb_lite | No |
| `DB_USER` | Database user | cmdb_user | No |
| `DB_PASSWORD` | Database password | - | Yes |
| `DB_SSL_MODE` | Database SSL mode | disable | No |
| `SERVER_PORT` | Server port | 8080 | No |
| `SERVER_HOST` | Server host | 0.0.0.0 | No |
| `JWT_SECRET` | JWT secret key | - | Yes |
| `JWT_EXPIRATION` | JWT expiration time | 24h | No |
| `LOG_LEVEL` | Log level (debug, info, warn, error) | info | No |
| `LOG_FORMAT` | Log format (text, json) | text | No |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limit requests per minute | 100 | No |
| `RATE_LIMIT_BURST` | Rate limit burst | 200 | No |

### Frontend Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `VITE_API_URL` | Backend API URL | http://localhost:8080 | Yes |
| `VITE_ENVIRONMENT` | Environment (development, staging, production) | development | No |
| `VITE_ENABLE_GRAPH_FEATURES` | Enable graph features | true | No |
| `VITE_ENABLE_AUDIT_LOGS` | Enable audit logs | true | No |
| `VITE_ENABLE_USER_MANAGEMENT` | Enable user management | true | No |

### Database Configuration

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `POSTGRES_DB` | Database name | cmdb_lite | No |
| `POSTGRES_USER` | Database user | cmdb_user | No |
| `POSTGRES_PASSWORD` | Database password | - | Yes |
| `POSTGRES_INITDB_ARGS` | PostgreSQL initdb arguments | - | No |
| `POSTGRES_HOST_AUTH_METHOD` | PostgreSQL host authentication method | md5 | No |

## Monitoring and Logging

Monitoring and logging are essential for maintaining the health and performance of CMDB Lite.

### Monitoring

#### Application Metrics

CMDB Lite includes built-in support for Prometheus metrics:

- **Request Rate**: Number of requests per second
- **Response Time**: Average and percentiles of response times
- **Error Rate**: Percentage of requests that result in errors
- **Memory Usage**: Application memory consumption
- **CPU Usage**: Application CPU consumption

To enable metrics collection, configure Prometheus to scrape the `/metrics` endpoint of the backend.

#### Database Metrics

Monitor the following database metrics:

- **Connection Count**: Number of active database connections
- **Query Performance**: Slow query count and average query time
- **Database Size**: Size of the database and growth rate
- **Replication Lag**: For replicated databases, the lag between primary and replicas

#### System Metrics

Monitor the following system metrics:

- **CPU Usage**: System CPU usage
- **Memory Usage**: System memory usage
- **Disk Usage**: Available disk space and I/O performance
- **Network Traffic**: Incoming and outgoing network traffic

### Logging

#### Backend Logging

The backend uses structured logging for better observability:

```json
{
  "level": "info",
  "timestamp": "2023-01-01T12:00:00Z",
  "request_id": "abc123",
  "user_id": "user123",
  "method": "GET",
  "path": "/api/v1/cis",
  "status_code": 200,
  "duration_ms": 45,
  "message": "Request completed successfully"
}
```

#### Frontend Logging

The frontend logs errors and warnings to the browser console, and can be configured to send logs to a logging service.

#### Log Aggregation

For production deployments, consider using a log aggregation system such as:

- ELK Stack (Elasticsearch, Logstash, Kibana)
- Splunk
- Datadog
- Grafana Loki

### Alerting

Set up alerts for the following conditions:

1. **High Error Rate**: Alert if error rate exceeds 5% for 5 minutes
2. **High Response Time**: Alert if 95th percentile response time exceeds 1 second for 5 minutes
3. **High Memory Usage**: Alert if memory usage exceeds 90% for 5 minutes
4. **High CPU Usage**: Alert if CPU usage exceeds 90% for 5 minutes
5. **Database Connection Issues**: Alert if database connection errors exceed 1% for 5 minutes
6. **Disk Space Low**: Alert if disk space is below 20%
7. **Service Unavailable**: Alert if the service is unavailable for more than 1 minute

## Backup and Recovery

Regular backups are essential for protecting against data loss and ensuring business continuity.

### Database Backups

#### Full Backups

Perform daily full backups of the PostgreSQL database:

```bash
# Create a backup directory
mkdir -p /backups/cmdb-lite

# Perform a full backup
pg_dump -h localhost -U cmdb_user -d cmdb_lite -f /backups/cmdb-lite/cmdb-lite-$(date +%Y%m%d).sql

# Compress the backup
gzip /backups/cmdb-lite/cmdb-lite-$(date +%Y%m%d).sql
```

#### Point-in-Time Recovery

Configure PostgreSQL for point-in-time recovery:

1. Enable WAL archiving in `postgresql.conf`:

```ini
wal_level = replica
archive_mode = on
archive_command = 'cp %p /backups/cmdb-lite/wal/%f'
```

2. Create a recovery configuration file:

```ini
restore_command = 'cp /backups/cmdb-lite/wal/%f %p'
recovery_target_time = '2023-01-01 12:00:00'
```

### Configuration Backups

Back up the following configuration files:

- Environment variable configurations
- Nginx configuration files
- Systemd service files
- SSL certificates

### Backup Testing

Regularly test your backups to ensure they can be restored successfully:

1. **Schedule Regular Tests**: Test backups at least quarterly
2. **Document Procedures**: Document backup and recovery procedures
3. **Test Different Scenarios**: Test recovery from different types of failures
4. **Update Procedures**: Update procedures based on test results

### Recovery Procedures

#### Database Recovery

1. Stop the application:

```bash
# For Docker Compose
docker-compose stop backend

# For Kubernetes
kubectl scale deployment cmdb-lite-backend --replicas=0 -n cmdb-lite

# For Traditional VM
sudo systemctl stop cmdb-lite-backend
```

2. Restore the database:

```bash
# Drop the existing database
dropdb -h localhost -U cmdb_user cmdb_lite

# Create a new database
createdb -h localhost -U cmdb_user cmdb_lite

# Restore the backup
gunzip -c /backups/cmdb-lite/cmdb-lite-20230101.sql.gz | psql -h localhost -U cmdb_user -d cmdb_lite
```

3. Start the application:

```bash
# For Docker Compose
docker-compose start backend

# For Kubernetes
kubectl scale deployment cmdb-lite-backend --replicas=3 -n cmdb-lite

# For Traditional VM
sudo systemctl start cmdb-lite-backend
```

## Troubleshooting

This section provides guidance on troubleshooting common issues with CMDB Lite.

### Common Issues and Solutions

#### Application Fails to Start

**Symptoms**: The application fails to start with error messages

**Possible Causes**:
- Missing or invalid configuration
- Database connection issues
- Port already in use
- File permission issues

**Troubleshooting Steps**:
1. Check the application logs for error messages
2. Verify that all required environment variables are set correctly
3. Check that the database is accessible and the schema is up to date
4. Verify that the specified port is available
5. Check file permissions for the application directory

#### High Response Times

**Symptoms**: The application responds slowly to requests

**Possible Causes**:
- Database performance issues
- High resource usage (CPU, memory)
- Network latency
- Inefficient queries

**Troubleshooting Steps**:
1. Check resource usage (CPU, memory, disk I/O)
2. Analyze database performance (slow queries, connection pool usage)
3. Check network latency between application and database
4. Review application logs for errors or warnings
5. Profile the application to identify performance bottlenecks

#### Database Connection Errors

**Symptoms**: The application reports database connection errors

**Possible Causes**:
- Database is down or not accessible
- Database connection pool is exhausted
- Network issues between application and database
- Database authentication issues

**Troubleshooting Steps**:
1. Check that the database is running and accessible
2. Verify database connection parameters (host, port, username, password)
3. Check network connectivity between application and database
4. Review database connection pool settings
5. Check database logs for connection-related errors

### Diagnostic Tools

#### Application Logs

Application logs are the primary source of information for troubleshooting:

```bash
# For Docker Compose
docker-compose logs backend

# For Kubernetes
kubectl logs <backend-pod-name> -n cmdb-lite

# For Traditional VM
sudo journalctl -u cmdb-lite-backend -f
```

#### Database Metrics

Use PostgreSQL's built-in monitoring views to diagnose database issues:

```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check long-running queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query, state
FROM pg_stat_activity
WHERE (now() - pg_stat_activity.query_start) > interval '5 minutes';

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) as size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

#### System Metrics

Use system tools to monitor resource usage:

```bash
# Check CPU usage
top
htop

# Check memory usage
free -h
vmstat

# Check disk usage
df -h
du -sh

# Check network connections
netstat -tuln
ss -tuln
```

For more information on operating CMDB Lite, see the [Operations Guide](README.md).