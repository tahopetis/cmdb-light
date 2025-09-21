package services

import (
	"context"
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

// MockPasswordManager is a mock implementation of PasswordManager interface
type MockPasswordManager struct {
	mock.Mock
}

func (m *MockPasswordManager) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordManager) CheckPassword(password, hash string) bool {
	args := m.Called(password, hash)
	return args.Bool(0)
}

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
func createTestUser(username, password, role string) *models.User {
	return &models.User{
		ID:           uuid.New(),
		Username:     username,
		PasswordHash: password, // In tests, we'll use plain text for simplicity
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

func TestAuthService_Login_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockPasswordManager.On("CheckPassword", "password123", "hashedpassword").Return(true)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockJWTManager.On("GenerateAccessToken", testUser).Return("access-token", nil)
	mockJWTManager.On("GenerateRefreshToken", testUser).Return("refresh-token", nil)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "testuser", "password123")
	
	// Assert
	require.NoError(t, err)
	assert.Equal(t, "access-token", result.AccessToken)
	assert.Equal(t, "refresh-token", result.RefreshToken)
	assert.NotNil(t, result.User)
	assert.Equal(t, testUser.ID, result.User.ID)
	assert.Equal(t, testUser.Username, result.User.Username)
	assert.Equal(t, testUser.Role, result.User.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
	mockJWTManager.AssertExpectations(t)
}

func TestAuthService_Login_InvalidPassword(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockPasswordManager.On("CheckPassword", "wrongpassword", "hashedpassword").Return(false)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "testuser", "wrongpassword")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid username or password", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "nonexistent", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid username or password", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_RepositoryError(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, assert.AnError)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "testuser", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid username or password", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Login_RefreshTokenError(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockPasswordManager.On("CheckPassword", "password123", "hashedpassword").Return(true)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(assert.AnError)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "testuser", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "INTERNAL_ERROR", appErr.Code)
	assert.Equal(t, "Failed to login", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
}

func TestAuthService_Login_TokenGenerationError(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
	mockPasswordManager.On("CheckPassword", "password123", "hashedpassword").Return(true)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockJWTManager.On("GenerateAccessToken", testUser).Return("", assert.AnError)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.Login(context.Background(), "testuser", "password123")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "INTERNAL_ERROR", appErr.Code)
	assert.Equal(t, "Failed to login", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
	mockJWTManager.AssertExpectations(t)
}

func TestAuthService_RefreshToken_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	testRefreshToken := createTestRefreshToken(testUser.ID, "valid-refresh-token")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "valid-refresh-token").Return(testRefreshToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
	mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
	mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
	mockJWTManager.On("GenerateAccessToken", testUser).Return("new-access-token", nil)
	mockJWTManager.On("GenerateRefreshToken", testUser).Return("new-refresh-token", nil)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.RefreshToken(context.Background(), "valid-refresh-token")
	
	// Assert
	require.NoError(t, err)
	assert.Equal(t, "new-access-token", result.AccessToken)
	assert.Equal(t, "new-refresh-token", result.RefreshToken)
	assert.NotNil(t, result.User)
	assert.Equal(t, testUser.ID, result.User.ID)
	assert.Equal(t, testUser.Username, result.User.Username)
	assert.Equal(t, testUser.Role, result.User.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
	mockJWTManager.AssertExpectations(t)
}

func TestAuthService_RefreshToken_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "invalid-refresh-token").Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.RefreshToken(context.Background(), "invalid-refresh-token")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid refresh token", appErr.Message)
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_RefreshToken_UserNotFound(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	testRefreshToken := createTestRefreshToken(testUser.ID, "valid-refresh-token")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "valid-refresh-token").Return(testRefreshToken, nil)
	mockUserRepo.On("GetByID", mock.Anything, testUser.ID).Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.RefreshToken(context.Background(), "valid-refresh-token")
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid refresh token", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByID", mock.Anything, testUser.ID).Return(testUser, nil)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.ValidateToken(context.Background(), testUser.ID)
	
	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.ID, result.ID)
	assert.Equal(t, testUser.Username, result.Username)
	assert.Equal(t, testUser.Role, result.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_ValidateToken_UserNotFound(t *testing.T) {
	// Arrange
	testUserID := uuid.New()
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByID", mock.Anything, testUserID).Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.ValidateToken(context.Background(), testUserID)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "User not found", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_Logout_Success(t *testing.T) {
	// Arrange
	testUser := createTestUser("testuser", "hashedpassword", "user")
	testRefreshToken := createTestRefreshToken(testUser.ID, "valid-refresh-token")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "valid-refresh-token").Return(testRefreshToken, nil)
	mockRefreshTokenRepo.On("Delete", mock.Anything, testRefreshToken.ID).Return(nil)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	err := service.Logout(context.Background(), "valid-refresh-token")
	
	// Assert
	require.NoError(t, err)
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_Logout_InvalidToken(t *testing.T) {
	// Arrange
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockRefreshTokenRepo.On("GetByToken", mock.Anything, "invalid-refresh-token").Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	err := service.Logout(context.Background(), "invalid-refresh-token")
	
	// Assert
	assert.Error(t, err)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "UNAUTHORIZED", appErr.Code)
	assert.Equal(t, "Invalid refresh token", appErr.Message)
	
	// Verify mock expectations
	mockRefreshTokenRepo.AssertExpectations(t)
}

func TestAuthService_CreateUser_Success(t *testing.T) {
	// Arrange
	newUser := &models.User{
		Username: "newuser",
		Password: "password123",
		Role:     "user",
	}
	
	createdUser := createTestUser("newuser", "hashedpassword", "user")
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockPasswordManager.On("HashPassword", "password123").Return("hashedpassword", nil)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil)
	mockUserRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.CreateUser(context.Background(), newUser)
	
	// Assert
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createdUser.Username, result.Username)
	assert.Equal(t, createdUser.Role, result.Role)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
}

