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
	"github.com/stretchr/testify/require"
)

// MockRelationshipRepository is a mock implementation of RelationshipRepository
type MockRelationshipRepository struct {
	rel  *models.Relationship
	rels []*models.Relationship
	err  error
}

func (m *MockRelationshipRepository) Create(ctx context.Context, rel *models.Relationship) error {
	return m.err
}

func (m *MockRelationshipRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error) {
	if m.rel != nil && m.rel.ID == id {
		return m.rel, m.err
	}
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

func TestRelationshipHandler_CreateRelationship(t *testing.T) {
	// Setup test data
	relID := uuid.New()
	sourceID := uuid.New()
	targetID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relID,
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	validRelationship := models.Relationship{
		SourceID: sourceID,
		TargetID: targetID,
		Type:     "depends_on",
	}

	invalidRelationship := models.Relationship{
		SourceID: sourceID,
		TargetID: targetID,
		Type:     "", // Missing type
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockRelationshipRepository, *MockAuditLogRepository)
		requestBody    interface{}
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful relationship creation",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: validRelationship,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid relationship data",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: invalidRelationship,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Source ID, target ID, and type are required",
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: validRelationship,
			setupContext: func(req *http.Request) *http.Request {
				// No username in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "Invalid request body",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
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
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: assert.AnError,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: validRelationship,
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to create relationship",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			relHandler := NewRelationshipHandler(relRepo, auditRepo)
			
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
			
			req, err := http.NewRequest("POST", "/relationships", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Setup context
			req = tt.setupContext(req)
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			relHandler.CreateRelationship(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.Relationship
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert relationship in response
				assert.NotEmpty(t, response.ID)
				assert.Equal(t, validRelationship.SourceID, response.SourceID)
				assert.Equal(t, validRelationship.TargetID, response.TargetID)
				assert.Equal(t, validRelationship.Type, response.Type)
			}
		})
	}
}

func TestRelationshipHandler_GetRelationship(t *testing.T) {
	// Setup test data
	relID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relID,
		SourceID:  uuid.New(),
		TargetID:  uuid.New(),
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockRelationshipRepository, *MockAuditLogRepository)
		urlParams      map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful relationship retrieval",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: testRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": relID.String(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID format",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "Relationship not found",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: sqlx.ErrNotFound,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": relID.String(),
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Relationship not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			relHandler := NewRelationshipHandler(relRepo, auditRepo)
			
			// Create request with URL parameters
			req, err := http.NewRequest("GET", "/relationships/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/relationships/{id}", relHandler.GetRelationship).Methods("GET")
			
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
				var response models.Relationship
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert relationship in response
				assert.Equal(t, testRelationship.ID, response.ID)
				assert.Equal(t, testRelationship.SourceID, response.SourceID)
				assert.Equal(t, testRelationship.TargetID, response.TargetID)
				assert.Equal(t, testRelationship.Type, response.Type)
			}
		})
	}
}

func TestRelationshipHandler_GetAllRelationships(t *testing.T) {
	// Setup test data
	rel1 := &models.Relationship{
		ID:        uuid.New(),
		SourceID:  uuid.New(),
		TargetID:  uuid.New(),
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}
	
	rel2 := &models.Relationship{
		ID:        uuid.New(),
		SourceID:  uuid.New(),
		TargetID:  uuid.New(),
		Type:      "hosts",
		CreatedAt: time.Now(),
	}
	
	allRelationships := []*models.Relationship{rel1, rel2}

	tests := []struct {
		name           string
		setupMock      func() (*MockRelationshipRepository, *MockAuditLogRepository)
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful relationship retrieval",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rels: allRelationships,
					err:  nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Repository error",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: assert.AnError,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get relationships",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			relHandler := NewRelationshipHandler(relRepo, auditRepo)
			
			// Create request
			req, err := http.NewRequest("GET", "/relationships", nil)
			require.NoError(t, err)
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			relHandler.GetAllRelationships(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response []*models.Relationship
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert data in response
				assert.Equal(t, tt.expectedCount, len(response))
			}
		})
	}
}

func TestRelationshipHandler_UpdateRelationship(t *testing.T) {
	// Setup test data
	relID := uuid.New()
	existingRelationship := &models.Relationship{
		ID:        relID,
		SourceID:  uuid.New(),
		TargetID:  uuid.New(),
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	updatedRelationship := models.Relationship{
		SourceID: uuid.New(),
		TargetID: uuid.New(),
		Type:     "hosts",
	}

	invalidRelationship := models.Relationship{
		SourceID: uuid.New(),
		TargetID: uuid.New(),
		Type:     "", // Missing type
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockRelationshipRepository, *MockAuditLogRepository)
		requestBody    interface{}
		urlParams      map[string]string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful relationship update",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: existingRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: updatedRelationship,
			urlParams: map[string]string{
				"id": relID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid relationship data",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: existingRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: invalidRelationship,
			urlParams: map[string]string{
				"id": relID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Source ID, target ID, and type are required",
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: existingRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: updatedRelationship,
			urlParams: map[string]string{
				"id": relID.String(),
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
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: updatedRelationship,
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
			name: "Relationship not found",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: sqlx.ErrNotFound,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			requestBody: updatedRelationship,
			urlParams: map[string]string{
				"id": relID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Relationship not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			relHandler := NewRelationshipHandler(relRepo, auditRepo)
			
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
			
			req, err := http.NewRequest("PUT", "/relationships/"+tt.urlParams["id"], bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Setup context
			req = tt.setupContext(req)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/relationships/{id}", relHandler.UpdateRelationship).Methods("PUT")
			
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
				var response models.Relationship
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert relationship in response
				assert.Equal(t, relID, response.ID)
				assert.Equal(t, updatedRelationship.SourceID, response.SourceID)
				assert.Equal(t, updatedRelationship.TargetID, response.TargetID)
				assert.Equal(t, updatedRelationship.Type, response.Type)
			}
		})
	}
}

func TestRelationshipHandler_DeleteRelationship(t *testing.T) {
	// Setup test data
	relID := uuid.New()
	testRelationship := &models.Relationship{
		ID:        relID,
		SourceID:  uuid.New(),
		TargetID:  uuid.New(),
		Type:      "depends_on",
		CreatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockRelationshipRepository, *MockAuditLogRepository)
		urlParams      map[string]string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful relationship deletion",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: testRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": relID.String(),
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
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					rel: testRelationship,
					err: nil,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": relID.String(),
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
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
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
			name: "Relationship not found",
			setupMock: func() (*MockRelationshipRepository, *MockAuditLogRepository) {
				relRepo := &MockRelationshipRepository{
					err: sqlx.ErrNotFound,
				}
				auditRepo := &MockAuditLogRepository{}
				return relRepo, auditRepo
			},
			urlParams: map[string]string{
				"id": relID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with username
				ctx := context.WithValue(req.Context(), middleware.UsernameKey, "testuser")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Relationship not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			relRepo, auditRepo := tt.setupMock()
			
			// Create handler
			relHandler := NewRelationshipHandler(relRepo, auditRepo)
			
			// Create request
			req, err := http.NewRequest("DELETE", "/relationships/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Setup context
			req = tt.setupContext(req)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/relationships/{id}", relHandler.DeleteRelationship).Methods("DELETE")
			
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
				assert.Equal(t, "Relationship deleted successfully", response["message"])
			}
		})
	}
}