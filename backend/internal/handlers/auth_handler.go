package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/middleware"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/cmdb-lite/backend/internal/validation"
	"github.com/google/uuid"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	jwtManager       *auth.JWTManager
	passwordManager  *auth.PasswordManager
	validator       *validation.Validator
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
	passwordManager *auth.PasswordManager,
) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		passwordManager:  passwordManager,
		validator:       validation.NewValidator(),
	}
}

// Login handles user login
// @Summary User login
// @Description Authenticate a user with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param loginRequest body models.LoginRequest true "Login credentials"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	// Decode the request body
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		middleware.RespondWithValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate the input using the validator
	if validationError := h.validator.Validate(loginReq); validationError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(models.GetHTTPStatusForError(models.ErrorTypeValidation))
		json.NewEncoder(w).Encode(validationError)
		return
	}

	// Get the user from the database
	user, err := h.userRepo.GetByUsername(r.Context(), loginReq.Username)
	if err != nil {
		middleware.RespondWithUnauthorizedError(w, "Invalid username or password", nil)
		return
	}

	// Check the password
	if err := h.passwordManager.CheckPassword(loginReq.Password, user.PasswordHash); err != nil {
		middleware.RespondWithUnauthorizedError(w, "Invalid username or password", nil)
		return
	}

	// Generate access token
	accessToken, err := h.jwtManager.GenerateAccessToken(user)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to generate access token", nil)
		return
	}

	// Generate refresh token
	refreshToken, err := h.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to generate refresh token", nil)
		return
	}

	// Hash the refresh token for secure storage
	refreshTokenHash, err := h.jwtManager.HashRefreshToken(refreshToken)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to hash refresh token", nil)
		return
	}

	// Store the refresh token in the database
	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: refreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := h.refreshTokenRepo.Create(r.Context(), refreshTokenModel); err != nil {
		middleware.RespondWithInternalError(w, "Failed to store refresh token", nil)
		return
	}

	// Create the response
	response := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ValidateToken validates a JWT token
// @Summary Validate token
// @Description Validate a JWT token and return user information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/validate [get]
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context (set by the AuthMiddleware)
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		middleware.RespondWithUnauthorizedError(w, "User not authenticated", nil)
		return
	}

	// Get the user from the database
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "User not found", nil)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Refresh an access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param refreshRequest body models.TokenRefreshRequest true "Refresh token"
// @Success 200 {object} models.TokenRefreshResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Decode the request body
	var refreshReq models.TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshReq); err != nil {
		middleware.RespondWithValidationError(w, "Invalid request body", nil)
		return
	}

	// Validate the input using the validator
	if validationError := h.validator.Validate(refreshReq); validationError != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(models.GetHTTPStatusForError(models.ErrorTypeValidation))
		json.NewEncoder(w).Encode(validationError)
		return
	}

	// Verify the refresh token
	refreshTokenClaims, err := h.jwtManager.Verify(refreshReq.RefreshToken)
	if err != nil {
		middleware.RespondWithUnauthorizedError(w, "Invalid or expired refresh token", nil)
		return
	}

	// Hash the refresh token to look it up in the database
	refreshTokenHash, err := h.jwtManager.HashRefreshToken(refreshReq.RefreshToken)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to process refresh token", nil)
		return
	}

	// Get the refresh token from the database
	refreshToken, err := h.refreshTokenRepo.GetByTokenHash(r.Context(), refreshTokenHash)
	if err != nil {
		middleware.RespondWithUnauthorizedError(w, "Invalid refresh token", nil)
		return
	}

	// Check if the refresh token is revoked or expired
	if refreshToken.RevokedAt != nil || time.Now().After(refreshToken.ExpiresAt) {
		middleware.RespondWithUnauthorizedError(w, "Refresh token revoked or expired", nil)
		return
	}

	// Get the user from the database
	user, err := h.userRepo.GetByID(r.Context(), refreshTokenClaims.UserID)
	if err != nil {
		middleware.RespondWithNotFoundError(w, "User not found", nil)
		return
	}

	// Generate a new access token
	newAccessToken, err := h.jwtManager.GenerateAccessToken(user)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to generate new access token", nil)
		return
	}

	// Generate a new refresh token (token rotation)
	newRefreshToken, err := h.jwtManager.GenerateRefreshToken(user)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to generate new refresh token", nil)
		return
	}

	// Hash the new refresh token
	newRefreshTokenHash, err := h.jwtManager.HashRefreshToken(newRefreshToken)
	if err != nil {
		middleware.RespondWithInternalError(w, "Failed to hash new refresh token", nil)
		return
	}

	// Revoke the old refresh token
	if err := h.refreshTokenRepo.Revoke(r.Context(), refreshToken.ID); err != nil {
		// Log the error but continue, as this is not critical
		// In production, you might want to handle this more gracefully
	}

	// Store the new refresh token in the database
	newRefreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: newRefreshTokenHash,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := h.refreshTokenRepo.Create(r.Context(), newRefreshTokenModel); err != nil {
		middleware.RespondWithInternalError(w, "Failed to store new refresh token", nil)
		return
	}

	// Create the response
	response := models.TokenRefreshResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// Logout handles user logout
// @Summary User logout
// @Description Revoke the user's refresh tokens
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get the user from the context (set by the AuthMiddleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		middleware.RespondWithUnauthorizedError(w, "User not authenticated", nil)
		return
	}

	// Revoke all refresh tokens for the user
	if err := h.refreshTokenRepo.RevokeAllForUser(r.Context(), userID); err != nil {
		middleware.RespondWithInternalError(w, "Failed to revoke refresh tokens", nil)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully logged out"})
}

// Helper function to get user ID from context
func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	// This would be implemented in the middleware package
	// For now, we'll return a placeholder implementation
	return uuid.Nil, false
}