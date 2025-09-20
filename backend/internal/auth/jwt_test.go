package auth

import (
	"testing"
	"time"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTManager_Generate(t *testing.T) {
	// Setup test data
	secretKey := "test-secret-key"
	tokenDuration := time.Hour
	userID := uuid.New()
	testUser := &models.User{
		ID:       userID,
		Username: "testuser",
		Role:     "user",
	}

	tests := []struct {
		name          string
		setupJWT      func() *JWTManager
		expectedError bool
	}{
		{
			name: "Successful token generation",
			setupJWT: func() *JWTManager {
				return NewJWTManager(secretKey, tokenDuration)
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup JWT manager
			jwtManager := tt.setupJWT()
			
			// Call the method
			token, err := jwtManager.Generate(testUser)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
				
				// Verify token can be parsed
				claims, err := jwtManager.Verify(token)
				assert.NoError(t, err)
				require.NotNil(t, claims)
				
				assert.Equal(t, testUser.ID, claims.UserID)
				assert.Equal(t, testUser.Username, claims.Username)
				assert.Equal(t, testUser.Role, claims.Role)
				
				// Verify expiration time
				assert.True(t, claims.ExpiresAt.After(time.Now()))
			}
		})
	}
}

func TestJWTManager_Verify(t *testing.T) {
	// Setup test data
	secretKey := "test-secret-key"
	tokenDuration := time.Hour
	userID := uuid.New()
	testUser := &models.User{
		ID:       userID,
		Username: "testuser",
		Role:     "user",
	}

	tests := []struct {
		name          string
		setupToken    func(*JWTManager) string
		expectedError bool
		errorType     error
	}{
		{
			name: "Successful token verification",
			setupToken: func(jwtManager *JWTManager) string {
				token, _ := jwtManager.Generate(testUser)
				return token
			},
			expectedError: false,
		},
		{
			name: "Invalid token",
			setupToken: func(jwtManager *JWTManager) string {
				return "invalid-token"
			},
			expectedError: true,
		},
		{
			name: "Token with invalid signing method",
			setupToken: func(jwtManager *JWTManager) string {
				// Create a token with a different signing method
				claims := UserClaims{
					UserID:   testUser.ID,
					Username: testUser.Username,
					Role:     testUser.Role,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenDuration)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						NotBefore: jwt.NewNumericDate(time.Now()),
					},
				}
				
				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				// This will fail because we don't have a private key, but that's fine for the test
				return token.Raw
			},
			expectedError: true,
		},
		{
			name: "Token with invalid claims",
			setupToken: func(jwtManager *JWTManager) string {
				// Create a token with invalid claims
				claims := jwt.MapClaims{
					"sub": "user",
					"exp": time.Now().Add(tokenDuration).Unix(),
				}
				
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(secretKey))
				return tokenString
			},
			expectedError: true,
			errorType:     jwt.ErrInvalidType,
		},
		{
			name: "Expired token",
			setupToken: func(jwtManager *JWTManager) string {
				// Create a token that's already expired
				claims := UserClaims{
					UserID:   testUser.ID,
					Username: testUser.Username,
					Role:     testUser.Role,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
						IssuedAt:  jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
						NotBefore: jwt.NewNumericDate(time.Now().Add(-time.Hour * 2)),
					},
				}
				
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(secretKey))
				return tokenString
			},
			expectedError: true,
			errorType:     jwt.ErrTokenExpired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup JWT manager
			jwtManager := NewJWTManager(secretKey, tokenDuration)
			
			// Setup token
			token := tt.setupToken(jwtManager)
			
			// Call the method
			claims, err := jwtManager.Verify(token)
			
			// Assert results
			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType != nil {
					assert.Contains(t, err.Error(), tt.errorType.Error())
				}
			} else {
				assert.NoError(t, err)
				require.NotNil(t, claims)
				
				assert.Equal(t, testUser.ID, claims.UserID)
				assert.Equal(t, testUser.Username, claims.Username)
				assert.Equal(t, testUser.Role, claims.Role)
			}
		})
	}
}

func TestUserClaims_StandardClaims(t *testing.T) {
	// Setup test data
	userID := uuid.New()
	
	claims := &UserClaims{
		UserID:   userID,
		Username: "testuser",
		Role:     "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	
	// Verify that the standard claims are properly initialized
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.NotNil(t, claims.NotBefore)
}