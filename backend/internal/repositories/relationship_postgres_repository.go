package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// RelationshipPostgresRepository implements the RelationshipRepository interface for PostgreSQL
type RelationshipPostgresRepository struct {
	db *sqlx.DB
}

// NewRelationshipPostgresRepository creates a new RelationshipPostgresRepository
func NewRelationshipPostgresRepository(db *sqlx.DB) *RelationshipPostgresRepository {
	return &RelationshipPostgresRepository{db: db}
}

// Create creates a new relationship in the database
func (r *RelationshipPostgresRepository) Create(ctx context.Context, relationship *models.Relationship) error {
	query := `
		INSERT INTO relationships (id, source_id, target_id, type, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		relationship.ID,
		relationship.SourceID,
		relationship.TargetID,
		relationship.Type,
		relationship.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves a relationship by ID
func (r *RelationshipPostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		WHERE id = $1
	`

	var relationship models.Relationship
	err := r.db.GetContext(ctx, &relationship, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("relationship not found")
		}
		return nil, err
	}

	return &relationship, nil
}

// GetBySourceCI retrieves relationships by source CI ID
func (r *RelationshipPostgresRepository) GetBySourceCI(ctx context.Context, sourceCIID uuid.UUID) ([]*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		WHERE source_id = $1
		ORDER BY created_at DESC
	`

	var relationships []*models.Relationship
	err := r.db.SelectContext(ctx, &relationships, query, sourceCIID)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetByTargetCI retrieves relationships by target CI ID
func (r *RelationshipPostgresRepository) GetByTargetCI(ctx context.Context, targetCIID uuid.UUID) ([]*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		WHERE target_id = $1
		ORDER BY created_at DESC
	`

	var relationships []*models.Relationship
	err := r.db.SelectContext(ctx, &relationships, query, targetCIID)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetBySourceAndTarget retrieves relationships by source and target CI IDs
func (r *RelationshipPostgresRepository) GetBySourceAndTarget(ctx context.Context, sourceCIID, targetCIID uuid.UUID) ([]*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		WHERE source_id = $1 AND target_id = $2
		ORDER BY created_at DESC
	`

	var relationships []*models.Relationship
	err := r.db.SelectContext(ctx, &relationships, query, sourceCIID, targetCIID)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetByType retrieves relationships by type
func (r *RelationshipPostgresRepository) GetByType(ctx context.Context, relationshipType string) ([]*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		WHERE type = $1
		ORDER BY created_at DESC
	`

	var relationships []*models.Relationship
	err := r.db.SelectContext(ctx, &relationships, query, relationshipType)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// GetAll retrieves all relationships from the database
func (r *RelationshipPostgresRepository) GetAll(ctx context.Context) ([]*models.Relationship, error) {
	query := `
		SELECT id, source_id, target_id, type, created_at
		FROM relationships
		ORDER BY created_at DESC
	`

	var relationships []*models.Relationship
	err := r.db.SelectContext(ctx, &relationships, query)
	if err != nil {
		return nil, err
	}

	return relationships, nil
}

// Update updates a relationship in the database
func (r *RelationshipPostgresRepository) Update(ctx context.Context, relationship *models.Relationship) error {
	query := `
		UPDATE relationships
		SET source_id = $2, target_id = $3, type = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		relationship.ID,
		relationship.SourceID,
		relationship.TargetID,
		relationship.Type,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("relationship not found")
	}

	return nil
}

// Delete deletes a relationship from the database
func (r *RelationshipPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM relationships WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("relationship not found")
	}

	return nil
}

// DeleteBySourceCI deletes all relationships for a source CI
func (r *RelationshipPostgresRepository) DeleteBySourceCI(ctx context.Context, sourceCIID uuid.UUID) error {
	query := `DELETE FROM relationships WHERE source_id = $1`

	_, err := r.db.ExecContext(ctx, query, sourceCIID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteByTargetCI deletes all relationships for a target CI
func (r *RelationshipPostgresRepository) DeleteByTargetCI(ctx context.Context, targetCIID uuid.UUID) error {
	query := `DELETE FROM relationships WHERE target_id = $1`

	_, err := r.db.ExecContext(ctx, query, targetCIID)
	if err != nil {
		return err
	}

	return nil
}
