package handlers

import (
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

// MockAuditLogRepository is a mock implementation of AuditLogRepository
type MockAuditLogRepository struct {
	auditLog *models.AuditLog
	auditLogs []*models.AuditLog
	err error
}

func (m *MockAuditLogRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	return m.err
}

func (m *MockAuditLogRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	if m.auditLog != nil && m.auditLog.ID == id {
		return m.auditLog, m.err
	}
	return nil, m.err
}

func (m *MockAuditLogRepository) GetAll(ctx context.Context) ([]*models.AuditLog, error) {
	return m.auditLogs, m.err
}

func (m *MockAuditLogRepository) GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error) {
	return m.auditLogs, m.err
}

func (m *MockAuditLogRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error) {
	return m.auditLogs, m.err
}

func (m *MockAuditLogRepository) GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error) {
	return m.auditLogs, m.err
}

func (m *MockAuditLogRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.err
}

func TestAuditLogHandler_GetAuditLog(t *testing.T) {
	// Setup test data
	auditLogID := uuid.New()
	testAuditLog := &models.AuditLog{
		ID:         auditLogID,
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI", "type": "server"},
	}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		urlParams      map[string]string
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful audit log retrieval",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLog: testAuditLog,
					err:      nil,
				}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid ID format",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "Audit log not found",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: sqlx.ErrNotFound,
				}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Audit log not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request with URL parameters
			req, err := http.NewRequest("GET", "/audit-logs/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/audit-logs/{id}", auditHandler.GetAuditLog).Methods("GET")
			
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
				var response models.AuditLog
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert audit log in response
				assert.Equal(t, testAuditLog.ID, response.ID)
				assert.Equal(t, testAuditLog.EntityType, response.EntityType)
				assert.Equal(t, testAuditLog.EntityID, response.EntityID)
				assert.Equal(t, testAuditLog.Action, response.Action)
				assert.Equal(t, testAuditLog.ChangedBy, response.ChangedBy)
			}
		})
	}
}

func TestAuditLogHandler_GetAllAuditLogs(t *testing.T) {
	// Setup test data
	auditLog1 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
	}
	
	auditLog2 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "update",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 2", "type": "database"},
	}
	
	allAuditLogs := []*models.AuditLog{auditLog1, auditLog2}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		queryParams    map[string]string
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful audit log retrieval",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval with pagination",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			queryParams: map[string]string{
				"page":  "1",
				"limit": "1",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Successful audit log retrieval with entity type filter",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			queryParams: map[string]string{
				"entity_type": "configuration_item",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval with entity ID filter",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			queryParams: map[string]string{
				"entity_id": auditLog1.EntityID.String(),
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval with changed by filter",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			queryParams: map[string]string{
				"changed_by": "testuser",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Invalid entity ID format",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			queryParams: map[string]string{
				"entity_id": "invalid-uuid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid entity ID format",
		},
		{
			name: "Repository error",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: assert.AnError,
				}
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get audit logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request with query parameters
			req, err := http.NewRequest("GET", "/audit-logs", nil)
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
			auditHandler.GetAllAuditLogs(rr, req)
			
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
				if tt.queryParams["page"] != "" {
					assert.Equal(t, tt.queryParams["page"], pagination["page"])
				} else {
					assert.Equal(t, 1, int(pagination["page"].(float64)))
				}
				
				if tt.queryParams["limit"] != "" {
					assert.Equal(t, tt.queryParams["limit"], pagination["limit"])
				} else {
					assert.Equal(t, 10, int(pagination["limit"].(float64)))
				}
				
				assert.Equal(t, 2, int(pagination["total"].(float64)))
			}
		})
	}
}

