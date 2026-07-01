ALTER TABLE lessons
ADD COLUMN IF NOT EXISTS schedule_id UUID REFERENCES schedules(id) ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS idx_lessons_schedule_date_unique
ON lessons(schedule_id, lesson_date)
WHERE schedule_id IS NOT NULL;
