# CMDB Lite Developer Guide

Welcome to the CMDB Lite Developer Guide! This guide provides information for developers who want to contribute to the CMDB Lite project or integrate with its API.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Development Environment Setup](#development-environment-setup)
- [Code Organization](#code-organization)
- [Backend Development](#backend-development)
- [Frontend Development](#frontend-development)
- [Database Development](#database-development)
- [Testing](#testing)
- [Contributing Guidelines](#contributing-guidelines)
- [Getting Help](#getting-help)

## Architecture Overview

CMDB Lite follows a three-tier architecture with clear separation of concerns:

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Frontend (Vue 3)                           │
├─────────────────────────────────────────────────────────────┤
│                     Backend (Go)                             │
├─────────────────────────────────────────────────────────────┤
│                    Database (PostgreSQL)                     │
└─────────────────────────────────────────────────────────────┘
```

### Frontend Architecture

The frontend is built with Vue 3 and follows a component-based architecture:

- **Vue 3**: Progressive JavaScript framework for building user interfaces
- **Pinia**: State management library for Vue
- **Vue Router**: Official routing library for Vue
- **TailwindCSS**: Utility-first CSS framework for styling
- **Axios**: HTTP client for API communication
- **D3.js**: Data visualization library for graph rendering

### Backend Architecture

The backend is built with Go and follows a clean architecture pattern:

- **Go**: Programming language known for performance and simplicity
- **Gin**: HTTP web framework for building RESTful APIs
- **GORM**: ORM library for database operations
- **JWT**: Authentication using JSON Web Tokens
- **PostgreSQL**: Relational database for data persistence

### Database Schema

The database uses a relational model with the following main entities:

- **Users**: Authentication and authorization
- **Configuration Items**: Core entities managed by the CMDB
- **Relationships**: Connections between configuration items
- **Audit Logs**: Track all changes for compliance and troubleshooting

For a detailed ERD, see the [Technical Specification Document](../../TSD.md).

## Development Environment Setup

### Prerequisites

Before setting up your development environment, ensure you have the following installed:

- **Go** (version 1.18 or higher)
- **Node.js** (version 16 or higher)
- **PostgreSQL** (version 12 or higher)
- **Git** (version 2.20 or higher)
- **Docker** and **Docker Compose** (version 20.10 or higher)

### Cloning the Repository

```bash
git clone https://github.com/yourorg/cmdb-lite.git
cd cmdb-lite
```

### Setting Up the Backend

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Edit the `.env` file with your local configuration:
   ```env
   # Environment
   ENVIRONMENT=development

   # Database
   DB_HOST=localhost
   DB_PORT=5432
   DB_NAME=cmdb_lite
   DB_USER=cmdb_user
   DB_PASSWORD=cmdb_password

   # Server
   SERVER_PORT=8080
   JWT_SECRET=your-secret-key-here

   # Logging
   LOG_LEVEL=debug
   ```

4. Install Go dependencies:
   ```bash
   go mod download
   ```

5. Run database migrations:
   ```bash
   go run cmd/migrate.go
   ```

6. Start the backend server:
   ```bash
   go run cmd/main.go
   ```

   The backend API will be available at http://localhost:8080.

### Setting Up the Frontend

1. Navigate to the frontend directory:
   ```bash
   cd ../frontend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Edit the `.env` file with your local configuration:
   ```env
   VITE_API_URL=http://localhost:8080
   VITE_ENVIRONMENT=development
   ```

4. Install Node.js dependencies:
   ```bash
   npm install
   ```

5. Start the frontend development server:
   ```bash
   npm run dev
   ```

   The frontend will be available at http://localhost:3000.

### Setting Up the Database

1. Ensure PostgreSQL is running on your system.

2. Create a database user and database:
   ```sql
   CREATE USER cmdb_user WITH PASSWORD 'cmdb_password';
   CREATE DATABASE cmdb_lite OWNER cmdb_user;
   ```

3. Run the database migrations (already done in the backend setup step).

### Using Docker Compose for Development

For a simpler setup, you can use Docker Compose to run all services:

1. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

2. Start the development environment:
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

3. Access the services:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Database Admin: http://localhost:8081

## Code Organization

### Project Structure

```
cmdb-lite/
├── backend/                 # Go backend application
│   ├── cmd/                # Application entry points
│   │   ├── main.go         # Main application entry point
│   │   └── migrate.go      # Database migration tool
│   ├── internal/           # Internal application packages
│   │   ├── auth/           # Authentication logic
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection and setup
│   │   ├── handlers/       # HTTP request handlers
│   │   ├── middleware/     # HTTP middleware
│   │   ├── models/         # Data models
│   │   ├── repositories/   # Data access layer
│   │   └── router/         # Route definitions
│   ├── migrations/         # Database migration files
│   ├── go.mod              # Go module definition
│   └── Dockerfile          # Docker configuration
├── frontend/               # Vue 3 frontend application
│   ├── public/             # Static assets
│   ├── src/                # Source code
│   │   ├── assets/         # CSS and other assets
│   │   ├── components/     # Vue components
│   │   ├── router/         # Vue router configuration
│   │   ├── services/       # API services
│   │   ├── stores/         # Pinia stores
│   │   ├── utils/          # Utility functions
│   │   ├── views/          # Vue views
│   │   ├── App.vue         # Root Vue component
│   │   └── main.js         # Application entry point
│   ├── package.json        # Node.js dependencies
│   └── vite.config.js      # Vite configuration
├── database/               # Database schema and migrations
│   ├── migrations/         # Database migration files
│   ├── schema/             # Database schema files
│   └── seeds/              # Database seed files
└── docs/                   # Documentation
```

### Backend Code Organization

The backend follows a clean architecture pattern with clear separation of concerns:

- **cmd**: Contains the application entry points.
- **internal**: Contains all the application code that is private to the project.
  - **auth**: Handles authentication and authorization logic.
  - **config**: Manages application configuration.
  - **database**: Handles database connection and setup.
  - **handlers**: Contains HTTP request handlers for each endpoint.
  - **middleware**: Contains HTTP middleware for authentication, logging, etc.
  - **models**: Contains data models and business logic.
  - **repositories**: Contains data access layer for each entity.
  - **router**: Contains route definitions and middleware setup.

### Frontend Code Organization

The frontend follows a component-based architecture:

- **src/assets**: Contains CSS and other static assets.
- **src/components**: Contains reusable Vue components.
  - **forms**: Contains form components for data input.
  - **graphs**: Contains graph visualization components.
  - **layout**: Contains layout components like header and sidebar.
  - **ui**: Contains UI components like modals and notifications.
- **src/router**: Contains Vue router configuration.
- **src/services**: Contains API services for backend communication.
- **src/stores**: Contains Pinia stores for state management.
- **src/utils**: Contains utility functions.
- **src/views**: Contains Vue views for each page.

## Backend Development

### Adding a New API Endpoint

To add a new API endpoint:

1. Define the route in `internal/router/router.go`:
   ```go
   // Example for a new endpoint
   apiGroup.GET("/new-endpoint", handlers.NewEndpointHandler)
   ```

2. Create the handler function in `internal/handlers/new_handler.go`:
   ```go
   func NewEndpointHandler(c *gin.Context) {
       // Handler logic here
   }
   ```

3. Add any necessary models in `internal/models/`:
   ```go
   type NewModel struct {
       ID   uuid.UUID `json:"id"`
       Name string    `json:"name"`
       // Other fields
   }
   ```

4. Add any necessary repository methods in `internal/repositories/`:
   ```go
   func (r *NewModelRepository) Create(model *models.NewModel) error {
       // Repository logic here
   }
   ```

5. Write tests for your new endpoint in `internal/handlers/new_handler_test.go`.

### Working with the Database

CMDB Lite uses GORM as the ORM for database operations. Here are some common operations:

#### Creating a Record

```go
user := models.User{
    Username: "testuser",
    Email:    "test@example.com",
}

result := database.DB.Create(&user)
if result.Error != nil {
    // Handle error
}
```

#### Querying Records

```go
// Get a single record
var user models.User
result := database.DB.First(&user, "id = ?", userID)
if result.Error != nil {
    // Handle error
}

// Get multiple records with conditions
var users []models.User
result := database.DB.Where("active = ?", true).Find(&users)
if result.Error != nil {
    // Handle error
}
```

#### Updating Records

```go
var user models.User
database.DB.First(&user, userID)

user.Email = "new-email@example.com"
result := database.DB.Save(&user)
if result.Error != nil {
    // Handle error
}
```

#### Deleting Records

```go
var user models.User
database.DB.First(&user, userID)

result := database.DB.Delete(&user)
if result.Error != nil {
    // Handle error
}
```

### Authentication and Authorization

CMDB Lite uses JWT for authentication. To protect an endpoint:

1. Add the authentication middleware to the route:
   ```go
   apiGroup.GET("/protected-endpoint", middleware.Authenticate(), handlers.ProtectedEndpointHandler)
   ```

2. Access the user information from the context:
   ```go
   func ProtectedEndpointHandler(c *gin.Context) {
       user, exists := c.Get("user")
       if !exists {
           c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
           return
       }
       
       // Use the user information
       userModel := user.(*models.User)
       // Handler logic here
   }
   ```

### Error Handling

CMDB Lite uses a standardized error response format:

```go
c.JSON(http.StatusBadRequest, gin.H{
    "error": "Error message",
    "details": "Additional error details",
})
```

For validation errors, use the HTTP status code 422 Unprocessable Entity:

```go
c.JSON(http.StatusUnprocessableEntity, gin.H{
    "error": "Validation failed",
    "details": validationErrors,
})
```

## Frontend Development

### Adding a New Vue Component

To add a new Vue component:

1. Create a new file in the appropriate directory under `src/components/`:
   ```vue
   <template>
     <div class="new-component">
       <!-- Component template here -->
     </div>
   </template>

   <script>
   export default {
     name: 'NewComponent',
     props: {
       // Define props here
     },
     data() {
       return {
         // Define component data here
       }
     },
     methods: {
       // Define component methods here
     }
   }
   </script>

   <style scoped>
   /* Component styles here */
   </style>
   ```

2. Import and use the component in other components or views:
   ```vue
   <script>
   import NewComponent from '@/components/NewComponent.vue'

   export default {
     components: {
       NewComponent
     }
   }
   </script>

   <template>
     <div>
       <NewComponent />
     </div>
   </template>
   ```

### Adding a New View

To add a new view:

1. Create a new file in `src/views/`:
   ```vue
   <template>
     <div class="new-view">
       <!-- View template here -->
     </div>
   </template>

   <script>
   export default {
     name: 'NewView',
     data() {
       return {
         // Define view data here
       }
     },
     methods: {
       // Define view methods here
     },
     created() {
       // Code to run when the view is created
     }
   }
   </script>

   <style scoped>
   /* View styles here */
   </style>
   ```

2. Add the route to `src/router/index.js`:
   ```javascript
   import NewView from '@/views/NewView.vue'

   const routes = [
     // Other routes
     {
       path: '/new-view',
       name: 'NewView',
       component: NewView,
       meta: { requiresAuth: true }
     }
   ]
   ```

### Working with Pinia Stores

CMDB Lite uses Pinia for state management. To use a store:

1. Access the store in a component:
   ```javascript
   import { useAuthStore } from '@/stores/auth'

   export default {
     setup() {
       const authStore = useAuthStore()
       
       return {
         authStore
       }
     }
   }
   ```

2. Use store actions and getters:
   ```javascript
   // Call an action
   this.authStore.login(username, password)

   // Access a getter
   const isAuthenticated = this.authStore.isAuthenticated
   ```

### Making API Calls

CMDB Lite uses Axios for API calls. To make an API call:

1. Use the API service:
   ```javascript
   import { ciService } from '@/services/ciService'

   export default {
     methods: {
       async getCIs() {
         try {
           const response = await ciService.getAll()
           this.cis = response.data
         } catch (error) {
           console.error('Error fetching CIs:', error)
         }
       }
     }
   }
   ```

2. Create a new API service if needed:
   ```javascript
   import api from '@/services/api'

   export const newService = {
     getAll() {
       return api.get('/new-endpoint')
     },
     
     create(data) {
       return api.post('/new-endpoint', data)
     },
     
     update(id, data) {
       return api.put(`/new-endpoint/${id}`, data)
     },
     
     delete(id) {
       return api.delete(`/new-endpoint/${id}`)
     }
   }
   ```

## Database Development

### Creating a New Migration

To create a new database migration:

1. Create a new migration file in `backend/migrations/` with the format `YYYYMMDDHHMMSS_description.up.sql` and `YYYYMMDDHHMMSS_description.down.sql`.

2. Write the SQL for the migration in the `.up.sql` file:
   ```sql
   -- Example: Add a new table
   CREATE TABLE new_table (
     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
     name VARCHAR(255) NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

3. Write the SQL to revert the migration in the `.down.sql` file:
   ```sql
   -- Example: Drop the new table
   DROP TABLE new_table;
   ```

4. Run the migration:
   ```bash
   go run cmd/migrate.go up
   ```

### Working with GORM Models

CMDB Lite uses GORM models to represent database entities. Here's an example model:

```go
package models

import (
    "time"
    
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ExampleModel struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
    Name      string    `gorm:"type:varchar(255);not null" json:"name"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ExampleModel) TableName() string {
    return "example_models"
}
```

### Database Seeding

To seed the database with initial data:

1. Create a seed file in `database/seeds/` with the format `YYYYMMDDHHMMSS_description.sql`.

2. Write the SQL to insert the seed data:
   ```sql
   -- Example: Insert seed data
   INSERT INTO users (id, username, email, role, created_at, updated_at) VALUES
     (gen_random_uuid(), 'admin', 'admin@example.com', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
   ```

3. Run the seed script:
   ```bash
   psql -U cmdb_user -d cmdb_lite -f database/seeds/YYYYMMDDHHMMSS_description.sql
   ```

## Testing

### Backend Testing

CMDB Lite uses Go's built-in testing framework for backend tests. To run the backend tests:

```bash
cd backend
go test ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

To run tests for a specific package:

```bash
go test ./internal/handlers
```

### Frontend Testing

CMDB Lite uses Vitest for frontend testing. To run the frontend tests:

```bash
cd frontend
npm test
```

To run tests with coverage:

```bash
npm run test:coverage
```

### End-to-End Testing

CMDB Lite uses Playwright for end-to-end testing. To run the E2E tests:

```bash
npm run test:e2e
```

### Writing Tests

#### Backend Test Example

```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestExampleHandler(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.GET("/example", ExampleHandler)
    
    // Create a test request
    req, _ := http.NewRequest("GET", "/example", nil)
    w := httptest.NewRecorder()
    
    // Perform the request
    router.ServeHTTP(w, req)
    
    // Assert the response
    assert.Equal(t, http.StatusOK, w.Code)
    assert.Contains(t, w.Body.String(), "expected response")
}
```

#### Frontend Test Example

```javascript
import { render, screen } from '@testing-library/vue'
import ExampleComponent from '@/components/ExampleComponent.vue'

describe('ExampleComponent', () => {
  it('renders properly', () => {
    render(ExampleComponent, {
      props: {
        title: 'Test Title'
      }
    })
    
    expect(screen.getByText('Test Title')).toBeInTheDocument()
  })
})
```

## Contributing Guidelines

We welcome contributions to CMDB Lite! Please follow these guidelines when contributing:

### Code Style

- **Backend**: Follow the Go standard code style as defined by `gofmt`.
- **Frontend**: Follow the Vue Style Guide and use ESLint with the provided configuration.

### Pull Request Process

1. Fork the repository.
2. Create a feature branch from `main`.
3. Make your changes.
4. Write tests for your changes.
5. Ensure all tests pass.
6. Update documentation if needed.
7. Submit a pull request with a clear description of your changes.

### Commit Message Format

We use a conventional commit format for commit messages:

```
type(scope): description

[optional body]

[optional footer(s)]
```

Types include:
- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `test`: Adding missing tests or correcting existing tests
- `chore`: Changes to the build process or auxiliary tools

Example:
```
feat(auth): add JWT token refresh functionality

Add the ability to refresh JWT tokens without requiring users to log in again.
This improves the user experience by reducing the frequency of login prompts.

Closes #123
```

## Getting Help

If you need help with development:

- Check the [API Documentation](api.md) for detailed API reference.
- Review existing code for examples of how to implement features.
- Check the [GitHub Issues](https://github.com/yourorg/cmdb-lite/issues) for similar problems.
- Create a new issue with your question if it hasn't been addressed before.

For more information on contributing, see the [Contribution Guide](../project/contributing.md).