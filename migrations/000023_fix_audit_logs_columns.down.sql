DROP INDEX IF EXISTS idx_audit_logs_organization_id;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_entity;
DROP INDEX IF EXISTS idx_audit_logs_created_at;

ALTER TABLE audit_logs
DROP COLUMN IF EXISTS description;

ALTER TABLE audit_logs
DROP COLUMN IF EXISTS metadata;

ALTER TABLE audit_logs
DROP COLUMN IF EXISTS ip_address;

ALTER TABLE audit_logs
DROP COLUMN IF EXISTS user_agent;
