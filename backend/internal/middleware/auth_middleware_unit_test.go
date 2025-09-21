package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockJWTManager is a mock implementation of JWTManager interface
type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateAccessToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) Verify(accessToken string) (*auth.UserClaims, error) {
	args := m.Called(accessToken)
	return args.Get(0).(*auth.UserClaims), args.Error(1)
}

func (m *MockJWTManager) HashRefreshToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) VerifyRefreshTokenHash(refreshToken, hash string) bool {
	args := m.Called(refreshToken, hash)
	return args.Bool(0)
}

// Helper functions to create test data
func createTestUser(username, role string) *models.User {
	return &models.User{
		ID:        uuid.New(),
		Username:  username,
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func createTestToken(userID uuid.UUID, username, role string) string {
	// Create a simple token for testing
	return "test-token-" + userID.String()
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "user")
	testClaims := &auth.UserClaims{
		UserID:   testUser.ID,
		Username: testUser.Username,
		Role:     testUser.Role,
		RegisteredClaims: auth.RegisteredClaims{
			ExpiresAt: auth.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	mockJWTManager := &MockJWTManager{}
	mockJWTManager.On("Verify", "valid-token").Return(testClaims, nil)

	middleware := AuthMiddleware(mockJWTManager)

	// Create a handler that the middleware will wrap
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if user info is in context
		userID := r.Context().Value("user_id")
		username := r.Context().Value("username")
		role := r.Context().Value("role")

		assert.NotNil(t, userID)
		assert.Equal(t, testUser.ID, userID)
		assert.NotNil(t, username)
		assert.Equal(t, testUser.Username, username)
		assert.NotNil(t, role)
		assert.Equal(t, testUser.Role, role)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("protected resource"))
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer valid-token")

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "protected resource", rr.Body.String())

	// Verify mock expectations
	mockJWTManager.AssertExpectations(t)
}

func TestAuthMiddleware_MissingToken(t *testing.T) {
	// Arrange
	mockJWTManager := &MockJWTManager{}
	middleware := AuthMiddleware(mockJWTManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	require.NoError(t, err)
	// No Authorization header

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Missing authorization token", response.Message)
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	// Arrange
	mockJWTManager := &MockJWTManager{}
	middleware := AuthMiddleware(mockJWTManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "InvalidFormat token") // Missing "Bearer"

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid authorization format", response.Message)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	// Arrange
	mockJWTManager := &MockJWTManager{}
	mockJWTManager.On("Verify", "invalid-token").Return(nil, assert.AnError)

	middleware := AuthMiddleware(mockJWTManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid token", response.Message)

	// Verify mock expectations
	mockJWTManager.AssertExpectations(t)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "user")
	testClaims := &auth.UserClaims{
		UserID:   testUser.ID,
		Username: testUser.Username,
		Role:     testUser.Role,
		RegisteredClaims: auth.RegisteredClaims{
			ExpiresAt: auth.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired
		},
	}

	mockJWTManager := &MockJWTManager{}
	mockJWTManager.On("Verify", "expired-token").Return(testClaims, nil)

	middleware := AuthMiddleware(mockJWTManager)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/protected", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer expired-token")

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Token is expired", response.Message)

	// Verify mock expectations
	mockJWTManager.AssertExpectations(t)
}

func TestRBACMiddleware_AdminRole(t *testing.T) {
	// Arrange
	testUser := createTestUser("admin", "admin")

	// Create context with user info
	ctx := context.WithValue(context.Background(), "user_id", testUser.ID)
	ctx = context.WithValue(ctx, "username", testUser.Username)
	ctx = context.WithValue(ctx, "role", testUser.Role)

	middleware := RBACMiddleware("admin")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("admin resource"))
	})

	req, err := http.NewRequest("GET", "/admin", nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "admin resource", rr.Body.String())
}

func TestRBACMiddleware_ViewerRole(t *testing.T) {
	// Arrange
	testUser := createTestUser("viewer", "viewer")

	// Create context with user info
	ctx := context.WithValue(context.Background(), "user_id", testUser.ID)
	ctx = context.WithValue(ctx, "username", testUser.Username)
	ctx = context.WithValue(ctx, "role", testUser.Role)

	middleware := RBACMiddleware("admin", "viewer")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("viewer resource"))
	})

	req, err := http.NewRequest("GET", "/viewer", nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "viewer resource", rr.Body.String())
}

func TestRBACMiddleware_InsufficientRole(t *testing.T) {
	// Arrange
	testUser := createTestUser("viewer", "viewer")

	// Create context with user info
	ctx := context.WithValue(context.Background(), "user_id", testUser.ID)
	ctx = context.WithValue(ctx, "username", testUser.Username)
	ctx = context.WithValue(ctx, "role", testUser.Role)

	middleware := RBACMiddleware("admin")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/admin", nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "FORBIDDEN", response.Code)
	assert.Equal(t, "Insufficient permissions", response.Message)
}

