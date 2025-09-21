# CMDB Lite Unit Testing Guide

This guide provides comprehensive instructions for writing unit tests for the CMDB Lite application. The project uses Go for the backend and follows established testing patterns.

## Table of Contents
- [Testing Framework and Dependencies](#testing-framework-and-dependencies)
- [Test Structure and Organization](#test-structure-and-organization)
- [Running Tests](#running-tests)
- [Writing Unit Tests](#writing-unit-tests)
- [Mocking Dependencies](#mocking-dependencies)
- [Testing Handlers](#testing-handlers)
- [Testing Repositories](#testing-repositories)
- [Testing Services](#testing-services)
- [Testing Middleware](#testing-middleware)
- [Integration Testing](#integration-testing)
- [Best Practices](#best-practices)
- [Test Examples](#test-examples)

## Testing Framework and Dependencies

### Core Testing Stack
- **Go Testing Package**: Built-in `testing` package
- **Testify**: Assertion library (`github.com/stretchr/testify`)
- **SQLMock**: Database mocking (`github.com/DATA-DOG/go-sqlmock`)
- **TestContainers**: Integration testing with real containers (`github.com/testcontainers/testcontainers-go`)

### Key Dependencies
```go
// Testing and assertions
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)

// HTTP testing
import (
    "net/http"
    "net/http/httptest"
    "github.com/gorilla/mux"
)

// Database mocking
import (
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/testcontainers/testcontainers-go/modules/postgres"
)
```

## Test Structure and Organization

### Directory Structure
```
backend/
├── internal/
│   ├── handlers/
│   │   ├── ci_handler.go
│   │   └── ci_handler_test.go
│   ├── repositories/
│   │   ├── ci_repository.go
│   │   └── ci_repository_test.go
│   ├── services/
│   │   ├── auth_service.go
│   │   └── auth_service_test.go
│   └── middleware/
│       ├── auth.go
│       └── auth_test.go
├── internal/testutils/
│   ├── testutils.go
│   └── mocks.go
└── test/
    ├── integration/
    │   └── ci_integration_test.go
    └── testdata/
        └── fixtures.json
```

### Test File Naming Convention
- Unit tests: `{filename}_test.go`
- Integration tests: `{filename}_integration_test.go`
- Test utilities: `testutils.go`, `mocks.go`

## Running Tests

### Basic Test Commands
```bash
# Run all unit tests
make test-unit

# Run all tests with coverage
make test-coverage

# Run specific test file
go test -v ./internal/handlers/ci_handler_test.go

# Run tests with verbose output
go test -v ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Test Categories
```bash
# Run only unit tests (short tests)
go test -v -short ./...

# Run integration tests
go test -v -run=Integration ./...

# Run tests with specific pattern
go test -v -run="TestCIHandler_Create" ./...
```

## Writing Unit Tests

### Basic Test Structure
```go
package handlers

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCIHandler_CreateCI(t *testing.T) {
    // Arrange: Setup test data and mocks
    testCI := &models.CI{
        Name: "Test CI",
        Type: "server",
    }
    
    // Act: Call the function being tested
    result, err := handler.CreateCI(testCI)
    
    // Assert: Verify the results
    require.NoError(t, err)
    assert.Equal(t, testCI.Name, result.Name)
    assert.Equal(t, testCI.Type, result.Type)
}
```

### Test Table Pattern
```go
func TestCIHandler_CreateCI(t *testing.T) {
    tests := []struct {
        name           string
        input          *models.CI
        expectedErr    error
        expectedStatus int
    }{
        {
            name: "Valid CI",
            input: &models.CI{Name: "Test CI", Type: "server"},
            expectedErr: nil,
            expectedStatus: http.StatusCreated,
        },
        {
            name: "Missing Name",
            input: &models.CI{Name: "", Type: "server"},
            expectedErr: errors.New("name is required"),
            expectedStatus: http.StatusBadRequest,
        },
        {
            name: "Missing Type",
            input: &models.CI{Name: "Test CI", Type: ""},
            expectedErr: errors.New("type is required"),
            expectedStatus: http.StatusBadRequest,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            handler := NewCIHandler(mockRepo, mockRelRepo, mockAuditRepo)
            
            // Act
            result, err := handler.CreateCI(tt.input)
            
            // Assert
            if tt.expectedErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.expectedErr.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.input.Name, result.Name)
                assert.Equal(t, tt.input.Type, result.Type)
            }
        })
    }
}
```

## Mocking Dependencies

### Creating Mock Repositories
```go
// MockCIRepository implements the CIRepository interface
type MockCIRepository struct {
    mock.Mock
}

func (m *MockCIRepository) Create(ctx context.Context, ci *models.CI) error {
    args := m.Called(ctx, ci)
    return args.Error(0)
}

func (m *MockCIRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CI, error) {
    args := m.Called(ctx, id)
    return args.Get(0).(*models.CI), args.Error(1)
}

func (m *MockCIRepository) GetAll(ctx context.Context) ([]*models.CI, error) {
    args := m.Called(ctx)
    return args.Get(0).([]*models.CI), args.Error(1)
}

// Setup mock expectations
func setupMockRepository() *MockCIRepository {
    mock := &MockCIRepository{}
    
    // Setup expected calls and return values
    mock.On("Create", mock.Anything, mock.AnythingOfType("*models.CI")).
        Return(nil)
    
    mock.On("GetByID", mock.Anything, mock.AnythingOfType("uuid.UUID")).
        Return(&models.CI{ID: uuid.New(), Name: "Test CI"}, nil)
    
    mock.On("GetAll", mock.Anything).
        Return([]*models.CI{{ID: uuid.New(), Name: "Test CI"}}, nil)
    
    return mock
}
```

### Using SQLMock for Database Testing
```go
func TestCIRepository_Create(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    sqlxDB := sqlx.NewDb(db, "sqlmock")
    repo := NewCIPostgresRepository(sqlxDB)

    testCI := &models.CI{
        ID:        uuid.New(),
        Name:      "Test CI",
        Type:      "server",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Expect the INSERT query
    mock.ExpectExec(`INSERT INTO configuration_items`).
        WithArgs(testCI.ID, testCI.Name, testCI.Type, mock.Anything(), mock.Anything(), mock.Anything(), mock.Anything()).
        WillReturnResult(sqlmock.NewResult(1, 1))

    // Call the method
    err = repo.Create(context.Background(), testCI)

    // Assert
    assert.NoError(t, err)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

## Testing Handlers

### HTTP Handler Testing Pattern
```go
func TestCIHandler_CreateCI_HTTP(t *testing.T) {
    // Setup mocks
    mockCIRepo := &MockCIRepository{}
    mockRelRepo := &MockRelationshipRepository{}
    mockAuditRepo := &MockAuditLogRepository{}
    
    handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

    // Create test request
    testCI := models.CI{
        Name: "Test CI",
        Type: "server",
    }
    
    requestBody, err := json.Marshal(testCI)
    require.NoError(t, err)

    req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(requestBody))
    require.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")

    // Add authentication context
    ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
    req = req.WithContext(ctx)

    // Create response recorder
    rr := httptest.NewRecorder()

    // Call handler
    handler.CreateCI(rr, req)

    // Assert response
    assert.Equal(t, http.StatusCreated, rr.Code)
    
    var response models.CI
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.Equal(t, testCI.Name, response.Name)
    assert.Equal(t, testCI.Type, response.Type)
    assert.NotEmpty(t, response.ID)
}
```

### Testing with URL Parameters
```go
func TestCIHandler_GetCI_HTTP(t *testing.T) {
    // Setup
    testCI := &models.CI{
        ID:   uuid.New(),
        Name: "Test CI",
        Type: "server",
    }
    
    mockCIRepo := &MockCIRepository{}
    mockCIRepo.On("GetByID", mock.Anything, testCI.ID).Return(testCI, nil)
    
    handler := NewCIHandler(mockCIRepo, &MockRelationshipRepository{}, &MockAuditLogRepository{})

    // Create request with URL parameters
    req, err := http.NewRequest("GET", "/api/v1/cis/"+testCI.ID.String(), nil)
    require.NoError(t, err)

    // Setup router with parameters
    router := mux.NewRouter()
    router.HandleFunc("/api/v1/cis/{id}", handler.GetCI).Methods("GET")

    // Create response recorder
    rr := httptest.NewRecorder()

    // Serve request
    router.ServeHTTP(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    
    var response models.CI
    err = json.Unmarshal(rr.Body.Bytes(), &response)
    require.NoError(t, err)
    
    assert.Equal(t, testCI.ID, response.ID)
    assert.Equal(t, testCI.Name, response.Name)
}
```

## Testing Repositories

### Repository Test with SQLMock
```go
func TestCIPostgresRepository_GetByID(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    sqlxDB := sqlx.NewDb(db, "sqlmock")
    repo := NewCIPostgresRepository(sqlxDB)

    testID := uuid.New()
    expectedCI := &models.CI{
        ID:        testID,
        Name:      "Test CI",
        Type:      "server",
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Expect the SELECT query
    rows := sqlmock.NewRows([]string{"id", "name", "type", "attributes", "tags", "created_at", "updated_at"}).
        AddRow(expectedCI.ID, expectedCI.Name, expectedCI.Type, []byte{}, []byte{}, expectedCI.CreatedAt, expectedCI.UpdatedAt)

    mock.ExpectQuery(`SELECT (.+) FROM configuration_items WHERE id = \$1`).
        WithArgs(testID).
        WillReturnRows(rows)

    // Call the method
    result, err := repo.GetByID(context.Background(), testID)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedCI.ID, result.ID)
    assert.Equal(t, expectedCI.Name, result.Name)
    assert.Equal(t, expectedCI.Type, result.Type)
    assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCIPostgresRepository_GetByID_NotFound(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    sqlxDB := sqlx.NewDb(db, "sqlmock")
    repo := NewCIPostgresRepository(sqlxDB)

    testID := uuid.New()

    // Expect the SELECT query to return no rows
    mock.ExpectQuery(`SELECT (.+) FROM configuration_items WHERE id = \$1`).
        WithArgs(testID).
        WillReturnError(sql.ErrNoRows)

    // Call the method
    result, err := repo.GetByID(context.Background(), testID)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

## Testing Services

### Service Layer Testing
```go
func TestAuthService_Login(t *testing.T) {
    tests := []struct {
        name        string
        username    string
        password    string
        setupMock   func(*MockUserRepository, *MockPasswordManager)
        expectError bool
        errorType   models.ErrorType
    }{
        {
            name:     "Successful login",
            username: "testuser",
            password: "validpassword",
            setupMock: func(mockRepo *MockUserRepository, mockPM *MockPasswordManager) {
                user := &models.User{
                    ID:           uuid.New(),
                    Username:     "testuser",
                    PasswordHash: "hashedpassword",
                    Role:         "user",
                }
                mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)
                mockPM.On("CheckPassword", "validpassword", "hashedpassword").Return(true)
            },
            expectError: false,
        },
        {
            name:     "Invalid password",
            username: "testuser",
            password: "wrongpassword",
            setupMock: func(mockRepo *MockUserRepository, *MockPasswordManager) {
                user := &models.User{
                    ID:           uuid.New(),
                    Username:     "testuser",
                    PasswordHash: "hashedpassword",
                    Role:         "user",
                }
                mockRepo.On("GetByUsername", mock.Anything, "testuser").Return(user, nil)
                mockPM.On("CheckPassword", "wrongpassword", "hashedpassword").Return(false)
            },
            expectError: true,
            errorType:   models.ErrorTypeUnauthorized,
        },
        {
            name:     "User not found",
            username: "nonexistent",
            password: "anypassword",
            setupMock: func(mockRepo *MockUserRepository, *MockPasswordManager) {
                mockRepo.On("GetByUsername", mock.Anything, "nonexistent").
                    Return(nil, sqlx.ErrNotFound)
            },
            expectError: true,
            errorType:   models.ErrorTypeUnauthorized,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockRepo := &MockUserRepository{}
            mockPM := &MockPasswordManager{}
            jwtManager := auth.NewJWTManager("test-secret", time.Hour, time.Hour*24)
            
            tt.setupMock(mockRepo, mockPM)
            
            service := NewAuthService(mockRepo, mockPM, jwtManager)

            // Act
            result, err := service.Login(context.Background(), tt.username, tt.password)

            // Assert
            if tt.expectError {
                assert.Error(t, err)
                if tt.errorType != "" {
                    var appErr *models.ErrorResponse
                    assert.ErrorAs(t, err, &appErr)
                    assert.Equal(t, string(tt.errorType), appErr.Code)
                }
            } else {
                assert.NoError(t, err)
                assert.NotEmpty(t, result.AccessToken)
                assert.NotEmpty(t, result.RefreshToken)
                assert.NotNil(t, result.User)
            }

            mockRepo.AssertExpectations(t)
            mockPM.AssertExpectations(t)
        })
    }
}
```

## Testing Middleware

### Authentication Middleware Testing
```go
func TestAuthMiddleware(t *testing.T) {
    tests := []struct {
        name           string
        setupRequest   func() *http.Request
        expectedStatus int
        expectedError  string
    }{
        {
            name: "Valid token",
            setupRequest: func() *http.Request {
                token := generateValidToken()
                req, _ := http.NewRequest("GET", "/protected", nil)
                req.Header.Set("Authorization", "Bearer "+token)
                return req
            },
            expectedStatus: http.StatusOK,
        },
        {
            name: "Missing token",
            setupRequest: func() *http.Request {
                req, _ := http.NewRequest("GET", "/protected", nil)
                return req
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Missing authorization token",
        },
        {
            name: "Invalid token format",
            setupRequest: func() *http.Request {
                req, _ := http.NewRequest("GET", "/protected", nil)
                req.Header.Set("Authorization", "InvalidFormat token")
                return req
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Invalid authorization format",
        },
        {
            name: "Expired token",
            setupRequest: func() *http.Request {
                token := generateExpiredToken()
                req, _ := http.NewRequest("GET", "/protected", nil)
                req.Header.Set("Authorization", "Bearer "+token)
                return req
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "Token is expired",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            jwtManager := auth.NewJWTManager("test-secret", time.Minute, time.Hour*24)
            middleware := AuthMiddleware(jwtManager)
            
            // Create a handler that the middleware will wrap
            handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("protected resource"))
            })

            // Act
            req := tt.setupRequest()
            rr := httptest.NewRecorder()
            
            // Wrap the handler with middleware
            wrappedHandler := middleware(handler)
            wrappedHandler.ServeHTTP(rr, req)

            // Assert
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedError != "" {
                var response map[string]interface{}
                err := json.Unmarshal(rr.Body.Bytes(), &response)
                require.NoError(t, err)
                assert.Contains(t, response["message"], tt.expectedError)
            }
        })
    }
}
```

## Integration Testing

### Integration Testing with TestContainers
```go
func TestCIHandler_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }

    // Setup PostgreSQL container
    ctx := context.Background()
    pgContainer, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15-alpine"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
    )
    require.NoError(t, err)
    defer pgContainer.Terminate(ctx)

    // Get connection string
    connStr, err := pgContainer.ConnectionString(ctx)
    require.NoError(t, err)

    // Connect to database
    db, err := sqlx.Connect("postgres", connStr)
    require.NoError(t, err)
    defer db.Close()

    // Run migrations
    err = runMigrations(db)
    require.NoError(t, err)

    // Setup real repositories
    ciRepo := repositories.NewCIPostgresRepository(db)
    relRepo := repositories.NewRelationshipPostgresRepository(db)
    auditRepo := repositories.NewAuditLogPostgresRepository(db)

    // Setup handler with real dependencies
    handler := handlers.NewCIHandler(ciRepo, relRepo, auditRepo)

    // Test the full flow
    t.Run("Create and Get CI", func(t *testing.T) {
        // Create CI
        createCI := models.CI{
            Name:  "Integration Test CI",
            Type:  "server",
            Tags:  []string{"integration", "test"},
        }

        req, err := http.NewRequest("POST", "/api/v1/cis", createTestJSONRequest(createCI))
        require.NoError(t, err)
        
        // Add authentication context
        ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
        req = req.WithContext(ctx)

        rr := httptest.NewRecorder()
        handler.CreateCI(rr, req)

        assert.Equal(t, http.StatusCreated, rr.Code)

        var createdCI models.CI
        err = json.Unmarshal(rr.Body.Bytes(), &createdCI)
        require.NoError(t, err)
        assert.NotEmpty(t, createdCI.ID)

        // Get CI
        req, err = http.NewRequest("GET", "/api/v1/cis/"+createdCI.ID.String(), nil)
        require.NoError(t, err)

        rr = httptest.NewRecorder()
        
        // Setup router for parameter extraction
        router := mux.NewRouter()
        router.HandleFunc("/api/v1/cis/{id}", handler.GetCI).Methods("GET")
        router.ServeHTTP(rr, req)

        assert.Equal(t, http.StatusOK, rr.Code)

        var retrievedCI models.CI
        err = json.Unmarshal(rr.Body.Bytes(), &retrievedCI)
        require.NoError(t, err)
        assert.Equal(t, createdCI.ID, retrievedCI.ID)
        assert.Equal(t, createdCI.Name, retrievedCI.Name)
        assert.Equal(t, createdCI.Type, retrievedCI.Type)
    })
}

// Helper functions
func createTestJSONRequest(data interface{}) *bytes.Buffer {
    jsonData, err := json.Marshal(data)
    if err != nil {
        panic(err)
    }
    return bytes.NewBuffer(jsonData)
}

func runMigrations(db *sqlx.DB) error {
    // Run database migrations here
    // This could involve reading migration files and executing them
    return nil
}
```

## Best Practices

### 1. Test Organization
```go
// Good: Clear test structure with Arrange-Act-Assert
func TestCIHandler_CreateCI_Success(t *testing.T) {
    // Arrange
    mockRepo := setupMockRepository()
    handler := NewCIHandler(mockRepo, mockRelRepo, mockAuditRepo)
    testCI := &models.CI{Name: "Test CI", Type: "server"}
    
    // Act
    result, err := handler.CreateCI(testCI)
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, testCI.Name, result.Name)
}

// Bad: Mixed concerns and unclear structure
func TestCIHandler_CreateCI(t *testing.T) {
    handler := NewCIHandler(mockRepo, mockRelRepo, mockAuditRepo)
    result, err := handler.CreateCI(&models.CI{Name: "Test CI", Type: "server"})
    assert.NoError(t, err)
    // Missing assertions and setup
}
```

### 2. Mock Usage
```go
// Good: Proper mock setup and verification
func TestCIRepository_Create(t *testing.T) {
    mock := &MockCIRepository{}
    mock.On("Create", mock.Anything, mock.AnythingOfType("*models.CI")).
        Return(nil)
    
    repo := NewCIService(mock)
    err := repo.CreateCI(context.Background(), &models.CI{Name: "Test"})
    
    assert.NoError(t, err)
    mock.AssertExpectations(t) // Verify mock expectations
}

// Bad: No mock verification
func TestCIRepository_Create(t *testing.T) {
    mock := &MockCIRepository{}
    mock.On("Create", mock.Anything, mock.Anything).Return(nil)
    
    repo := NewCIService(mock)
    err := repo.CreateCI(context.Background(), &models.CI{Name: "Test"})
    
    assert.NoError(t, err)
    // Missing mock.AssertExpectations(t)
}
```

### 3. Error Testing
```go
// Good: Comprehensive error testing
func TestCIHandler_CreateCI_ValidationErrors(t *testing.T) {
    tests := []struct {
        name     string
        ci       *models.CI
        expected string
    }{
        {
            name:     "Empty name",
            ci:       &models.CI{Name: "", Type: "server"},
            expected: "name is required",
        },
        {
            name:     "Empty type",
            ci:       &models.CI{Name: "Test", Type: ""},
            expected: "type is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            handler := NewCIHandler(mockRepo, mockRelRepo, mockAuditRepo)
            
            _, err := handler.CreateCI(tt.ci)
            
            assert.Error(t, err)
            assert.Contains(t, err.Error(), tt.expected)
        })
    }
}
```

### 4. Test Data Management
```go
// Good: Centralized test data creation
func createTestCI(name, ciType string) *models.CI {
    return &models.CI{
        ID:        uuid.New(),
        Name:      name,
        Type:      ciType,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
}

func createTestUser(username, role string) *models.User {
    return &models.User{
        ID:        uuid.New(),
        Username:  username,
        Role:      role,
        CreatedAt: time.Now(),
    }
}

// Usage in tests
func TestCIHandler_GetCI(t *testing.T) {
    testCI := createTestCI("Test CI", "server")
    // ... test implementation
}
```

## Test Examples

### Complete Handler Test Example
```go
package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/cmdb-lite/backend/internal/middleware"
    "github.com/cmdb-lite/backend/internal/models"
    "github.com/cmdb-lite/backend/internal/repositories"
    "github.com/cmdb-lite/backend/internal/testutils"
    "github.com/google/uuid"
    "github.com/gorilla/mux"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestCIHandler_CreateCI_Comprehensive(t *testing.T) {
    // Test data
    validCI := models.CI{
        Name: "Test Server",
        Type: "server",
        Attributes: models.JSONBMap{
            "cpu":    "4 cores",
            "memory": "16GB",
            "storage": "500GB SSD",
        },
        Tags: []string{"production", "linux"},
    }

    tests := []struct {
        name           string
        setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
        requestBody    interface{}
        setupContext   func(*http.Request) *http.Request
        expectedStatus int
        expectedError  string
        validateFunc   func(*testing.T, []byte)
    }{
        {
            name: "Successful CI creation with full data",
            setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
                ciRepo := &MockCIRepository{err: nil}
                relRepo := &MockRelationshipRepository{}
                auditRepo := &MockAuditLogRepository{err: nil}
                return ciRepo, relRepo, auditRepo
            },
            requestBody: validCI,
            setupContext: func(req *http.Request) *http.Request {
                ctx := context.WithValue(req.Context(), middleware.UsernameKey, "admin")
                return req.WithContext(ctx)
            },
            expectedStatus: http.StatusCreated,
            validateFunc: func(t *testing.T, body []byte) {
                var response models.CI
                err := json.Unmarshal(body, &response)
                require.NoError(t, err)
                
                assert.NotEmpty(t, response.ID)
                assert.Equal(t, validCI.Name, response.Name)
                assert.Equal(t, validCI.Type, response.Type)
                assert.Equal(t, validCI.Attributes, response.Attributes)
                assert.Equal(t, validCI.Tags, response.Tags)
                assert.NotZero(t, response.CreatedAt)
                assert.NotZero(t, response.UpdatedAt)
            },
        },
        {
            name: "Validation error - missing required fields",
            setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
                return &MockCIRepository{}, &MockRelationshipRepository{}, &MockAuditLogRepository{}
            },
            requestBody: models.CI{Name: "", Type: "server"},
            setupContext: func(req *http.Request) *http.Request {
                ctx := context.WithValue(req.Context(), middleware.UsernameKey, "admin")
                return req.WithContext(ctx)
            },
            expectedStatus: http.StatusBadRequest,
            expectedError:  "Validation failed",
            validateFunc: func(t *testing.T, body []byte) {
                var response models.ErrorResponse
                err := json.Unmarshal(body, &response)
                require.NoError(t, err)
                
                assert.Equal(t, "VALIDATION_ERROR", response.Code)
                assert.Contains(t, response.Details.(map[string]interface{}), "name")
            },
        },
        {
            name: "Authentication error - missing user context",
            setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
                return &MockCIRepository{}, &MockRelationshipRepository{}, &MockAuditLogRepository{}
            },
            requestBody: validCI,
            setupContext: func(req *http.Request) *http.Request {
                // No user context
                return req
            },
            expectedStatus: http.StatusUnauthorized,
            expectedError:  "User not authenticated",
        },
        {
            name: "Repository error - database failure",
            setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
                ciRepo := &MockCIRepository{err: assert.AnError}
                relRepo := &MockRelationshipRepository{}
                auditRepo := &MockAuditLogRepository{err: nil}
                return ciRepo, relRepo, auditRepo
            },
            requestBody: validCI,
            setupContext: func(req *http.Request) *http.Request {
                ctx := context.WithValue(req.Context(), middleware.UsernameKey, "admin")
                return req.WithContext(ctx)
            },
            expectedStatus: http.StatusInternalServerError,
            expectedError:  "Failed to create CI",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            ciRepo, relRepo, auditRepo := tt.setupMock()
            handler := NewCIHandler(ciRepo, relRepo, auditRepo)
            
            // Create request
            var reqBody []byte
            var err error
            
            switch v := tt.requestBody.(type) {
            case string:
                reqBody = []byte(v)
            default:
                reqBody, err = json.Marshal(v)
                require.NoError(t, err)
            }
            
            req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(reqBody))
            require.NoError(t, err)
            req.Header.Set("Content-Type", "application/json")
            
            // Setup context
            req = tt.setupContext(req)
            
            // Create response recorder
            rr := httptest.NewRecorder()
            
            // Act
            handler.CreateCI(rr, req)
            
            // Assert
            assert.Equal(t, tt.expectedStatus, rr.Code)
            
            if tt.expectedError != "" {
                assert.Contains(t, rr.Body.String(), tt.expectedError)
            }
            
            if tt.validateFunc != nil {
                tt.validateFunc(t, rr.Body.Bytes())
            }
        })
    }
}
```

This comprehensive testing guide provides everything needed to write effective unit tests for the CMDB Lite application, covering all aspects from basic test structure to advanced integration testing patterns.
