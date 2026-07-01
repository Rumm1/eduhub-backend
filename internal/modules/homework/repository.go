package homework

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

func (r *Repository) Create(ctx context.Context, homework Homework) (Homework, error) {
	var groupHomeworkEnabled bool
	var subjectHomeworkEnabled bool

	err := r.db.QueryRow(ctx, `
SELECT
l.group_id,
COALESCE(l.teacher_id::text, ''),
COALESCE(g.homework_enabled, true),
COALESCE(s.homework_enabled, true)
FROM lessons l
JOIN groups g ON g.id = l.group_id
JOIN subjects s ON s.id = l.subject_id
WHERE l.id = $1
  AND l.organization_id = $2
`, homework.LessonID, homework.OrganizationID).Scan(
		&homework.GroupID,
		&homework.TeacherID,
		&groupHomeworkEnabled,
		&subjectHomeworkEnabled,
	)
	if err != nil {
		return Homework{}, ErrLessonNotFound
	}

	if homework.TeacherID == "" {
		return Homework{}, ErrTeacherNotFound
	}

	if !groupHomeworkEnabled || !subjectHomeworkEnabled {
		return Homework{}, ErrHomeworkDisabled
	}

	var dueDate interface{}
	if homework.DueDate != "" {
		dueDate = homework.DueDate
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO homeworks (
id,
organization_id,
group_id,
lesson_id,
teacher_id,
title,
description,
due_date
)
VALUES ($1, $2, $3, $4, $5::uuid, $6, $7, $8::date)
RETURNING
id,
organization_id,
group_id,
lesson_id,
teacher_id::text,
title,
COALESCE(description, ''),
COALESCE(due_date::text, '')
`,
		homework.ID,
		homework.OrganizationID,
		homework.GroupID,
		homework.LessonID,
		homework.TeacherID,
		homework.Title,
		homework.Description,
		dueDate,
	).Scan(
		&homework.ID,
		&homework.OrganizationID,
		&homework.GroupID,
		&homework.LessonID,
		&homework.TeacherID,
		&homework.Title,
		&homework.Description,
		&homework.DueDate,
	)
	if err != nil {
		return Homework{}, err
	}

	return homework, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Homework, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
group_id,
lesson_id,
teacher_id::text,
title,
COALESCE(description, ''),
COALESCE(due_date::text, '')
FROM homeworks
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	homeworks := make([]Homework, 0)

	for rows.Next() {
		var item Homework

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.GroupID,
			&item.LessonID,
			&item.TeacherID,
			&item.Title,
			&item.Description,
			&item.DueDate,
		); err != nil {
			return nil, err
		}

		homeworks = append(homeworks, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (r *Repository) ListByLessonID(ctx context.Context, organizationID uuid.UUID, lessonID uuid.UUID) ([]Homework, error) {
	var lessonExists bool

	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM lessons
WHERE id = $1
  AND organization_id = $2
)
`, lessonID, organizationID).Scan(&lessonExists)
	if err != nil {
		return nil, err
	}

	if !lessonExists {
		return nil, ErrLessonNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
group_id,
lesson_id,
teacher_id::text,
title,
COALESCE(description, ''),
COALESCE(due_date::text, '')
FROM homeworks
WHERE organization_id = $1
  AND lesson_id = $2
ORDER BY created_at DESC
`, organizationID, lessonID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	homeworks := make([]Homework, 0)

	for rows.Next() {
		var item Homework

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.GroupID,
			&item.LessonID,
			&item.TeacherID,
			&item.Title,
			&item.Description,
			&item.DueDate,
		); err != nil {
			return nil, err
		}

		homeworks = append(homeworks, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return homeworks, nil
}
