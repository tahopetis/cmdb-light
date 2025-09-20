#!/bin/bash

# CMDB Lite Database Migration Script for Containerized Environment
# This script helps with running database migrations in Docker containers

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

# Function to create database if it doesn't exist
create_database() {
    print_step "Checking if database exists..."
    
    # Check if database exists
    DB_EXISTS=$(docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -tAc "SELECT 1 FROM pg_database WHERE datname='\${DB_NAME:-cmdb_lite}'" 2>/dev/null || echo "0")
    
    if [ "$DB_EXISTS" != "1" ]; then
        print_status "Database \${DB_NAME:-cmdb_lite} does not exist. Creating..."
        docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -c "CREATE DATABASE \${DB_NAME:-cmdb_lite};" > /dev/null 2>&1
        print_status "Database \${DB_NAME:-cmdb_lite} created successfully."
    else
        print_status "Database \${DB_NAME:-cmdb_lite} already exists."
    fi
}

# Function to apply schema files
apply_schema() {
    print_step "Applying database schema..."
    
    # Apply schema files
    docker-compose -f $COMPOSE_FILE exec -T db bash -c "
        for file in \$(find /docker-entrypoint-initdb.d/schema -name \"*.sql\" | sort); do
            echo \"Applying schema file: \$file\"
            psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} -f \"\$file\"
        done
    "
    
    print_status "Database schema applied successfully."
}

# Function to apply migrations
apply_migrations() {
    print_step "Applying database migrations..."
    
    # Apply migration files
    docker-compose -f $COMPOSE_FILE exec -T db bash -c "
        for file in \$(find /docker-entrypoint-initdb.d/migrations -name \"*.sql\" | sort); do
            echo \"Applying migration file: \$file\"
            psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} -f \"\$file\"
        done
    "
    
    print_status "Database migrations applied successfully."
}

# Function to apply seed data
apply_seeds() {
    print_step "Applying seed data..."
    
    # Apply seed files
    docker-compose -f $COMPOSE_FILE exec -T db bash -c "
        for file in \$(find /docker-entrypoint-initdb.d/seeds -name \"*.sql\" | sort); do
            echo \"Applying seed file: \$file\"
            psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} -f \"\$file\"
        done
    "
    
    print_status "Seed data applied successfully."
}

# Function to run all migrations
run_all_migrations() {
    print_step "Running all database migrations..."
    create_database
    apply_schema
    apply_migrations
    apply_seeds
    print_status "All database migrations completed successfully."
}

# Function to create a new migration
create_migration() {
    if [ -z "$2" ]; then
        print_error "Please provide a migration name."
        print_status "Usage: $0 create <migration-name>"
        exit 1
    fi
    
    MIGRATION_NAME=$2
    TIMESTAMP=$(date +%Y%m%d%H%M%S)
    MIGRATION_FILE="./database/migrations/${TIMESTAMP}_${MIGRATION_NAME}.sql"
    
    print_step "Creating new migration: $MIGRATION_NAME"
    
    # Create migration file
    cat > "$MIGRATION_FILE" << EOF
-- Migration: $MIGRATION_NAME
-- Created at: $(date)

-- Add your SQL here

EOF
    
    print_status "Migration file created: $MIGRATION_FILE"
}

# Function to show migration status
show_migration_status() {
    print_step "Database Migration Status:"
    
    # Show database information
    docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} -c "
        SELECT 
            schemaname as schema,
            tablename as table,
            tableowner as owner
        FROM 
            pg_tables 
        WHERE 
            schemaname = 'public'
        ORDER BY 
            tablename;
    "
}

# Main script logic
case "$1" in
    "create")
        create_migration "$@"
        ;;
    "status")
        check_docker
        wait_for_db
        show_migration_status
        ;;
    "schema")
        check_docker
        wait_for_db
        create_database
        apply_schema
        ;;
    "migrate")
        check_docker
        wait_for_db
        create_database
        apply_schema
        apply_migrations
        ;;
    "seed")
        check_docker
        wait_for_db
        create_database
        apply_schema
        apply_migrations
        apply_seeds
        ;;
    "all")
        check_docker
        wait_for_db
        run_all_migrations
        ;;
    *)
        echo "CMDB Lite Database Migration Script"
        echo ""
        echo "Usage: $0 {create|status|schema|migrate|seed|all}"
        echo ""
        echo "  create   - Create a new migration file"
        echo "  status   - Show database migration status"
        echo "  schema   - Apply database schema only"
        echo "  migrate  - Apply database schema and migrations"
        echo "  seed     - Apply database schema, migrations, and seed data"
        echo "  all      - Apply all database migrations (schema, migrations, and seeds)"
        echo ""
        echo "Environment: $ENVIRONMENT"
        echo "Compose file: $COMPOSE_FILE"
        exit 1
        ;;
esac