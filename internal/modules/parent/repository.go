package parent

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context, organizationID uuid.UUID) ([]Parent, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
full_name,
COALESCE(phone, ''),
COALESCE(email, ''),
created_at::text,
updated_at::text
FROM parents
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Parent, 0)

	for rows.Next() {
		var item Parent

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.FullName,
			&item.Phone,
			&item.Email,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, organizationID uuid.UUID, parentID uuid.UUID) (Parent, error) {
	var item Parent

	err := r.db.QueryRow(ctx, `
SELECT
id,
organization_id,
full_name,
COALESCE(phone, ''),
COALESCE(email, ''),
created_at::text,
updated_at::text
FROM parents
WHERE id = $1
  AND organization_id = $2
`, parentID, organizationID).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.FullName,
		&item.Phone,
		&item.Email,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Parent{}, ErrParentNotFound
		}

		return Parent{}, err
	}

	return item, nil
}

func (r *Repository) Create(ctx context.Context, item Parent) (Parent, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO parents (
id,
organization_id,
full_name,
phone,
email
)
VALUES ($1, $2, $3, $4, $5)
RETURNING
id,
organization_id,
full_name,
COALESCE(phone, ''),
COALESCE(email, ''),
created_at::text,
updated_at::text
`,
		item.ID,
		item.OrganizationID,
		item.FullName,
		item.Phone,
		item.Email,
	).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.FullName,
		&item.Phone,
		&item.Email,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return Parent{}, err
	}

	return item, nil
}

func (r *Repository) Update(ctx context.Context, organizationID uuid.UUID, parentID uuid.UUID, input Parent) (Parent, error) {
	var item Parent

	err := r.db.QueryRow(ctx, `
UPDATE parents
SET
full_name = COALESCE(NULLIF($3, ''), full_name),
phone = COALESCE(NULLIF($4, ''), phone),
email = COALESCE(NULLIF($5, ''), email),
updated_at = now()
WHERE id = $1
  AND organization_id = $2
RETURNING
id,
organization_id,
full_name,
COALESCE(phone, ''),
COALESCE(email, ''),
created_at::text,
updated_at::text
`, parentID, organizationID, input.FullName, input.Phone, input.Email).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.FullName,
		&item.Phone,
		&item.Email,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return Parent{}, ErrParentNotFound
		}

		return Parent{}, err
	}

	return item, nil
}

func (r *Repository) Delete(ctx context.Context, organizationID uuid.UUID, parentID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
DELETE FROM parents
WHERE id = $1
  AND organization_id = $2
`, parentID, organizationID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrParentNotFound
	}

	return nil
}

func (r *Repository) AttachStudent(
	ctx context.Context,
	organizationID uuid.UUID,
	parentID uuid.UUID,
	studentID uuid.UUID,
	relation string,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var parentExists bool
	err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM parents
WHERE id = $1
  AND organization_id = $2
)
`, parentID, organizationID).Scan(&parentExists)
	if err != nil {
		return err
	}

	if !parentExists {
		return ErrParentNotFound
	}

	var studentExists bool
	err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM students
WHERE id = $1
  AND organization_id = $2
)
`, studentID, organizationID).Scan(&studentExists)
	if err != nil {
		return err
	}

	if !studentExists {
		return ErrStudentNotFound
	}

	_, err = tx.Exec(ctx, `
DELETE FROM student_parents
WHERE student_id = $1
  AND parent_id = $2
`, studentID, parentID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO student_parents (
student_id,
parent_id,
relation
)
VALUES ($1, $2, $3)
`, studentID, parentID, relation)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) DetachStudent(
	ctx context.Context,
	organizationID uuid.UUID,
	parentID uuid.UUID,
	studentID uuid.UUID,
) error {
	tag, err := r.db.Exec(ctx, `
DELETE FROM student_parents sp
USING parents p, students s
WHERE sp.parent_id = p.id
  AND sp.student_id = s.id
  AND sp.parent_id = $1
  AND sp.student_id = $2
  AND p.organization_id = $3
  AND s.organization_id = $3
`, parentID, studentID, organizationID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrStudentNotFound
	}

	return nil
}

func (r *Repository) ListStudents(ctx context.Context, organizationID uuid.UUID, parentID uuid.UUID) ([]Student, error) {
	var parentExists bool
	err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM parents
WHERE id = $1
  AND organization_id = $2
)
`, parentID, organizationID).Scan(&parentExists)
	if err != nil {
		return nil, err
	}

	if !parentExists {
		return nil, ErrParentNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
s.id,
s.branch_id,
s.full_name,
COALESCE(s.phone, ''),
s.status,
COALESCE(sp.relation, '')
FROM students s
JOIN student_parents sp ON sp.student_id = s.id
WHERE sp.parent_id = $1
  AND s.organization_id = $2
ORDER BY s.full_name ASC
`, parentID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Student, 0)

	for rows.Next() {
		var item Student

		if err := rows.Scan(
			&item.ID,
			&item.BranchID,
			&item.FullName,
			&item.Phone,
			&item.Status,
			&item.Relation,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
