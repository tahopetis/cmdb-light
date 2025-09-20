# CMDB Lite Testing Guide

This document provides information about the testing strategy and practices used in the CMDB Lite project. It covers unit testing, integration testing, end-to-end testing, and how to run tests.

## Table of Contents

- [Testing Strategy](#testing-strategy)
- [Backend Testing](#backend-testing)
  - [Unit Testing](#unit-testing)
  - [Integration Testing](#integration-testing)
  - [Test Coverage](#test-coverage)
- [Frontend Testing](#frontend-testing)
  - [Component Testing](#component-testing)
  - [View Testing](#view-testing)
  - [Store Testing](#store-testing)
- [End-to-End Testing](#end-to-end-testing)
- [Test Data Management](#test-data-management)
- [Running Tests](#running-tests)
  - [Backend Tests](#backend-tests)
  - [Frontend Tests](#frontend-tests)
  - [End-to-End Tests](#end-to-end-tests)
  - [All Tests](#all-tests)
- [Writing Tests](#writing-tests)
  - [Backend Test Examples](#backend-test-examples)
  - [Frontend Test Examples](#frontend-test-examples)
  - [E2E Test Examples](#e2e-test-examples)
- [Continuous Integration](#continuous-integration)

## Testing Strategy

CMDB Lite follows a comprehensive testing strategy that includes multiple levels of testing to ensure the quality and reliability of the application:

- **Unit Testing**: Tests individual components and functions in isolation.
- **Integration Testing**: Tests how multiple components work together.
- **End-to-End Testing**: Tests the entire application flow from a user's perspective.

### Testing Pyramid

We follow the testing pyramid model, with a larger number of unit tests, fewer integration tests, and even fewer end-to-end tests:

```
    E2E Tests
   /           \
  / Integration \
 /               \
-------------------
|   Unit Tests   |
-------------------
```

This approach ensures that we have fast feedback from unit tests while still verifying that the application works as a whole.

## Backend Testing

The backend uses Go's built-in testing framework along with additional libraries for mocking and assertions.

### Unit Testing

Unit tests focus on testing individual functions and methods in isolation. We use the following libraries for unit testing:

- **testing**: Go's built-in testing package
- **testify**: Provides assertions and mock functionality
- **gomock**: Go mocking framework
- **sqlmock**: Mock SQL database driver

### Integration Testing

Integration tests verify that multiple components work together correctly. For the backend, this primarily means testing the interaction between handlers, services, and the database.

We use **Testcontainers** to spin up real database instances for integration tests, ensuring that our tests run against a real database rather than in-memory mocks.

### Test Coverage

We aim for a minimum of 80% test coverage for all backend code. Coverage is measured using Go's built-in coverage tool.

To generate a coverage report:

```bash
cd backend
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

This will generate an HTML report showing which lines of code are covered by tests.

## Frontend Testing

The frontend uses Vitest for testing, along with Vue Test Utils for testing Vue components.

### Component Testing

Component tests focus on testing individual Vue components in isolation. We test:

- Component rendering
- User interactions (clicks, form submissions, etc.)
- Props and events
- Computed properties and methods

### View Testing

View tests verify that entire views (pages) work correctly, including:

- Navigation between views
- Data fetching and display
- User authentication and authorization

### Store Testing

Store tests verify that Pinia stores work correctly, including:

- State management
- Actions and mutations
- Getters

## End-to-End Testing

End-to-end (E2E) tests verify that the entire application works correctly from a user's perspective. We use Playwright for E2E testing, which allows us to:

- Automate browser interactions
- Test across multiple browsers (Chrome, Firefox, Safari)
- Test responsive design on different screen sizes
- Test authentication flows

## Test Data Management

### Backend Test Data

For backend tests, we use a combination of approaches:

- **Factories**: Generate test data with realistic values
- **Fixtures**: Predefined data for common test scenarios
- **Seed Scripts**: Populate the test database with initial data

### Frontend Test Data

For frontend tests, we use:

- **Mock Services**: Mock API responses to test components without hitting a real backend
- **Test Factories**: Generate test data for components
- **Storybook**: For visual testing of components in isolation (future feature)

## Running Tests

### Backend Tests

To run all backend tests:

```bash
cd backend
go test ./...
```

To run tests for a specific package:

```bash
go test ./internal/handlers
```

To run tests with verbose output:

```bash
go test -v ./...
```

To run tests with coverage:

```bash
go test -cover ./...
```

To run a specific test:

```bash
go test -run TestSpecificFunction ./internal/handlers
```

### Frontend Tests

To run all frontend tests:

```bash
cd frontend
npm test
```

To run tests in watch mode:

```bash
npm run test:watch
```

To run tests with coverage:

```bash
npm run test:coverage
```

To run tests for a specific file:

```bash
npm test ExampleComponent.spec.js
```

### End-to-End Tests

To run all E2E tests:

```bash
npm run test:e2e
```

To run E2E tests in headed mode (visible browser):

```bash
npm run test:e2e:headed
```

To run E2E tests for a specific browser:

```bash
npm run test:e2e -- --browser=firefox
```

To run a specific E2E test file:

```bash
npm run test:e2e -- example.spec.js
```

### All Tests

To run all tests (backend, frontend, and E2E):

```bash
npm run test:all
```

This command is defined in the root package.json and runs all test suites in sequence.

## Writing Tests

### Backend Test Examples

#### Handler Test Example

```go
package handlers

import (
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "github.com/yourorg/cmdb-lite/internal/mocks"
    "github.com/yourorg/cmdb-lite/internal/models"
)

func TestGetCIHandler(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    
    // Create a mock repository
    mockRepo := new(mocks.CIRepository)
    
    // Create a test CI
    testCI := models.ConfigurationItem{
        ID:   "123e4567-e89b-12d3-a456-426614174000",
        Name: "Test CI",
        Type: "Server",
    }
    
    // Set up expectations
    mockRepo.On("GetByID", testCI.ID).Return(&testCI, nil)
    
    // Create a handler with the mock repository
    handler := NewCIHandler(mockRepo)
    
    // Create a router
    router := gin.Default()
    router.GET("/cis/:id", handler.GetCI)
    
    // Create a test request
    req, _ := http.NewRequest("GET", "/cis/"+testCI.ID, nil)
    w := httptest.NewRecorder()
    
    // Perform the request
    router.ServeHTTP(w, req)
    
    // Assert the response
    assert.Equal(t, http.StatusOK, w.Code)
    
    // Assert that the expectations were met
    mockRepo.AssertExpectations(t)
}
```

#### Service Test Example

```go
package services

import (
    "testing"
    
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "github.com/yourorg/cmdb-lite/internal/mocks"
    "github.com/yourorg/cmdb-lite/internal/models"
)

func TestCIService_CreateCI(t *testing.T) {
    // Create a mock repository
    mockRepo := new(mocks.CIRepository)
    
    // Create test data
    testCI := &models.ConfigurationItem{
        Name: "Test CI",
        Type: "Server",
    }
    
    // Set up expectations
    mockRepo.On("Create", testCI).Return(nil)
    
    // Create a service with the mock repository
    service := NewCIService(mockRepo)
    
    // Call the method under test
    err := service.CreateCI(testCI)
    
    // Assert the result
    assert.NoError(t, err)
    
    // Assert that the expectations were met
    mockRepo.AssertExpectations(t)
}
```

#### Repository Test Example with sqlmock

```go
package repositories

import (
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    
    "github.com/yourorg/cmdb-lite/internal/models"
)

func TestCIRepository_GetByID(t *testing.T) {
    // Create a mock database connection
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("Error creating mock database: %v", err)
    }
    defer db.Close()
    
    // Create a GORM database instance with the mock
    gormDB, err := gorm.Open(postgres.New(postgres.Config{
        Conn: db,
    }), &gorm.Config{})
    if err != nil {
        t.Fatalf("Error creating GORM database: %v", err)
    }
    
    // Create a repository with the mock database
    repo := NewCIRepository(gormDB)
    
    // Test data
    testID := "123e4567-e89b-12d3-a456-426614174000"
    
    // Set up expectations
    rows := sqlmock.NewRows([]string{"id", "name", "type", "created_at", "updated_at"}).
        AddRow(testID, "Test CI", "Server", "2023-01-01T00:00:00Z", "2023-01-01T00:00:00Z")
    
    mock.ExpectQuery(`SELECT \* FROM "configuration_items" WHERE id = \$1`).
        WithArgs(testID).
        WillReturnRows(rows)
    
    // Call the method under test
    ci, err := repo.GetByID(testID)
    
    // Assert the result
    assert.NoError(t, err)
    assert.NotNil(t, ci)
    assert.Equal(t, testID, ci.ID)
    assert.Equal(t, "Test CI", ci.Name)
    assert.Equal(t, "Server", ci.Type)
    
    // Assert that all expectations were met
    assert.NoError(t, mock.ExpectationsWereMet())
}
```

### Frontend Test Examples

#### Component Test Example

```javascript
import { render, screen, fireEvent } from '@testing-library/vue'
import { describe, it, expect, vi } from 'vitest'
import CIForm from '@/components/forms/CIForm.vue'
import { useCIStore } from '@/stores/ciStore'

// Mock the CI store
vi.mock('@/stores/ciStore', () => ({
  useCIStore: vi.fn()
}))

describe('CIForm', () => {
  it('renders form fields correctly', () => {
    // Setup mock store
    const mockStore = {
      ci: null,
      loading: false,
      error: null,
      createCI: vi.fn(),
      updateCI: vi.fn()
    }
    useCIStore.mockReturnValue(mockStore)
    
    // Render the component
    render(CIForm)
    
    // Assert that form fields are rendered
    expect(screen.getByLabelText('Name')).toBeInTheDocument()
    expect(screen.getByLabelText('Type')).toBeInTheDocument()
    expect(screen.getByText('Save')).toBeInTheDocument()
  })
  
  it('submits form with correct data', async () => {
    // Setup mock store
    const mockStore = {
      ci: null,
      loading: false,
      error: null,
      createCI: vi.fn()
    }
    useCIStore.mockReturnValue(mockStore)
    
    // Render the component
    render(CIForm)
    
    // Fill in form fields
    await fireEvent.update(screen.getByLabelText('Name'), 'Test CI')
    await fireEvent.update(screen.getByLabelText('Type'), 'Server')
    
    // Submit the form
    await fireEvent.click(screen.getByText('Save'))
    
    // Assert that createCI was called with correct data
    expect(mockStore.createCI).toHaveBeenCalledWith({
      name: 'Test CI',
      type: 'Server',
      attributes: {},
      tags: []
    })
  })
  
  it('shows loading state when submitting', async () => {
    // Setup mock store
    const mockStore = {
      ci: null,
      loading: true,
      error: null,
      createCI: vi.fn()
    }
    useCIStore.mockReturnValue(mockStore)
    
    // Render the component
    render(CIForm)
    
    // Assert that save button is disabled and shows loading text
    const saveButton = screen.getByText('Saving...')
    expect(saveButton).toBeInTheDocument()
    expect(saveButton).toBeDisabled()
  })
})
```

#### View Test Example

```javascript
import { render, screen, waitFor } from '@testing-library/vue'
import { describe, it, expect, vi } from 'vitest'
import DashboardView from '@/views/DashboardView.vue'
import { useCIStore } from '@/stores/ciStore'
import { useAuthStore } from '@/stores/authStore'

// Mock stores
vi.mock('@/stores/ciStore', () => ({
  useCIStore: vi.fn()
}))

vi.mock('@/stores/authStore', () => ({
  useAuthStore: vi.fn()
}))

describe('DashboardView', () => {
  it('renders dashboard with CI statistics', async () => {
    // Setup mock stores
    const mockCIStore = {
      cis: [
        { id: '1', name: 'CI 1', type: 'Server' },
        { id: '2', name: 'CI 2', type: 'Application' },
        { id: '3', name: 'CI 3', type: 'Database' }
      ],
      loading: false,
      error: null,
      fetchCIs: vi.fn()
    }
    
    const mockAuthStore = {
      user: { id: '1', username: 'testuser', role: 'admin' },
      isAuthenticated: true
    }
    
    useCIStore.mockReturnValue(mockCIStore)
    useAuthStore.mockReturnValue(mockAuthStore)
    
    // Render the view
    render(DashboardView)
    
    // Wait for async operations to complete
    await waitFor(() => {
      // Assert that dashboard title is rendered
      expect(screen.getByText('Dashboard')).toBeInTheDocument()
      
      // Assert that CI statistics are rendered
      expect(screen.getByText('Total CIs: 3')).toBeInTheDocument()
      
      // Assert that CI type breakdown is rendered
      expect(screen.getByText('Servers: 1')).toBeInTheDocument()
      expect(screen.getByText('Applications: 1')).toBeInTheDocument()
      expect(screen.getByText('Databases: 1')).toBeInTheDocument()
    })
    
    // Assert that fetchCIs was called
    expect(mockCIStore.fetchCIs).toHaveBeenCalled()
  })
  
  it('shows loading state while fetching data', () => {
    // Setup mock stores
    const mockCIStore = {
      cis: [],
      loading: true,
      error: null,
      fetchCIs: vi.fn()
    }
    
    const mockAuthStore = {
      user: { id: '1', username: 'testuser', role: 'admin' },
      isAuthenticated: true
    }
    
    useCIStore.mockReturnValue(mockCIStore)
    useAuthStore.mockReturnValue(mockAuthStore)
    
    // Render the view
    render(DashboardView)
    
    // Assert that loading indicator is shown
    expect(screen.getByText('Loading...')).toBeInTheDocument()
  })
  
  it('shows error message when fetch fails', async () => {
    // Setup mock stores
    const mockCIStore = {
      cis: [],
      loading: false,
      error: 'Failed to fetch CIs',
      fetchCIs: vi.fn()
    }
    
    const mockAuthStore = {
      user: { id: '1', username: 'testuser', role: 'admin' },
      isAuthenticated: true
    }
    
    useCIStore.mockReturnValue(mockCIStore)
    useAuthStore.mockReturnValue(mockAuthStore)
    
    // Render the view
    render(DashboardView)
    
    // Wait for async operations to complete
    await waitFor(() => {
      // Assert that error message is shown
      expect(screen.getByText('Failed to fetch CIs')).toBeInTheDocument()
    })
  })
})
```

#### Store Test Example

```javascript
import { setActivePinia, createPinia } from 'pinia'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { useCIStore } from '@/stores/ciStore'
import { ciService } from '@/services/ciService'

// Mock the CI service
vi.mock('@/services/ciService', () => ({
  ciService: {
    getAll: vi.fn(),
    create: vi.fn(),
    update: vi.fn(),
    delete: vi.fn()
  }
}))

describe('CI Store', () => {
  beforeEach(() => {
    // Create a fresh pinia instance for each test
    setActivePinia(createPinia())
  })
  
  it('fetches CIs and updates state', async () => {
    // Setup mock response
    const mockCIs = [
      { id: '1', name: 'CI 1', type: 'Server' },
      { id: '2', name: 'CI 2', type: 'Application' }
    ]
    ciService.getAll.mockResolvedValue({ data: mockCIs })
    
    // Get the store instance
    const store = useCIStore()
    
    // Call the action
    await store.fetchCIs()
    
    // Assert that the service was called
    expect(ciService.getAll).toHaveBeenCalled()
    
    // Assert that the state was updated
    expect(store.cis).toEqual(mockCIs)
    expect(store.loading).toBe(false)
    expect(store.error).toBe(null)
  })
  
  it('handles error when fetching CIs', async () => {
    // Setup mock error
    const errorMessage = 'Failed to fetch CIs'
    ciService.getAll.mockRejectedValue(new Error(errorMessage))
    
    // Get the store instance
    const store = useCIStore()
    
    // Call the action
    await store.fetchCIs()
    
    // Assert that the service was called
    expect(ciService.getAll).toHaveBeenCalled()
    
    // Assert that the error state was set
    expect(store.cis).toEqual([])
    expect(store.loading).toBe(false)
    expect(store.error).toBe(errorMessage)
  })
  
  it('creates a new CI and updates state', async () => {
    // Setup mock response
    const newCI = { id: '3', name: 'CI 3', type: 'Database' }
    ciService.create.mockResolvedValue({ data: newCI })
    
    // Get the store instance
    const store = useCIStore()
    
    // Set initial state
    store.cis = [
      { id: '1', name: 'CI 1', type: 'Server' },
      { id: '2', name: 'CI 2', type: 'Application' }
    ]
    
    // Call the action
    await store.createCI(newCI)
    
    // Assert that the service was called
    expect(ciService.create).toHaveBeenCalledWith(newCI)
    
    // Assert that the state was updated
    expect(store.cis).toHaveLength(3)
    expect(store.cis).toContainEqual(newCI)
  })
})
```

### E2E Test Examples

#### Authentication Flow Test

```javascript
import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test('successful login and logout', async ({ page }) => {
    // Navigate to login page
    await page.goto('/login')
    
    // Fill in login form
    await page.fill('[data-testid="username"]', 'admin')
    await page.fill('[data-testid="password"]', 'password')
    
    // Click login button
    await page.click('[data-testid="login-button"]')
    
    // Verify redirect to dashboard
    await expect(page).toHaveURL('/dashboard')
    
    // Verify user is logged in
    await expect(page.locator('[data-testid="user-menu"]')).toBeVisible()
    
    // Logout
    await page.click('[data-testid="user-menu"]')
    await page.click('[data-testid="logout-button"]')
    
    // Verify redirect to login page
    await expect(page).toHaveURL('/login')
  })
  
  test('shows error for invalid credentials', async ({ page }) => {
    // Navigate to login page
    await page.goto('/login')
    
    // Fill in login form with invalid credentials
    await page.fill('[data-testid="username"]', 'invalid')
    await page.fill('[data-testid="password"]', 'invalid')
    
    // Click login button
    await page.click('[data-testid="login-button"]')
    
    // Verify error message is shown
    await expect(page.locator('[data-testid="error-message"]')).toBeVisible()
    await expect(page.locator('[data-testid="error-message"]')).toHaveText('Invalid username or password')
    
    // Verify user is not redirected
    await expect(page).toHaveURL('/login')
  })
})
```

#### CRUD Operations Test

```javascript
import { test, expect } from '@playwright/test'

test.describe('CI CRUD Operations', () => {
  test.beforeEach(async ({ page }) => {
    // Login before each test
    await page.goto('/login')
    await page.fill('[data-testid="username"]', 'admin')
    await page.fill('[data-testid="password"]', 'password')
    await page.click('[data-testid="login-button"]')
    
    // Verify redirect to dashboard
    await expect(page).toHaveURL('/dashboard')
  })
  
  test('create, read, update, and delete a CI', async ({ page }) => {
    // Navigate to CI list
    await page.click('[data-testid="cis-menu-item"]')
    await expect(page).toHaveURL('/cis')
    
    // Click "Create New CI" button
    await page.click('[data-testid="create-ci-button"]')
    
    // Fill in CI form
    await page.fill('[data-testid="ci-name"]', 'Test CI')
    await page.selectOption('[data-testid="ci-type"]', 'Server')
    await page.fill('[data-testid="ci-attribute-ip_address"]', '192.168.1.100')
    
    // Click save button
    await page.click('[data-testid="save-ci-button"]')
    
    // Verify redirect to CI list
    await expect(page).toHaveURL('/cis')
    
    // Verify CI was created
    await expect(page.locator('text=Test CI')).toBeVisible()
    
    // Click on the CI to view details
    await page.click('text=Test CI')
    
    // Verify CI details page
    await expect(page.locator('[data-testid="ci-details"]')).toBeVisible()
    await expect(page.locator('[data-testid="ci-name-value"]')).toHaveText('Test CI')
    await expect(page.locator('[data-testid="ci-type-value"]')).toHaveText('Server')
    await expect(page.locator('[data-testid="ci-attribute-ip_address"]')).toHaveText('192.168.1.100')
    
    // Click edit button
    await page.click('[data-testid="edit-ci-button"]')
    
    // Update CI form
    await page.fill('[data-testid="ci-name"]', 'Updated Test CI')
    await page.fill('[data-testid="ci-attribute-ip_address"]', '192.168.1.101')
    
    // Click save button
    await page.click('[data-testid="save-ci-button"]')
    
    // Verify CI was updated
    await expect(page.locator('[data-testid="ci-name-value"]')).toHaveText('Updated Test CI')
    await expect(page.locator('[data-testid="ci-attribute-ip_address"]')).toHaveText('192.168.1.101')
    
    // Click delete button
    await page.click('[data-testid="delete-ci-button"]')
    
    // Confirm deletion in dialog
    await page.click('[data-testid="confirm-delete-button"]')
    
    // Verify redirect to CI list
    await expect(page).toHaveURL('/cis')
    
    // Verify CI was deleted
    await expect(page.locator('text=Updated Test CI')).not.toBeVisible()
  })
})
```

## Continuous Integration

CMDB Lite uses GitHub Actions for continuous integration. The CI pipeline runs on every push and pull request, and includes the following steps:

1. **Code Checkout**: Check out the code from the repository.
2. **Setup Dependencies**: Set up Go, Node.js, and other dependencies.
3. **Linting**: Run linters to ensure code quality.
4. **Backend Tests**: Run all backend tests with coverage.
5. **Frontend Tests**: Run all frontend tests with coverage.
6. **E2E Tests**: Run all end-to-end tests.
7. **Build**: Build the application.
8. **Upload Coverage**: Upload coverage reports to a coverage service.

### Test Coverage Requirements

The CI pipeline enforces minimum test coverage requirements:

- **Backend**: Minimum 80% coverage
- **Frontend**: Minimum 70% coverage

If coverage falls below these thresholds, the pipeline will fail.

### Test Results

Test results are published as artifacts in GitHub Actions, allowing developers to download and review detailed test reports.

For more information on contributing to tests, see the [Developer Guide](README.md) and the [Contribution Guide](../project/contributing.md).