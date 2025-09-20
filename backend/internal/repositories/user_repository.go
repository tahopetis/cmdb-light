package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// Error constants
var (
	ErrUserNotFound = errors.New("user not found")
)

// UserPostgresRepository implements the UserRepository interface for PostgreSQL
type UserPostgresRepository struct {
	db *sqlx.DB
}

// NewUserPostgresRepository creates a new UserPostgresRepository
func NewUserPostgresRepository(db *sqlx.DB) *UserPostgresRepository {
	return &UserPostgresRepository{db: db}
}

// Create creates a new user in the database
func (r *UserPostgresRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (id, username, password_hash, role, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// GetByID retrieves a user by ID
func (r *UserPostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at
		FROM users
		WHERE id = $1
	`
	
	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserPostgresRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, password_hash, role, created_at
		FROM users
		WHERE username = $1
	`
	
	var user models.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	
	return &user, nil
}

// Update updates a user in the database
func (r *UserPostgresRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $2, password_hash = $3, role = $4
		WHERE id = $1
	`
	
	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.PasswordHash,
		user.Role,
	)
	
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}

// Delete deletes a user from the database
func (r *UserPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return ErrUserNotFound
	}
	
	return nil
}