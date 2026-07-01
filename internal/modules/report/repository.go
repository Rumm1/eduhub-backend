package report

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

func (r *Repository) GetTeacherSchedule(
	ctx context.Context,
	organizationID uuid.UUID,
	teacherID uuid.UUID,
	fromDate string,
	toDate string,
) (TeacherScheduleReport, error) {
	var teacherName string

	err := r.db.QueryRow(ctx, `
SELECT full_name
FROM users
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, teacherID, organizationID).Scan(&teacherName)
	if err != nil {
		return TeacherScheduleReport{}, ErrTeacherNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
l.id::text,
l.lesson_date::text,
l.start_time::text,
l.end_time::text,
ROUND((EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600)::numeric, 2)::text AS hours,
COALESCE(l.topic, ''),
l.status,
g.id::text,
g.name,
b.id::text,
b.name,
s.id::text,
s.name,
COALESCE(l.planned_teacher_id::text, ''),
COALESCE(planned_user.full_name, ''),
COALESCE(l.actual_teacher_id::text, ''),
COALESCE(actual_user.full_name, ''),
(
COALESCE(l.actual_teacher_id, l.teacher_id) <> COALESCE(l.planned_teacher_id, l.teacher_id)
) AS is_substitution,
CASE
WHEN COALESCE(l.actual_teacher_id, l.teacher_id) = $2 THEN 'actual'
WHEN l.planned_teacher_id = $2 THEN 'planned_only'
ELSE 'unknown'
END AS teacher_role_in_lesson,
COALESCE(l.substitution_reason, '')
FROM lessons l
JOIN groups g ON g.id = l.group_id
JOIN branches b ON b.id = l.branch_id
JOIN subjects s ON s.id = l.subject_id
LEFT JOIN users planned_user ON planned_user.id = l.planned_teacher_id
LEFT JOIN users actual_user ON actual_user.id = l.actual_teacher_id
WHERE l.organization_id = $1
  AND l.lesson_date >= $3::date
  AND l.lesson_date <= $4::date
  AND (
COALESCE(l.actual_teacher_id, l.teacher_id) = $2
OR l.planned_teacher_id = $2
  )
ORDER BY l.lesson_date ASC, l.start_time ASC
`, organizationID, teacherID, fromDate, toDate)
	if err != nil {
		return TeacherScheduleReport{}, err
	}
	defer rows.Close()

	items := make([]TeacherScheduleItem, 0)
	actualLessons := 0
	plannedOnlyLessons := 0
	substitutions := 0
	totalActualHours := "0"

	for rows.Next() {
		var item TeacherScheduleItem

		if err := rows.Scan(
			&item.LessonID,
			&item.LessonDate,
			&item.StartTime,
			&item.EndTime,
			&item.Hours,
			&item.Topic,
			&item.Status,
			&item.GroupID,
			&item.GroupName,
			&item.BranchID,
			&item.BranchName,
			&item.SubjectID,
			&item.SubjectName,
			&item.PlannedTeacherID,
			&item.PlannedTeacherName,
			&item.ActualTeacherID,
			&item.ActualTeacherName,
			&item.IsSubstitution,
			&item.TeacherRoleInLesson,
			&item.SubstitutionReason,
		); err != nil {
			return TeacherScheduleReport{}, err
		}

		if item.TeacherRoleInLesson == "actual" {
			actualLessons++
		}

		if item.TeacherRoleInLesson == "planned_only" {
			plannedOnlyLessons++
		}

		if item.IsSubstitution {
			substitutions++
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return TeacherScheduleReport{}, err
	}

	err = r.db.QueryRow(ctx, `
SELECT
ROUND(COALESCE(SUM(EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600), 0)::numeric, 2)::text
FROM lessons l
WHERE l.organization_id = $1
  AND l.lesson_date >= $3::date
  AND l.lesson_date <= $4::date
  AND COALESCE(l.actual_teacher_id, l.teacher_id) = $2
  AND l.status <> 'cancelled'
`, organizationID, teacherID, fromDate, toDate).Scan(&totalActualHours)
	if err != nil {
		return TeacherScheduleReport{}, err
	}

	return TeacherScheduleReport{
		TeacherID:          teacherID.String(),
		TeacherName:        teacherName,
		FromDate:           fromDate,
		ToDate:             toDate,
		TotalLessons:       len(items),
		ActualLessons:      actualLessons,
		PlannedOnlyLessons: plannedOnlyLessons,
		Substitutions:      substitutions,
		TotalActualHours:   totalActualHours,
		Items:              items,
	}, nil
}
