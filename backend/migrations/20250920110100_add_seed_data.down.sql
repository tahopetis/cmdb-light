-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

-- Delete audit logs
DELETE FROM audit_logs;

-- Delete relationships
DELETE FROM relationships;

-- Delete configuration items
DELETE FROM configuration_items;

-- Delete users
DELETE FROM users;