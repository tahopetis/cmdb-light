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

// MockCIRepository is a mock implementation of CIRepository
type MockCIRepository struct {
	ci  *models.CI
	cis []*models.CI
	err error
}

func (m *MockCIRepository) Create(ctx context.Context, ci *models.CI) error {
	return m.err
}

func (m *MockCIRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CI, error) {
	if m.ci != nil && m.ci.ID == id {
		return m.ci, m.err
	}
	return nil, m.err
}

func (m *MockCIRepository) GetAll(ctx context.Context) ([]*models.CI, error) {
	return m.cis, m.err
}

func (m *MockCIRepository) Update(ctx context.Context, ci *models.CI) error {
	return m.err
}

func (m *MockCIRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

// MockRelationshipRepository is a mock implementation of RelationshipRepository
type MockRelationshipRepository struct {
	rels []*models.Relationship
	err  error
}

func (m *MockRelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	return m.err
}

func (m *MockRelationshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error) {
	return nil, m.err
}

func (m *MockRelationshipRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	return m.rels, m.err
}

func (m *MockRelationshipRepository) GetBySourceID(ctx context.Context, id uuid.UUID) ([]*models.Relationship, error) {
	return m.rels, m.err
}

func (m *MockRelationshipRepository) GetByTargetID(ctx context.Context, id uuid.UUID) ([]*models.Relationship, error) {
	return m.rels, m.err
}

func (m *MockRelationshipRepository) Update(ctx context.Context, rel *models.Relationship) error {
	return m.err
}

func (m *MockRelationshipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

// MockAuditLogRepository is a mock implementation of AuditLogRepository
type MockAuditLogRepository struct {
	err error
}

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return m.err
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	return nil, m.err
}

func (m *MockAuditLogRepository) GetAll(ctx context.Context) ([]*models.AuditLog, error) {
	return nil, m.err
}

func (m *MockAuditLogRepository) GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error) {
	return nil, m.err
}

func (m *MockAuditLogRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error) {
	return nil, m.err
}

func (m *MockAuditLogRepository) GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error) {
	return nil, m.err
}

func (m *MockAuditLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func TestCIHandler_CreateCI(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:        ciID,
		Name:      "Test CI",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	validCI := models.CI{
		Name: "Test CI",
		Type: "server",
	}

	invalidCI := models.CI{
		Name: "", // Missing name
		Type: "server",
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		requestBody    interface{}
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful CI creation",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: validCI,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid CI data",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: invalidCI,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name and type are required",
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: validCI,
			setupContext: func(req *http.Request) *http.Request {
				// No username in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "Invalid request body",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: "invalid-json",
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "Repository error",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: assert.AnError,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: validCI,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create CI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
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
			
			req, err := http.NewRequest("POST", "/cis", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Setup context
			req = tt.setupContext(req)
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			ciHandler.CreateCI(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.CI
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert CI in response
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, validCI.Name, response.Name)
				assert.Equal(t, validCI.Type, response.Type)
			}
		})
	}
}

func TestCIHandler_GetCI(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:        ciID,
		Name:      "Test CI",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		urlParams      map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful CI retrieval",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  testCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID format",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "CI not found",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: sqlx.ErrNotFound,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CI not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
			// Create request with URL parameters
			req, err := http.NewRequest("GET", "/cis/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/cis/{id}", ciHandler.GetCI).Methods("GET")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Serve request
			router.ServeHTTP(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.CI
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert CI in response
				assert.Equal(t, testCI.ID, response.ID)
				assert.Equal(t, testCI.Name, response.Name)
				assert.Equal(t, testCI.Type, response.Type)
			}
		})
	}
}

func TestCIHandler_GetAllCIs(t *testing.T) {
	// Setup test data
	ci1 := &models.CI{
		ID:        uuid.New(),
		Name:      "Test CI 1",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	ci2 := &models.CI{
		ID:        uuid.New(),
		Name:      "Test CI 2",
		Type:      "database",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	allCIs := []*models.CI{ci1, ci2}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		queryParams    map[string]string
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful CI retrieval",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					cis: allCIs,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful CI retrieval with pagination",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					cis: allCIs,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			queryParams: map[string]string{
				"page":  "1",
				"limit": "1",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Repository error",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: assert.AnError,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get CIs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
			// Create request with query parameters
			req, err := http.NewRequest("GET", "/cis", nil)
			require.NoError(t, err)
			
			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			ciHandler.GetAllCIs(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert data in response
				data := response["data"].([]interface{})
				assert.Equal(t, tt.expectedCount, len(data))
				
				// Assert pagination info
				pagination := response["pagination"].(map[string]interface{})
				assert.Equal(t, tt.queryParams["page"], pagination["page"])
				assert.Equal(t, tt.queryParams["limit"], pagination["limit"])
				assert.Equal(t, 2, int(pagination["total"].(float64)))
			}
		})
	}
}

