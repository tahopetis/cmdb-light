-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Drop triggers
DROP TRIGGER IF EXISTS update_configuration_items_updated_at ON configuration_items;

-- Drop trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_audit_logs_changed_by;
DROP INDEX IF EXISTS idx_audit_logs_entity_id;
DROP INDEX IF EXISTS idx_audit_logs_entity_type;
DROP INDEX IF EXISTS idx_relationships_type;
DROP INDEX IF EXISTS idx_relationships_target_id;
DROP INDEX IF EXISTS idx_relationships_source_id;
DROP INDEX IF EXISTS idx_configuration_items_tags;
DROP INDEX IF EXISTS idx_configuration_items_type;

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS relationships;
DROP TABLE IF EXISTS configuration_items;
DROP TABLE IF EXISTS users;

-- Drop UUID extension
DROP EXTENSION IF EXISTS "uuid-ossp";