func TestAuthService_CreateUser_UsernameExists(t *testing.T) {
	// Arrange
	existingUser := createTestUser("existinguser", "hashedpassword", "user")
	newUser := &models.User{
		Username: "existinguser",
		Password: "password123",
		Role:     "user",
	}
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockUserRepo.On("GetByUsername", mock.Anything, "existinguser").Return(existingUser, nil)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.CreateUser(context.Background(), newUser)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "CONFLICT", appErr.Code)
	assert.Equal(t, "Username already exists", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
}

func TestAuthService_CreateUser_PasswordHashError(t *testing.T) {
	// Arrange
	newUser := &models.User{
		Username: "newuser",
		Password: "password123",
		Role:     "user",
	}
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockPasswordManager.On("HashPassword", "password123").Return("", assert.AnError)
	mockUserRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, repositories.ErrNotFound)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.CreateUser(context.Background(), newUser)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "INTERNAL_ERROR", appErr.Code)
	assert.Equal(t, "Failed to create user", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
}

func TestAuthService_CreateUser_RepositoryError(t *testing.T) {
	// Arrange
	newUser := &models.User{
		Username: "newuser",
		Password: "password123",
		Role:     "user",
	}
	
	mockUserRepo := &MockUserRepository{}
	mockRefreshTokenRepo := &MockRefreshTokenRepository{}
	mockPasswordManager := &MockPasswordManager{}
	mockJWTManager := &MockJWTManager{}
	
	// Setup mock expectations
	mockPasswordManager.On("HashPassword", "password123").Return("hashedpassword", nil)
	mockUserRepo.On("GetByUsername", mock.Anything, "newuser").Return(nil, repositories.ErrNotFound)
	mockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(assert.AnError)
	
	service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)
	
	// Act
	result, err := service.CreateUser(context.Background(), newUser)
	
	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	
	var appErr *models.ErrorResponse
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, "INTERNAL_ERROR", appErr.Code)
	assert.Equal(t, "Failed to create user", appErr.Message)
	
	// Verify mock expectations
	mockUserRepo.AssertExpectations(t)
	mockPasswordManager.AssertExpectations(t)
}

// Table-driven test for login scenarios
func TestAuthService_Login_TableDriven(t *testing.T) {
	testUser := createTestUser("testuser", "hashedpassword", "user")
	
	tests := []struct {
		name           string
		setupMock      func() (*MockUserRepository, *MockRefreshTokenRepository, *MockPasswordManager, *MockJWTManager)
		username       string
		password       string
		expectError    bool
		expectedError  string
	}{
		{
			name: "Successful login",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository, *MockPasswordManager, *MockJWTManager) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				mockPasswordManager := &MockPasswordManager{}
				mockJWTManager := &MockJWTManager{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
				mockPasswordManager.On("CheckPassword", "password123", "hashedpassword").Return(true)
				mockRefreshTokenRepo.On("DeleteByUserID", mock.Anything, testUser.ID).Return(nil)
				mockRefreshTokenRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.RefreshToken")).Return(nil)
				mockJWTManager.On("GenerateAccessToken", testUser).Return("access-token", nil)
				mockJWTManager.On("GenerateRefreshToken", testUser).Return("refresh-token", nil)
				
				return mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager
			},
			username:    "testuser",
			password:    "password123",
			expectError: false,
		},
		{
			name: "Invalid password",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository, *MockPasswordManager, *MockJWTManager) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				mockPasswordManager := &MockPasswordManager{}
				mockJWTManager := &MockJWTManager{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(testUser, nil)
				mockPasswordManager.On("CheckPassword", "wrongpassword", "hashedpassword").Return(false)
				
				return mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager
			},
			username:       "testuser",
			password:       "wrongpassword",
			expectError:    true,
			expectedError:  "Invalid username or password",
		},
		{
			name: "User not found",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository, *MockPasswordManager, *MockJWTManager) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				mockPasswordManager := &MockPasswordManager{}
				mockJWTManager := &MockJWTManager{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "nonexistent").Return(nil, repositories.ErrNotFound)
				
				return mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager
			},
			username:       "nonexistent",
			password:       "password123",
			expectError:    true,
			expectedError:  "Invalid username or password",
		},
		{
			name: "Repository error",
			setupMock: func() (*MockUserRepository, *MockRefreshTokenRepository, *MockPasswordManager, *MockJWTManager) {
				mockUserRepo := &MockUserRepository{}
				mockRefreshTokenRepo := &MockRefreshTokenRepository{}
				mockPasswordManager := &MockPasswordManager{}
				mockJWTManager := &MockJWTManager{}
				
				mockUserRepo.On("GetByUsername", mock.Anything, "testuser").Return(nil, assert.AnError)
				
				return mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager
			},
			username:       "testuser",
			password:       "password123",
			expectError:    true,
			expectedError:  "Invalid username or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager := tt.setupMock()
			service := NewAuthService(mockUserRepo, mockRefreshTokenRepo, mockPasswordManager, mockJWTManager)

			// Act
			result, err := service.Login(context.Background(), tt.username, tt.password)

			// Assert
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
				assert.NotNil(t, result.User)
			}

			// Verify mock expectations
			mockUserRepo.AssertExpectations(t)
			mockRefreshTokenRepo.AssertExpectations(t)
			mockPasswordManager.AssertExpectations(t)
			mockJWTManager.AssertExpectations(t)
		})
	}
}
