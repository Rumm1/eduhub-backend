package teacher

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

func (r *Repository) Create(
	ctx context.Context,
	organizationID uuid.UUID,
	teacher Teacher,
	subjectIDs []uuid.UUID,
) (Teacher, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Teacher{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var userExists bool
	err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM users
WHERE id = $1 AND organization_id = $2 AND status = 'active'
)
`, teacher.UserID, organizationID).Scan(&userExists)
	if err != nil {
		return Teacher{}, err
	}

	if !userExists {
		return Teacher{}, ErrUserNotFound
	}

	var hasTeacherRole bool
	err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = $1
  AND r.organization_id = $2
  AND r.code = 'TEACHER'
)
`, teacher.UserID, organizationID).Scan(&hasTeacherRole)
	if err != nil {
		return Teacher{}, err
	}

	if !hasTeacherRole {
		return Teacher{}, ErrUserIsNotTeacher
	}

	for _, subjectID := range subjectIDs {
		var subjectExists bool

		err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM subjects
WHERE id = $1 AND organization_id = $2 AND status = 'active'
)
`, subjectID, organizationID).Scan(&subjectExists)
		if err != nil {
			return Teacher{}, err
		}

		if !subjectExists {
			return Teacher{}, ErrSubjectNotFound
		}
	}

	_, err = tx.Exec(ctx, `
INSERT INTO teacher_profiles (
user_id,
organization_id,
bio,
experience_years,
employment_type,
hourly_rate,
fixed_salary
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (user_id) DO UPDATE SET
bio = EXCLUDED.bio,
experience_years = EXCLUDED.experience_years,
employment_type = EXCLUDED.employment_type,
hourly_rate = EXCLUDED.hourly_rate,
fixed_salary = EXCLUDED.fixed_salary,
updated_at = now()
`,
		teacher.UserID,
		organizationID,
		teacher.Bio,
		teacher.ExperienceYears,
		teacher.EmploymentType,
		teacher.HourlyRate,
		teacher.FixedSalary,
	)
	if err != nil {
		return Teacher{}, err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM teacher_subjects
WHERE teacher_id = $1
`, teacher.UserID)
	if err != nil {
		return Teacher{}, err
	}

	for _, subjectID := range subjectIDs {
		_, err = tx.Exec(ctx, `
INSERT INTO teacher_subjects (teacher_id, subject_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, teacher.UserID, subjectID)
		if err != nil {
			return Teacher{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Teacher{}, err
	}

	return r.GetByUserID(ctx, organizationID, teacher.UserID)
}

func (r *Repository) GetByUserID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (Teacher, error) {
	var teacher Teacher

	err := r.db.QueryRow(ctx, `
SELECT
tp.user_id,
tp.organization_id,
u.email,
u.full_name,
COALESCE(u.phone, ''),
COALESCE(tp.bio, ''),
tp.experience_years,
COALESCE(tp.employment_type, ''),
COALESCE(tp.hourly_rate, 0)::float8,
COALESCE(tp.fixed_salary, 0)::float8
FROM teacher_profiles tp
JOIN users u ON u.id = tp.user_id
WHERE tp.organization_id = $1 AND tp.user_id = $2
`, organizationID, userID).Scan(
		&teacher.UserID,
		&teacher.OrganizationID,
		&teacher.Email,
		&teacher.FullName,
		&teacher.Phone,
		&teacher.Bio,
		&teacher.ExperienceYears,
		&teacher.EmploymentType,
		&teacher.HourlyRate,
		&teacher.FixedSalary,
	)
	if err != nil {
		return Teacher{}, err
	}

	return teacher, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Teacher, error) {
	rows, err := r.db.Query(ctx, `
SELECT
tp.user_id,
tp.organization_id,
u.email,
u.full_name,
COALESCE(u.phone, ''),
COALESCE(tp.bio, ''),
tp.experience_years,
COALESCE(tp.employment_type, ''),
COALESCE(tp.hourly_rate, 0)::float8,
COALESCE(tp.fixed_salary, 0)::float8
FROM teacher_profiles tp
JOIN users u ON u.id = tp.user_id
WHERE tp.organization_id = $1
ORDER BY tp.created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teachers := make([]Teacher, 0)

	for rows.Next() {
		var teacher Teacher

		if err := rows.Scan(
			&teacher.UserID,
			&teacher.OrganizationID,
			&teacher.Email,
			&teacher.FullName,
			&teacher.Phone,
			&teacher.Bio,
			&teacher.ExperienceYears,
			&teacher.EmploymentType,
			&teacher.HourlyRate,
			&teacher.FixedSalary,
		); err != nil {
			return nil, err
		}

		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return teachers, nil
}

func (r *Repository) GetSubjectsByTeacherID(ctx context.Context, teacherID uuid.UUID) ([]Subject, error) {
	rows, err := r.db.Query(ctx, `
SELECT
s.id,
s.name
FROM subjects s
JOIN teacher_subjects ts ON ts.subject_id = s.id
WHERE ts.teacher_id = $1
ORDER BY s.name
`, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := make([]Subject, 0)

	for rows.Next() {
		var subject Subject

		if err := rows.Scan(&subject.ID, &subject.Name); err != nil {
			return nil, err
		}

		subjects = append(subjects, subject)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}
