package repositories

import (
	"context"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// CIRepository defines the interface for CI (Configuration Item) repository operations
type CIRepository interface {
	// Create creates a new CI in the database
	Create(ctx context.Context, ci *models.CI) error

	// GetByID retrieves a CI by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.CI, error)

	// GetByName retrieves a CI by name
	GetByName(ctx context.Context, name string) (*models.CI, error)

	// GetByType retrieves CIs by type
	GetByType(ctx context.Context, ciType string) ([]*models.CI, error)

	// GetAll retrieves all CIs from the database
	GetAll(ctx context.Context) ([]*models.CI, error)

	// Update updates a CI in the database
	Update(ctx context.Context, ci *models.CI) error

	// Delete deletes a CI from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByStatus retrieves CIs by status
	GetByStatus(ctx context.Context, status string) ([]*models.CI, error)
}
