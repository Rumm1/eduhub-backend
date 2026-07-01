DROP INDEX IF EXISTS idx_lessons_actual_teacher_id;
DROP INDEX IF EXISTS idx_lessons_planned_teacher_id;

ALTER TABLE lessons
DROP COLUMN IF EXISTS substitution_reason,
DROP COLUMN IF EXISTS actual_teacher_id,
DROP COLUMN IF EXISTS planned_teacher_id;
