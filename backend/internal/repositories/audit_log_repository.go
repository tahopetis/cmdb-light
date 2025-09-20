package repositories

import (
	"context"
	"errors"
	
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// Error constants
var (
	ErrAuditLogNotFound = errors.New("audit log not found")
)

// AuditLogRepository defines the interface for audit log repository operations
type AuditLogRepository interface {
	Create(ctx context.Context, auditLog *models.AuditLog) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error)
	GetAll(ctx context.Context) ([]*models.AuditLog, error)
	GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error)
	GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error)
	GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error)
	Delete(ctx context.Context, id uuid.UUID) error
}