func TestRBACMiddleware_MissingRole(t *testing.T) {
	// Arrange
	// Create context without role
	ctx := context.Background()

	middleware := RBACMiddleware("admin")

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("GET", "/admin", nil)
	require.NoError(t, err)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusForbidden, rr.Code)

	var response models.ErrorResponse
	err = models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.Equal(t, "FORBIDDEN", response.Code)
	assert.Equal(t, "Insufficient permissions", response.Message)
}

func TestCORSMiddleware(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		CORSAllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		CORSAllowCredentials: true,
	}

	middleware := CORS(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("CORS test"))
	})

	tests := []struct {
		name           string
		origin         string
		method         string
		expectedOrigin string
	}{
		{
			name:           "Allowed origin",
			origin:         "http://localhost:3000",
			method:         "GET",
			expectedOrigin: "http://localhost:3000",
		},
		{
			name:           "Another allowed origin",
			origin:         "https://example.com",
			method:         "POST",
			expectedOrigin: "https://example.com",
		},
		{
			name:           "Disallowed origin",
			origin:         "http://malicious.com",
			method:         "GET",
			expectedOrigin: "", // Should not be set
		},
		{
			name:           "No origin",
			origin:         "",
			method:         "GET",
			expectedOrigin: "", // Should not be set
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, "/test", nil)
			require.NoError(t, err)
			
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}

			rr := httptest.NewRecorder()

			// Act
			wrappedHandler := middleware(handler)
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			if tt.expectedOrigin != "" {
				assert.Equal(t, tt.expectedOrigin, rr.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, rr.Header().Get("Access-Control-Allow-Origin"))
			}

			// Check other CORS headers
			assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
			assert.Equal(t, "Accept, Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
			assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
			assert.Equal(t, "86400", rr.Header().Get("Access-Control-Max-Age"))
		})
	}
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	// Arrange
	cfg := &config.Config{
		CORSAllowedOrigins:   []string{"http://localhost:3000"},
		CORSAllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		CORSAllowedHeaders:   []string{"Accept", "Content-Type", "Authorization"},
		CORSAllowCredentials: true,
	}

	middleware := CORS(cfg)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("should not reach here"))
	})

	req, err := http.NewRequest("OPTIONS", "/test", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "http://localhost:3000", rr.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET, POST, PUT, DELETE, OPTIONS", rr.Header().Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "Accept, Content-Type, Authorization", rr.Header().Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", rr.Header().Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "86400", rr.Header().Get("Access-Control-Max-Age"))

	// Body should be empty for preflight requests
	assert.Empty(t, rr.Body.String())
}

func TestRecoveryMiddleware(t *testing.T) {
	// Arrange
	middleware := Recovery()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a panic
		panic("test panic")
	})

	req, err := http.NewRequest("GET", "/test", nil)
	require.NoError(t, err)

	rr := httptest.NewRecorder()

	// Act
	wrappedHandler := middleware(handler)
	wrappedHandler.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Equal(t, "Internal Server Error\n", rr.Body.String())
}

// Table-driven test for RBAC middleware
func TestRBACMiddleware_TableDriven(t *testing.T) {
	tests := []struct {
		name          string
		userRole      string
		allowedRoles  []string
		expectedStatus int
	}{
		{
			name:          "Admin accessing admin resource",
			userRole:      "admin",
			allowedRoles:  []string{"admin"},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Viewer accessing viewer resource",
			userRole:      "viewer",
			allowedRoles:  []string{"viewer"},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Admin accessing viewer resource",
			userRole:      "admin",
			allowedRoles:  []string{"viewer"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:          "Viewer accessing admin resource",
			userRole:      "viewer",
			allowedRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:          "Admin accessing multi-role resource",
			userRole:      "admin",
			allowedRoles:  []string{"admin", "viewer"},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "Viewer accessing multi-role resource",
			userRole:      "viewer",
			allowedRoles:  []string{"admin", "viewer"},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "User accessing admin resource",
			userRole:      "user",
			allowedRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:          "No role accessing protected resource",
			userRole:      "",
			allowedRoles:  []string{"admin"},
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			testUser := createTestUser("testuser", tt.userRole)

			// Create context with user info
			ctx := context.WithValue(context.Background(), "user_id", testUser.ID)
			ctx = context.WithValue(ctx, "username", testUser.Username)
			if tt.userRole != "" {
				ctx = context.WithValue(ctx, "role", tt.userRole)
			}

			middleware := RBACMiddleware(tt.allowedRoles...)

			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("protected resource"))
			})

			req, err := http.NewRequest("GET", "/protected", nil)
			require.NoError(t, err)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			// Act
			wrappedHandler := middleware(handler)
			wrappedHandler.ServeHTTP(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "protected resource", rr.Body.String())
			} else {
				var response models.ErrorResponse
				err := models.UnmarshalErrorResponse(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, "FORBIDDEN", response.Code)
			}
		})
	}
}
