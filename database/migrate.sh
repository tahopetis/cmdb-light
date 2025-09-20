#!/bin/bash

# Migration script for CMDB Lite database
# This script helps apply database schema and migrations

set -e

# Default values
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-cmdb_lite}
DB_USER=${DB_USER:-cmdb_user}
DB_PASSWORD=${DB_PASSWORD:-cmdb_password}

# Paths
SCHEMA_PATH=${SCHEMA_PATH:-./schema}
MIGRATIONS_PATH=${MIGRATIONS_PATH:-./migrations}
SEEDS_PATH=${SEEDS_PATH:-./seeds}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

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

# Check if psql is installed
if ! command -v psql &> /dev/null; then
    print_error "psql is not installed. Please install PostgreSQL client tools."
    exit 1
fi

# Function to apply schema files
apply_schema() {
    print_status "Applying database schema..."
    
    for file in $(find "$SCHEMA_PATH" -name "*.sql" | sort); do
        print_status "Applying schema file: $file"
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    done
    
    print_status "Database schema applied successfully."
}

# Function to apply migrations
apply_migrations() {
    print_status "Applying database migrations..."
    
    for file in $(find "$MIGRATIONS_PATH" -name "*.sql" | sort); do
        print_status "Applying migration file: $file"
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    done
    
    print_status "Database migrations applied successfully."
}

# Function to apply seed data
apply_seeds() {
    print_status "Applying seed data..."
    
    for file in $(find "$SEEDS_PATH" -name "*.sql" | sort); do
        print_status "Applying seed file: $file"
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$file"
    done
    
    print_status "Seed data applied successfully."
}

# Function to create database if it doesn't exist
create_database() {
    print_status "Checking if database exists..."
    
    # Check if database exists
    DB_EXISTS=$(PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")
    
    if [ "$DB_EXISTS" != "1" ]; then
        print_status "Database $DB_NAME does not exist. Creating..."
        PGPASSWORD=$DB_PASSWORD psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -c "CREATE DATABASE $DB_NAME;"
        print_status "Database $DB_NAME created successfully."
    else
        print_status "Database $DB_NAME already exists."
    fi
}

# Main script logic
case "$1" in
    "create")
        create_database
        ;;
    "schema")
        create_database
        apply_schema
        ;;
    "migrate")
        create_database
        apply_schema
        apply_migrations
        ;;
    "seed")
        create_database
        apply_schema
        apply_migrations
        apply_seeds
        ;;
    *)
        echo "Usage: $0 {create|schema|migrate|seed}"
        echo ""
        echo "  create  - Create the database"
        echo "  schema  - Apply database schema"
        echo "  migrate - Apply database schema and migrations"
        echo "  seed    - Apply database schema, migrations, and seed data"
        exit 1
        ;;
esac

print_status "Database operation completed successfully."