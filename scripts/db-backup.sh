#!/bin/bash

# CMDB Lite Database Backup Script
# This script helps with creating and restoring database backups

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Environment variables
ENVIRONMENT=${ENVIRONMENT:-development}
COMPOSE_FILE="docker-compose.yml"
if [ "$ENVIRONMENT" = "development" ]; then
    COMPOSE_FILE="docker-compose.dev.yml"
fi

# Default backup directory
BACKUP_DIR="${BACKUP_DIR:-./backups}"

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Function to check if Docker is running
check_docker() {
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker."
        exit 1
    fi
}

# Function to wait for database to be ready
wait_for_db() {
    print_step "Waiting for database to be ready..."
    
    # Maximum number of attempts
    MAX_ATTEMPTS=30
    ATTEMPT=0
    
    # Wait for database to be ready
    until docker-compose -f $COMPOSE_FILE exec -T db pg_isready -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} > /dev/null 2>&1; do
        ATTEMPT=$((ATTEMPT+1))
        if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
            print_error "Database is not ready after $MAX_ATTEMPTS attempts. Exiting."
            exit 1
        fi
        print_status "Waiting for database... (Attempt $ATTEMPT/$MAX_ATTEMPTS)"
        sleep 2
    done
    
    print_status "Database is ready."
}

# Function to create a backup
create_backup() {
    print_step "Creating database backup..."
    
    # Create backup directory if it doesn't exist
    mkdir -p $BACKUP_DIR
    
    # Generate timestamp
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_FILE="${BACKUP_DIR}/cmdb_lite_backup_${TIMESTAMP}.sql"
    
    # Create backup
    docker-compose -f $COMPOSE_FILE exec -T db pg_dump -U \${DB_USER:-cmdb_user} \${DB_NAME:-cmdb_lite} > $BACKUP_FILE
    
    # Compress backup
    gzip $BACKUP_FILE
    COMPRESSED_FILE="${BACKUP_FILE}.gz"
    
    # Create checksum
    sha256sum $COMPRESSED_FILE > "${COMPRESSED_FILE}.sha256"
    
    print_status "Database backup created: $COMPRESSED_FILE"
    print_status "Checksum: $(cat ${COMPRESSED_FILE}.sha256 | cut -d' ' -f1)"
}

