package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/testutils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	user *models.User
	err  error
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	return nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return m.user, m.err
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	return m.user, m.err
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	if m.user != nil {
		return []*models.User{m.user}, nil
	}
	return []*models.User{}, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestAuthHandler_Login(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	password := "test-password-123"
	hashedPassword, _ := auth.NewPasswordManager().HashPassword(password)
	
	testUser := &models.User{
		ID:           userID,
		Username:     "testuser",
		PasswordHash: hashedPassword,
		Role:         "user",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	validLoginReq := models.LoginRequest{
		Username: "testuser",
		Password: password,
	}

	invalidLoginReq := models.LoginRequest{
		Username: "testuser",
		Password: "wrong-password",
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager)
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful login",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{
					user: testUser,
					err:  nil,
				}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			requestBody:    validLoginReq,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid password",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{
					user: testUser,
					err:  nil,
				}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			requestBody:    invalidLoginReq,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid username or password",
		},
		{
			name: "User not found",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{
					user: nil,
					err:  sqlx.ErrNotFound,
				}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			requestBody:    validLoginReq,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid username or password",
		},
		{
			name: "Invalid request body",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			requestBody:    "invalid-json",
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid request body",
		},
		{
			name: "Missing username or password",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			requestBody:    models.LoginRequest{Username: "", Password: "password"},
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo, jwtManager, passwordManager := tt.setupMock()
			
			// Create handler
			authHandler := NewAuthHandler(userRepo, jwtManager, passwordManager)
			
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
			
			req, err := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			authHandler.Login(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.LoginResponse
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert token and user in response
				assert.NotEmpty(t, response.Token)
				assert.Equal(t, testUser.ID, response.User.ID)
				assert.Equal(t, testUser.Username, response.User.Username)
				assert.Equal(t, testUser.Role, response.User.Role)
			}
		})
	}
}

func TestAuthHandler_ValidateToken(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	testUser := &models.User{
		ID:        userID,
		Username:  "testuser",
		Role:      "user",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tests := []struct {
		name           string
		setupMock      func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager)
		setupContext   func(*http.Request) *http.Request
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful token validation",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{
					user: testUser,
					err:  nil,
				}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with user ID
				ctx := context.WithValue(req.Context(), "user_id", userID)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not authenticated",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			setupContext: func(req *http.Request) *http.Request {
				// No user ID in context
				return req
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not authenticated",
		},
		{
			name: "User not found",
			setupMock: func() (*MockUserRepository, *auth.JWTManager, *auth.PasswordManager) {
				userRepo := &MockUserRepository{
					user: nil,
					err:  sqlx.ErrNotFound,
				}
				jwtManager := auth.NewJWTManager("test-secret-key", time.Hour)
				passwordManager := auth.NewPasswordManager()
				return userRepo, jwtManager, passwordManager
			},
			setupContext: func(req *http.Request) *http.Request {
				// Mock context with user ID
				ctx := context.WithValue(req.Context(), "user_id", userID)
				return req.WithContext(ctx)
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "User not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			userRepo, jwtManager, passwordManager := tt.setupMock()
			
			// Create handler
			authHandler := NewAuthHandler(userRepo, jwtManager, passwordManager)
			
			// Create request
			req, err := http.NewRequest("GET", "/auth/validate", nil)
			require.NoError(t, err)
			
			// Setup context
			req = tt.setupContext(req)
			
			// Create response recorder
			rr := httptest.NewRecorder()
			
			// Call handler
			authHandler.ValidateToken(rr, req)
			
			// Assert status code
			assert.Equal(t, tt.expectedStatus, rr.Code)
			
			// Assert error message if expected
			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else {
				// Parse response body
				var response models.User
				err = json.Unmarshal(rr.Body.Bytes(), &response)
				require.NoError(t, err)
				
				// Assert user in response
				assert.Equal(t, testUser.ID, response.ID)
				assert.Equal(t, testUser.Username, response.Username)
				assert.Equal(t, testUser.Role, response.Role)
			}
		})
	}
}

func TestGetUserIDFromContext(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	
	tests := []struct {
		name          string
		setupContext  func() context.Context
		expectedID    uuid.UUID
		expectedFound bool
	}{
		{
			name: "User ID in context",
			setupContext: func() context.Context {
				return context.WithValue(context.Background(), "user_id", userID)
			},
			expectedID:    userID,
			expectedFound: true,
		},
		{
			name: "No user ID in context",
			setupContext: func() context.Context {
				return context.Background()
			},
			expectedID:    uuid.Nil,
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup context
			ctx := tt.setupContext()
			
			// Create request with context
			req := &http.Request{}
			req = req.WithContext(ctx)
			
			// Call function
			id, found := GetUserIDFromContext(req)
			
			// Assert results
			assert.Equal(t, tt.expectedFound, found)
			assert.Equal(t, tt.expectedID, id)
		})
	}
}