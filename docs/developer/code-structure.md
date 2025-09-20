# CMDB Lite Code Structure

This document provides an overview of the CMDB Lite codebase organization, explaining the purpose of each directory and file, and how they relate to each other.

## Table of Contents

- [Project Structure Overview](#project-structure-overview)
- [Backend Structure](#backend-structure)
  - [cmd](#cmd)
  - [internal](#internal)
  - [migrations](#migrations)
  - [Configuration Files](#configuration-files)
- [Frontend Structure](#frontend-structure)
  - [public](#public)
  - [src](#src)
  - [Configuration Files](#configuration-files-1)
- [Database Structure](#database-structure)
  - [migrations](#migrations-1)
  - [schema](#schema)
  - [seeds](#seeds)
- [Deployment Structure](#deployment-structure)
  - [Docker](#docker)
  - [Kubernetes](#kubernetes)
  - [Scripts](#scripts)
- [Documentation Structure](#documentation-structure)
- [Key Architectural Patterns](#key-architectural-patterns)
  - [Backend Patterns](#backend-patterns)
  - [Frontend Patterns](#frontend-patterns)

## Project Structure Overview

```
cmdb-lite/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   ├── internal/           # Internal application packages
│   ├── migrations/         # Database migration files
│   ├── go.mod              # Go module definition
│   ├── go.sum              # Go module checksums
│   ├── .air.toml           # Air configuration for hot reload
│   └── Dockerfile          # Docker configuration
├── frontend/               # Vue 3 frontend application
│   ├── public/             # Static assets
│   ├── src/                # Source code
│   ├── package.json        # Node.js dependencies
│   ├── vite.config.js      # Vite configuration
│   ├── tailwind.config.js  # TailwindCSS configuration
│   ├── postcss.config.js   # PostCSS configuration
│   ├── .env.example        # Environment variables example
│   └── Dockerfile          # Docker configuration
├── database/               # Database schema and migrations
│   ├── migrations/         # Database migration files
│   ├── schema/             # Database schema files
│   ├── seeds/              # Database seed files
│   ├── .env.example        # Environment variables example
│   └── migrate.sh          # Migration script
├── k8s/                    # Kubernetes deployment configurations
│   └── helm/               # Helm chart for Kubernetes deployment
├── scripts/                # Deployment and utility scripts
├── docs/                   # Documentation
├── .env.example            # Environment variables example
├── docker-compose.yml      # Docker Compose configuration
├── docker-compose.dev.yml  # Docker Compose development configuration
├── .gitignore              # Git ignore rules
├── README.md               # Project README
├── PRD.md                  # Product Requirements Document
├── TSD.md                  # Technical Specifications Document
└── DEPLOYMENT.md           # Deployment Guide
```

## Backend Structure

The backend is a Go application that provides the REST API for CMDB Lite.

### cmd

The `cmd` directory contains the application entry points.

```
backend/cmd/
├── main.go         # Main application entry point
├── migrate.go      # Database migration tool
└── server/         # Server-related commands
    └── server.go   # Server implementation
```

- **main.go**: The main entry point for the application. It initializes the configuration, database, router, and starts the server.
- **migrate.go**: Command-line tool for running database migrations.
- **server/server.go**: Contains the server implementation, including HTTP server setup and graceful shutdown.

### internal

The `internal` directory contains all the application code that is private to the project. This follows the Go standard project layout.

```
backend/internal/
├── auth/           # Authentication logic
├── config/         # Configuration management
├── database/       # Database connection and setup
├── handlers/       # HTTP request handlers
├── middleware/     # HTTP middleware
├── models/         # Data models
├── repositories/   # Data access layer
└── router/         # Route definitions
```

#### auth

Contains authentication and authorization logic.

```
backend/internal/auth/
├── auth.go         # Authentication service
├── jwt.go          # JWT token handling
└── middleware.go   # Authentication middleware
```

- **auth.go**: Provides authentication services, including user authentication and token generation.
- **jwt.go**: Contains JWT token creation, validation, and parsing logic.
- **middleware.go**: Implements authentication middleware for protecting routes.

#### config

Contains configuration management logic.

```
backend/internal/config/
└── config.go       # Configuration loading and management
```

- **config.go**: Handles loading configuration from environment variables, configuration files, and default values.

#### database

Contains database connection and setup logic.

```
backend/internal/database/
├── database.go     # Database connection setup
├── migrations.go   # Database migration functions
└── models.go       # GORM model definitions
```

- **database.go**: Sets up the database connection using GORM.
- **migrations.go**: Provides functions for running database migrations.
- **models.go**: Contains GORM model definitions for all entities.

#### handlers

Contains HTTP request handlers for each endpoint.

```
backend/internal/handlers/
├── auth_handler.go     # Authentication handlers
├── ci_handler.go       # Configuration item handlers
├── relationship_handler.go  # Relationship handlers
├── audit_handler.go   # Audit log handlers
└── user_handler.go    # User management handlers
```

- **auth_handler.go**: Handles authentication endpoints (login, logout, refresh).
- **ci_handler.go**: Handles configuration item CRUD operations.
- **relationship_handler.go**: Handles relationship CRUD operations.
- **audit_handler.go**: Handles audit log retrieval.
- **user_handler.go**: Handles user management operations.

#### middleware

Contains HTTP middleware for cross-cutting concerns.

```
backend/internal/middleware/
├── auth.go          # Authentication middleware
├── cors.go          # CORS middleware
├── logging.go       # Request logging middleware
├── recovery.go      # Panic recovery middleware
└── rate_limit.go    # Rate limiting middleware
```

- **auth.go**: Middleware for authenticating requests using JWT tokens.
- **cors.go**: Middleware for handling Cross-Origin Resource Sharing (CORS).
- **logging.go**: Middleware for logging HTTP requests and responses.
- **recovery.go**: Middleware for recovering from panics and returning appropriate error responses.
- **rate_limit.go**: Middleware for rate limiting API requests.

#### models

Contains data models and business logic.

```
backend/internal/models/
├── user.go          # User model
├── ci.go            # Configuration item model
├── relationship.go  # Relationship model
├── audit.go         # Audit log model
└── base.go          # Base model with common fields
```

- **user.go**: Defines the User model and related business logic.
- **ci.go**: Defines the Configuration Item model and related business logic.
- **relationship.go**: Defines the Relationship model and related business logic.
- **audit.go**: Defines the Audit Log model and related business logic.
- **base.go**: Defines a base model with common fields (ID, created_at, updated_at).

#### repositories

Contains the data access layer for each entity.

```
backend/internal/repositories/
├── user_repository.go      # User data access
├── ci_repository.go        # Configuration item data access
├── relationship_repository.go  # Relationship data access
└── audit_repository.go     # Audit log data access
```

- **user_repository.go**: Provides data access methods for User entities.
- **ci_repository.go**: Provides data access methods for Configuration Item entities.
- **relationship_repository.go**: Provides data access methods for Relationship entities.
- **audit_repository.go**: Provides data access methods for Audit Log entities.

#### router

Contains route definitions and middleware setup.

```
backend/internal/router/
└── router.go       # Route definitions and middleware setup
```

- **router.go**: Defines all API routes and sets up middleware chains.

### migrations

Contains database migration files.

```
backend/migrations/
├── 20230101000000_create_users_table.up.sql
├── 20230101000000_create_users_table.down.sql
├── 20230101000001_create_cis_table.up.sql
├── 20230101000001_create_cis_table.down.sql
├── 20230101000002_create_relationships_table.up.sql
├── 20230101000002_create_relationships_table.down.sql
└── ...
```

Each migration consists of two files:
- **.up.sql**: Contains the SQL to apply the migration.
- **.down.sql**: Contains the SQL to revert the migration.

### Configuration Files

```
backend/
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── .air.toml           # Air configuration for hot reload
└── Dockerfile          # Docker configuration
```

- **go.mod**: Defines the Go module name and its dependencies.
- **go.sum**: Contains checksums for the dependencies to ensure integrity.
- **.air.toml**: Configuration for Air, a tool for live reloading Go applications.
- **Dockerfile**: Instructions for building a Docker image for the backend.

## Frontend Structure

The frontend is a Vue 3 application that provides the user interface for CMDB Lite.

### public

Contains static assets that will be served as-is.

```
frontend/public/
├── favicon.ico         # Website favicon
├── logo.png           # Application logo
└── index.html         # HTML template
```

- **favicon.ico**: The favicon for the application.
- **logo.png**: The application logo.
- **index.html**: The HTML template that will be used to serve the application.

### src

Contains the source code for the frontend application.

```
frontend/src/
├── assets/         # CSS and other assets
├── components/     # Vue components
├── router/         # Vue router configuration
├── services/       # API services
├── stores/         # Pinia stores
├── utils/          # Utility functions
├── views/          # Vue views
├── App.vue         # Root Vue component
└── main.js         # Application entry point
```

#### assets

Contains CSS and other static assets.

```
frontend/src/assets/
├── main.css        # Main CSS file
└── logo.svg        # SVG version of the logo
```

- **main.css**: Contains global CSS styles for the application.
- **logo.svg**: SVG version of the application logo.

#### components

Contains reusable Vue components.

```
frontend/src/components/
├── forms/          # Form components
│   ├── CIForm.vue          # Configuration item form
│   ├── RelationshipForm.vue  # Relationship form
│   └── UserForm.vue         # User form
├── graphs/         # Graph visualization components
│   ├── CIGraph.vue          # Configuration item graph
│   └── GraphControls.vue    # Graph control panel
├── layout/         # Layout components
│   ├── Header.vue          # Page header
│   ├── Sidebar.vue         # Navigation sidebar
│   └── Footer.vue          # Page footer
└── ui/             # UI components
    ├── Button.vue          # Button component
    ├── Card.vue            # Card component
    ├── Modal.vue           # Modal component
    ├── Notification.vue    # Notification component
    └── DataTable.vue       # Data table component
```

- **forms/**: Components for data input forms.
- **graphs/**: Components for graph visualization.
- **layout/**: Components for page layout.
- **ui/**: Reusable UI components.

#### router

Contains Vue router configuration.

```
frontend/src/router/
└── index.js       # Router configuration
```

- **index.js**: Defines all routes for the application, including route guards for authentication.

#### services

Contains API services for backend communication.

```
frontend/src/services/
├── api.js         # API client configuration
├── authService.js # Authentication API service
├── ciService.js   # Configuration item API service
├── relationshipService.js  # Relationship API service
├── auditService.js # Audit log API service
└── userService.js  # User API service
```

- **api.js**: Configures the Axios HTTP client with base URL, interceptors, and error handling.
- **authService.js**: Provides methods for authentication API calls.
- **ciService.js**: Provides methods for configuration item API calls.
- **relationshipService.js**: Provides methods for relationship API calls.
- **auditService.js**: Provides methods for audit log API calls.
- **userService.js**: Provides methods for user API calls.

#### stores

Contains Pinia stores for state management.

```
frontend/src/stores/
├── authStore.js   # Authentication state
├── ciStore.js     # Configuration item state
├── relationshipStore.js  # Relationship state
├── auditStore.js  # Audit log state
└── userStore.js   # User state
```

- **authStore.js**: Manages authentication state and related actions.
- **ciStore.js**: Manages configuration item state and related actions.
- **relationshipStore.js**: Manages relationship state and related actions.
- **auditStore.js**: Manages audit log state and related actions.
- **userStore.js**: Manages user state and related actions.

#### utils

Contains utility functions.

```
frontend/src/utils/
├── formatDate.js  # Date formatting utilities
├── validators.js  # Validation utilities
└── helpers.js     # General helper functions
```

- **formatDate.js**: Utilities for formatting dates.
- **validators.js**: Validation functions for form inputs.
- **helpers.js**: General helper functions used throughout the application.

#### views

Contains Vue views for each page.

```
frontend/src/views/
├── DashboardView.vue    # Dashboard page
├── CIListView.vue       # Configuration item list page
├── CIDetailView.vue     # Configuration item detail page
├── GraphView.vue        # Graph visualization page
├── AuditLogView.vue     # Audit log page
├── UserListView.vue     # User list page
├── UserDetailView.vue   # User detail page
├── LoginView.vue        # Login page
├── SettingsView.vue     # Settings page
└── NotFoundView.vue     # 404 page
```

- **DashboardView.vue**: The dashboard page showing an overview of CMDB data.
- **CIListView.vue**: The page for listing and searching configuration items.
- **CIDetailView.vue**: The page for viewing and editing a single configuration item.
- **GraphView.vue**: The page for visualizing relationships between configuration items.
- **AuditLogView.vue**: The page for viewing audit logs.
- **UserListView.vue**: The page for listing and managing users.
- **UserDetailView.vue**: The page for viewing and editing a single user.
- **LoginView.vue**: The login page.
- **SettingsView.vue**: The settings page.
- **NotFoundView.vue**: The 404 page.

#### App.vue and main.js

```
frontend/src/
├── App.vue         # Root Vue component
└── main.js         # Application entry point
```

- **App.vue**: The root Vue component that contains the application layout.
- **main.js**: The entry point for the application, which initializes Vue, plugins, and mounts the app.

### Configuration Files

```
frontend/
├── package.json        # Node.js dependencies
├── vite.config.js      # Vite configuration
├── tailwind.config.js  # TailwindCSS configuration
├── postcss.config.js   # PostCSS configuration
├── .env.example        # Environment variables example
└── Dockerfile          # Docker configuration
```

- **package.json**: Defines Node.js dependencies and scripts.
- **vite.config.js**: Configuration for Vite, the build tool and development server.
- **tailwind.config.js**: Configuration for TailwindCSS, the utility-first CSS framework.
- **postcss.config.js**: Configuration for PostCSS, a tool for transforming CSS.
- **.env.example**: Example environment variables for the frontend.
- **Dockerfile**: Instructions for building a Docker image for the frontend.

## Database Structure

The database directory contains schema, migrations, and seed files for the PostgreSQL database.

### migrations

Contains database migration files.

```
database/migrations/
├── 20230101000000_create_users_table.up.sql
├── 20230101000000_create_users_table.down.sql
├── 20230101000001_create_cis_table.up.sql
├── 20230101000001_create_cis_table.down.sql
├── 20230101000002_create_relationships_table.up.sql
├── 20230101000002_create_relationships_table.down.sql
└── ...
```

These are the same migration files as in `backend/migrations/`, but stored in a central location for database management.

### schema

Contains database schema files.

```
database/schema/
├── schema.sql      # Complete database schema
└── erd.png         # Entity-relationship diagram
```

- **schema.sql**: Contains the complete database schema, useful for setting up a new database.
- **erd.png**: An entity-relationship diagram showing the relationships between tables.

### seeds

Contains database seed files.

```
database/seeds/
├── 20230101000000_seed_users.sql
├── 20230101000001_seed_ci_types.sql
└── ...
```

- **seed_users.sql**: SQL to seed the users table with initial data.
- **seed_ci_types.sql**: SQL to seed the CI types table with initial data.

## Deployment Structure

The deployment directory contains configurations for deploying CMDB Lite using Docker and Kubernetes.

### Docker

```
├── backend/Dockerfile          # Backend Docker configuration
├── frontend/Dockerfile         # Frontend Docker configuration
├── docker-compose.yml         # Docker Compose configuration
└── docker-compose.dev.yml     # Docker Compose development configuration
```

- **backend/Dockerfile**: Instructions for building a Docker image for the backend.
- **frontend/Dockerfile**: Instructions for building a Docker image for the frontend.
- **docker-compose.yml**: Docker Compose configuration for production deployment.
- **docker-compose.dev.yml**: Docker Compose configuration for development deployment.

### Kubernetes

```
k8s/
└── helm/               # Helm chart for Kubernetes deployment
    ├── Chart.yaml     # Helm chart metadata
    ├── values.yaml    # Default Helm values
    ├── templates/     # Kubernetes resource templates
    │   ├── backend-deployment.yaml
    │   ├── backend-service.yaml
    │   ├── frontend-deployment.yaml
    │   ├── frontend-service.yaml
    │   ├── ingress.yaml
    │   └── secrets.yaml
    └── README.md      # Helm chart documentation
```

- **Chart.yaml**: Metadata for the Helm chart.
- **values.yaml**: Default values for the Helm chart.
- **templates/**: Kubernetes resource templates.
- **README.md**: Documentation for the Helm chart.

### Scripts

```
scripts/
├── deploy.sh          # Main deployment script
├── db-migrate.sh      # Database migration script
├── db-backup.sh       # Database backup script
└── init.sh            # Initialization script
```

- **deploy.sh**: Script for deploying the application.
- **db-migrate.sh**: Script for running database migrations.
- **db-backup.sh**: Script for backing up the database.
- **init.sh**: Script for initializing the environment.

## Documentation Structure

```
docs/
├── README.md          # Documentation index
├── user/              # User documentation
│   ├── README.md      # User guide
│   ├── features.md    # Features overview
│   └── faq.md         # Frequently asked questions
├── developer/         # Developer documentation
│   ├── README.md      # Developer guide
│   ├── api.md         # API documentation
│   ├── testing.md     # Testing guide
│   └── code-structure.md  # Code structure
├── operator/          # Operator documentation
│   ├── README.md      # Operations guide
│   ├── deployment.md  # Deployment options
│   ├── monitoring.md  # Monitoring and troubleshooting
│   └── backup-recovery.md  # Backup and recovery
└── project/           # Project documentation
    ├── README.md      # Project overview
    ├── contributing.md  # Contribution guide
    └── license.md     # License information
```

- **README.md**: The main documentation index.
- **user/**: Documentation for end users.
- **developer/**: Documentation for developers.
- **operator/**: Documentation for system operators.
- **project/**: Documentation for project maintainers.

## Key Architectural Patterns

### Backend Patterns

The backend follows several architectural patterns to ensure maintainability, testability, and scalability:

#### Clean Architecture

The backend follows the Clean Architecture pattern, which emphasizes separation of concerns and dependency inversion:

- **Entities**: Core business objects (models)
- **Use Cases**: Business logic (services)
- **Interface Adapters**: Convert data between external formats and internal formats (handlers, repositories)
- **Frameworks & Drivers**: External tools and frameworks (database, web framework)

#### Repository Pattern

The Repository Pattern is used to abstract the data layer:

- **Repositories**: Define interfaces for data access operations
- **Implementations**: Provide concrete implementations for different data sources
- **Services**: Use repository interfaces for data access, making them independent of the data source

#### Dependency Injection

Dependency Injection is used to manage dependencies between components:

- **Constructor Injection**: Dependencies are provided through constructors
- **Interface Segregation**: Components depend on abstractions, not concrete implementations
- **Inversion of Control**: The framework manages the lifecycle of dependencies

#### Middleware Pattern

The Middleware Pattern is used for cross-cutting concerns:

- **Authentication**: Verifies user identity
- **Authorization**: Checks user permissions
- **Logging**: Logs requests and responses
- **Rate Limiting**: Limits request rates
- **CORS**: Handles cross-origin requests

### Frontend Patterns

The frontend follows several architectural patterns to ensure maintainability, testability, and scalability:

#### Component-Based Architecture

The frontend is built using a component-based architecture:

- **Reusable Components**: Small, reusable components for UI elements
- **Container Components**: Components that manage state and pass data to presentational components
- **Presentational Components**: Components that focus on rendering UI based on props

#### State Management

Pinia is used for state management:

- **Stores**: Manage state for different domains (auth, CIs, relationships, etc.)
- **Actions**: Define methods for modifying state
- **Getters**: Define computed properties based on state
- **Modules**: Organize state into logical modules

#### Service Layer

A service layer is used for API communication:

- **API Services**: Encapsulate API calls and error handling
- **Interceptors**: Handle request/response transformations and error handling
- **Models**: Define data structures for API responses

#### Router Guards

Router guards are used for navigation control:

- **Authentication Guards**: Protect routes that require authentication
- **Authorization Guards**: Protect routes based on user roles
- **Data Resolvers**: Pre-fetch data before navigating to a route

#### Composition API

The Composition API is used for organizing component logic:

- **Reactive State**: Use `ref` and `reactive` for reactive state
- **Computed Properties**: Use `computed` for derived state
- **Lifecycle Hooks**: Use lifecycle hooks for component lifecycle events
- **Composables**: Reuse logic across components using composables

These architectural patterns ensure that the codebase is maintainable, testable, and scalable, making it easier to add new features and fix bugs.