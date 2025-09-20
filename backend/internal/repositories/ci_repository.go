package repositories

import (
	"context"
	"errors"
	
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// Error constants
var (
	ErrCINotFound = errors.New("configuration item not found")
)

// CIRepository defines the interface for CI repository operations
type CIRepository interface {
	Create(ctx context.Context, ci *models.CI) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.CI, error)
	GetAll(ctx context.Context) ([]*models.CI, error)
	Update(ctx context.Context, ci *models.CI) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetByType(ctx context.Context, ciType string) ([]*models.CI, error)
}

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// CITypeRepository defines the interface for CI type repository operations
type CITypeRepository interface {
	Create(ctx context.Context, ciType *models.CIType) error
	GetByID(ctx context.Context, id uint) (*models.CIType, error)
	GetAll(ctx context.Context) ([]*models.CIType, error)
	Update(ctx context.Context, ciType *models.CIType) error
	Delete(ctx context.Context, id uint) error
}

// RelationshipRepository defines the interface for relationship repository operations
type RelationshipRepository interface {
	Create(ctx context.Context, relationship *models.Relationship) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error)
	GetAll(ctx context.Context) ([]*models.Relationship, error)
	GetBySourceID(ctx context.Context, sourceID uuid.UUID) ([]*models.Relationship, error)
	GetByTargetID(ctx context.Context, targetID uuid.UUID) ([]*models.Relationship, error)
	Update(ctx context.Context, relationship *models.Relationship) error
	Delete(ctx context.Context, id uuid.UUID) error
}

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