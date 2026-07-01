package attendance

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

func (r *Repository) MarkLessonAttendance(
	ctx context.Context,
	organizationID uuid.UUID,
	lessonID uuid.UUID,
	markedBy uuid.UUID,
	items []Attendance,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var groupID uuid.UUID

	err = tx.QueryRow(ctx, `
SELECT group_id
FROM lessons
WHERE id = $1
  AND organization_id = $2
`, lessonID, organizationID).Scan(&groupID)
	if err != nil {
		return ErrLessonNotFound
	}

	for _, item := range items {
		var studentExists bool

		err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM group_students gs
JOIN students s ON s.id = gs.student_id
WHERE gs.group_id = $1
  AND gs.student_id = $2
  AND gs.status = 'active'
  AND s.organization_id = $3
  AND s.status = 'active'
)
`, groupID, item.StudentID, organizationID).Scan(&studentExists)
		if err != nil {
			return err
		}

		if !studentExists {
			return ErrStudentNotInLessonGroup
		}

		_, err = tx.Exec(ctx, `
INSERT INTO attendance (
id,
lesson_id,
student_id,
status,
reason,
comment,
marked_by,
marked_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, now())
ON CONFLICT (lesson_id, student_id)
DO UPDATE SET
status = EXCLUDED.status,
reason = EXCLUDED.reason,
comment = EXCLUDED.comment,
marked_by = EXCLUDED.marked_by,
marked_at = now(),
updated_at = now()
`,
			uuid.New(),
			lessonID,
			item.StudentID,
			item.Status,
			item.Reason,
			item.Comment,
			markedBy,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListByLessonID(
	ctx context.Context,
	organizationID uuid.UUID,
	lessonID uuid.UUID,
) ([]Attendance, error) {
	var groupID uuid.UUID

	err := r.db.QueryRow(ctx, `
SELECT group_id
FROM lessons
WHERE id = $1
  AND organization_id = $2
`, lessonID, organizationID).Scan(&groupID)
	if err != nil {
		return nil, ErrLessonNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
COALESCE(a.id::text, ''),
l.id,
s.id,
s.full_name,
COALESCE(a.status, 'unmarked'),
COALESCE(a.reason, ''),
COALESCE(a.comment, ''),
COALESCE(a.marked_by::text, ''),
COALESCE(a.marked_at::text, '')
FROM lessons l
JOIN group_students gs ON gs.group_id = l.group_id AND gs.status = 'active'
JOIN students s ON s.id = gs.student_id AND s.status = 'active'
LEFT JOIN attendance a ON a.lesson_id = l.id AND a.student_id = s.id
WHERE l.id = $1
  AND l.organization_id = $2
  AND s.organization_id = $2
ORDER BY s.full_name ASC
`, lessonID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]Attendance, 0)

	for rows.Next() {
		var item Attendance

		if err := rows.Scan(
			&item.ID,
			&item.LessonID,
			&item.StudentID,
			&item.StudentFullName,
			&item.Status,
			&item.Reason,
			&item.Comment,
			&item.MarkedBy,
			&item.MarkedAt,
		); err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
