# CMDB Lite

A lightweight Configuration Management Database (CMDB) application built with Go backend, Vue 3 frontend, and PostgreSQL database.

## Documentation

For detailed information about CMDB Lite, please refer to our documentation:

- [User Guide](docs/user/user-guide.md) - For end users of the application
- [Developer Documentation](docs/developer/README.md) - For developers contributing to the project
- [Operator Documentation](docs/operator/README.md) - For system operators and administrators

## Features

- Configuration Item (CI) Management
- CI Type Management
- Relationship Mapping between CIs
- User Authentication and Authorization
- RESTful API
- Responsive Web Interface

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: Vue 3 with TailwindCSS
- **Database**: PostgreSQL
- **Containerization**: Docker and Docker Compose

## Prerequisites

- Docker and Docker Compose
- Node.js (for local frontend development)
- Go (for local backend development)
- PostgreSQL (for local database development)

## Quick Start with Docker

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/cmdb-lite.git
   cd cmdb-lite
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Start the application with Docker Compose:
   ```bash
   docker-compose up -d
   ```

4. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database Admin: http://localhost:8081 (username: cmdb_user, password: cmdb_password)

## Development Setup

### Backend Development

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

4. Run the backend:
   ```bash
   go run cmd/main.go
   ```

### Frontend Development

1. Navigate to the frontend directory:
   ```bash
   cd frontend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Install dependencies:
   ```bash
   npm install
   ```

4. Run the frontend:
   ```bash
   npm run serve
   ```

### Database Setup

1. Navigate to the database directory:
   ```bash
   cd database
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Create and seed the database:
   ```bash
   ./migrate.sh seed
   ```

## Development with Docker Compose

For development with hot reload:

1. Start the development environment:
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

2. Access the application:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database Admin: http://localhost:8081

## API Documentation

The backend API follows RESTful conventions. The base URL is `/api`.

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
cmdb-lite/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry point
│   ├── internal/           # Internal application packages
│   │   ├── config/         # Configuration management
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── models/         # Data models
│   │   └── repositories/   # Data access layer
│   ├── go.mod              # Go module definition
│   ├── go.sum              # Go module checksums
│   └── Dockerfile          # Docker configuration for production
├── frontend/               # Vue 3 frontend application
│   ├── public/             # Static assets
│   ├── src/                # Source code
│   │   ├── assets/         # CSS and other assets
│   │   ├── components/     # Vue components
│   │   ├── router/         # Vue router configuration
│   │   ├── stores/         # Pinia stores
│   │   ├── utils/          # Utility functions
│   │   ├── views/          # Vue views
│   │   ├── App.vue         # Root Vue component
│   │   └── main.js         # Application entry point
│   ├── package.json        # Node.js dependencies
│   ├── vite.config.js      # Vite configuration
│   └── Dockerfile          # Docker configuration for production
├── database/               # Database schema and migrations
│   ├── migrations/         # Database migration files
│   ├── schema/             # Database schema files
│   ├── seeds/              # Database seed files
│   └── migrate.sh          # Migration script
├── docs/                   # Documentation
│   ├── user/               # User documentation
│   ├── developer/          # Developer documentation
│   └── operator/           # Operator documentation
├── docker-compose.yml     # Docker Compose configuration for production
├── docker-compose.dev.yml # Docker Compose configuration for development
└── README.md               # This file
```

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details on how to contribute to this project.

## Roadmap

Please see [ROADMAP.md](ROADMAP.md) for information about the planned features and improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.