func TestAuditLogHandler_GetAuditLogsByEntityType(t *testing.T) {
	// Setup test data
	auditLog1 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
	}
	
	auditLog2 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "relationship",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"source_id": uuid.New(), "target_id": uuid.New(), "type": "depends_on"},
	}
	
	allAuditLogs := []*models.AuditLog{auditLog1, auditLog2}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		urlParams      map[string]string
		queryParams    map[string]string
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful audit log retrieval by entity type",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"entity_type": "configuration_item",
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval by entity type with pagination",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"entity_type": "configuration_item",
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
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: assert.AnError,
				}
			},
			urlParams: map[string]string{
				"entity_type": "configuration_item",
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get audit logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request with URL and query parameters
			req, err := http.NewRequest("GET", "/audit-logs/entity-type/"+tt.urlParams["entity_type"], nil)
			require.NoError(t, err)
			
			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/audit-logs/entity-type/{entity_type}", auditHandler.GetAuditLogsByEntityType).Methods("GET")
			
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
				
				// Assert data in response
				data := response["data"].([]interface{})
				assert.Equal(t, tt.expectedCount, len(data))
				
				// Assert pagination info
				pagination := response["pagination"].(map[string]interface{})
				if tt.queryParams["page"] != "" {
					assert.Equal(t, tt.queryParams["page"], pagination["page"])
				} else {
					assert.Equal(t, 1, int(pagination["page"].(float64)))
				}
				
				if tt.queryParams["limit"] != "" {
					assert.Equal(t, tt.queryParams["limit"], pagination["limit"])
				} else {
					assert.Equal(t, 10, int(pagination["limit"].(float64)))
				}
				
				assert.Equal(t, 2, int(pagination["total"].(float64)))
			}
		})
	}
}

func TestAuditLogHandler_GetAuditLogsByEntityID(t *testing.T) {
	// Setup test data
	entityID := uuid.New()
	auditLog1 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   entityID,
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
	}
	
	auditLog2 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   entityID,
		Action:     "update",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 1", "type": "database"},
	}
	
	allAuditLogs := []*models.AuditLog{auditLog1, auditLog2}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		urlParams      map[string]string
		queryParams    map[string]string
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful audit log retrieval by entity ID",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"entity_id": entityID.String(),
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval by entity ID with pagination",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"entity_id": entityID.String(),
			},
			queryParams: map[string]string{
				"page":  "1",
				"limit": "1",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Invalid entity ID format",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			urlParams: map[string]string{
				"entity_id": "invalid-uuid",
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid entity ID format",
		},
		{
			name: "Repository error",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: assert.AnError,
				}
			},
			urlParams: map[string]string{
				"entity_id": entityID.String(),
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get audit logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request with URL and query parameters
			req, err := http.NewRequest("GET", "/audit-logs/entity-id/"+tt.urlParams["entity_id"], nil)
			require.NoError(t, err)
			
			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/audit-logs/entity-id/{entity_id}", auditHandler.GetAuditLogsByEntityID).Methods("GET")
			
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
				
				// Assert data in response
				data := response["data"].([]interface{})
				assert.Equal(t, tt.expectedCount, len(data))
				
				// Assert pagination info
				pagination := response["pagination"].(map[string]interface{})
				if tt.queryParams["page"] != "" {
					assert.Equal(t, tt.queryParams["page"], pagination["page"])
				} else {
					assert.Equal(t, 1, int(pagination["page"].(float64)))
				}
				
				if tt.queryParams["limit"] != "" {
					assert.Equal(t, tt.queryParams["limit"], pagination["limit"])
				} else {
					assert.Equal(t, 10, int(pagination["limit"].(float64)))
				}
				
				assert.Equal(t, 2, int(pagination["total"].(float64)))
			}
		})
	}
}

