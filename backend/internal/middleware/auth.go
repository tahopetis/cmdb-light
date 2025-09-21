package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/google/uuid"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	UserContextKey ContextKey = "user"
)

// AuthMiddleware creates a middleware for authentication
func AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				RespondWithUnauthorizedError(w, "Authorization header required", nil)
				return
			}

			// Check if the Authorization header has the correct format
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				RespondWithUnauthorizedError(w, "Invalid authorization format", nil)
				return
			}

			// Extract the token
			token := parts[1]

			// Verify the token
			claims, err := jwtManager.Verify(token)
			if err != nil {
				RespondWithUnauthorizedError(w, "Invalid or expired token", nil)
				return
			}

			// Add the user information to the context
			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts the user claims from the context
func GetUserFromContext(ctx context.Context) (*auth.UserClaims, bool) {
	user, ok := ctx.Value(UserContextKey).(*auth.UserClaims)
	return user, ok
}

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return user.UserID, true
}

// GetUserRoleFromContext extracts the user role from the context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	return user.Role, true
}

// GetUsernameFromContext extracts the username from the context
func GetUsernameFromContext(ctx context.Context) (string, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return "", false
	}
	return user.Username, true
}