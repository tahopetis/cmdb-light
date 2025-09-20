# CMDB Lite Backup and Recovery Guide

This document provides comprehensive guidance on backing up and recovering CMDB Lite data and configurations in production environments.

## Table of Contents

- [Backup Strategy Overview](#backup-strategy-overview)
- [Database Backups](#database-backups)
  - [Full Backups](#full-backups)
  - [Incremental Backups](#incremental-backups)
  - [Point-in-Time Recovery](#point-in-time-recovery)
  - [Automated Backup Scripts](#automated-backup-scripts)
- [Configuration Backups](#configuration-backups)
  - [Application Configuration](#application-configuration)
  - [Infrastructure Configuration](#infrastructure-configuration)
  - [Secrets Management](#secrets-management)
- [File System Backups](#file-system-backups)
  - [Application Files](#application-files)
  - [Static Assets](#static-assets)
  - [Log Files](#log-files)
- [Backup Storage](#backup-storage)
  - [Local Storage](#local-storage)
  - [Cloud Storage](#cloud-storage)
  - [Offsite Storage](#offsite-storage)
- [Backup Verification](#backup-verification)
  - [Automated Verification](#automated-verification)
  - [Manual Verification](#manual-verification)
  - [Recovery Testing](#recovery-testing)
- [Recovery Procedures](#recovery-procedures)
  - [Database Recovery](#database-recovery)
  - [Configuration Recovery](#configuration-recovery)
  - [Complete System Recovery](#complete-system-recovery)
- [Disaster Recovery](#disaster-recovery)
  - [Disaster Recovery Plan](#disaster-recovery-plan)
  - [Recovery Time Objectives](#recovery-time-objectives)
  - [Recovery Point Objectives](#recovery-point-objectives)
- [Backup Security](#backup-security)
  - [Encryption](#encryption)
  - [Access Control](#access-control)
  - [Audit Logging](#audit-logging)
- [Backup Maintenance](#backup-maintenance)
  - [Backup Rotation](#backup-rotation)
  - [Backup Cleanup](#backup-cleanup)
  - [Backup Monitoring](#backup-monitoring)

## Backup Strategy Overview

A comprehensive backup strategy is essential for protecting CMDB Lite data and ensuring business continuity. The strategy should include:

1. **Regular Backups**: Schedule regular backups to minimize data loss
2. **Multiple Backup Types**: Use different backup types for different recovery scenarios
3. **Offsite Storage**: Store backups in multiple locations to protect against local disasters
4. **Backup Verification**: Regularly verify that backups can be restored successfully
5. **Documented Procedures**: Document backup and recovery procedures for reference during emergencies

### Backup Types

| Backup Type | Description | Frequency | Retention | Recovery Time |
|-------------|-------------|-----------|-----------|---------------|
| Full Database | Complete database backup | Daily | 30 days | Moderate |
| Incremental | Changes since last backup | Hourly | 7 days | Fast |
| Point-in-Time | Database state at specific time | Continuous | 14 days | Fast |
| Configuration | Application and infrastructure configuration | Weekly | 90 days | Fast |
| File System | Application files, static assets, logs | Daily | 30 days | Moderate |

### Backup Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   CMDB Lite     │    │   Backup Agent  │    │   Backup Server │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ Database  │  │───▶│  │ Collector │  │───▶│  │ Storage   │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ Config    │  │───▶│  │ Scheduler │  │───▶│  │ Scheduler │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ Files     │  │───▶│  │ Encryptor │  │───▶│  │ Encryptor │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │                      │                      │
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Local Storage│    │   Cloud Storage │    │  Offsite Storage│
│                 │    │                 │    │                 │
│  ┌───────────┐  │    │  ┌───────────┐  │    │  ┌───────────┐  │
│  │ Disk/NAS  │  │    │  │ S3/GCS    │  │    │  │ Tape/DR   │  │
│  └───────────┘  │    │  └───────────┘  │    │  └───────────┘  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Database Backups

Database backups are the most critical component of the backup strategy, as they contain all the configuration item data, relationships, and user information.

### Full Backups

Full backups capture the entire database at a specific point in time. They are the foundation of the backup strategy and are typically performed daily.

#### Using pg_dump

```bash
#!/bin/bash

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
BACKUP_DIR="/backups/cmdb-lite/database"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cmdb-lite-full_$DATE.sql"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Perform backup
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Log backup
echo "$(date): Full database backup completed: $BACKUP_FILE.gz" >> $BACKUP_DIR/backup.log

# Verify backup
if [ -f "$BACKUP_FILE.gz" ]; then
    echo "$(date): Backup verification successful" >> $BACKUP_DIR/backup.log
else
    echo "$(date): Backup verification failed" >> $BACKUP_DIR/backup.log
    # Send alert
    send_alert "Database backup verification failed"
fi
```

#### Using pg_dumpall

```bash
#!/bin/bash

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
BACKUP_DIR="/backups/cmdb-lite/database"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cmdb-lite-all_$DATE.sql"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Perform backup of all databases
PGPASSWORD=$DB_PASSWORD pg_dumpall -h $DB_HOST -p $DB_PORT -U $DB_USER -f $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Log backup
echo "$(date): Full database cluster backup completed: $BACKUP_FILE.gz" >> $BACKUP_DIR/backup.log
```

#### Using Custom Format

```bash
#!/bin/bash

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
BACKUP_DIR="/backups/cmdb-lite/database"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cmdb-lite-custom_$DATE.dump"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Perform backup in custom format
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -F c -f $BACKUP_FILE

# Log backup
echo "$(date): Custom format database backup completed: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
```

### Incremental Backups

Incremental backups capture only the changes since the last backup, reducing storage requirements and backup time. In PostgreSQL, this is typically achieved using Write-Ahead Logging (WAL).

#### Configuring WAL Archiving

1. Configure PostgreSQL for WAL archiving:

```ini
# /etc/postgresql/13/main/postgresql.conf
wal_level = replica
archive_mode = on
archive_command = 'test ! -f /backups/cmdb-lite/wal/%f && cp %p /backups/cmdb-lite/wal/%f'
archive_timeout = 60
max_wal_senders = 3
```

2. Create the WAL archive directory:

```bash
mkdir -p /backups/cmdb-lite/wal
chown -R postgres:postgres /backups/cmdb-lite/wal
chmod 700 /backups/cmdb-lite/wal
```

3. Restart PostgreSQL:

```bash
sudo systemctl restart postgresql
```

#### WAL Backup Script

```bash
#!/bin/bash

# Configuration
WAL_DIR="/backups/cmdb-lite/wal"
BACKUP_DIR="/backups/cmdb-lite/database"
DATE=$(date +%Y%m%d_%H%M%S)
WAL_BACKUP_FILE="$BACKUP_DIR/cmdb-lite-wal_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Create WAL backup
tar -czf $WAL_BACKUP_FILE -C $(dirname $WAL_DIR) $(basename $WAL_DIR)

# Log backup
echo "$(date): WAL backup completed: $WAL_BACKUP_FILE" >> $BACKUP_DIR/backup.log

# Clean up old WAL files (keep last 7 days)
find $WAL_DIR -name "*.wal" -type f -mtime +7 -delete

# Log cleanup
echo "$(date): Cleaned up old WAL files" >> $BACKUP_DIR/backup.log
```

### Point-in-Time Recovery

Point-in-Time Recovery (PITR) allows you to restore the database to any specific point in time, which is essential for recovering from accidental data corruption or deletion.

#### Configuring PITR

1. Ensure WAL archiving is configured as described in the previous section.

2. Create a recovery configuration file:

```bash
# /backups/cmdb-lite/recovery.conf
restore_command = 'cp /backups/cmdb-lite/wal/%f %p'
recovery_target_time = '2023-01-01 12:00:00'
recovery_target_action = 'promote'
```

#### PITR Backup Script

```bash
#!/bin/bash

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
BACKUP_DIR="/backups/cmdb-lite/database"
WAL_DIR="/backups/cmdb-lite/wal"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cmdb-lite-pitr_$DATE.sql"
WAL_BACKUP_FILE="$BACKUP_DIR/cmdb-lite-wal_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Perform base backup
PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $BACKUP_FILE

# Compress backup
gzip $BACKUP_FILE

# Create WAL backup
tar -czf $WAL_BACKUP_FILE -C $(dirname $WAL_DIR) $(basename $WAL_DIR)

# Log backup
echo "$(date): PITR backup completed: $BACKUP_FILE.gz, $WAL_BACKUP_FILE" >> $BACKUP_DIR/backup.log

# Create recovery configuration
cat > $BACKUP_DIR/recovery_$DATE.conf << EOF
restore_command = 'cp $WAL_DIR/%f %p'
recovery_target_time = '$(date -d "1 hour ago" +"%Y-%m-%d %H:%M:%S")'
recovery_target_action = 'promote'
EOF

# Log recovery configuration
echo "$(date): Recovery configuration created: $BACKUP_DIR/recovery_$DATE.conf" >> $BACKUP_DIR/backup.log
```

### Automated Backup Scripts

Automating backups ensures they are performed consistently and reduces the risk of human error.

#### Full Backup Automation

```bash
#!/bin/bash

# /usr/local/bin/cmdb-lite-backup-full.sh

# Configuration
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
BACKUP_DIR="/backups/cmdb-lite/database"
RETENTION_DAYS=30
LOG_FILE="$BACKUP_DIR/backup.log"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Perform backup
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/cmdb-lite-full_$DATE.sql"

log "Starting full database backup"

# Perform backup
if PGPASSWORD=$DB_PASSWORD pg_dump -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $BACKUP_FILE; then
    log "Full database backup completed: $BACKUP_FILE"
    
    # Compress backup
    if gzip $BACKUP_FILE; then
        log "Backup compressed: $BACKUP_FILE.gz"
    else
        log "Backup compression failed"
        send_alert "CMDB Lite backup compression failed"
        exit 1
    fi
    
    # Verify backup
    if [ -f "$BACKUP_FILE.gz" ]; then
        log "Backup verification successful"
    else
        log "Backup verification failed"
        send_alert "CMDB Lite backup verification failed"
        exit 1
    fi
    
    # Clean up old backups
    find $BACKUP_DIR -name "cmdb-lite-full_*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
    log "Cleaned up backups older than $RETENTION_DAYS days"
    
    # Upload to cloud storage (optional)
    if command -v aws &> /dev/null; then
        aws s3 cp $BACKUP_FILE.gz s3://your-backup-bucket/cmdb-lite/database/
        log "Backup uploaded to S3"
    fi
    
    log "Full database backup process completed successfully"
else
    log "Full database backup failed"
    send_alert "CMDB Lite full database backup failed"
    exit 1
fi
```

#### Incremental Backup Automation

```bash
#!/bin/bash

# /usr/local/bin/cmdb-lite-backup-incremental.sh

# Configuration
WAL_DIR="/backups/cmdb-lite/wal"
BACKUP_DIR="/backups/cmdb-lite/database"
RETENTION_DAYS=7
LOG_FILE="$BACKUP_DIR/backup.log"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Perform backup
DATE=$(date +%Y%m%d_%H%M%S)
WAL_BACKUP_FILE="$BACKUP_DIR/cmdb-lite-wal_$DATE.tar.gz"

log "Starting incremental database backup"

# Create WAL backup
if tar -czf $WAL_BACKUP_FILE -C $(dirname $WAL_DIR) $(basename $WAL_DIR); then
    log "Incremental database backup completed: $WAL_BACKUP_FILE"
    
    # Verify backup
    if [ -f "$WAL_BACKUP_FILE" ]; then
        log "Backup verification successful"
    else
        log "Backup verification failed"
        send_alert "CMDB Lite incremental backup verification failed"
        exit 1
    fi
    
    # Clean up old WAL files
    find $WAL_DIR -name "*.wal" -type f -mtime +$RETENTION_DAYS -delete
    log "Cleaned up WAL files older than $RETENTION_DAYS days"
    
    # Upload to cloud storage (optional)
    if command -v aws &> /dev/null; then
        aws s3 cp $WAL_BACKUP_FILE s3://your-backup-bucket/cmdb-lite/database/
        log "Backup uploaded to S3"
    fi
    
    log "Incremental database backup process completed successfully"
else
    log "Incremental database backup failed"
    send_alert "CMDB Lite incremental database backup failed"
    exit 1
fi
```

#### Scheduling Backups with Cron

```bash
# Edit crontab
crontab -e

# Add the following lines to schedule backups
# Full backup every day at 2 AM
0 2 * * * /usr/local/bin/cmdb-lite-backup-full.sh

# Incremental backup every hour at 30 minutes past the hour
30 * * * * /usr/local/bin/cmdb-lite-backup-incremental.sh
```

## Configuration Backups

Configuration backups ensure that you can restore the application and infrastructure settings in case of a failure.

### Application Configuration

#### Environment Variables

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
DATE=$(date +%Y%m%d_%H%M%S)
ENV_FILE="/opt/cmdb-lite/.env"
BACKUP_FILE="$BACKUP_DIR/env_$DATE.txt"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup environment variables
if [ -f "$ENV_FILE" ]; then
    cp $ENV_FILE $BACKUP_FILE
    echo "$(date): Environment variables backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Environment variables backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Environment variables file not found: $ENV_FILE" >> $BACKUP_DIR/backup.log
fi
```

#### Application Configuration Files

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
DATE=$(date +%Y%m%d_%H%M%S)
CONFIG_DIR="/etc/cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/config_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup configuration files
if [ -d "$CONFIG_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $CONFIG_DIR) $(basename $CONFIG_DIR)
    echo "$(date): Configuration files backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Configuration files backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Configuration directory not found: $CONFIG_DIR" >> $BACKUP_DIR/backup.log
fi
```

### Infrastructure Configuration

#### Docker Compose Configuration

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
DATE=$(date +%Y%m%d_%H%M%S)
COMPOSE_DIR="/opt/cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/docker-compose_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup Docker Compose files
if [ -d "$COMPOSE_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $COMPOSE_DIR) $(basename $COMPOSE_DIR)/docker-compose*.yml
    echo "$(date): Docker Compose files backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Docker Compose files backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Docker Compose directory not found: $COMPOSE_DIR" >> $BACKUP_DIR/backup.log
fi
```

#### Kubernetes Configuration

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
DATE=$(date +%Y%m%d_%H%M%S)
NAMESPACE="cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/kubernetes_$DATE.yaml"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup Kubernetes resources
kubectl get all,configmap,secret,ingress -n $NAMESPACE -o yaml > $BACKUP_FILE
echo "$(date): Kubernetes resources backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log

# Encrypt backup
if command -v gpg &> /dev/null; then
    gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
    rm $BACKUP_FILE
    echo "$(date): Kubernetes resources backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
fi
```

#### Nginx Configuration

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
DATE=$(date +%Y%m%d_%H%M%S)
NGINX_DIR="/etc/nginx"
BACKUP_FILE="$BACKUP_DIR/nginx_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup Nginx configuration
if [ -d "$NGINX_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $NGINX_DIR) $(basename $NGINX_DIR)/sites-available $(basename $NGINX_DIR)/sites-enabled
    echo "$(date): Nginx configuration backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Nginx configuration backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Nginx directory not found: $NGINX_DIR" >> $BACKUP_DIR/backup.log
fi
```

### Secrets Management

#### HashiCorp Vault

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/secrets"
DATE=$(date +%Y%m%d_%H%M%S)
VAULT_ADDR="https://vault.example.com"
VAULT_TOKEN="your-vault-token"
BACKUP_FILE="$BACKUP_DIR/vault_$DATE.json"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup Vault secrets
export VAULT_ADDR=$VAULT_ADDR
export VAULT_TOKEN=$VAULT_TOKEN

# List all secrets
vault kv list -format=json cmdb-lite > $BACKUP_DIR/secrets_list_$DATE.json

# Backup each secret
for path in $(vault kv list cmdb-lite | grep -E '^  [^/]+$' | sed 's/^  //'); do
    vault kv get -format=json cmdb-lite/$path > $BACKUP_DIR/secret_${path}_$DATE.json
done

# Create archive
tar -czf $BACKUP_FILE -C $BACKUP_DIR secrets_list_$DATE.json secret_*_$DATE.json
echo "$(date): Vault secrets backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log

# Clean up individual files
rm $BACKUP_DIR/secrets_list_$DATE.json $BACKUP_DIR/secret_*_$DATE.json

# Encrypt backup
if command -v gpg &> /dev/null; then
    gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
    rm $BACKUP_FILE
    echo "$(date): Vault secrets backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
fi
```

#### Kubernetes Secrets

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/secrets"
DATE=$(date +%Y%m%d_%H%M%S)
NAMESPACE="cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/kubernetes-secrets_$DATE.yaml"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup Kubernetes secrets
kubectl get secrets -n $NAMESPACE -o yaml > $BACKUP_FILE
echo "$(date): Kubernetes secrets backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log

# Encrypt backup
if command -v gpg &> /dev/null; then
    gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
    rm $BACKUP_FILE
    echo "$(date): Kubernetes secrets backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
fi
```

## File System Backups

File system backups ensure that you can restore application files, static assets, and logs in case of a failure.

### Application Files

#### Backend Application Files

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/files"
DATE=$(date +%Y%m%d_%H%M%S)
APP_DIR="/opt/cmdb-lite/backend"
BACKUP_FILE="$BACKUP_DIR/backend_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup backend application files
if [ -d "$APP_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $APP_DIR) $(basename $APP_DIR)
    echo "$(date): Backend application files backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Backend application files backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Backend application directory not found: $APP_DIR" >> $BACKUP_DIR/backup.log
fi
```

#### Frontend Application Files

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/files"
DATE=$(date +%Y%m%d_%H%M%S)
APP_DIR="/opt/cmdb-lite/frontend"
BACKUP_FILE="$BACKUP_DIR/frontend_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup frontend application files
if [ -d "$APP_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $APP_DIR) $(basename $APP_DIR)
    echo "$(date): Frontend application files backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Frontend application files backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Frontend application directory not found: $APP_DIR" >> $BACKUP_DIR/backup.log
fi
```

### Static Assets

#### Web Assets

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/files"
DATE=$(date +%Y%m%d_%H%M%S)
WEB_DIR="/var/www/cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/web-assets_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup web assets
if [ -d "$WEB_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $WEB_DIR) $(basename $WEB_DIR)
    echo "$(date): Web assets backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Web assets backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): Web assets directory not found: $WEB_DIR" >> $BACKUP_DIR/backup.log
fi
```

#### SSL Certificates

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/files"
DATE=$(date +%Y%m%d_%H%M%S)
SSL_DIR="/etc/letsencrypt/live/your-domain.com"
BACKUP_FILE="$BACKUP_DIR/ssl-certificates_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup SSL certificates
if [ -d "$SSL_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $SSL_DIR) $(basename $SSL_DIR)
    echo "$(date): SSL certificates backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): SSL certificates backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): SSL certificates directory not found: $SSL_DIR" >> $BACKUP_DIR/backup.log
fi
```

### Log Files

#### Application Logs

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/logs"
DATE=$(date +%Y%m%d_%H%M%S)
LOG_DIR="/var/log/cmdb-lite"
BACKUP_FILE="$BACKUP_DIR/application-logs_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup application logs
if [ -d "$LOG_DIR" ]; then
    tar -czf $BACKUP_FILE -C $(dirname $LOG_DIR) $(basename $LOG_DIR)
    echo "$(date): Application logs backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): Application logs backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
    
    # Rotate logs
    find $LOG_DIR -name "*.log" -type f -mtime +30 -delete
    echo "$(date): Rotated application logs older than 30 days" >> $BACKUP_DIR/backup.log
else
    echo "$(date): Application logs directory not found: $LOG_DIR" >> $BACKUP_DIR/backup.log
fi
```

#### System Logs

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/logs"
DATE=$(date +%Y%m%d_%H%M%S)
SYSTEM_LOG_DIR="/var/log"
BACKUP_FILE="$BACKUP_DIR/system-logs_$DATE.tar.gz"

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Backup system logs
if [ -d "$SYSTEM_LOG_DIR" ]; then
    tar -czf $BACKUP_FILE -C $SYSTEM_LOG_DIR auth.log kern.log syslog
    echo "$(date): System logs backed up: $BACKUP_FILE" >> $BACKUP_DIR/backup.log
    
    # Encrypt backup
    if command -v gpg &> /dev/null; then
        gpg --batch --yes --passphrase "your-encryption-password" --cipher-algo AES256 --symmetric --output $BACKUP_FILE.gpg $BACKUP_FILE
        rm $BACKUP_FILE
        echo "$(date): System logs backup encrypted: $BACKUP_FILE.gpg" >> $BACKUP_DIR/backup.log
    fi
else
    echo "$(date): System logs directory not found: $SYSTEM_LOG_DIR" >> $BACKUP_DIR/backup.log
fi
```

## Backup Storage

Choosing the right storage solution for backups is critical for ensuring data durability and availability.

### Local Storage

#### Disk Storage

Local disk storage is the simplest option for storing backups, but it lacks redundancy and protection against local disasters.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
MOUNT_POINT="/mnt/backup-disk"

# Check if backup disk is mounted
if mountpoint -q $MOUNT_POINT; then
    echo "$(date): Backup disk is mounted at $MOUNT_POINT" >> $BACKUP_DIR/backup.log
    
    # Create backup directory on backup disk
    mkdir -p $MOUNT_POINT/cmdb-lite
    
    # Copy backups to backup disk
    rsync -av $BACKUP_DIR/ $MOUNT_POINT/cmdb-lite/
    echo "$(date): Backups copied to backup disk" >> $BACKUP_DIR/backup.log
else
    echo "$(date): Backup disk is not mounted at $MOUNT_POINT" >> $BACKUP_DIR/backup.log
    
    # Try to mount backup disk
    if mount /dev/sdb1 $MOUNT_POINT; then
        echo "$(date): Backup disk mounted successfully" >> $BACKUP_DIR/backup.log
        
        # Create backup directory on backup disk
        mkdir -p $MOUNT_POINT/cmdb-lite
        
        # Copy backups to backup disk
        rsync -av $BACKUP_DIR/ $MOUNT_POINT/cmdb-lite/
        echo "$(date): Backups copied to backup disk" >> $BACKUP_DIR/backup.log
    else
        echo "$(date): Failed to mount backup disk" >> $BACKUP_DIR/backup.log
        # Send alert
        send_alert "Failed to mount backup disk"
    fi
fi
```

#### Network Attached Storage (NAS)

NAS provides network-accessible storage with redundancy features like RAID.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
NAS_MOUNT_POINT="/mnt/nas"
NAS_SERVER="nas.example.com"
NAS_SHARE="backups"
NAS_USER="backup-user"
NAS_PASSWORD="backup-password"

# Check if NAS is mounted
if mountpoint -q $NAS_MOUNT_POINT; then
    echo "$(date): NAS is mounted at $NAS_MOUNT_POINT" >> $BACKUP_DIR/backup.log
    
    # Create backup directory on NAS
    mkdir -p $NAS_MOUNT_POINT/cmdb-lite
    
    # Copy backups to NAS
    rsync -av $BACKUP_DIR/ $NAS_MOUNT_POINT/cmdb-lite/
    echo "$(date): Backups copied to NAS" >> $BACKUP_DIR/backup.log
else
    echo "$(date): NAS is not mounted at $NAS_MOUNT_POINT" >> $BACKUP_DIR/backup.log
    
    # Try to mount NAS
    if mount -t cifs //$NAS_SERVER/$NAS_SHARE $NAS_MOUNT_POINT -o username=$NAS_USER,password=$NAS_PASSWORD; then
        echo "$(date): NAS mounted successfully" >> $BACKUP_DIR/backup.log
        
        # Create backup directory on NAS
        mkdir -p $NAS_MOUNT_POINT/cmdb-lite
        
        # Copy backups to NAS
        rsync -av $BACKUP_DIR/ $NAS_MOUNT_POINT/cmdb-lite/
        echo "$(date): Backups copied to NAS" >> $BACKUP_DIR/backup.log
    else
        echo "$(date): Failed to mount NAS" >> $BACKUP_DIR/backup.log
        # Send alert
        send_alert "Failed to mount NAS"
    fi
fi
```

### Cloud Storage

#### Amazon S3

Amazon S3 provides durable, scalable object storage with high availability.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
S3_BUCKET="your-backup-bucket"
S3_PREFIX="cmdb-lite"

# Check if AWS CLI is installed
if command -v aws &> /dev/null; then
    echo "$(date): AWS CLI is installed" >> $BACKUP_DIR/backup.log
    
    # Sync backups to S3
    if aws s3 sync $BACKUP_DIR/ s3://$S3_BUCKET/$S3_PREFIX/; then
        echo "$(date): Backups synced to S3" >> $BACKUP_DIR/backup.log
    else
        echo "$(date): Failed to sync backups to S3" >> $BACKUP_DIR/backup.log
        # Send alert
        send_alert "Failed to sync backups to S3"
    fi
else
    echo "$(date): AWS CLI is not installed" >> $BACKUP_DIR/backup.log
    # Send alert
    send_alert "AWS CLI is not installed"
fi
```

#### Google Cloud Storage

Google Cloud Storage provides durable, scalable object storage with high availability.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
GCS_BUCKET="your-backup-bucket"
GCS_PREFIX="cmdb-lite"

# Check if gsutil is installed
if command -v gsutil &> /dev/null; then
    echo "$(date): gsutil is installed" >> $BACKUP_DIR/backup.log
    
    # Sync backups to GCS
    if gsutil -m rsync -r $BACKUP_DIR/ gs://$GCS_BUCKET/$GCS_PREFIX/; then
        echo "$(date): Backups synced to GCS" >> $BACKUP_DIR/backup.log
    else
        echo "$(date): Failed to sync backups to GCS" >> $BACKUP_DIR/backup.log
        # Send alert
        send_alert "Failed to sync backups to GCS"
    fi
else
    echo "$(date): gsutil is not installed" >> $BACKUP_DIR/backup.log
    # Send alert
    send_alert "gsutil is not installed"
fi
```

### Offsite Storage

#### Tape Storage

Tape storage provides a cost-effective solution for long-term archival of backups.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
TAPE_DEVICE="/dev/nst0"
DATE=$(date +%Y%m%d_%H%M%S)
LOG_FILE="$BACKUP_DIR/backup.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Check if tape device is available
if [ -e "$TAPE_DEVICE" ]; then
    echo "$(date): Tape device is available: $TAPE_DEVICE" >> $LOG_FILE
    
    # Rewind tape
    mt -f $TAPE_DEVICE rewind
    echo "$(date): Tape rewound" >> $LOG_FILE
    
    # Write backups to tape
    if tar -czf $TAPE_DEVICE $BACKUP_DIR; then
        echo "$(date): Backups written to tape" >> $LOG_FILE
        
        # Eject tape
        mt -f $TAPE_DEVICE offline
        echo "$(date): Tape ejected" >> $LOG_FILE
    else
        echo "$(date): Failed to write backups to tape" >> $LOG_FILE
        send_alert "Failed to write backups to tape"
    fi
else
    echo "$(date): Tape device is not available: $TAPE_DEVICE" >> $LOG_FILE
    send_alert "Tape device is not available"
fi
```

#### Offsite Server

An offsite server provides a remote location for storing backups, protecting against local disasters.

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
OFFSITE_SERVER="backup.example.com"
OFFSITE_USER="backup-user"
OFFSITE_DIR="/backups/cmdb-lite"
SSH_KEY="/home/backup-user/.ssh/id_rsa"
LOG_FILE="$BACKUP_DIR/backup.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Check if SSH key is available
if [ -f "$SSH_KEY" ]; then
    echo "$(date): SSH key is available: $SSH_KEY" >> $LOG_FILE
    
    # Create backup directory on offsite server
    ssh -i $SSH_KEY $OFFSITE_USER@$OFFSITE_SERVER "mkdir -p $OFFSITE_DIR"
    echo "$(date): Backup directory created on offsite server" >> $LOG_FILE
    
    # Sync backups to offsite server
    if rsync -av -e "ssh -i $SSH_KEY" $BACKUP_DIR/ $OFFSITE_USER@$OFFSITE_SERVER:$OFFSITE_DIR/; then
        echo "$(date): Backups synced to offsite server" >> $LOG_FILE
    else
        echo "$(date): Failed to sync backups to offsite server" >> $LOG_FILE
        send_alert "Failed to sync backups to offsite server"
    fi
else
    echo "$(date): SSH key is not available: $SSH_KEY" >> $LOG_FILE
    send_alert "SSH key is not available"
fi
```

## Backup Verification

Regularly verifying that backups can be restored successfully is essential for ensuring the effectiveness of the backup strategy.

### Automated Verification

#### Database Backup Verification

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/database"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite_verification"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
LOG_FILE="$BACKUP_DIR/verification.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Verification Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "cmdb-lite-full_*.sql.gz" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No backup found for verification"
    send_alert "No backup found for verification"
    exit 1
fi

log "Starting verification of backup: $LATEST_BACKUP"

# Create verification database
dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME 2>/dev/null || true
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME

# Restore backup to verification database
if gunzip -c $LATEST_BACKUP | psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; then
    log "Backup restored successfully to verification database"
    
    # Count tables in verification database
    TABLE_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM information_schema.tables WHERE table_schema = 'public';")
    log "Table count in verification database: $TABLE_COUNT"
    
    # Count rows in key tables
    CI_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM configuration_items;")
    log "Configuration item count in verification database: $CI_COUNT"
    
    RELATIONSHIP_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM relationships;")
    log "Relationship count in verification database: $RELATIONSHIP_COUNT"
    
    USER_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -t -c "SELECT count(*) FROM users;")
    log "User count in verification database: $USER_COUNT"
    
    # Check if counts are reasonable
    if [ "$TABLE_COUNT" -gt 0 ] && [ "$CI_COUNT" -gt 0 ] && [ "$RELATIONSHIP_COUNT" -gt 0 ] && [ "$USER_COUNT" -gt 0 ]; then
        log "Backup verification successful"
    else
        log "Backup verification failed: Table or row counts are zero"
        send_alert "CMDB Lite backup verification failed: Table or row counts are zero"
    fi
else
    log "Backup restoration failed"
    send_alert "CMDB Lite backup restoration failed"
fi

# Drop verification database
dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME

log "Backup verification process completed"
```

#### Configuration Backup Verification

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
LOG_FILE="$BACKUP_DIR/verification.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Verification Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "env_*.txt.gpg" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No configuration backup found for verification"
    send_alert "No configuration backup found for verification"
    exit 1
fi

log "Starting verification of configuration backup: $LATEST_BACKUP"

# Create temporary directory for verification
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Decrypt backup
if gpg --batch --yes --passphrase "your-encryption-password" --output $TEMP_DIR/env.txt --decrypt $LATEST_BACKUP; then
    log "Configuration backup decrypted successfully"
    
    # Check if environment file is valid
    if [ -f "$TEMP_DIR/env.txt" ] && [ -s "$TEMP_DIR/env.txt" ]; then
        log "Configuration backup is valid"
        
        # Count environment variables
        ENV_VAR_COUNT=$(wc -l < $TEMP_DIR/env.txt)
        log "Environment variable count: $ENV_VAR_COUNT"
        
        # Check for required environment variables
        REQUIRED_VARS=("DB_HOST" "DB_PORT" "DB_NAME" "DB_USER" "DB_PASSWORD" "JWT_SECRET")
        MISSING_VARS=0
        
        for var in "${REQUIRED_VARS[@]}"; do
            if ! grep -q "^$var=" $TEMP_DIR/env.txt; then
                log "Required environment variable missing: $var"
                MISSING_VARS=$((MISSING_VARS + 1))
            fi
        done
        
        if [ $MISSING_VARS -eq 0 ]; then
            log "All required environment variables are present"
            log "Configuration backup verification successful"
        else
            log "Configuration backup verification failed: $MISSING_VARS required environment variables are missing"
            send_alert "CMDB Lite configuration backup verification failed: Required environment variables are missing"
        fi
    else
        log "Configuration backup is invalid"
        send_alert "CMDB Lite configuration backup is invalid"
    fi
else
    log "Configuration backup decryption failed"
    send_alert "CMDB Lite configuration backup decryption failed"
fi

log "Configuration backup verification process completed"
```

### Manual Verification

Manual verification involves periodically performing a complete restoration of backups in a test environment to ensure that all components can be restored successfully.

#### Manual Verification Checklist

1. **Database Restoration**:
   - [ ] Restore the latest full backup
   - [ ] Apply incremental backups
   - [ ] Verify data integrity
   - [ ] Test application functionality

2. **Configuration Restoration**:
   - [ ] Restore application configuration
   - [ ] Restore infrastructure configuration
   - [ ] Restore secrets
   - [ ] Verify configuration validity

3. **File System Restoration**:
   - [ ] Restore application files
   - [ ] Restore static assets
   - [ ] Restore SSL certificates
   - [ ] Verify file integrity

4. **Application Testing**:
   - [ ] Start the application
   - [ ] Verify all services are running
   - [ ] Test all application features
   - [ ] Verify user access

#### Manual Verification Procedure

1. **Prepare Test Environment**:
   - Set up a test environment that mirrors production
   - Ensure sufficient resources are available
   - Document the test environment configuration

2. **Restore Backups**:
   - Follow the documented restoration procedures
   - Document any issues encountered
   - Verify all components are restored correctly

3. **Test Application**:
   - Perform comprehensive testing of all features
   - Test with realistic data volumes
   - Verify performance metrics
   - Document any issues encountered

4. **Document Results**:
   - Create a verification report
   - Include any issues encountered and resolutions
   - Update backup and restoration procedures if necessary
   - Share results with the operations team

### Recovery Testing

Recovery testing involves simulating various failure scenarios to ensure that the recovery procedures are effective and that the team can restore the system within the required timeframes.

#### Failure Scenarios

1. **Database Failure**:
   - Simulate database corruption
   - Test database restoration procedures
   - Verify data integrity after restoration

2. **Application Failure**:
   - Simulate application corruption
   - Test application restoration procedures
   - Verify application functionality after restoration

3. **Infrastructure Failure**:
   - Simulate server failure
   - Test infrastructure restoration procedures
   - Verify system availability after restoration

4. **Complete System Failure**:
   - Simulate complete system failure
   - Test complete system restoration procedures
   - Verify all components are functioning correctly after restoration

#### Recovery Testing Procedure

1. **Plan the Test**:
   - Define the failure scenario to test
   - Determine the scope of the test
   - Identify the team members involved
   - Schedule the test during a maintenance window

2. **Prepare the Environment**:
   - Create a backup of the current system
   - Ensure rollback procedures are in place
   - Prepare monitoring tools to track the recovery process

3. **Execute the Test**:
   - Introduce the failure according to the scenario
   - Follow the documented recovery procedures
   - Document any deviations from the procedures
   - Measure the time taken to recover

4. **Evaluate the Results**:
   - Verify that the system is fully recovered
   - Compare the recovery time with the objectives
   - Identify any issues with the recovery procedures
   - Document lessons learned

5. **Update Procedures**:
   - Update the recovery procedures based on lessons learned
   - Train the team on any changes
   - Schedule regular recovery tests

## Recovery Procedures

Recovery procedures are documented steps for restoring CMDB Lite after a failure. These procedures should be tested regularly and updated as needed.

### Database Recovery

#### Full Database Recovery

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/database"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
LOG_FILE="/var/log/cmdb-lite/recovery.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Recovery Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Stop the application
log "Stopping CMDB Lite application"
sudo systemctl stop cmdb-lite-backend
sudo systemctl stop cmdb-lite-frontend

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "cmdb-lite-full_*.sql.gz" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No backup found for recovery"
    send_alert "No backup found for CMDB Lite recovery"
    exit 1
fi

log "Starting recovery from backup: $LATEST_BACKUP"

# Drop the existing database
log "Dropping existing database"
dropdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME

# Create a new database
log "Creating new database"
createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME

# Restore the backup
log "Restoring backup"
if gunzip -c $LATEST_BACKUP | psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME; then
    log "Backup restored successfully"
    
    # Start the application
    log "Starting CMDB Lite application"
    sudo systemctl start cmdb-lite-backend
    sudo systemctl start cmdb-lite-frontend
    
    # Verify application is running
    if systemctl is-active --quiet cmdb-lite-backend && systemctl is-active --quiet cmdb-lite-frontend; then
        log "CMDB Lite application started successfully"
        send_alert "CMDB Lite database recovery completed successfully"
    else
        log "Failed to start CMDB Lite application"
        send_alert "Failed to start CMDB Lite application after database recovery"
    fi
else
    log "Backup restoration failed"
    send_alert "CMDB Lite backup restoration failed"
    exit 1
fi

log "Database recovery process completed"
```

#### Point-in-Time Recovery

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/database"
WAL_DIR="/backups/cmdb-lite/wal"
DB_HOST="localhost"
DB_PORT="5432"
DB_NAME="cmdb_lite"
DB_USER="cmdb_user"
DB_PASSWORD="secure_password"
RECOVERY_TIME="2023-01-01 12:00:00"
LOG_FILE="/var/log/cmdb-lite/recovery.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Recovery Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Stop the application
log "Stopping CMDB Lite application"
sudo systemctl stop cmdb-lite-backend
sudo systemctl stop cmdb-lite-frontend
sudo systemctl stop postgresql

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "cmdb-lite-full_*.sql.gz" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No backup found for recovery"
    send_alert "No backup found for CMDB Lite recovery"
    exit 1
fi

log "Starting point-in-time recovery from backup: $LATEST_BACKUP to time: $RECOVERY_TIME"

# Create a temporary directory for recovery
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Extract backup to temporary directory
log "Extracting backup to temporary directory"
mkdir -p $TEMP_DIR/data
gunzip -c $LATEST_BACKUP | tar -xf - -C $TEMP_DIR/data

# Remove existing data directory
log "Removing existing data directory"
sudo rm -rf /var/lib/postgresql/13/main/*

# Restore base backup
log "Restoring base backup"
sudo cp -r $TEMP_DIR/data/* /var/lib/postgresql/13/main/
sudo chown -R postgres:postgres /var/lib/postgresql/13/main

# Create recovery configuration
log "Creating recovery configuration"
sudo tee /var/lib/postgresql/13/main/recovery.conf > /dev/null << EOF
restore_command = 'cp $WAL_DIR/%f %p'
recovery_target_time = '$RECOVERY_TIME'
recovery_target_action = 'promote'
EOF

# Start PostgreSQL
log "Starting PostgreSQL"
sudo systemctl start postgresql

# Wait for PostgreSQL to start
sleep 10

# Check if recovery is in progress
if sudo -u postgres psql -c "SELECT pg_is_in_recovery();" | grep -q "t"; then
    log "Recovery is in progress"
    
    # Wait for recovery to complete
    while sudo -u postgres psql -c "SELECT pg_is_in_recovery();" | grep -q "t"; do
        log "Waiting for recovery to complete..."
        sleep 5
    done
    
    log "Recovery completed"
    
    # Start the application
    log "Starting CMDB Lite application"
    sudo systemctl start cmdb-lite-backend
    sudo systemctl start cmdb-lite-frontend
    
    # Verify application is running
    if systemctl is-active --quiet cmdb-lite-backend && systemctl is-active --quiet cmdb-lite-frontend; then
        log "CMDB Lite application started successfully"
        send_alert "CMDB Lite point-in-time recovery completed successfully"
    else
        log "Failed to start CMDB Lite application"
        send_alert "Failed to start CMDB Lite application after point-in-time recovery"
    fi
else
    log "Recovery is not in progress or failed"
    send_alert "CMDB Lite point-in-time recovery failed"
    exit 1
fi

log "Point-in-time recovery process completed"
```

### Configuration Recovery

#### Application Configuration Recovery

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
APP_DIR="/opt/cmdb-lite"
LOG_FILE="/var/log/cmdb-lite/recovery.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Recovery Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Stop the application
log "Stopping CMDB Lite application"
sudo systemctl stop cmdb-lite-backend
sudo systemctl stop cmdb-lite-frontend

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "env_*.txt.gpg" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No configuration backup found for recovery"
    send_alert "No configuration backup found for CMDB Lite recovery"
    exit 1
fi

log "Starting configuration recovery from backup: $LATEST_BACKUP"

# Create a temporary directory for recovery
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Decrypt backup
log "Decrypting configuration backup"
if gpg --batch --yes --passphrase "your-encryption-password" --output $TEMP_DIR/env.txt --decrypt $LATEST_BACKUP; then
    log "Configuration backup decrypted successfully"
    
    # Restore environment file
    log "Restoring environment file"
    cp $TEMP_DIR/env.txt $APP_DIR/.env
    
    # Set permissions
    chown cmdb:cmdb $APP_DIR/.env
    chmod 600 $APP_DIR/.env
    
    log "Environment file restored successfully"
    
    # Start the application
    log "Starting CMDB Lite application"
    sudo systemctl start cmdb-lite-backend
    sudo systemctl start cmdb-lite-frontend
    
    # Verify application is running
    if systemctl is-active --quiet cmdb-lite-backend && systemctl is-active --quiet cmdb-lite-frontend; then
        log "CMDB Lite application started successfully"
        send_alert "CMDB Lite configuration recovery completed successfully"
    else
        log "Failed to start CMDB Lite application"
        send_alert "Failed to start CMDB Lite application after configuration recovery"
    fi
else
    log "Configuration backup decryption failed"
    send_alert "CMDB Lite configuration backup decryption failed"
    exit 1
fi

log "Configuration recovery process completed"
```

#### Infrastructure Configuration Recovery

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/config"
LOG_FILE="/var/log/cmdb-lite/recovery.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Recovery Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Stop the application
log "Stopping CMDB Lite application"
sudo systemctl stop cmdb-lite-backend
sudo systemctl stop cmdb-lite-frontend
sudo systemctl stop nginx

# Find the latest backup
LATEST_BACKUP=$(find $BACKUP_DIR -name "docker-compose_*.tar.gz.gpg" -type f | sort | tail -n 1)

if [ -z "$LATEST_BACKUP" ]; then
    log "No infrastructure configuration backup found for recovery"
    send_alert "No infrastructure configuration backup found for CMDB Lite recovery"
    exit 1
fi

log "Starting infrastructure configuration recovery from backup: $LATEST_BACKUP"

# Create a temporary directory for recovery
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Decrypt backup
log "Decrypting infrastructure configuration backup"
if gpg --batch --yes --passphrase "your-encryption-password" --output $TEMP_DIR/docker-compose.tar.gz --decrypt $LATEST_BACKUP; then
    log "Infrastructure configuration backup decrypted successfully"
    
    # Extract backup
    log "Extracting infrastructure configuration backup"
    tar -xzf $TEMP_DIR/docker-compose.tar.gz -C $TEMP_DIR
    
    # Restore Docker Compose files
    log "Restoring Docker Compose files"
    cp $TEMP_DIR/docker-compose*.yml /opt/cmdb-lite/
    
    # Set permissions
    chown cmdb:cmdb /opt/cmdb-lite/docker-compose*.yml
    
    log "Docker Compose files restored successfully"
    
    # Find the latest Nginx backup
    LATEST_NGINX_BACKUP=$(find $BACKUP_DIR -name "nginx_*.tar.gz.gpg" -type f | sort | tail -n 1)
    
    if [ -n "$LATEST_NGINX_BACKUP" ]; then
        log "Starting Nginx configuration recovery from backup: $LATEST_NGINX_BACKUP"
        
        # Decrypt Nginx backup
        log "Decrypting Nginx configuration backup"
        if gpg --batch --yes --passphrase "your-encryption-password" --output $TEMP_DIR/nginx.tar.gz --decrypt $LATEST_NGINX_BACKUP; then
            log "Nginx configuration backup decrypted successfully"
            
            # Extract Nginx backup
            log "Extracting Nginx configuration backup"
            tar -xzf $TEMP_DIR/nginx.tar.gz -C $TEMP_DIR
            
            # Restore Nginx configuration
            log "Restoring Nginx configuration"
            cp -r $TEMP_DIR/sites-available/* /etc/nginx/sites-available/
            cp -r $TEMP_DIR/sites-enabled/* /etc/nginx/sites-enabled/
            
            # Test Nginx configuration
            log "Testing Nginx configuration"
            if nginx -t; then
                log "Nginx configuration is valid"
            else
                log "Nginx configuration is invalid"
                send_alert "CMDB Lite Nginx configuration is invalid after recovery"
            fi
        else
            log "Nginx configuration backup decryption failed"
            send_alert "CMDB Lite Nginx configuration backup decryption failed"
        fi
    fi
    
    # Start the services
    log "Starting services"
    sudo systemctl start nginx
    sudo systemctl start cmdb-lite-backend
    sudo systemctl start cmdb-lite-frontend
    
    # Verify services are running
    if systemctl is-active --quiet nginx && systemctl is-active --quiet cmdb-lite-backend && systemctl is-active --quiet cmdb-lite-frontend; then
        log "All services started successfully"
        send_alert "CMDB Lite infrastructure configuration recovery completed successfully"
    else
        log "Failed to start all services"
        send_alert "Failed to start all services after infrastructure configuration recovery"
    fi
else
    log "Infrastructure configuration backup decryption failed"
    send_alert "CMDB Lite infrastructure configuration backup decryption failed"
    exit 1
fi

log "Infrastructure configuration recovery process completed"
```

### Complete System Recovery

#### Complete System Recovery Procedure

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
LOG_FILE="/var/log/cmdb-lite/recovery.log"

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Recovery Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

log "Starting complete system recovery"

# 1. Recover infrastructure configuration
log "Step 1: Recovering infrastructure configuration"
sudo /usr/local/bin/cmdb-lite-recover-infrastructure-config.sh

# 2. Recover application files
log "Step 2: Recovering application files"
sudo /usr/local/bin/cmdb-lite-recover-application-files.sh

# 3. Recover database
log "Step 3: Recovering database"
sudo /usr/local/bin/cmdb-lite-recover-database.sh

# 4. Recover configuration
log "Step 4: Recovering application configuration"
sudo /usr/local/bin/cmdb-lite-recover-application-config.sh

# 5. Recover static assets
log "Step 5: Recovering static assets"
sudo /usr/local/bin/cmdb-lite-recover-static-assets.sh

# 6. Recover SSL certificates
log "Step 6: Recovering SSL certificates"
sudo /usr/local/bin/cmdb-lite-recover-ssl-certificates.sh

# 7. Start services
log "Step 7: Starting services"
sudo systemctl start nginx
sudo systemctl start cmdb-lite-backend
sudo systemctl start cmdb-lite-frontend

# 8. Verify services are running
log "Step 8: Verifying services are running"
if systemctl is-active --quiet nginx && systemctl is-active --quiet cmdb-lite-backend && systemctl is-active --quiet cmdb-lite-frontend; then
    log "All services started successfully"
    send_alert "CMDB Lite complete system recovery completed successfully"
else
    log "Failed to start all services"
    send_alert "Failed to start all services after complete system recovery"
    exit 1
fi

# 9. Test application functionality
log "Step 9: Testing application functionality"
# Add application testing commands here

log "Complete system recovery process completed"
```

## Disaster Recovery

Disaster recovery involves planning for and recovering from major disasters that could significantly impact the business.

### Disaster Recovery Plan

A disaster recovery plan (DRP) outlines the steps to recover CMDB Lite after a major disaster. The plan should include:

1. **Disaster Recovery Team**: Identify the team members responsible for disaster recovery
2. **Recovery Objectives**: Define recovery time objectives (RTO) and recovery point objectives (RPO)
3. **Recovery Sites**: Identify primary and secondary recovery sites
4. **Recovery Procedures**: Document the procedures for recovering CMDB Lite
5. **Communication Plan**: Outline how to communicate with stakeholders during a disaster
6. **Testing Plan**: Schedule regular testing of the disaster recovery plan

#### Disaster Recovery Team

| Role | Responsibilities | Contact Information |
|------|-----------------|---------------------|
| Disaster Recovery Coordinator | Overall coordination of disaster recovery efforts | coordinator@your-domain.com |
| Database Administrator | Database recovery and verification | dba@your-domain.com |
| System Administrator | System recovery and configuration | sysadmin@your-domain.com |
| Application Administrator | Application recovery and testing | appadmin@your-domain.com |
| Communications Officer | Communication with stakeholders | comms@your-domain.com |

### Recovery Time Objectives

Recovery Time Objective (RTO) is the maximum acceptable time to recover a system after a disaster.

| Component | RTO | Recovery Priority |
|-----------|-----|-------------------|
| Database | 4 hours | Critical |
| Backend Application | 4 hours | Critical |
| Frontend Application | 4 hours | Critical |
| Configuration | 2 hours | High |
| Static Assets | 8 hours | Medium |
| Logs | 24 hours | Low |

### Recovery Point Objectives

Recovery Point Objective (RPO) is the maximum acceptable amount of data loss, measured in time.

| Component | RPO | Data Criticality |
|-----------|-----|------------------|
| Database | 1 hour | Critical |
| Configuration | 24 hours | High |
| Static Assets | 24 hours | Medium |
| Logs | 7 days | Low |

## Backup Security

Backup security is essential for protecting sensitive data and ensuring the integrity of backups.

### Encryption

Encrypting backups protects sensitive data from unauthorized access.

#### Using GPG for Encryption

```bash
#!/bin/bash

# Configuration
BACKUP_FILE="/backups/cmdb-lite/database/cmdb-lite-full_20230101_120000.sql"
ENCRYPTED_FILE="$BACKUP_FILE.gpg"
PASSPHRASE="your-encryption-password"

# Encrypt backup
gpg --batch --yes --passphrase "$PASSPHRASE" --cipher-algo AES256 --symmetric --output $ENCRYPTED_FILE $BACKUP_FILE

# Verify encryption
if [ -f "$ENCRYPTED_FILE" ]; then
    echo "Backup encrypted successfully: $ENCRYPTED_FILE"
    
    # Remove original file
    rm $BACKUP_FILE
    echo "Original backup file removed"
else
    echo "Backup encryption failed"
    exit 1
fi
```

#### Using OpenSSL for Encryption

```bash
#!/bin/bash

# Configuration
BACKUP_FILE="/backups/cmdb-lite/database/cmdb-lite-full_20230101_120000.sql"
ENCRYPTED_FILE="$BACKUP_FILE.enc"
PASSPHRASE="your-encryption-password"

# Encrypt backup
openssl enc -aes-256-cbc -salt -in $BACKUP_FILE -out $ENCRYPTED_FILE -k $PASSPHRASE

# Verify encryption
if [ -f "$ENCRYPTED_FILE" ]; then
    echo "Backup encrypted successfully: $ENCRYPTED_FILE"
    
    # Remove original file
    rm $BACKUP_FILE
    echo "Original backup file removed"
else
    echo "Backup encryption failed"
    exit 1
fi
```

### Access Control

Implementing access controls ensures that only authorized personnel can access and restore backups.

#### File Permissions

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
BACKUP_USER="backup-user"
BACKUP_GROUP="backup"

# Set ownership of backup directory
chown -R $BACKUP_USER:$BACKUP_GROUP $BACKUP_DIR

# Set permissions for backup directory
chmod -R 750 $BACKUP_DIR

# Set permissions for backup files
find $BACKUP_DIR -type f -name "*.sql" -exec chmod 640 {} \;
find $BACKUP_DIR -type f -name "*.sql.gz" -exec chmod 640 {} \;
find $BACKUP_DIR -type f -name "*.dump" -exec chmod 640 {} \;
find $BACKUP_DIR -type f -name "*.tar.gz" -exec chmod 640 {} \;
find $BACKUP_DIR -type f -name "*.gpg" -exec chmod 640 {} \;

echo "Backup directory permissions set"
```

#### sudo Configuration

```bash
# /etc/sudoers.d/cmdb-lite-backup

# Allow backup-user to run backup scripts as root without password
backup-user ALL=(root) NOPASSWD: /usr/local/bin/cmdb-lite-backup-*.sh

# Allow backup-user to stop and start services
backup-user ALL=(root) NOPASSWD: /bin/systemctl stop cmdb-lite-*
backup-user ALL=(root) NOPASSWD: /bin/systemctl start cmdb-lite-*
backup-user ALL=(root) NOPASSWD: /bin/systemctl restart cmdb-lite-*

# Allow backup-user to access PostgreSQL
backup-user ALL=(postgres) NOPASSWD: /usr/bin/pg_dump, /usr/bin/pg_dumpall, /usr/bin/psql, /usr/bin/createdb, /usr/bin/dropdb
```

### Audit Logging

Audit logging tracks who accessed backups and what actions they performed.

#### Using auditd

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"

# Install auditd if not already installed
if ! command -v auditctl &> /dev/null; then
    echo "Installing auditd..."
    sudo apt-get install -y auditd audispd-plugins
fi

# Add audit rules for backup directory
sudo auditctl -w $BACKUP_DIR -p rwxa -k cmdb-lite-backup

# Add audit rules for backup scripts
sudo auditctl -w /usr/local/bin/cmdb-lite-backup-*.sh -p x -k cmdb-lite-backup-script

# Add audit rules for backup restoration
sudo auditctl -w /usr/local/bin/cmdb-lite-recover-*.sh -p x -k cmdb-lite-recover-script

echo "Audit rules added for backup directory and scripts"
```

#### Custom Audit Logging

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
AUDIT_LOG="$BACKUP_DIR/audit.log"
USER=$(whoami)
TIMESTAMP=$(date +"%Y-%m-%d %H:%M:%S")

# Function to log audit events
log_audit() {
    local action=$1
    local target=$2
    local result=$3
    
    echo "$TIMESTAMP $USER $action $target $result" >> $AUDIT_LOG
}

# Example usage
if [ "$1" = "backup" ]; then
    log_audit "BACKUP" "$2" "SUCCESS"
elif [ "$1" = "restore" ]; then
    log_audit "RESTORE" "$2" "SUCCESS"
else
    log_audit "UNKNOWN" "$2" "FAILURE"
fi
```

## Backup Maintenance

Regular maintenance of backups ensures that they remain effective and don't consume excessive storage.

### Backup Rotation

Implementing a backup rotation strategy ensures that you have multiple recovery points while managing storage usage.

#### Grandfather-Father-Son (GFS) Rotation

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite/database"
DAILY_RETENTION=7
WEEKLY_RETENTION=4
MONTHLY_RETENTION=12
LOG_FILE="$BACKUP_DIR/rotation.log"

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Get current day of week (0-6, where 0 is Sunday)
DAY_OF_WEEK=$(date +%w)

# Get current day of month (1-31)
DAY_OF_MONTH=$(date +%d)

# Get current month (1-12)
MONTH=$(date +%m)

log "Starting backup rotation"

# Rotate daily backups (keep last 7 days)
log "Rotating daily backups"
find $BACKUP_DIR -name "cmdb-lite-daily_*.sql.gz" -type f -mtime +$DAILY_RETENTION -delete
log "Daily backups older than $DAILY_RETENTION days removed"

# Rotate weekly backups (keep last 4 weeks)
if [ "$DAY_OF_WEEK" = "0" ]; then
    log "Creating weekly backup"
    LATEST_DAILY=$(find $BACKUP_DIR -name "cmdb-lite-daily_*.sql.gz" -type f | sort | tail -n 1)
    if [ -n "$LATEST_DAILY" ]; then
        WEEKLY_BACKUP="$BACKUP_DIR/cmdb-lite-weekly_$(date +%Y%m%d).sql.gz"
        cp $LATEST_DAILY $WEEKLY_BACKUP
        log "Weekly backup created: $WEEKLY_BACKUP"
    fi
    
    log "Rotating weekly backups"
    find $BACKUP_DIR -name "cmdb-lite-weekly_*.sql.gz" -type f -mtime +$((WEEKLY_RETENTION * 7)) -delete
    log "Weekly backups older than $WEEKLY_RETENTION weeks removed"
fi

# Rotate monthly backups (keep last 12 months)
if [ "$DAY_OF_MONTH" = "01" ]; then
    log "Creating monthly backup"
    LATEST_DAILY=$(find $BACKUP_DIR -name "cmdb-lite-daily_*.sql.gz" -type f | sort | tail -n 1)
    if [ -n "$LATEST_DAILY" ]; then
        MONTHLY_BACKUP="$BACKUP_DIR/cmdb-lite-monthly_$(date +%Y%m).sql.gz"
        cp $LATEST_DAILY $MONTHLY_BACKUP
        log "Monthly backup created: $MONTHLY_BACKUP"
    fi
    
    log "Rotating monthly backups"
    find $BACKUP_DIR -name "cmdb-lite-monthly_*.sql.gz" -type f -mtime +$((MONTHLY_RETENTION * 30)) -delete
    log "Monthly backups older than $MONTHLY_RETENTION months removed"
fi

log "Backup rotation completed"
```

### Backup Cleanup

Regular cleanup of old backups ensures that storage space is not wasted.

#### Automated Backup Cleanup

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
RETENTION_DAYS=30
LOG_FILE="$BACKUP_DIR/cleanup.log"

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Cleanup Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

log "Starting backup cleanup"

# Clean up old database backups
log "Cleaning up old database backups"
DB_BACKUP_COUNT=$(find $BACKUP_DIR/database -name "*.sql.gz" -type f -mtime +$RETENTION_DAYS | wc -l)
find $BACKUP_DIR/database -name "*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
log "Removed $DB_BACKUP_COUNT old database backups"

# Clean up old configuration backups
log "Cleaning up old configuration backups"
CONFIG_BACKUP_COUNT=$(find $BACKUP_DIR/config -name "*.gpg" -type f -mtime +$RETENTION_DAYS | wc -l)
find $BACKUP_DIR/config -name "*.gpg" -type f -mtime +$RETENTION_DAYS -delete
log "Removed $CONFIG_BACKUP_COUNT old configuration backups"

# Clean up old file backups
log "Cleaning up old file backups"
FILE_BACKUP_COUNT=$(find $BACKUP_DIR/files -name "*.tar.gz.gpg" -type f -mtime +$RETENTION_DAYS | wc -l)
find $BACKUP_DIR/files -name "*.tar.gz.gpg" -type f -mtime +$RETENTION_DAYS -delete
log "Removed $FILE_BACKUP_COUNT old file backups"

# Clean up old log backups
log "Cleaning up old log backups"
LOG_BACKUP_COUNT=$(find $BACKUP_DIR/logs -name "*.tar.gz.gpg" -type f -mtime +$RETENTION_DAYS | wc -l)
find $BACKUP_DIR/logs -name "*.tar.gz.gpg" -type f -mtime +$RETENTION_DAYS -delete
log "Removed $LOG_BACKUP_COUNT old log backups"

# Check available disk space
DISK_USAGE=$(df $BACKUP_DIR | awk 'NR==2 {print $5}' | sed 's/%//')
log "Disk usage after cleanup: $DISK_USAGE%"

# Send alert if disk usage is still high
if [ "$DISK_USAGE" -gt 80 ]; then
    send_alert "CMDB Lite backup disk usage is still high after cleanup: $DISK_USAGE%"
fi

log "Backup cleanup completed"
```

### Backup Monitoring

Monitoring backups ensures that they are running successfully and that any issues are detected promptly.

#### Backup Monitoring Script

```bash
#!/bin/bash

# Configuration
BACKUP_DIR="/backups/cmdb-lite"
LOG_FILE="$BACKUP_DIR/monitoring.log"
MAX_AGE_HOURS=24

# Function to log messages
log() {
    local message=$1
    echo "$(date): $message" >> $LOG_FILE
}

# Function to send alerts
send_alert() {
    local message=$1
    # Send email
    echo "$message" | mail -s "CMDB Lite Backup Monitoring Alert" ops-team@your-domain.com
    # Send to Slack
    curl -X POST -H 'Content-type: application/json' --data "{\"text\":\"$message\"}" https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
}

log "Starting backup monitoring"

# Check if backup directory exists
if [ ! -d "$BACKUP_DIR" ]; then
    log "Backup directory does not exist: $BACKUP_DIR"
    send_alert "CMDB Lite backup directory does not exist: $BACKUP_DIR"
    exit 1
fi

# Check for recent database backups
log "Checking for recent database backups"
RECENT_DB_BACKUP=$(find $BACKUP_DIR/database -name "*.sql.gz" -type f -mtime -$((MAX_AGE_HOURS/24)) | sort | tail -n 1)
if [ -z "$RECENT_DB_BACKUP" ]; then
    log "No recent database backup found (within $MAX_AGE_HOURS hours)"
    send_alert "CMDB Lite: No recent database backup found (within $MAX_AGE_HOURS hours)"
else
    log "Recent database backup found: $RECENT_DB_BACKUP"
    
    # Check backup size
    BACKUP_SIZE=$(du -h $RECENT_DB_BACKUP | cut -f1)
    log "Database backup size: $BACKUP_SIZE"
    
    # Check if backup is empty
    if [ $(gzip -l $RECENT_DB_BACKUP | awk 'NR==2 {print $2}') -eq 0 ]; then
        log "Database backup is empty: $RECENT_DB_BACKUP"
        send_alert "CMDB Lite: Database backup is empty: $RECENT_DB_BACKUP"
    fi
fi

# Check for recent configuration backups
log "Checking for recent configuration backups"
RECENT_CONFIG_BACKUP=$(find $BACKUP_DIR/config -name "*.gpg" -type f -mtime -$((MAX_AGE_HOURS/24)) | sort | tail -n 1)
if [ -z "$RECENT_CONFIG_BACKUP" ]; then
    log "No recent configuration backup found (within $MAX_AGE_HOURS hours)"
    send_alert "CMDB Lite: No recent configuration backup found (within $MAX_AGE_HOURS hours)"
else
    log "Recent configuration backup found: $RECENT_CONFIG_BACKUP"
fi

# Check for recent file backups
log "Checking for recent file backups"
RECENT_FILE_BACKUP=$(find $BACKUP_DIR/files -name "*.tar.gz.gpg" -type f -mtime -$((MAX_AGE_HOURS/24)) | sort | tail -n 1)
if [ -z "$RECENT_FILE_BACKUP" ]; then
    log "No recent file backup found (within $MAX_AGE_HOURS hours)"
    send_alert "CMDB Lite: No recent file backup found (within $MAX_AGE_HOURS hours)"
else
    log "Recent file backup found: $RECENT_FILE_BACKUP"
fi

# Check backup log for errors
log "Checking backup log for errors"
if [ -f "$BACKUP_DIR/database/backup.log" ]; then
    ERROR_COUNT=$(grep -i "error\|failed" $BACKUP_DIR/database/backup.log | wc -l)
    if [ "$ERROR_COUNT" -gt 0 ]; then
        log "Found $ERROR_COUNT errors in backup log"
        send_alert "CMDB Lite: Found $ERROR_COUNT errors in backup log"
    else
        log "No errors found in backup log"
    fi
else
    log "Backup log not found: $BACKUP_DIR/database/backup.log"
fi

# Check available disk space
log "Checking available disk space"
DISK_USAGE=$(df $BACKUP_DIR | awk 'NR==2 {print $5}' | sed 's/%//')
log "Disk usage: $DISK_USAGE%"

# Send alert if disk usage is high
if [ "$DISK_USAGE" -gt 80 ]; then
    send_alert "CMDB Lite backup disk usage is high: $DISK_USAGE%"
fi

log "Backup monitoring completed"
```

For more information on operating CMDB Lite, see the [Operations Guide](README.md) and the [Monitoring and Troubleshooting Guide](monitoring.md).