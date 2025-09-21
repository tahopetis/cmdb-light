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
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCIRepository is a mock implementation of CIRepository interface
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

func (m *MockCIRepository) Update(ctx context.Context, ci *models.CI) error {
	args := m.Called(ctx, ci)
	return args.Error(0)
}

func (m *MockCIRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRelationshipRepository is a mock implementation of RelationshipRepository interface
type MockRelationshipRepository struct {
	mock.Mock
}

func (m *MockRelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	args := m.Called(ctx, rel)
	return args.Error(0)
}

func (m *MockRelationshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Relationship), args.Error(1)
}

func (m *MockRelationshipRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Relationship), args.Error(1)
}

func (m *MockRelationshipRepository) GetBySourceID(ctx context.Context, id uuid.UUID) ([]*models.Relationship, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*models.Relationship), args.Error(1)
}

func (m *MockRelationshipRepository) GetByTargetID(ctx context.Context, id uuid.UUID) ([]*models.Relationship, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]*models.Relationship), args.Error(1)
}

func (m *MockRelationshipRepository) Update(ctx context.Context, rel *models.Relationship) error {
	args := m.Called(ctx, rel)
	return args.Error(0)
}

func (m *MockRelationshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAuditLogRepository is a mock implementation of AuditLogRepository interface
type MockAuditLogRepository struct {
	mock.Mock
}

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	args := m.Called(ctx, auditLog)
	return args.Error(0)
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetAll(ctx context.Context) ([]*models.AuditLog, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error) {
	args := m.Called(ctx, entityType)
	return args.Get(0).([]*models.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error) {
	args := m.Called(ctx, entityID)
	return args.Get(0).([]*models.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error) {
	args := m.Called(ctx, changedBy)
	return args.Get(0).([]*models.AuditLog), args.Error(1)
}

func (m *MockAuditLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helper functions to create test data
func createTestCI(name, ciType string) *models.CI {
	return &models.CI{
		ID:        uuid.New(),
		Name:      name,
		Type:      ciType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestRelationship(sourceID, targetID uuid.UUID, relType string) *models.Relationship {
	return &models.Relationship{
		ID:        uuid.New(),
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      relType,
		CreatedAt: time.Now(),
	}
}

func createTestAuditLog(entityType string, entityID uuid.UUID, action, changedBy string) *models.AuditLog {
	return &models.AuditLog{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		ChangedBy:  changedBy,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{},
	}
}

func TestCIHandler_CreateCI_Success(t *testing.T) {
	// Arrange
	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CI")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	testCI := models.CI{
		Name: "Test Server",
		Type: "server",
		Attributes: models.JSONBMap{
			"cpu":    "4 cores",
			"memory": "16GB",
		},
		Tags: []string{"production", "linux"},
	}

	requestBody, err := json.Marshal(testCI)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add authentication context
	ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	handler.CreateCI(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)

	var response models.CI
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.NotEmpty(t, response.ID)
	assert.Equal(t, testCI.Name, response.Name)
	assert.Equal(t, testCI.Type, response.Type)
	assert.Equal(t, testCI.Attributes, response.Attributes)
	assert.Equal(t, testCI.Tags, response.Tags)

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCIHandler_CreateCI_ValidationError(t *testing.T) {
	// Arrange
	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	// Test with missing required fields
	invalidCI := models.CI{
		Name: "", // Missing name
		Type: "server",
	}

	requestBody, err := json.Marshal(invalidCI)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add authentication context
	ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	handler.CreateCI(rr, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "VALIDATION_ERROR", response.Code)
	assert.Contains(t, response.Details.(map[string]interface{}), "name")
}

func TestCIHandler_CreateCI_Unauthorized(t *testing.T) {
	// Arrange
	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	testCI := models.CI{
		Name: "Test Server",
		Type: "server",
	}

	requestBody, err := json.Marshal(testCI)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// No authentication context
	rr := httptest.NewRecorder()

	// Act
	handler.CreateCI(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "User not authenticated", response.Message)
}

func TestCIHandler_GetCI_Success(t *testing.T) {
	// Arrange
	testCI := createTestCI("Test Server", "server")

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetByID", mock.Anything, testCI.ID).Return(testCI, nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	req, err := http.NewRequest("GET", "/api/v1/cis/"+testCI.ID.String(), nil)
	require.NoError(t, err)

	// Setup router with parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/cis/{id}", handler.GetCI).Methods("GET")

	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.CI
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, testCI.ID, response.ID)
	assert.Equal(t, testCI.Name, response.Name)
	assert.Equal(t, testCI.Type, response.Type)

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
}

func TestCIHandler_GetCI_NotFound(t *testing.T) {
	// Arrange
	testID := uuid.New()

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetByID", mock.Anything, testID).Return(nil, repositories.ErrNotFound)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	req, err := http.NewRequest("GET", "/api/v1/cis/"+testID.String(), nil)
	require.NoError(t, err)

	// Setup router with parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/cis/{id}", handler.GetCI).Methods("GET")

	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "NOT_FOUND", response.Code)
	assert.Equal(t, "CI not found", response.Message)

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
}

func TestCIHandler_GetAllCIs_Success(t *testing.T) {
	// Arrange
	testCIs := []*models.CI{
		createTestCI("Server 1", "server"),
		createTestCI("Database 1", "database"),
		createTestCI("Application 1", "application"),
	}

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetAll", mock.Anything).Return(testCIs, nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	req, err := http.NewRequest("GET", "/api/v1/cis", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	handler.GetAllCIs(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check data
	data := response["data"].([]interface{})
	assert.Equal(t, len(testCIs), len(data))

	// Check pagination
	pagination := response["pagination"].(map[string]interface{})
	assert.Equal(t, 1, int(pagination["page"].(float64)))
	assert.Equal(t, 10, int(pagination["limit"].(float64)))
	assert.Equal(t, len(testCIs), int(pagination["total"].(float64)))

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
}

func TestCIHandler_UpdateCI_Success(t *testing.T) {
	// Arrange
	existingCI := createTestCI("Old Name", "server")
	updatedCI := models.CI{
		Name: "Updated Name",
		Type: "database",
		Attributes: models.JSONBMap{
			"cpu":    "8 cores",
			"memory": "32GB",
		},
		Tags: []string{"production", "critical"},
	}

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetByID", mock.Anything, existingCI.ID).Return(existingCI, nil)
	mockCIRepo.On("Update", mock.Anything, mock.AnythingOfType("*models.CI")).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	requestBody, err := json.Marshal(updatedCI)
	require.NoError(t, err)

	req, err := http.NewRequest("PUT", "/api/v1/cis/"+existingCI.ID.String(), bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Add authentication context
	ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
	req = req.WithContext(ctx)

	// Setup router with parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/cis/{id}", handler.UpdateCI).Methods("PUT")

	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response models.CI
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, existingCI.ID, response.ID)
	assert.Equal(t, updatedCI.Name, response.Name)
	assert.Equal(t, updatedCI.Type, response.Type)
	assert.Equal(t, updatedCI.Attributes, response.Attributes)
	assert.Equal(t, updatedCI.Tags, response.Tags)

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCIHandler_DeleteCI_Success(t *testing.T) {
	// Arrange
	testCI := createTestCI("Test Server", "server")

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetByID", mock.Anything, testCI.ID).Return(testCI, nil)
	mockCIRepo.On("Delete", mock.Anything, testCI.ID).Return(nil)
	mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	req, err := http.NewRequest("DELETE", "/api/v1/cis/"+testCI.ID.String(), nil)
	require.NoError(t, err)

	// Add authentication context
	ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
	req = req.WithContext(ctx)

	// Setup router with parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/cis/{id}", handler.DeleteCI).Methods("DELETE")

	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "CI deleted successfully", response["message"])

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
	mockAuditRepo.AssertExpectations(t)
}

func TestCIHandler_GetCIGraph_Success(t *testing.T) {
	// Arrange
	testCI := createTestCI("Test Server", "server")
	sourceRel := createTestRelationship(testCI.ID, uuid.New(), "depends_on")
	targetRel := createTestRelationship(uuid.New(), testCI.ID, "hosts")

	mockCIRepo := &MockCIRepository{}
	mockRelRepo := &MockRelationshipRepository{}
	mockAuditRepo := &MockAuditLogRepository{}

	// Setup mock expectations
	mockCIRepo.On("GetByID", mock.Anything, testCI.ID).Return(testCI, nil)
	mockRelRepo.On("GetBySourceID", mock.Anything, testCI.ID).Return([]*models.Relationship{sourceRel}, nil)
	mockRelRepo.On("GetByTargetID", mock.Anything, testCI.ID).Return([]*models.Relationship{targetRel}, nil)

	handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

	req, err := http.NewRequest("GET", "/api/v1/cis/"+testCI.ID.String()+"/graph", nil)
	require.NoError(t, err)

	// Setup router with parameters
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/cis/{id}/graph", handler.GetCIGraph).Methods("GET")

	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	// Check node
	node := response["node"].(map[string]interface{})
	assert.Equal(t, testCI.ID.String(), node["id"])
	assert.Equal(t, testCI.Name, node["name"])
	assert.Equal(t, testCI.Type, node["type"])

	// Check relationships
	relationships := response["relationships"].(map[string]interface{})
	outgoing := relationships["outgoing"].([]interface{})
	incoming := relationships["incoming"].([]interface{})

	assert.Equal(t, 1, len(outgoing))
	assert.Equal(t, 1, len(incoming))

	// Verify mock expectations
	mockCIRepo.AssertExpectations(t)
	mockRelRepo.AssertExpectations(t)
}

// Table-driven test for multiple scenarios
func TestCIHandler_CreateCI_TableDriven(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		requestBody    models.CI
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful creation",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				mockCIRepo := &MockCIRepository{}
				mockRelRepo := &MockRelationshipRepository{}
				mockAuditRepo := &MockAuditLogRepository{}
				
				mockCIRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CI")).Return(nil)
				mockAuditRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.AuditLog")).Return(nil)
				
				return mockCIRepo, mockRelRepo, mockAuditRepo
			},
			requestBody: models.CI{
				Name: "Test CI",
				Type: "server",
			},
			setupContext: func(req *http.Request) *http.Request {
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Validation error - missing name",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				return &MockCIRepository{}, &MockRelationshipRepository{}, &MockAuditLogRepository{}
			},
			requestBody: models.CI{
				Name: "",
				Type: "server",
			},
			setupContext: func(req *http.Request) *http.Request {
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Validation failed",
		},
		{
			name: "Repository error",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				mockCIRepo := &MockCIRepository{}
				mockRelRepo := &MockRelationshipRepository{}
				mockAuditRepo := &MockAuditLogRepository{}
				
				mockCIRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.CI")).Return(assert.AnError)
				
				return mockCIRepo, mockRelRepo, mockAuditRepo
			},
			requestBody: models.CI{
				Name: "Test CI",
				Type: "server",
			},
			setupContext: func(req *http.Request) *http.Request {
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create CI",
		},
		{
			name: "Unauthorized - missing context",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				return &MockCIRepository{}, &MockRelationshipRepository{}, &MockAuditLogRepository{}
			},
			requestBody: models.CI{
				Name: "Test CI",
				Type: "server",
			},
			setupContext: func(req *http.Request) *http.Request {
				// No authentication context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockCIRepo, mockRelRepo, mockAuditRepo := tt.setupMock()
			handler := NewCIHandler(mockCIRepo, mockRelRepo, mockAuditRepo)

			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/cis", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Setup context
			req = tt.setupContext(req)

			rr := httptest.NewRecorder()

			// Act
			handler.CreateCI(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			}

			// Verify mock expectations
			mockCIRepo.AssertExpectations(t)
			mockRelRepo.AssertExpectations(t)
			mockAuditRepo.AssertExpectations(t)
		})
	}
}
