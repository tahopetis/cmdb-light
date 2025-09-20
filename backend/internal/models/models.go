package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Role         string    `json:"role" db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// CI represents a Configuration Item
type CI struct {
	ID        uuid.UUID      `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Type      string         `json:"type" db:"type"`
	Attributes JSONBMap      `json:"attributes" db:"attributes"`
	Tags      []string       `json:"tags" db:"tags"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}

// Relationship represents a relationship between CIs
type Relationship struct {
	ID        uuid.UUID `json:"id" db:"id"`
	SourceID  uuid.UUID `json:"source_id" db:"source_id"`
	TargetID  uuid.UUID `json:"target_id" db:"target_id"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID         uuid.UUID `json:"id" db:"id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   uuid.UUID `json:"entity_id" db:"entity_id"`
	Action     string    `json:"action" db:"action"`
	ChangedBy  string    `json:"changed_by" db:"changed_by"`
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