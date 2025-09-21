package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/cmdb-lite/backend/internal/auth"
	"github.com/cmdb-lite/backend/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	Create(ctx context.Context, refreshToken *models.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.RefreshToken, error)
	Revoke(ctx context.Context, tokenID uuid.UUID) error
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
	CleanExpired(ctx context.Context) error
}

// RefreshTokenPostgresRepository implements RefreshTokenRepository for PostgreSQL
type RefreshTokenPostgresRepository struct {
	db *sqlx.DB
}

// NewRefreshTokenPostgresRepository creates a new RefreshTokenPostgresRepository
func NewRefreshTokenPostgresRepository(db *sqlx.DB) *RefreshTokenPostgresRepository {
	return &RefreshTokenPostgresRepository{db: db}
}

// Create creates a new refresh token in the database
func (r *RefreshTokenPostgresRepository) Create(ctx context.Context, refreshToken *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at, revoked_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(ctx, query,
		refreshToken.ID,
		refreshToken.UserID,
		refreshToken.TokenHash,
		refreshToken.ExpiresAt,
		refreshToken.CreatedAt,
		refreshToken.RevokedAt,
	)
	return err
}

// GetByTokenHash retrieves a refresh token by its hash
func (r *RefreshTokenPostgresRepository) GetByTokenHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token_hash = $1
	`
	var refreshToken models.RefreshToken
	err := r.db.GetContext(ctx, &refreshToken, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &refreshToken, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *RefreshTokenPostgresRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.RefreshToken, error) {
	query := `
		SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1
	`
	var refreshTokens []*models.RefreshToken
	err := r.db.SelectContext(ctx, &refreshTokens, query, userID)
	if err != nil {
		return nil, err
	}
	return refreshTokens, nil
}

// Revoke revokes a refresh token by setting the revoked_at timestamp
func (r *RefreshTokenPostgresRepository) Revoke(ctx context.Context, tokenID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE id = $2 AND revoked_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query, now, tokenID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenPostgresRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	now := time.Now()
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE user_id = $2 AND revoked_at IS NULL
	`
	_, err := r.db.ExecContext(ctx, query, now, userID)
	return err
}

// CleanExpired deletes all expired refresh tokens
func (r *RefreshTokenPostgresRepository) CleanExpired(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1 OR (revoked_at IS NOT NULL AND revoked_at < $2)
	`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now)
	return err
}