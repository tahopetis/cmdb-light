package repositories

import (
	"context"
	"errors"
	
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// Error constants
var (
	ErrRelationshipNotFound = errors.New("relationship not found")
)

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