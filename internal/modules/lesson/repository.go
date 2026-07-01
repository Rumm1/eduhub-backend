package lesson

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, lesson Lesson) (Lesson, error) {
	err := r.db.QueryRow(ctx, `
SELECT
branch_id,
subject_id,
COALESCE(teacher_id::text, '')
FROM groups
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, lesson.GroupID, lesson.OrganizationID).Scan(
		&lesson.BranchID,
		&lesson.SubjectID,
		&lesson.TeacherID,
	)
	if err != nil {
		return Lesson{}, ErrGroupNotFound
	}

	var teacherID interface{}
	if lesson.TeacherID != "" {
		teacherID = lesson.TeacherID
		lesson.PlannedTeacherID = lesson.TeacherID
		lesson.ActualTeacherID = lesson.TeacherID
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO lessons (
id,
organization_id,
branch_id,
group_id,
teacher_id,
planned_teacher_id,
actual_teacher_id,
subject_id,
lesson_date,
start_time,
end_time,
topic,
status
)
VALUES ($1, $2, $3, $4, $5::uuid, $5::uuid, $5::uuid, $6, $7::date, $8::time, $9::time, $10, $11)
RETURNING
id,
organization_id,
branch_id,
group_id,
COALESCE(teacher_id::text, ''),
COALESCE(planned_teacher_id::text, ''),
COALESCE(actual_teacher_id::text, ''),
subject_id,
COALESCE(schedule_id::text, ''),
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status,
COALESCE(substitution_reason, '')
`,
		lesson.ID,
		lesson.OrganizationID,
		lesson.BranchID,
		lesson.GroupID,
		teacherID,
		lesson.SubjectID,
		lesson.LessonDate,
		lesson.StartTime,
		lesson.EndTime,
		lesson.Topic,
		lesson.Status,
	).Scan(
		&lesson.ID,
		&lesson.OrganizationID,
		&lesson.BranchID,
		&lesson.GroupID,
		&lesson.TeacherID,
		&lesson.PlannedTeacherID,
		&lesson.ActualTeacherID,
		&lesson.SubjectID,
		&lesson.ScheduleID,
		&lesson.LessonDate,
		&lesson.StartTime,
		&lesson.EndTime,
		&lesson.Topic,
		&lesson.Status,
		&lesson.SubstitutionReason,
	)
	if err != nil {
		return Lesson{}, err
	}

	return lesson, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Lesson, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
group_id,
COALESCE(teacher_id::text, ''),
COALESCE(planned_teacher_id::text, ''),
COALESCE(actual_teacher_id::text, ''),
subject_id,
COALESCE(schedule_id::text, ''),
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status,
COALESCE(substitution_reason, '')
FROM lessons
WHERE organization_id = $1
ORDER BY lesson_date DESC, start_time DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lessons := make([]Lesson, 0)

	for rows.Next() {
		var item Lesson

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.BranchID,
			&item.GroupID,
			&item.TeacherID,
			&item.PlannedTeacherID,
			&item.ActualTeacherID,
			&item.SubjectID,
			&item.ScheduleID,
			&item.LessonDate,
			&item.StartTime,
			&item.EndTime,
			&item.Topic,
			&item.Status,
			&item.SubstitutionReason,
		); err != nil {
			return nil, err
		}

		lessons = append(lessons, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return lessons, nil
}

func (r *Repository) UpdateActualTeacher(
	ctx context.Context,
	organizationID uuid.UUID,
	lessonID uuid.UUID,
	actualTeacherID uuid.UUID,
	reason string,
) (Lesson, error) {
	var teacherExists bool

	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM teacher_profiles
WHERE user_id = $1
  AND organization_id = $2
)
`, actualTeacherID, organizationID).Scan(&teacherExists)
	if err != nil {
		return Lesson{}, err
	}

	if !teacherExists {
		return Lesson{}, ErrActualTeacherNotFound
	}

	var lesson Lesson

	err = r.db.QueryRow(ctx, `
UPDATE lessons
SET
teacher_id = $3,
actual_teacher_id = $3,
substitution_reason = $4,
updated_at = now()
WHERE id = $1
  AND organization_id = $2
RETURNING
id,
organization_id,
branch_id,
group_id,
COALESCE(teacher_id::text, ''),
COALESCE(planned_teacher_id::text, ''),
COALESCE(actual_teacher_id::text, ''),
subject_id,
COALESCE(schedule_id::text, ''),
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status,
COALESCE(substitution_reason, '')
`, lessonID, organizationID, actualTeacherID, reason).Scan(
		&lesson.ID,
		&lesson.OrganizationID,
		&lesson.BranchID,
		&lesson.GroupID,
		&lesson.TeacherID,
		&lesson.PlannedTeacherID,
		&lesson.ActualTeacherID,
		&lesson.SubjectID,
		&lesson.ScheduleID,
		&lesson.LessonDate,
		&lesson.StartTime,
		&lesson.EndTime,
		&lesson.Topic,
		&lesson.Status,
		&lesson.SubstitutionReason,
	)
	if err != nil {
		return Lesson{}, ErrLessonNotFound
	}

	return lesson, nil
}
