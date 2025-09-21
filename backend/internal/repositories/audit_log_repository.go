package repositories

import (
	"context"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// AuditLogRepository defines the interface for audit log repository operations
type AuditLogRepository interface {
	// Create creates a new audit log in the database
	Create(ctx context.Context, auditLog *models.AuditLog) error

	// GetByID retrieves an audit log by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error)

	// GetAll retrieves all audit logs from the database
	GetAll(ctx context.Context) ([]*models.AuditLog, error)

	// GetByEntityType retrieves audit logs by entity type
	GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error)

	// GetByEntityID retrieves audit logs by entity ID
	GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error)

	// GetByChangedBy retrieves audit logs by the user who made the change
	GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error)

	// Delete deletes an audit log from the database
	Delete(ctx context.Context, id uuid.UUID) error
}
