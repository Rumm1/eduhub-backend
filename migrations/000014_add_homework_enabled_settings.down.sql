ALTER TABLE groups
DROP COLUMN IF EXISTS homework_enabled;

ALTER TABLE subjects
DROP COLUMN IF EXISTS homework_enabled;