# Function to list backups
list_backups() {
    print_step "Available backups:"
    
    if [ ! -d "$BACKUP_DIR" ] || [ -z "$(ls -A $BACKUP_DIR)" ]; then
        print_warning "No backups found in $BACKUP_DIR"
        return
    fi
    
    # List backups with details
    for file in $(ls -t $BACKUP_DIR/*.sql.gz 2>/dev/null); do
        if [ -f "$file" ]; then
            FILENAME=$(basename "$file")
            FILESIZE=$(du -h "$file" | cut -f1)
            FILEDATE=$(date -r "$file" "+%Y-%m-%d %H:%M:%S")
            CHECKSUM=""
            if [ -f "${file}.sha256" ]; then
                CHECKSUM=$(cat "${file}.sha256" | cut -d' ' -f1)
            fi
            echo "  File: $FILENAME"
            echo "  Size: $FILESIZE"
            echo "  Date: $FILEDATE"
            if [ -n "$CHECKSUM" ]; then
                echo "  Checksum: $CHECKSUM"
            fi
            echo ""
        fi
    done
}

# Function to restore from a backup
restore_backup() {
    if [ -z "$2" ]; then
        print_error "Please provide the backup file path."
        print_status "Usage: $0 restore <backup-file>"
        exit 1
    fi
    
    BACKUP_FILE=$2
    
    # Check if backup file exists
    if [ ! -f "$BACKUP_FILE" ]; then
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    # Check if backup file is compressed
    if [[ "$BACKUP_FILE" == *.gz ]]; then
        # Check checksum if available
        if [ -f "${BACKUP_FILE}.sha256" ]; then
            print_step "Verifying backup checksum..."
            CURRENT_CHECKSUM=$(sha256sum "$BACKUP_FILE" | cut -d' ' -f1)
            EXPECTED_CHECKSUM=$(cat "${BACKUP_FILE}.sha256" | cut -d' ' -f1)
            
            if [ "$CURRENT_CHECKSUM" != "$EXPECTED_CHECKSUM" ]; then
                print_error "Backup checksum verification failed!"
                print_status "Expected: $EXPECTED_CHECKSUM"
                print_status "Actual: $CURRENT_CHECKSUM"
                exit 1
            fi
            print_status "Backup checksum verified successfully."
        fi
        
        # Decompress backup
        print_step "Decompressing backup..."
        TEMP_FILE="${BACKUP_FILE%.gz}"
        gunzip -c "$BACKUP_FILE" > "$TEMP_FILE"
        BACKUP_FILE="$TEMP_FILE"
    fi
    
    print_step "Restoring database from $BACKUP_FILE..."
    
    # Stop the backend service to prevent conflicts
    docker-compose -f $COMPOSE_FILE stop backend
    
    # Wait for database to be ready
    wait_for_db
    
    # Drop existing database and create a new one
    docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -c "DROP DATABASE IF EXISTS \${DB_NAME:-cmdb_lite};" > /dev/null 2>&1
    docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -c "CREATE DATABASE \${DB_NAME:-cmdb_lite};" > /dev/null 2>&1
    
    # Restore database
    docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} < "$BACKUP_FILE"
    
    # Clean up temporary file if it was created
    if [[ "$2" == *.gz ]]; then
        rm -f "$BACKUP_FILE"
    fi
    
    # Start the backend service
    docker-compose -f $COMPOSE_FILE start backend
    
    print_status "Database restored successfully."
}

# Function to clean old backups
clean_backups() {
    # Default number of backups to keep
    KEEP=${2:-5}
    
    print_step "Cleaning old backups, keeping last $KEEP backups..."
    
    if [ ! -d "$BACKUP_DIR" ] || [ -z "$(ls -A $BACKUP_DIR)" ]; then
        print_status "No backups found in $BACKUP_DIR"
        return
    fi
    
    # Count backups
    BACKUP_COUNT=$(ls -1 $BACKUP_DIR/*.sql.gz 2>/dev/null | wc -l)
    
    if [ $BACKUP_COUNT -le $KEEP ]; then
        print_status "Only $BACKUP_COUNT backups found, nothing to clean."
        return
    fi
    
    # Remove old backups
    ls -t $BACKUP_DIR/*.sql.gz 2>/dev/null | tail -n +$(($KEEP + 1)) | xargs -r rm -f
    ls -t $BACKUP_DIR/*.sha256 2>/dev/null | tail -n +$(($KEEP + 1)) | xargs -r rm -f
    
    print_status "Old backups cleaned successfully."
}

# Function to schedule automatic backups
schedule_backup() {
    if [ -z "$2" ]; then
        print_error "Please provide a schedule in cron format."
        print_status "Usage: $0 schedule '<cron-expression>'"
        print_status "Example: $0 schedule '0 2 * * *' (daily at 2 AM)"
        exit 1
    fi
    
    CRON_SCHEDULE=$2
    
    print_step "Scheduling automatic backups with cron schedule: $CRON_SCHEDULE"
    
    # Create cron job
    CRON_JOB="$CRON_SCHEDULE cd $(pwd) && ./scripts/db-backup.sh create > /dev/null 2>&1"
    
    # Check if crontab is available
    if ! command -v crontab &> /dev/null; then
        print_error "crontab is not installed. Please install cron to schedule automatic backups."
        exit 1
    fi
    
    # Add to crontab
    (crontab -l 2>/dev/null; echo "$CRON_JOB") | crontab -
    
    print_status "Automatic backup scheduled successfully."
    print_status "To view scheduled jobs, run: crontab -l"
}

# Main script logic
case "$1" in
    "create")
        check_docker
        wait_for_db
        create_backup
        ;;
    "list")
        list_backups
        ;;
    "restore")
        check_docker
        restore_backup "$@"
        ;;
    "clean")
        clean_backups "$@"
        ;;
    "schedule")
        schedule_backup "$@"
        ;;
    *)
        echo "CMDB Lite Database Backup Script"
        echo ""
        echo "Usage: $0 {create|list|restore|clean|schedule}"
        echo ""
        echo "  create    - Create a new database backup"
        echo "  list      - List all available backups"
        echo "  restore   - Restore database from a backup file"
        echo "  clean     - Clean old backups (keep last N backups)"
        echo "  schedule  - Schedule automatic backups using cron"
        echo ""
        echo "Environment: $ENVIRONMENT"
        echo "Compose file: $COMPOSE_FILE"
        echo "Backup directory: $BACKUP_DIR"
        exit 1
        ;;
esac