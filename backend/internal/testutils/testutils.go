package testutils

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
)

// MockDB is a mock implementation of sqlx.DB for testing
type MockDB struct {
	ShouldError bool
	ErrorMsg     string
}

// NewMockDB creates a new MockDB
func NewMockDB(shouldError bool, errorMsg string) *MockDB {
	return &MockDB{
		ShouldError: shouldError,
		ErrorMsg:     errorMsg,
	}
}

// GetContext is a helper function to create a context with a user ID
func GetContext(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), "user_id", userID)
}

// GetAuthenticatedContext is a helper function to create a context with authenticated user
func GetAuthenticatedContext(username, role string) context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.UsernameKey, username)
	ctx = context.WithValue(ctx, middleware.UserRoleKey, role)
	return ctx
}

// CreateTestUser creates a test user for testing purposes
func CreateTestUser(username, password, role string) *models.User {
	userID := uuid.New()
	hashedPassword, _ := auth.NewPasswordManager().HashPassword(password)
	
	return &models.User{
		ID:           userID,
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateTestCI creates a test CI for testing purposes
func CreateTestCI(name, ciType string) *models.CI {
	return &models.CI{
		ID:        uuid.New(),
		Name:      name,
		Type:      ciType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestRelationship creates a test relationship for testing purposes
func CreateTestRelationship(sourceID, targetID uuid.UUID, relType string) *models.Relationship {
	return &models.Relationship{
		ID:        uuid.New(),
		SourceID:  sourceID,
		TargetID:  targetID,
		Type:      relType,
		CreatedAt: time.Now(),
	}
}

// CreateTestAuditLog creates a test audit log for testing purposes
func CreateTestAuditLog(entityType string, entityID uuid.UUID, action, changedBy string) *models.AuditLog {
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

// MakeTestRequest is a helper function to create an HTTP request for testing
func MakeTestRequest(method, url string, body interface{}, headers map[string]string) (*http.Request, error) {
	var reqBody []byte
	var err error
	
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = []byte(v)
		default:
			reqBody, err = json.Marshal(body)
			if err != nil {
				return nil, err
			}
		}
	}
	
	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}
	
	// Set headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	
	// Set default content type if not provided
	if _, ok := headers["Content-Type"]; !ok && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req, nil
}

// ExecuteTestRequest is a helper function to execute an HTTP request and return the response
func ExecuteTestRequest(handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// ParseJSONResponse is a helper function to parse a JSON response
func ParseJSONResponse(rr *httptest.ResponseRecorder, target interface{}) error {
	return json.Unmarshal(rr.Body.Bytes(), target)
}

// AssertErrorResponse is a helper function to assert that a response contains an error
func AssertErrorResponse(t require.TestingT, rr *httptest.ResponseRecorder, expectedStatusCode int, expectedErrorMsg string) {
	require.Equal(t, expectedStatusCode, rr.Code)
	require.Contains(t, rr.Body.String(), expectedErrorMsg)
}

// AssertSuccessResponse is a helper function to assert that a response is successful
func AssertSuccessResponse(t require.TestingT, rr *httptest.ResponseRecorder, expectedStatusCode int, target interface{}) {
	require.Equal(t, expectedStatusCode, rr.Code)
	require.NoError(t, ParseJSONResponse(rr, target))
}