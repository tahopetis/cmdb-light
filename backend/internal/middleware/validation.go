package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/cmdb-lite/backend/internal/validation"
)

// Validator is an interface for validation operations
type Validator interface {
	Validate(s interface{}) *models.ErrorResponse
}

// ValidationMiddleware creates a middleware that validates request bodies
func ValidationMiddleware(validator Validator, model interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only validate POST, PUT, and PATCH requests
			if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodPatch {
				next.ServeHTTP(w, r)
				return
			}

			// Decode the request body into the model
			if err := json.NewDecoder(r.Body).Decode(model); err != nil {
				RespondWithValidationError(w, "Invalid request body", nil)
				return
			}

			// Validate the model
			if validationError := validator.Validate(model); validationError != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(models.GetHTTPStatusForError(models.ErrorTypeValidation))
				json.NewEncoder(w).Encode(validationError)
				return
			}

			// Store the validated model in the request context
			ctx := SetValidatedModel(r.Context(), model)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateStruct is a helper function to validate a struct directly
func ValidateStruct(validator Validator, s interface{}) *models.ErrorResponse {
	return validator.Validate(s)
}

// contextKey is a private type for context keys
type contextKey string

const (
	validatedModelKey contextKey = "validatedModel"
)

// SetValidatedModel stores the validated model in the context
func SetValidatedModel(ctx context.Context, model interface{}) context.Context {
	return context.WithValue(ctx, validatedModelKey, model)
}

// GetValidatedModel retrieves the validated model from the context
func GetValidatedModel(ctx context.Context, model interface{}) bool {
	val, ok := ctx.Value(validatedModelKey).(interface{})
	if !ok {
		return false
	}
	
	// Use reflection to set the model value
	rv := reflect.ValueOf(model)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return false
	}
	
	rv.Elem().Set(reflect.ValueOf(val))
	return true
}

// NewValidator creates a new validator instance
func NewValidator() Validator {
	return validation.NewValidator()
}