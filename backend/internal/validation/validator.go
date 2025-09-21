package validation

import (
	"fmt"
	"reflect"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/cmdb-lite/backend/internal/models"
)

// Validator is a wrapper around the validator.v10 validator
type Validator struct {
	validator *validator.Validate
}

// NewValidator creates a new Validator instance
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

// Validate validates a struct and returns a standardized error response if validation fails
func (v *Validator) Validate(s interface{}) *models.ErrorResponse {
	err := v.validator.Struct(s)
	if err == nil {
		return nil
	}

	// Convert validator error to our standardized error format
	validationErrors := v.formatValidationErrors(err)
	return models.NewErrorResponse(
		models.ErrorTypeValidation,
		"Validation failed",
		validationErrors,
	)
}

// formatValidationErrors converts validator errors to a more readable format
func (v *Validator) formatValidationErrors(err error) map[string]interface{} {
	errors := make(map[string]interface{})
	
	// Type assert to validator.ValidationErrors
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		errors["general"] = err.Error()
		return errors
	}

	// Process each validation error
	for _, e := range validationErrors {
		fieldName := e.Field()
		
		// Get a more readable field name by converting from CamelCase to snake_case
		fieldName = toSnakeCase(fieldName)
		
		// Create a descriptive error message
		message := v.getErrorMessage(e)
		
		// Add to errors map
		if _, exists := errors[fieldName]; !exists {
			errors[fieldName] = []string{}
		}
		if fieldErrors, ok := errors[fieldName].([]string); ok {
			errors[fieldName] = append(fieldErrors, message)
		} else {
			errors[fieldName] = []string{message}
		}
	}
	
	return errors
}

// getErrorMessage generates a user-friendly error message for a validation error
func (v *Validator) getErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", e.Field())
	case "min":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at least %s characters", e.Field(), e.Param())
		}
		return fmt.Sprintf("%s must be at least %s", e.Field(), e.Param())
	case "max":
		if e.Kind() == reflect.String {
			return fmt.Sprintf("%s must be at most %s characters", e.Field(), e.Param())
		}
		return fmt.Sprintf("%s must be at most %s", e.Field(), e.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", e.Field())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", e.Field())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s is invalid (failed on %s)", e.Field(), e.Tag())
	}
}

// toSnakeCase converts a string from CamelCase to snake_case
func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// RegisterCustomValidation registers a custom validation function
func (v *Validator) RegisterCustomValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}