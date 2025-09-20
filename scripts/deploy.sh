#!/bin/bash

# CMDB Lite Deployment Script
# This script helps with deploying CMDB Lite using Docker Compose

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

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose."
        exit 1
    fi
}

# Function to create .env file if it doesn't exist
setup_env() {
    print_step "Setting up environment variables..."
    
    if [ ! -f .env ]; then
        print_status "Creating .env file from template..."
        if [ -f .env.example ]; then
            cp .env.example .env
            print_status ".env file created successfully. Please review and update it with your values."
        else
            print_warning "No .env.example found. Creating a basic .env file."
            cat > .env << EOF
# Environment
ENVIRONMENT=${ENVIRONMENT}

# Database
DB_HOST=db
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=cmdb_password

# Backend
SERVER_PORT=8080
JWT_SECRET=your-secret-key-change-this-in-production

# Frontend
FRONTEND_PORT=3000

# Adminer
ADMINER_PORT=8081
EOF
            print_status "Basic .env file created. Please review and update it with your values."
        fi
    else
        print_status ".env file already exists."
    fi
}

# Function to start services
start_services() {
    print_step "Starting CMDB Lite services..."
    
    # Build and start services
    docker-compose -f $COMPOSE_FILE up --build -d
    
    print_status "Services started successfully."
    print_status "Frontend: http://localhost:${FRONTEND_PORT:-3000}"
    print_status "Backend API: http://localhost:${SERVER_PORT:-8080}"
    if [ "$ENVIRONMENT" = "development" ]; then
        print_status "Adminer: http://localhost:${ADMINER_PORT:-8081}"
    fi
}

# Function to stop services
stop_services() {
    print_step "Stopping CMDB Lite services..."
    docker-compose -f $COMPOSE_FILE down
    print_status "Services stopped successfully."
}

# Function to restart services
restart_services() {
    print_step "Restarting CMDB Lite services..."
    docker-compose -f $COMPOSE_FILE restart
    print_status "Services restarted successfully."
}

# Function to show logs
show_logs() {
    print_step "Showing logs for CMDB Lite services..."
    if [ -n "$2" ]; then
        docker-compose -f $COMPOSE_FILE logs -f "$2"
    else
        docker-compose -f $COMPOSE_FILE logs -f
    fi
}

# Function to run database migrations
run_migrations() {
    print_step "Running database migrations..."
    
    # Wait for database to be ready
    print_status "Waiting for database to be ready..."
    docker-compose -f $COMPOSE_FILE exec -T db bash -c "until pg_isready -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite}; do sleep 1; done"
    
    # Run migrations
    docker-compose -f $COMPOSE_FILE exec -T db bash -c "cd /docker-entrypoint-initdb.d && ./migrate.sh seed"
    
    print_status "Database migrations completed successfully."
}

# Function to backup database
backup_database() {
    print_step "Creating database backup..."
    
    TIMESTAMP=$(date +%Y%m%d_%H%M%S)
    BACKUP_DIR="./backups"
    BACKUP_FILE="${BACKUP_DIR}/cmdb_lite_backup_${TIMESTAMP}.sql"
    
    # Create backup directory if it doesn't exist
    mkdir -p $BACKUP_DIR
    
    # Create backup
    docker-compose -f $COMPOSE_FILE exec -T db pg_dump -U \${DB_USER:-cmdb_user} \${DB_NAME:-cmdb_lite} > $BACKUP_FILE
    
    print_status "Database backup created: $BACKUP_FILE"
}

# Function to restore database
restore_database() {
    if [ -z "$2" ]; then
        print_error "Please provide the backup file path."
        print_status "Usage: $0 restore <backup-file>"
        exit 1
    fi
    
    BACKUP_FILE=$2
    if [ ! -f "$BACKUP_FILE" ]; then
        print_error "Backup file not found: $BACKUP_FILE"
        exit 1
    fi
    
    print_step "Restoring database from $BACKUP_FILE..."
    
    # Stop the backend service to prevent conflicts
    docker-compose -f $COMPOSE_FILE stop backend
    
    # Restore database
    docker-compose -f $COMPOSE_FILE exec -T db psql -U \${DB_USER:-cmdb_user} -d \${DB_NAME:-cmdb_lite} < $BACKUP_FILE
    
    # Start the backend service
    docker-compose -f $COMPOSE_FILE start backend
    
    print_status "Database restored successfully."
}

# Function to show status
show_status() {
    print_step "CMDB Lite Services Status:"
    docker-compose -f $COMPOSE_FILE ps
}

# Function to clean up unused resources
cleanup() {
    print_step "Cleaning up unused Docker resources..."
    docker system prune -f
    print_status "Cleanup completed."
}

# Main script logic
case "$1" in
    "setup")
        check_docker
        setup_env
        ;;
    "start")
        check_docker
        setup_env
        start_services
        ;;
    "stop")
        check_docker
        stop_services
        ;;
    "restart")
        check_docker
        restart_services
        ;;
    "logs")
        check_docker
        show_logs "$@"
        ;;
    "migrate")
        check_docker
        run_migrations
        ;;
    "backup")
        check_docker
        backup_database
        ;;
    "restore")
        check_docker
        restore_database "$@"
        ;;
    "status")
        check_docker
        show_status
        ;;
    "cleanup")
        check_docker
        cleanup
        ;;
    "reset")
        print_warning "This will stop all services and remove all containers, networks, and volumes."
        read -p "Are you sure you want to continue? (y/n): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            check_docker
            docker-compose -f $COMPOSE_FILE down -v --remove-orphans
            print_status "Environment reset completed."
        fi
        ;;
    *)
        echo "CMDB Lite Deployment Script"
        echo ""
        echo "Usage: $0 {setup|start|stop|restart|logs|migrate|backup|restore|status|cleanup|reset}"
        echo ""
        echo "  setup    - Set up environment variables"
        echo "  start    - Start all services"
        echo "  stop     - Stop all services"
        echo "  restart  - Restart all services"
        echo "  logs     - Show logs for services (optionally specify service name)"
        echo "  migrate  - Run database migrations"
        echo "  backup   - Create database backup"
        echo "  restore  - Restore database from backup"
        echo "  status   - Show status of all services"
        echo "  cleanup  - Clean up unused Docker resources"
        echo "  reset    - Reset environment (stops services and removes all data)"
        echo ""
        echo "Environment: $ENVIRONMENT"
        echo "Compose file: $COMPOSE_FILE"
        exit 1
        ;;
esac