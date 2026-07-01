package student

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

func (r *Repository) CreateWithParent(
	ctx context.Context,
	student Student,
	parent *Parent,
) (Student, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Student{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var branchExists bool
	err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM branches
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
)
`, student.BranchID, student.OrganizationID).Scan(&branchExists)
	if err != nil {
		return Student{}, err
	}

	if !branchExists {
		return Student{}, ErrBranchNotFound
	}

	var birthDate interface{}
	if student.BirthDate != "" {
		birthDate = student.BirthDate
	}

	err = tx.QueryRow(ctx, `
INSERT INTO students (
id,
organization_id,
branch_id,
full_name,
phone,
birth_date,
gender,
status,
source,
notes
)
VALUES ($1, $2, $3, $4, $5, $6::date, $7, $8, $9, $10)
RETURNING
id,
organization_id,
branch_id,
full_name,
COALESCE(phone, ''),
COALESCE(birth_date::text, ''),
COALESCE(gender, ''),
status,
COALESCE(source, ''),
COALESCE(notes, '')
`,
		student.ID,
		student.OrganizationID,
		student.BranchID,
		student.FullName,
		student.Phone,
		birthDate,
		student.Gender,
		student.Status,
		student.Source,
		student.Notes,
	).Scan(
		&student.ID,
		&student.OrganizationID,
		&student.BranchID,
		&student.FullName,
		&student.Phone,
		&student.BirthDate,
		&student.Gender,
		&student.Status,
		&student.Source,
		&student.Notes,
	)
	if err != nil {
		return Student{}, err
	}

	if parent != nil {
		_, err = tx.Exec(ctx, `
INSERT INTO parents (
id,
organization_id,
full_name,
phone,
email
)
VALUES ($1, $2, $3, $4, $5)
`,
			parent.ID,
			parent.OrganizationID,
			parent.FullName,
			parent.Phone,
			parent.Email,
		)
		if err != nil {
			return Student{}, err
		}

		_, err = tx.Exec(ctx, `
INSERT INTO student_parents (
student_id,
parent_id,
relation
)
VALUES ($1, $2, $3)
`, student.ID, parent.ID, parent.Relation)
		if err != nil {
			return Student{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Student{}, err
	}

	return student, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Student, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
full_name,
COALESCE(phone, ''),
COALESCE(birth_date::text, ''),
COALESCE(gender, ''),
status,
COALESCE(source, ''),
COALESCE(notes, '')
FROM students
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	students := make([]Student, 0)

	for rows.Next() {
		var item Student

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.BranchID,
			&item.FullName,
			&item.Phone,
			&item.BirthDate,
			&item.Gender,
			&item.Status,
			&item.Source,
			&item.Notes,
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

func (r *Repository) GetParentsByStudentID(ctx context.Context, studentID uuid.UUID) ([]Parent, error) {
	rows, err := r.db.Query(ctx, `
SELECT
p.id,
p.organization_id,
p.full_name,
COALESCE(p.phone, ''),
COALESCE(p.email, ''),
COALESCE(sp.relation, '')
FROM parents p
JOIN student_parents sp ON sp.parent_id = p.id
WHERE sp.student_id = $1
ORDER BY p.created_at DESC
`, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	parents := make([]Parent, 0)

	for rows.Next() {
		var parent Parent

		if err := rows.Scan(
			&parent.ID,
			&parent.OrganizationID,
			&parent.FullName,
			&parent.Phone,
			&parent.Email,
			&parent.Relation,
		); err != nil {
			return nil, err
		}

		parents = append(parents, parent)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return parents, nil
}
