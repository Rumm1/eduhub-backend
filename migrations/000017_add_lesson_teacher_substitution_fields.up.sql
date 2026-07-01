ALTER TABLE lessons
ADD COLUMN IF NOT EXISTS planned_teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS actual_teacher_id UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS substitution_reason TEXT;

UPDATE lessons
SET
planned_teacher_id = COALESCE(planned_teacher_id, teacher_id),
actual_teacher_id = COALESCE(actual_teacher_id, teacher_id)
WHERE teacher_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_lessons_planned_teacher_id ON lessons(planned_teacher_id);
CREATE INDEX IF NOT EXISTS idx_lessons_actual_teacher_id ON lessons(actual_teacher_id);
