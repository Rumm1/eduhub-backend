package group

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

func (r *Repository) Create(ctx context.Context, group Group) (Group, error) {
	var branchExists bool
	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM branches
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
)
`, group.BranchID, group.OrganizationID).Scan(&branchExists)
	if err != nil {
		return Group{}, err
	}

	if !branchExists {
		return Group{}, ErrBranchNotFound
	}

	var subjectExists bool
	err = r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM subjects
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
)
`, group.SubjectID, group.OrganizationID).Scan(&subjectExists)
	if err != nil {
		return Group{}, err
	}

	if !subjectExists {
		return Group{}, ErrSubjectNotFound
	}

	if group.TeacherID != "" {
		var teacherExists bool
		err = r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM teacher_profiles
WHERE user_id = $1
  AND organization_id = $2
)
`, group.TeacherID, group.OrganizationID).Scan(&teacherExists)
		if err != nil {
			return Group{}, err
		}

		if !teacherExists {
			return Group{}, ErrTeacherNotFound
		}
	}

	var teacherID interface{}
	if group.TeacherID != "" {
		teacherID = group.TeacherID
	}

	var startDate interface{}
	if group.StartDate != "" {
		startDate = group.StartDate
	}

	var endDate interface{}
	if group.EndDate != "" {
		endDate = group.EndDate
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO groups (
id,
organization_id,
branch_id,
subject_id,
teacher_id,
name,
level,
status,
max_students,
start_date,
end_date
)
VALUES ($1, $2, $3, $4, $5::uuid, $6, $7, $8, $9, $10::date, $11::date)
RETURNING
id,
organization_id,
branch_id,
subject_id,
COALESCE(teacher_id::text, ''),
name,
COALESCE(level, ''),
status,
max_students,
COALESCE(start_date::text, ''),
COALESCE(end_date::text, '')
`,
		group.ID,
		group.OrganizationID,
		group.BranchID,
		group.SubjectID,
		teacherID,
		group.Name,
		group.Level,
		group.Status,
		group.MaxStudents,
		startDate,
		endDate,
	).Scan(
		&group.ID,
		&group.OrganizationID,
		&group.BranchID,
		&group.SubjectID,
		&group.TeacherID,
		&group.Name,
		&group.Level,
		&group.Status,
		&group.MaxStudents,
		&group.StartDate,
		&group.EndDate,
	)
	if err != nil {
		return Group{}, err
	}

	return group, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Group, error) {
	rows, err := r.db.Query(ctx, `
SELECT
g.id,
g.organization_id,
g.branch_id,
g.subject_id,
COALESCE(g.teacher_id::text, ''),
g.name,
COALESCE(g.level, ''),
g.status,
g.max_students,
COALESCE(g.start_date::text, ''),
COALESCE(g.end_date::text, ''),
COUNT(gs.student_id) FILTER (WHERE gs.status = 'active') AS students_count
FROM groups g
LEFT JOIN group_students gs ON gs.group_id = g.id
WHERE g.organization_id = $1
GROUP BY g.id
ORDER BY g.created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := make([]Group, 0)

	for rows.Next() {
		var item Group

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.BranchID,
			&item.SubjectID,
			&item.TeacherID,
			&item.Name,
			&item.Level,
			&item.Status,
			&item.MaxStudents,
			&item.StartDate,
			&item.EndDate,
			&item.StudentsCount,
		); err != nil {
			return nil, err
		}

		groups = append(groups, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *Repository) AddStudent(ctx context.Context, organizationID uuid.UUID, groupID uuid.UUID, studentID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var groupBranchID uuid.UUID
	var maxStudents int

	err = tx.QueryRow(ctx, `
SELECT branch_id, max_students
FROM groups
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, groupID, organizationID).Scan(&groupBranchID, &maxStudents)
	if err != nil {
		return ErrGroupNotFound
	}

	var studentBranchID uuid.UUID
	err = tx.QueryRow(ctx, `
SELECT branch_id
FROM students
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, studentID, organizationID).Scan(&studentBranchID)
	if err != nil {
		return ErrStudentNotFound
	}

	if studentBranchID != groupBranchID {
		return ErrStudentBranchMismatch
	}

	var activeStudentsCount int
	err = tx.QueryRow(ctx, `
SELECT COUNT(*)
FROM group_students
WHERE group_id = $1
  AND status = 'active'
`, groupID).Scan(&activeStudentsCount)
	if err != nil {
		return err
	}

	if activeStudentsCount >= maxStudents {
		return ErrGroupIsFull
	}

	_, err = tx.Exec(ctx, `
INSERT INTO group_students (
group_id,
student_id,
status
)
VALUES ($1, $2, 'active')
ON CONFLICT (group_id, student_id)
DO UPDATE SET
status = 'active',
left_at = NULL
`, groupID, studentID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *Repository) ListStudents(ctx context.Context, organizationID uuid.UUID, groupID uuid.UUID) ([]GroupStudent, error) {
	var groupExists bool
	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM groups
WHERE id = $1
  AND organization_id = $2
)
`, groupID, organizationID).Scan(&groupExists)
	if err != nil {
		return nil, err
	}

	if !groupExists {
		return nil, ErrGroupNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
s.id,
s.full_name,
COALESCE(s.phone, ''),
gs.status,
gs.joined_at::text
FROM group_students gs
JOIN students s ON s.id = gs.student_id
WHERE gs.group_id = $1
  AND s.organization_id = $2
ORDER BY gs.joined_at DESC, s.full_name ASC
`, groupID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]GroupStudent, 0)

	for rows.Next() {
		var item GroupStudent

		if err := rows.Scan(
			&item.StudentID,
			&item.FullName,
			&item.Phone,
			&item.Status,
			&item.JoinedAt,
		); err != nil {
			return nil, err
		}

		students = append(students, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}
