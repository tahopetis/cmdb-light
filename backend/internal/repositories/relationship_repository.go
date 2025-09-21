package repositories

import (
	"context"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// RelationshipRepository defines the interface for relationship repository operations
type RelationshipRepository interface {
	// Create creates a new relationship in the database
	Create(ctx context.Context, relationship *models.Relationship) error

	// GetByID retrieves a relationship by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error)

	// GetBySourceCI retrieves relationships by source CI ID
	GetBySourceCI(ctx context.Context, sourceCIID uuid.UUID) ([]*models.Relationship, error)

	// GetByTargetCI retrieves relationships by target CI ID
	GetByTargetCI(ctx context.Context, targetCIID uuid.UUID) ([]*models.Relationship, error)

	// GetBySourceAndTarget retrieves relationships by source and target CI IDs
	GetBySourceAndTarget(ctx context.Context, sourceCIID, targetCIID uuid.UUID) ([]*models.Relationship, error)

	// GetByType retrieves relationships by type
	GetByType(ctx context.Context, relationshipType string) ([]*models.Relationship, error)

	// GetAll retrieves all relationships from the database
	GetAll(ctx context.Context) ([]*models.Relationship, error)

	// Update updates a relationship in the database
	Update(ctx context.Context, relationship *models.Relationship) error

	// Delete deletes a relationship from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteBySourceCI deletes all relationships for a source CI
	DeleteBySourceCI(ctx context.Context, sourceCIID uuid.UUID) error

	// DeleteByTargetCI deletes all relationships for a target CI
	DeleteByTargetCI(ctx context.Context, targetCIID uuid.UUID) error
}
