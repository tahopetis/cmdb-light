package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
		INSERT INTO users (id, username, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a user by ID
func (r *UserPostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		WHERE id = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserPostgresRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		WHERE username = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserPostgresRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		WHERE email = $1
	`

	var user models.User
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetAll retrieves all users from the database
func (r *UserPostgresRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at, last_login
		FROM users
		ORDER BY created_at DESC
	`

	var users []*models.User
	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Update updates a user in the database
func (r *UserPostgresRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $2, email = $3, password_hash = $4, role = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
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
		return errors.New("user not found")
	}

	return nil
}

// UpdateLastLogin updates the last login timestamp for a user
func (r *UserPostgresRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users SET last_login = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}
