package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/repositories"
	"github.com/google/uuid"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	userRepo       repositories.UserRepository
	jwtManager     *auth.JWTManager
	passwordManager *auth.PasswordManager
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(
	userRepo repositories.UserRepository,
	jwtManager *auth.JWTManager,
	passwordManager *auth.PasswordManager,
) *AuthHandler {
	return &AuthHandler{
		userRepo:       userRepo,
		jwtManager:     jwtManager,
		passwordManager: passwordManager,
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
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the input
	if loginReq.Username == "" || loginReq.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Get the user from the database
	user, err := h.userRepo.GetByUsername(r.Context(), loginReq.Username)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Check the password
	if err := h.passwordManager.CheckPassword(loginReq.Password, user.PasswordHash); err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := h.jwtManager.Generate(user)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create the response
	response := models.LoginResponse{
		Token: token,
		User:  user,
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
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	// Get the user from the database
	user, err := h.userRepo.GetByID(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Helper function to get user ID from context
func GetUserIDFromContext(r *http.Request) (uuid.UUID, bool) {
	// This would be implemented in the middleware package
	// For now, we'll return a placeholder implementation
	return uuid.Nil, false
}