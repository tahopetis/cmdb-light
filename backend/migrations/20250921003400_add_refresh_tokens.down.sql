-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Drop indexes
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_token_hash;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

-- Drop refresh tokens table
DROP TABLE IF EXISTS refresh_tokens;