func TestCIHandler_UpdateCI(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	existingCI := &models.CI{
		ID:        ciID,
		Name:      "Test CI",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	updatedCI := models.CI{
		Name: "Updated Test CI",
		Type: "database",
	}

	invalidCI := models.CI{
		Name: "", // Missing name
		Type: "database",
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		requestBody    interface{}
		urlParams      map[string]string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful CI update",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  existingCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: updatedCI,
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid CI data",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  existingCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: invalidCI,
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Name and type are required",
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  existingCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: updatedCI,
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// No username in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "Invalid ID format",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: updatedCI,
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "CI not found",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: sqlx.ErrNotFound,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			requestBody: updatedCI,
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CI not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
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
			
			req, err := http.NewRequest("PUT", "/cis/"+tt.urlParams["id"], bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Setup context
			req = tt.setupContext(req)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/cis/{id}", ciHandler.UpdateCI).Methods("PUT")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Serve request
			router.ServeHTTP(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.CI
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert CI in response
				assert.Equal(t, ciID, response.ID)
				assert.Equal(t, updatedCI.Name, response.Name)
				assert.Equal(t, updatedCI.Type, response.Type)
			}
		})
	}
}

func TestCIHandler_DeleteCI(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:        ciID,
		Name:      "Test CI",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		urlParams      map[string]string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful CI deletion",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  testCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  testCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// No username in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "Invalid ID format",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "CI not found",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: sqlx.ErrNotFound,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CI not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
			// Create request
			req, err := http.NewRequest("DELETE", "/cis/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Setup context
			req = tt.setupContext(req)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/cis/{id}", ciHandler.DeleteCI).Methods("DELETE")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Serve request
			router.ServeHTTP(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response map[string]string
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert message in response
				assert.Equal(t, "CI deleted successfully", response["message"])
			}
		})
	}
}

func TestCIHandler_GetCIGraph(t *testing.T) {
	// Setup test data
	ciID := uuid.New()
	testCI := &models.CI{
		ID:        ciID,
		Name:      "Test CI",
		Type:      "server",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	sourceRel := &models.Relationship{
		ID:       uuid.New(),
		SourceID: ciID,
		TargetID: uuid.New(),
		Type:     "depends_on",
	}
	
	targetRel := &models.Relationship{
		ID:       uuid.New(),
		SourceID: uuid.New(),
		TargetID: ciID,
		Type:     "hosts",
	}
	
	sourceRels := []*models.Relationship{sourceRel}
	targetRels := []*models.Relationship{targetRel}

	tests := []struct {
		name           string
		setupMock      func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository)
		urlParams      map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful CI graph retrieval",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					ci:  testCI,
					err: nil,
				}
				relRepo := &MockRelationshipRepository{
					rels: sourceRels,
					err:  nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID format",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "CI not found",
			setupMock: func() (*MockCIRepository, *MockRelationshipRepository, *MockAuditLogRepository) {
				ciRepo := &MockCIRepository{
					err: sqlx.ErrNotFound,
				}
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return ciRepo, relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": ciID.String(),
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "CI not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			ciRepo, relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			ciHandler := NewCIHandler(ciRepo, relRepo, auditRepo)
			
			// Create request
			req, err := http.NewRequest("GET", "/cis/"+tt.urlParams["id"]+"/graph", nil)
			require.NoError(t, err)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/cis/{id}/graph", ciHandler.GetCIGraph).Methods("GET")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Serve request
			router.ServeHTTP(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response map[string]interface{}
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert node in response
				node := response["node"].(map[string]interface{})
				assert.Equal(t, testCI.ID.String(), node["id"])
				assert.Equal(t, testCI.Name, node["name"])
				assert.Equal(t, testCI.Type, node["type"])
				
				// Assert relationships in response
				relationships := response["relationships"].(map[string]interface{})
				outgoing := relationships["outgoing"].([]interface{})
				incoming := relationships["incoming"].([]interface{})
				
				// Note: In this test, we're only checking that the relationships are present,
				// not their exact content, since the mock returns the same relationships for both
				// source and target queries.
				assert.Equal(t, 1, len(outgoing))
				assert.Equal(t, 1, len(incoming))
			}
		})
	}
}