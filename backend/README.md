# CMDB Lite Backend

This is the backend API for CMDB Lite, built with Go and the Gin framework.

## Features

- RESTful API for CMDB operations
- JWT-based authentication
- PostgreSQL database integration
- Configuration management
- Structured logging

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmdb-lite.git
   cd cmdb-lite/backend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Set up the database:
   ```bash
   cd ../database
   ./migrate.sh seed
   cd ../backend
   ```

5. Run the application:
   ```bash
   go run cmd/main.go
   ```

The API will be available at http://localhost:8080.

### Development with Hot Reload

1. Install Air for hot reload:
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. Run with Air:
   ```bash
   air
   ```

### Running with Docker

1. Build the Docker image:
   ```bash
   docker build -t cmdb-lite-backend .
   ```

2. Run the container:
   ```bash
   docker run -p 8080:8080 --env-file .env cmdb-lite-backend
   ```

## API Documentation

### Authentication

All API endpoints (except login) require authentication using JWT tokens.

### Endpoints

#### Configuration Items

- `GET /api/cis` - Get all configuration items
- `GET /api/cis/:id` - Get a specific configuration item
- `POST /api/cis` - Create a new configuration item
- `PUT /api/cis/:id` - Update a configuration item
- `DELETE /api/cis/:id` - Delete a configuration item

#### CI Types

- `GET /api/ci-types` - Get all CI types
- `GET /api/ci-types/:id` - Get a specific CI type
- `POST /api/ci-types` - Create a new CI type
- `PUT /api/ci-types/:id` - Update a CI type
- `DELETE /api/ci-types/:id` - Delete a CI type

#### Relationships

- `GET /api/relationships` - Get all relationships
- `GET /api/relationships/:id` - Get a specific relationship
- `POST /api/relationships` - Create a new relationship
- `PUT /api/relationships/:id` - Update a relationship
- `DELETE /api/relationships/:id` - Delete a relationship

## Project Structure

```
backend/
├── cmd/                    # Application entry point
│   └── main.go            # Main application file
├── internal/              # Internal application packages
│   ├── config/            # Configuration management
│   │   └── config.go      # Configuration handling
│   ├── handlers/          # HTTP request handlers
│   │   └── ci_handler.go  # CI-related handlers
│   ├── models/            # Data models
│   │   └── models.go      # Data model definitions
│   └── repositories/      # Data access layer
│       └── ci_repository.go # CI data access
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── .env.example           # Environment variables example
├── .air.toml              # Air configuration for hot reload
├── Dockerfile             # Docker configuration for production
└── Dockerfile.dev         # Docker configuration for development
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| SERVER_PORT | Port for the server to listen on | 8080 |
| DATABASE_URL | Database host URL | localhost:5432 |
| DATABASE_NAME | Database name | cmdb_lite |
| DATABASE_USER | Database username | cmdb_user |
| DATABASE_PASSWORD | Database password | cmdb_password |
| JWT_SECRET | Secret key for JWT signing | - |

## Testing

To run the tests:

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Commit your changes
6. Push to the branch
7. Create a Pull Request