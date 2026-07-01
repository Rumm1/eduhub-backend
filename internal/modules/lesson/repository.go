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
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO lessons (
id,
organization_id,
branch_id,
group_id,
teacher_id,
subject_id,
lesson_date,
start_time,
end_time,
topic,
status
)
VALUES ($1, $2, $3, $4, $5::uuid, $6, $7::date, $8::time, $9::time, $10, $11)
RETURNING
id,
organization_id,
branch_id,
group_id,
COALESCE(teacher_id::text, ''),
subject_id,
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status
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
		&lesson.SubjectID,
		&lesson.LessonDate,
		&lesson.StartTime,
		&lesson.EndTime,
		&lesson.Topic,
		&lesson.Status,
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
subject_id,
lesson_date::text,
start_time::text,
end_time::text,
COALESCE(topic, ''),
status
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
			&item.SubjectID,
			&item.LessonDate,
			&item.StartTime,
			&item.EndTime,
			&item.Topic,
			&item.Status,
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