func TestAuditLogHandler_GetAuditLogsByChangedBy(t *testing.T) {
	// Setup test data
	changedBy := "testuser"
	auditLog1 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  changedBy,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI 1", "type": "server"},
	}
	
	auditLog2 := &models.AuditLog{
		ID:         uuid.New(),
		EntityType: "relationship",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  changedBy,
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"source_id": uuid.New(), "target_id": uuid.New(), "type": "depends_on"},
	}
	
	allAuditLogs := []*models.AuditLog{auditLog1, auditLog2}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		urlParams      map[string]string
		queryParams    map[string]string
		expectedStatus int
		expectedError  string
		expectedCount  int
	}{
		{
			name: "Successful audit log retrieval by changed by",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"changed_by": changedBy,
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "Successful audit log retrieval by changed by with pagination",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					auditLogs: allAuditLogs,
					err:       nil,
				}
			},
			urlParams: map[string]string{
				"changed_by": changedBy,
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
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: assert.AnError,
				}
			},
			urlParams: map[string]string{
				"changed_by": changedBy,
			},
			queryParams:    map[string]string{},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to get audit logs",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request with URL and query parameters
			req, err := http.NewRequest("GET", "/audit-logs/changed-by/"+tt.urlParams["changed_by"], nil)
			require.NoError(t, err)
			
			// Add query parameters
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/audit-logs/changed-by/{changed_by}", auditHandler.GetAuditLogsByChangedBy).Methods("GET")
			
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
				
				// Assert data in response
				data := response["data"].([]interface{})
				assert.Equal(t, tt.expectedCount, len(data))
				
				// Assert pagination info
				pagination := response["pagination"].(map[string]interface{})
				if tt.queryParams["page"] != "" {
					assert.Equal(t, tt.queryParams["page"], pagination["page"])
				} else {
					assert.Equal(t, 1, int(pagination["page"].(float64)))
				}
				
				if tt.queryParams["limit"] != "" {
					assert.Equal(t, tt.queryParams["limit"], pagination["limit"])
				} else {
					assert.Equal(t, 10, int(pagination["limit"].(float64)))
				}
				
				assert.Equal(t, 2, int(pagination["total"].(float64)))
			}
		})
	}
}

func TestAuditLogHandler_DeleteAuditLog(t *testing.T) {
	// Setup test data
	auditLogID := uuid.New()
	testAuditLog := &models.AuditLog{
		ID:         auditLogID,
		EntityType: "configuration_item",
		EntityID:   uuid.New(),
		Action:     "create",
		ChangedBy:  "testuser",
		ChangedAt:  time.Now(),
		Details:    models.JSONBMap{"name": "Test CI", "type": "server"},
	}

	tests := []struct {
		name           string
		setupMock      func() *MockAuditLogRepository
		urlParams      map[string]string
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful audit log deletion by admin",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: nil,
				}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with admin role
				ctx := context.WithValue(req.Context(), middleware.UserRoleKey, "admin")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not authenticated",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// No role in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "Insufficient permissions",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with user role (not admin)
				ctx := context.WithValue(req.Context(), middleware.UserRoleKey, "user")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusForbidden,
			expectedError:  "Insufficient permissions",
		},
		{
			name: "Invalid ID format",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{}
			},
			urlParams: map[string]string{
				"id": "invalid-uuid",
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with admin role
				ctx := context.WithValue(req.Context(), middleware.UserRoleKey, "admin")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid ID format",
		},
		{
			name: "Repository error",
			setupMock: func() *MockAuditLogRepository {
				return &MockAuditLogRepository{
					err: assert.AnError,
				}
			},
			urlParams: map[string]string{
				"id": auditLogID.String(),
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with admin role
				ctx := context.WithValue(req.Context(), middleware.UserRoleKey, "admin")
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "Failed to delete audit log",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			auditRepo := tt.setupMock()
			
			// Create handler
			auditHandler := NewAuditLogHandler(auditRepo)
			
			// Create request
			req, err := http.NewRequest("DELETE", "/audit-logs/"+tt.urlParams["id"], nil)
			require.NoError(t, err)
			
			// Setup context
			req = tt.setupContext(req)
			
			// Set up mux router with parameters
			router := mux.NewRouter()
			router.HandleFunc("/audit-logs/{id}", auditHandler.DeleteAuditLog).Methods("DELETE")
			
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
				assert.Equal(t, "Audit log deleted successfully", response["message"])
			}
		})
	}
}