package repositories

import (
	"context"
	"errors"
	
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
)

// Error constants
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository defines the interface for user repository operations
type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}