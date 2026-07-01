DROP INDEX IF EXISTS idx_lessons_schedule_date_unique;

ALTER TABLE lessons
DROP COLUMN IF EXISTS schedule_id;
