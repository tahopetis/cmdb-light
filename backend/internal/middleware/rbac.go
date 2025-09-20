package middleware

import (
	"net/http"
)

// RBACMiddleware creates a middleware for role-based access control
func RBACMiddleware(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the user role from the context
			role, ok := GetUserRoleFromContext(r.Context())
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if the user has the required role
			if !hasRequiredRole(role, roles) {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			// User has the required role, proceed to the next handler
			next.ServeHTTP(w, r)
		})
	}
}

// AdminOnly creates a middleware that only allows admin users
func AdminOnly() func(http.Handler) http.Handler {
	return RBACMiddleware("admin")
}

// AdminOrViewer creates a middleware that allows admin and viewer users
func AdminOrViewer() func(http.Handler) http.Handler {
	return RBACMiddleware("admin", "viewer")
}

// hasRequiredRole checks if the user has any of the required roles
func hasRequiredRole(userRole string, requiredRoles []string) bool {
	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}
	return false
}