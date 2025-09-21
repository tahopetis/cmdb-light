-- +goose Up
-- Initial Schema for CMDB Lite
-- This file contains the initial database schema creation

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'viewer',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Configuration Items table
CREATE TABLE IF NOT EXISTS configuration_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    type VARCHAR(100) NOT NULL,
    attributes JSONB,
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Relationships table
CREATE TABLE IF NOT EXISTS relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    source_id UUID NOT NULL REFERENCES configuration_items(id) ON DELETE CASCADE,
    target_id UUID NOT NULL REFERENCES configuration_items(id) ON DELETE CASCADE,
    type VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(source_id, target_id, type)
);

-- Audit Logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    changed_by VARCHAR(50) NOT NULL,
    changed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    details JSONB
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_configuration_items_type ON configuration_items(type);
CREATE INDEX IF NOT EXISTS idx_configuration_items_tags ON configuration_items USING GIN(tags);
CREATE INDEX IF NOT EXISTS idx_relationships_source_id ON relationships(source_id);
CREATE INDEX IF NOT EXISTS idx_relationships_target_id ON relationships(target_id);
CREATE INDEX IF NOT EXISTS idx_relationships_type ON relationships(type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_changed_by ON audit_logs(changed_by);

-- +goose Down
-- Drop the initial schema for CMDB Lite

-- Drop indexes
DROP INDEX IF EXISTS idx_configuration_items_type;
DROP INDEX IF EXISTS idx_configuration_items_tags;
DROP INDEX IF NOT EXISTS idx_relationships_source_id;
DROP INDEX IF NOT EXISTS idx_relationships_target_id;
DROP INDEX IF NOT EXISTS idx_relationships_type;
DROP INDEX IF NOT EXISTS idx_audit_logs_entity_type;
DROP INDEX IF NOT EXISTS idx_audit_logs_entity_id;
DROP INDEX IF NOT EXISTS idx_audit_logs_changed_by;

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS relationships;
DROP TABLE IF EXISTS configuration_items;
DROP TABLE IF EXISTS users;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
