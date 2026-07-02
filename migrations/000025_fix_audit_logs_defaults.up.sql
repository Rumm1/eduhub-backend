CREATE EXTENSION IF NOT EXISTS pgcrypto;

ALTER TABLE audit_logs
ALTER COLUMN id SET DEFAULT gen_random_uuid();

ALTER TABLE audit_logs
ALTER COLUMN created_at SET DEFAULT now();

ALTER TABLE audit_logs
ALTER COLUMN metadata SET DEFAULT '{}'::jsonb;
