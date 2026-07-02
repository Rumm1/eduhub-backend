ALTER TABLE users
DROP COLUMN IF EXISTS avatar_path;

ALTER TABLE organizations
DROP COLUMN IF EXISTS logo_path;
