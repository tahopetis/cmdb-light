package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// AuditLogPostgresRepository implements the AuditLogRepository interface for PostgreSQL
type AuditLogPostgresRepository struct {
	db *sqlx.DB
}

// NewAuditLogPostgresRepository creates a new AuditLogPostgresRepository
func NewAuditLogPostgresRepository(db *sqlx.DB) *AuditLogPostgresRepository {
	return &AuditLogPostgresRepository{db: db}
}

// Create creates a new audit log in the database
func (r *AuditLogPostgresRepository) Create(ctx context.Context, auditLog *models.AuditLog) error {
	query := `
		INSERT INTO audit_logs (id, entity_type, entity_id, action, changed_by, changed_at, details)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.ExecContext(ctx, query,
		auditLog.ID,
		auditLog.EntityType,
		auditLog.EntityID,
		auditLog.Action,
		auditLog.ChangedBy,
		auditLog.ChangedAt,
		auditLog.Details,
	)
	
	if err != nil {
		return err
	}
	
	return nil
}

// GetByID retrieves an audit log by ID
func (r *AuditLogPostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, changed_by, changed_at, details
		FROM audit_logs
		WHERE id = $1
	`
	
	var auditLog models.AuditLog
	err := r.db.GetContext(ctx, &auditLog, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("audit log not found")
		}
		return nil, err
	}
	
	return &auditLog, nil
}

// GetAll retrieves all audit logs from the database
func (r *AuditLogPostgresRepository) GetAll(ctx context.Context) ([]*models.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, changed_by, changed_at, details
		FROM audit_logs
		ORDER BY changed_at DESC
	`
	
	var auditLogs []*models.AuditLog
	err := r.db.SelectContext(ctx, &auditLogs, query)
	if err != nil {
		return nil, err
	}
	
	return auditLogs, nil
}

// GetByEntityType retrieves audit logs by entity type
func (r *AuditLogPostgresRepository) GetByEntityType(ctx context.Context, entityType string) ([]*models.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, changed_by, changed_at, details
		FROM audit_logs
		WHERE entity_type = $1
		ORDER BY changed_at DESC
	`
	
	var auditLogs []*models.AuditLog
	err := r.db.SelectContext(ctx, &auditLogs, query, entityType)
	if err != nil {
		return nil, err
	}
	
	return auditLogs, nil
}

// GetByEntityID retrieves audit logs by entity ID
func (r *AuditLogPostgresRepository) GetByEntityID(ctx context.Context, entityID uuid.UUID) ([]*models.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, changed_by, changed_at, details
		FROM audit_logs
		WHERE entity_id = $1
		ORDER BY changed_at DESC
	`
	
	var auditLogs []*models.AuditLog
	err := r.db.SelectContext(ctx, &auditLogs, query, entityID)
	if err != nil {
		return nil, err
	}
	
	return auditLogs, nil
}

// GetByChangedBy retrieves audit logs by the user who made the change
func (r *AuditLogPostgresRepository) GetByChangedBy(ctx context.Context, changedBy string) ([]*models.AuditLog, error) {
	query := `
		SELECT id, entity_type, entity_id, action, changed_by, changed_at, details
		FROM audit_logs
		WHERE changed_by = $1
		ORDER BY changed_at DESC
	`
	
	var auditLogs []*models.AuditLog
	err := r.db.SelectContext(ctx, &auditLogs, query, changedBy)
	if err != nil {
		return nil, err
	}
	
	return auditLogs, nil
}

// Delete deletes an audit log from the database
func (r *AuditLogPostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM audit_logs WHERE id = $1`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	
	if rowsAffected == 0 {
		return errors.New("audit log not found")
	}
	
	return nil
}