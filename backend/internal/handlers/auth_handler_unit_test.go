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
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockUserRepository is a mock implementation of UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockRefreshTokenRepository is a mock implementation of RefreshTokenRepository interface
type MockRefreshTokenRepository struct {
	mock.Mock
}

func (m *MockRefreshTokenRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	args := m.Called(ctx, refreshToken)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.RefreshToken, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) GetByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.RefreshToken, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*models.RefreshToken), args.Error(1)
}

func (m *MockRefreshTokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRefreshTokenRepository) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

// Helper functions to create test data
func createTestUser(username, password, role string) *models.User {
	hashedPassword, _ := auth.NewPasswordManager().HashPassword(password)
	return &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: hashedPassword,
		Role:         role,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func createTestRefreshToken(userID uuid.UUID, token string) *models.RefreshToken {
	return &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}
}

func TestAuthHandler_Login_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "password123", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	loginRequest := models.LoginRequest{
		Username: "testuser",
		Password: "password123",
	}
	
	requestBody, err := json.Marshal(loginRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Login(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, testUser.ID, response.User.ID)
	assert.Equal(t, testUser.Username, response.User.Username)
	assert.Equal(t, testUser.Role, response.User.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidCredentials(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "password123", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	loginRequest := models.LoginRequest{
		Username: "testuser",
		Password: "wrongpassword", // Wrong password
	}
	
	requestBody, err := json.Marshal(loginRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Login(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid username or password", response.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repositories.ErrNotFound)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	loginRequest := models.LoginRequest{
		Username: "nonexistent",
		Password: "password123",
	}
	
	requestBody, err := json.Marshal(loginRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Login(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid username or password", response.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Login_InvalidRequestBody(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	// Invalid JSON
	req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid-json")))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Login(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusBadRequest, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "VALIDATION_ERROR", response.Code)
	assert.Equal(t, "Invalid request body", response.Message)
}

func TestAuthHandler_Login_MissingFields(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	tests := []struct {
		name        string
		loginRequest models.LoginRequest
		expectedMsg string
	}{
		{
			name: "Missing username",
			loginRequest: models.LoginRequest{
				Username: "",
				Password: "password123",
			},
			expectedMsg: "Username and password are required",
		},
		{
			name: "Missing password",
			loginRequest: models.LoginRequest{
				Username: "testuser",
				Password: "",
			},
			expectedMsg: "Username and password are required",
		},
		{
			name: "Both missing",
			loginRequest: models.LoginRequest{
				Username: "",
				Password: "",
			},
			expectedMsg: "Username and password are required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requestBody, err := json.Marshal(tt.loginRequest)
			require.NoError(t, err)
			
			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			
			rr := httptest.NewRecorder()
			
			// Act
			handler.Login(rr, req)
			
			// Assert
			assert.Equal(t, http.StatusBadRequest, rr.Code)
			
			var response models.ErrorResponse
			err = json.Unmarshal(rr.Body.Bytes(), &response)
			require.NoError(t, err)
			
			assert.Equal(t, "VALIDATION_ERROR", response.Code)
			assert.Equal(t, tt.expectedMsg, response.Message)
		})
	}
}

func TestAuthHandler_RefreshToken_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "password123", "user")
	testRefreshToken := createTestRefreshToken(testUser.ID, "valid-refresh-token")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "valid-refresh-token").Return(testRefreshToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	refreshRequest := models.RefreshTokenRequest{
		RefreshToken: "valid-refresh-token",
	}
	
	requestBody, err := json.Marshal(refreshRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.RefreshToken(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.LoginResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.NotEmpty(t, response.Token)
	assert.NotEmpty(t, response.RefreshToken)
	assert.NotNil(t, response.User)
	assert.Equal(t, testUser.ID, response.User.ID)
	assert.Equal(t, testUser.Username, response.User.Username)
	assert.Equal(t, testUser.Role, response.User.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "invalid-refresh-token").Return(nil, repositories.ErrNotFound)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	refreshRequest := models.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}
	
	requestBody, err := json.Marshal(refreshRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.RefreshToken(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid refresh token", response.Message)
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_ValidateToken_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "password123", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockUserRepo.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	req, err := http.NewRequest("GET", "/api/v1/auth/validate", nil)
	require.NoError(t, err)
	
	// Add user ID to context (simulating JWT middleware)
	ctx := context.WithValue(req.Context(), "user_id", testUser.ID)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.ValidateToken(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response models.User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, testUser.ID, response.ID)
	assert.Equal(t, testUser.Username, response.Username)
	assert.Equal(t, testUser.Role, response.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_ValidateToken_Unauthorized(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	req, err := http.NewRequest("GET", "/api/v1/auth/validate", nil)
	require.NoError(t, err)
	
	// No user ID in context
	rr := httptest.NewRecorder()
	
	// Act
	handler.ValidateToken(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "User not authenticated", response.Message)
}

func TestAuthHandler_ValidateToken_UserNotFound(t *testing.T) {
	// Arrange
	testUserID := uuid.New()
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockUserRepo.On("GetByID", mock.Anything, testUserID).Return(nil, repositories.ErrNotFound)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	req, err := http.NewRequest("GET", "/api/v1/auth/validate", nil)
	require.NoError(t, err)
	
	// Add user ID to context (simulating JWT middleware)
	ctx := context.WithValue(req.Context(), "user_id", testUserID)
	req = req.WithContext(ctx)
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.ValidateToken(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "User not found", response.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthHandler_Logout_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "password123", "user")
	testRefreshToken := createTestRefreshToken(testUser.ID, "valid-refresh-token")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "valid-refresh-token").Return(testRefreshToken, nil)
	mockRefreshTokenRepo.On("Delete", mock.Anything, testRefreshToken.ID).Return(nil)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	logoutRequest := models.LogoutRequest{
		RefreshToken: "valid-refresh-token",
	}
	
	requestBody, err := json.Marshal(logoutRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Logout(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	
	var response map[string]string
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "Successfully logged out", response["message"])
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthHandler_Logout_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	
	jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
	passwordManager := auth.NewPasswordManager()
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "invalid-refresh-token").Return(nil, repositories.ErrNotFound)
	
	handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)
	
	logoutRequest := models.LogoutRequest{
		RefreshToken: "invalid-refresh-token",
	}
	
	requestBody, err := json.Marshal(logoutRequest)
	require.NoError(t, err)
	
	req, err := http.NewRequest("POST", "/api/v1/auth/logout", bytes.NewBuffer(requestBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	
	rr := httptest.NewRecorder()
	
	// Act
	handler.Logout(rr, req)
	
	// Assert
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	
	var response models.ErrorResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	
	assert.Equal(t, "UNAUTHORIZED", response.Code)
	assert.Equal(t, "Invalid refresh token", response.Message)
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

// Table-driven test for login scenarios
func TestAuthHandler_Login_TableDriven(t *testing.T) {
	testUser := createTestUser("testuser", "password123", "user")
	
	tests := []struct {
		name           string
		setupMock      func() (*MockUserRepository, *MockRefreshTokenRepository)
		requestBody    models.LoginRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Successful login",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
				mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
				mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
				
				return mockUserRepo, mockRefreshTokenRepo
			},
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "User not found",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repositories.ErrNotFound)
				
				return mockUserRepo, mockRefreshTokenRepo
			},
			requestBody: models.LoginRequest{
				Username: "nonexistent",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid username or password",
		},
		{
			name: "Repository error",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, assert.AnError)
				
				return mockUserRepo, mockRefreshTokenRepo
			},
			requestBody: models.LoginRequest{
				Username: "testuser",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo, mockRefreshTokenRepo := tt.setupMock()
			
			jwtManager := auth.NewJWTManager("test-secret", 15*time.Minute, 24*time.Hour)
			passwordManager := auth.NewPasswordManager()
			
			handler := NewAuthHandler(mockUserRepo, mockRefreshTokenRepo, jwtManager, passwordManager)

			requestBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			// Act
			handler.Login(rr, req)

			// Assert
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			}

			// Verify mock expectations
			mockUserRepo.AssertExpectations(t)
			mockRefreshTokenRepo.AssertExpectations(t)
		})
	}
}
