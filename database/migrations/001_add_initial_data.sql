-- Migration to add initial data to CMDB Lite database

-- Insert default admin user (password should be changed in production)
-- Password is 'admin123' hashed with bcrypt
INSERT INTO users (username, password_hash, role) VALUES
('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin'),
('viewer', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'viewer');

-- Insert sample configuration items
INSERT INTO configuration_items (name, type, attributes, tags) VALUES
('Web Server 01', 'server', '{"cpu": "4 cores", "memory": "16GB", "storage": "500GB SSD", "os": "Ubuntu 20.04"}', ARRAY['web', 'production', 'linux']),
('Database Server', 'server', '{"cpu": "8 cores", "memory": "32GB", "storage": "1TB SSD", "os": "Ubuntu 20.04"}', ARRAY['database', 'production', 'linux']),
('Application Database', 'database', '{"type": "PostgreSQL", "version": "13.4", "size": "500GB"}', ARRAY['postgres', 'production', 'primary']),
('Web Application', 'application', '{"version": "1.2.0", "language": "JavaScript", "framework": "Vue.js"}', ARRAY['frontend', 'production', 'web']),
('Core Router', 'network', '{"model": "Cisco ISR 4000", "ip_address": "192.168.1.1", "mac_address": "00:1B:44:11:3A:B7"}', ARRAY['network', 'production', 'cisco']),
('Primary Storage', 'storage', '{"capacity": "10TB", "type": "NAS", "raid_level": "RAID 5"}', ARRAY['storage', 'production', 'nas']);

-- Insert sample relationships
-- Note: We need to reference the UUIDs of the configuration items we just inserted
-- For simplicity in seeding, we'll use a CTE to get the IDs
WITH ci_ids AS (
  SELECT id, name FROM configuration_items
)
INSERT INTO relationships (source_id, target_id, type)
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Web Server 01'),
  (SELECT id FROM ci_ids WHERE name = 'Web Application'),
  'hosts'
UNION ALL
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Database Server'),
  (SELECT id FROM ci_ids WHERE name = 'Application Database'),
  'hosts'
UNION ALL
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Web Application'),
  (SELECT id FROM ci_ids WHERE name = 'Application Database'),
  'connects_to'
UNION ALL
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Web Server 01'),
  (SELECT id FROM ci_ids WHERE name = 'Core Router'),
  'connects_to'
UNION ALL
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Database Server'),
  (SELECT id FROM ci_ids WHERE name = 'Core Router'),
  'connects_to'
UNION ALL
SELECT
  (SELECT id FROM ci_ids WHERE name = 'Database Server'),
  (SELECT id FROM ci_ids WHERE name = 'Primary Storage'),
  'connects_to';

-- Insert sample audit logs
WITH ci_ids AS (
  SELECT id, name FROM configuration_items
)
INSERT INTO audit_logs (entity_type, entity_id, action, changed_by, details)
SELECT
  'configuration_item',
  id,
  'create',
  'admin',
  jsonb_build_object('name', name, 'type', type)
FROM configuration_items
LIMIT 6;