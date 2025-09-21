package models

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id" validate:"uuid"`
	Username     string    `json:"username" db:"username" validate:"required,min=3,max=50"`
	Email        string    `json:"email" db:"email" validate:"required,email"`
	PasswordHash string    `json:"-" db:"password_hash" validate:"required,min=8"`
	Role         string    `json:"role" db:"role" validate:"required,oneof=admin user viewer"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

// RefreshToken represents a refresh token in the database
type RefreshToken struct {
	ID        uuid.UUID    `json:"id" db:"id" validate:"uuid"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id" validate:"required,uuid"`
	TokenHash string       `json:"-" db:"token_hash" validate:"required"`
	ExpiresAt time.Time    `json:"expires_at" db:"expires_at" validate:"required"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	RevokedAt *time.Time   `json:"revoked_at,omitempty" db:"revoked_at"`
}

// TokenRefreshRequest represents a token refresh request
type TokenRefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenRefreshResponse represents a token refresh response
type TokenRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required,oneof=admin user viewer"`
}

// CI represents a Configuration Item
type CI struct {
	ID        uuid.UUID      `json:"id" db:"id" validate:"uuid"`
	Name      string         `json:"name" db:"name" validate:"required,min=1,max=100"`
	Type      string         `json:"type" db:"type" validate:"required,min=1,max=50"`
	Attributes JSONBMap      `json:"attributes" db:"attributes"`
	Tags      []string       `json:"tags" db:"tags"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// Relationship represents a relationship between CIs
type Relationship struct {
	ID        uuid.UUID `json:"id" db:"id" validate:"uuid"`
	SourceID  uuid.UUID `json:"source_id" db:"source_id" validate:"required,uuid"`
	TargetID  uuid.UUID `json:"target_id" db:"target_id" validate:"required,uuid"`
	Type      string    `json:"type" db:"type" validate:"required,min=1,max=50"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uuid.UUID `json:"id" db:"id" validate:"uuid"`
	EntityType string    `json:"entity_type" db:"entity_type" validate:"required,min=1,max=50"`
	EntityID   uuid.UUID `json:"entity_id" db:"entity_id" validate:"required,uuid"`
	Action     string    `json:"action" db:"action" validate:"required,min=1,max=20,oneof=create update delete"`
	ChangedBy  string    `json:"changed_by" db:"changed_by" validate:"required,min=1,max=50"`
	ChangedAt  time.Time `json:"changed_at" db:"changed_at"`
	Details    JSONBMap  `json:"details" db:"details"`
}

// JSONBMap is a custom type for handling JSONB data
type JSONBMap map[string]interface{}

// Value implements the driver.Valuer interface for JSONBMap
func (j JSONBMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface for JSONBMap
func (j *JSONBMap) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &j)
}

// StringArray is a custom type for handling string arrays (PostgreSQL text[])
type StringArray []string

// Value implements the driver.Valuer interface for StringArray
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan implements the sql.Scanner interface for StringArray
func (s *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, &s)
}

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code      string      `json:"code"`        // Error code or type
	Message   string      `json:"message"`     // Human-readable error message
	Details   interface{} `json:"details"`     // Optional details field for additional context
	Timestamp time.Time   `json:"timestamp"`   // Timestamp of the error
}

// ErrorType represents different types of errors
type ErrorType string

const (
	// Validation errors (400 Bad Request)
	ErrorTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrorTypeInvalidInput   ErrorType = "INVALID_INPUT"
	ErrorTypeMissingParam   ErrorType = "MISSING_PARAMETER"
	
	// Authentication errors (401 Unauthorized)
	ErrorTypeUnauthorized   ErrorType = "UNAUTHORIZED"
	ErrorTypeInvalidToken  ErrorType = "INVALID_TOKEN"
	ErrorTypeTokenExpired  ErrorType = "TOKEN_EXPIRED"
	
	// Authorization errors (403 Forbidden)
	ErrorTypeForbidden     ErrorType = "FORBIDDEN"
	ErrorTypeInsufficientPermissions ErrorType = "INSUFFICIENT_PERMISSIONS"
	
	// Not found errors (404 Not Found)
	ErrorTypeNotFound      ErrorType = "NOT_FOUND"
	ErrorTypeUserNotFound  ErrorType = "USER_NOT_FOUND"
	ErrorTypeCINotFound    ErrorType = "CI_NOT_FOUND"
	ErrorTypeRelationshipNotFound ErrorType = "RELATIONSHIP_NOT_FOUND"
	ErrorTypeAuditLogNotFound ErrorType = "AUDIT_LOG_NOT_FOUND"
	
	// Server errors (500 Internal Server Error)
	ErrorTypeInternal      ErrorType = "INTERNAL_ERROR"
	ErrorTypeDatabase      ErrorType = "DATABASE_ERROR"
	ErrorTypeService       ErrorType = "SERVICE_ERROR"
)

// NewErrorResponse creates a new error response with the given parameters
func NewErrorResponse(errorType ErrorType, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Code:      string(errorType),
		Message:   message,
		Details:   details,
		Timestamp: time.Now(),
	}
}

// GetHTTPStatusForError returns the appropriate HTTP status code for a given error type
func GetHTTPStatusForError(errorType ErrorType) int {
	switch errorType {
	case ErrorTypeValidation, ErrorTypeInvalidInput, ErrorTypeMissingParam:
		return http.StatusBadRequest
	case ErrorTypeUnauthorized, ErrorTypeInvalidToken, ErrorTypeTokenExpired:
		return http.StatusUnauthorized
	case ErrorTypeForbidden, ErrorTypeInsufficientPermissions:
		return http.StatusForbidden
	case ErrorTypeNotFound, ErrorTypeUserNotFound, ErrorTypeCINotFound,
	     ErrorTypeRelationshipNotFound, ErrorTypeAuditLogNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}