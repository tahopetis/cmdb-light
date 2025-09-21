package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// CIPostgresRepository implements the CIRepository interface for PostgreSQL
type CIPostgresRepository struct {
	db *sqlx.DB
}

// NewCIPostgresRepository creates a new CIPostgresRepository
func NewCIPostgresRepository(db *sqlx.DB) *CIPostgresRepository {
	return &CIPostgresRepository{db: db}
}

// Create creates a new CI in the database
func (r *CIPostgresRepository) Create(ctx context.Context, ci *models.CI) error {
	query := `
		INSERT INTO configuration_items (id, name, type, attributes, tags, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		ci.ID,
		ci.Name,
		ci.Type,
		ci.Attributes,
		ci.Tags,
		ci.CreatedAt,
		ci.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a CI by ID
func (r *CIPostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.CI, error) {
	query := `
		SELECT id, name, type, attributes, tags, created_at, updated_at
		FROM configuration_items
		WHERE id = $1
	`

	var ci models.CI
	err := r.db.GetContext(ctx, &ci, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("CI not found")
		}
		return nil, err
	}

	return &ci, nil
}

// GetByName retrieves a CI by name
func (r *CIPostgresRepository) GetByName(ctx context.Context, name string) (*models.CI, error) {
	query := `
		SELECT id, name, type, attributes, tags, created_at, updated_at
		FROM configuration_items
		WHERE name = $1
	`

	var ci models.CI
	err := r.db.GetContext(ctx, &ci, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("CI not found")
		}
		return nil, err
	}

	return &ci, nil
}

// GetAll retrieves all CIs from the database
func (r *CIPostgresRepository) GetAll(ctx context.Context) ([]*models.CI, error) {
	query := `
		SELECT id, name, type, attributes, tags, created_at, updated_at
		FROM configuration_items
		ORDER BY created_at DESC
	`

	var cis []*models.CI
	err := r.db.SelectContext(ctx, &cis, query)
	if err != nil {
		return nil, err
	}

	return cis, nil
}

// GetByType retrieves CIs by type
func (r *CIPostgresRepository) GetByType(ctx context.Context, ciType string) ([]*models.CI, error) {
	query := `
		SELECT id, name, type, attributes, tags, created_at, updated_at
		FROM configuration_items
		WHERE type = $1
		ORDER BY created_at DESC
	`

	var cis []*models.CI
	err := r.db.SelectContext(ctx, &cis, query, ciType)
	if err != nil {
		return nil, err
	}

	return cis, nil
}

// GetByStatus retrieves CIs by status
func (r *CIPostgresRepository) GetByStatus(ctx context.Context, status string) ([]*models.CI, error) {
	query := `
		SELECT id, name, type, attributes, tags, created_at, updated_at
		FROM configuration_items
		WHERE attributes->>'status' = $1
		ORDER BY created_at DESC
	`

	var cis []*models.CI
	err := r.db.SelectContext(ctx, &cis, query, status)
	if err != nil {
		return nil, err
	}

	return cis, nil
}

// Update updates a CI in the database
func (r *CIPostgresRepository) Update(ctx context.Context, ci *models.CI) error {
	query := `
		UPDATE configuration_items
		SET name = $2, type = $3, attributes = $4, tags = $5, updated_at = $6
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		ci.ID,
		ci.Name,
		ci.Type,
		ci.Attributes,
		ci.Tags,
		ci.UpdatedAt,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("CI not found")
	}

	return nil
}

// Delete deletes a CI from the database
func (r *CIPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM configuration_items WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("CI not found")
	}

	return nil
}
