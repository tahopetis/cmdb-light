package auth

import (
	"errors"
	"time"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTManager manages JWT tokens
type JWTManager struct {
	secretKey           string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// UserClaims represents the claims in a JWT token
type UserClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	jwt.RegisteredClaims
}

// NewJWTManager creates a new JWTManager
func NewJWTManager(secretKey string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:           secretKey,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// GenerateAccessToken generates a new JWT access token for a user
func (manager *JWTManager) GenerateAccessToken(user *models.User) (string, error) {
	claims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.accessTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// GenerateRefreshToken generates a new JWT refresh token for a user
func (manager *JWTManager) GenerateRefreshToken(user *models.User) (string, error) {
	claims := UserClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(manager.refreshTokenDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(manager.secretKey))
}

// Verify verifies a JWT token and returns the claims
func (manager *JWTManager) Verify(accessToken string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// HashRefreshToken creates a hash of the refresh token for secure storage
func (manager *JWTManager) HashRefreshToken(refreshToken string) (string, error) {
	// Use a simple hash for now - in production, consider using bcrypt or similar
	// This is just for demonstration purposes
	return manager.hashString(refreshToken)
}

// VerifyRefreshTokenHash verifies that a refresh token matches its hash
func (manager *JWTManager) VerifyRefreshTokenHash(refreshToken, hash string) bool {
	tokenHash, err := manager.hashString(refreshToken)
	if err != nil {
		return false
	}
	return tokenHash == hash
}

// hashString is a helper function to hash a string
func (manager *JWTManager) hashString(s string) (string, error) {
	// Use a simple hash for now - in production, consider using bcrypt or similar
	// This is just for demonstration purposes
	h := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"data": s,
		"salt": time.Now().Unix(),
	})
	token, err := h.SignedString([]byte(manager.secretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}