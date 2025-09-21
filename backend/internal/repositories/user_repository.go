package repositories

import (
	"context"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	// Create creates a new user in the database
	Create(ctx context.Context, user *models.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*models.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*models.User, error)

	// GetAll retrieves all users from the database
	GetAll(ctx context.Context) ([]*models.User, error)

	// Update updates a user in the database
	Update(ctx context.Context, user *models.User) error

	// Delete deletes a user from the database
	Delete(ctx context.Context, id uuid.UUID) error

	// UpdateLastLogin updates the last login timestamp for a user
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}
