#!/bin/bash

# CMDB Lite Initialization Script
# This script helps with setting up the CMDB Lite application for the first time

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

    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker."
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
JWT_SECRET=$(openssl rand -base64 32)

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

    # Create development environment file if it doesn't exist
    if [ "$ENVIRONMENT" = "development" ] && [ ! -f backend/.env.development ]; then
        print_status "Creating backend development environment file..."
        mkdir -p backend
        cat > backend/.env.development << EOF
# Development Environment
ENVIRONMENT=development
GIN_MODE=debug

# Database
DB_HOST=db
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=cmdb_password

# Server
SERVER_PORT=8080
JWT_SECRET=$(openssl rand -base64 32)
EOF
        print_status "Backend development environment file created."
    fi

    # Create production environment file if it doesn't exist
    if [ "$ENVIRONMENT" = "production" ] && [ ! -f backend/.env.production ]; then
        print_status "Creating backend production environment file..."
        mkdir -p backend
        cat > backend/.env.production << EOF
# Production Environment
ENVIRONMENT=production
GIN_MODE=release

# Database
DB_HOST=db
DB_PORT=5432
DB_NAME=cmdb_lite
DB_USER=cmdb_user
DB_PASSWORD=cmdb_password

# Server
SERVER_PORT=8080
JWT_SECRET=$(openssl rand -base64 32)
EOF
        print_status "Backend production environment file created."
    fi

    # Create frontend environment file if it doesn't exist
    if [ ! -f frontend/.env ]; then
        print_status "Creating frontend environment file..."
        mkdir -p frontend
        cat > frontend/.env << EOF
VITE_API_URL=http://localhost:8080
VITE_ENVIRONMENT=${ENVIRONMENT}
EOF
        print_status "Frontend environment file created."
    fi
}

# Function to create necessary directories
setup_directories() {
    print_step "Setting up directories..."
    
    # Create backups directory
    mkdir -p ./backups
    
    # Create logs directory
    mkdir -p ./logs
    
    print_status "Directories created successfully."
}

# Function to set up scripts
setup_scripts() {
    print_step "Setting up scripts..."
    
    # Make scripts executable
    chmod +x ./scripts/*.sh
    
    print_status "Scripts made executable."
}

# Function to pull Docker images
pull_images() {
    print_step "Pulling Docker images..."
    
    docker-compose -f $COMPOSE_FILE pull
    
    print_status "Docker images pulled successfully."
}

# Function to start services
start_services() {
    print_step "Starting CMDB Lite services..."
    
    # Build and start services
    docker-compose -f $COMPOSE_FILE up --build -d
    
    print_status "Services started successfully."
}

# Function to wait for services to be ready
wait_for_services() {
    print_step "Waiting for services to be ready..."
    
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
    
    # Wait for backend to be ready
    ATTEMPT=0
    until curl -s http://localhost:${SERVER_PORT:-8080}/health > /dev/null 2>&1; do
        ATTEMPT=$((ATTEMPT+1))
        if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
            print_error "Backend is not ready after $MAX_ATTEMPTS attempts. Exiting."
            exit 1
        fi
        print_status "Waiting for backend... (Attempt $ATTEMPT/$MAX_ATTEMPTS)"
        sleep 2
    done
    
    print_status "Backend is ready."
    
    # Wait for frontend to be ready
    ATTEMPT=0
    until curl -s http://localhost:${FRONTEND_PORT:-3000} > /dev/null 2>&1; do
        ATTEMPT=$((ATTEMPT+1))
        if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
            print_error "Frontend is not ready after $MAX_ATTEMPTS attempts. Exiting."
            exit 1
        fi
        print_status "Waiting for frontend... (Attempt $ATTEMPT/$MAX_ATTEMPTS)"
        sleep 2
    done
    
    print_status "Frontend is ready."
}

# Function to run database migrations
run_migrations() {
    print_step "Running database migrations..."
    
    # Run migrations using the db-migrate script
    ./scripts/db-migrate.sh all
    
    print_status "Database migrations completed successfully."
}

# Function to create admin user
create_admin_user() {
    print_step "Creating admin user..."
    
    # Create admin user using the backend API
    ADMIN_EMAIL="admin@example.com"
    ADMIN_PASSWORD="admin123"
    
    # Create admin user
    curl -s -X POST http://localhost:${SERVER_PORT:-8080}/api/users \
        -H "Content-Type: application/json" \
        -d "{\"email\":\"$ADMIN_EMAIL\",\"password\":\"$ADMIN_PASSWORD\",\"name\":\"Admin User\",\"role\":\"admin\"}" > /dev/null
    
    print_status "Admin user created successfully."
    print_status "Email: $ADMIN_EMAIL"
    print_status "Password: $ADMIN_PASSWORD"
    print_warning "Please change the admin password after first login."
}

# Function to display setup completion message
display_completion() {
    print_step "CMDB Lite setup completed successfully!"
    echo ""
    print_status "Application URLs:"
    echo "  Frontend: http://localhost:${FRONTEND_PORT:-3000}"
    echo "  Backend API: http://localhost:${SERVER_PORT:-8080}"
    if [ "$ENVIRONMENT" = "development" ]; then
        echo "  Adminer (Database UI): http://localhost:${ADMINER_PORT:-8081}"
    fi
    echo ""
    print_status "Admin Login:"
    echo "  Email: admin@example.com"
    echo "  Password: admin123"
    echo ""
    print_warning "Please change the admin password after first login."
    echo ""
    print_status "To manage the application, use the following commands:"
    echo "  ./scripts/deploy.sh start   - Start services"
    echo "  ./scripts/deploy.sh stop    - Stop services"
    echo "  ./scripts/deploy.sh logs    - View logs"
    echo "  ./scripts/deploy.sh status  - Check status"
    echo ""
    print_status "To manage the database, use the following commands:"
    echo "  ./scripts/db-migrate.sh all    - Run migrations"
    echo "  ./scripts/db-backup.sh create  - Create backup"
    echo "  ./scripts/db-backup.sh restore - Restore from backup"
    echo ""
}

# Main script logic
print_step "CMDB Lite Initialization Script"
print_step "Environment: $ENVIRONMENT"
echo ""

check_docker
setup_env
setup_directories
setup_scripts
pull_images
start_services
wait_for_services
run_migrations
create_admin_user
display_completion