package importer

import (
	"context"
	"strings"

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

func (r *Repository) FindGroupByName(ctx context.Context, organizationID uuid.UUID, groupName string) (GroupLookup, error) {
	var result GroupLookup

	err := r.db.QueryRow(ctx, `
SELECT
g.id,
g.branch_id,
g.name,
b.name
FROM groups g
JOIN branches b ON b.id = g.branch_id
WHERE g.organization_id = $1
  AND LOWER(TRIM(g.name)) = LOWER(TRIM($2))
  AND g.status = 'active'
LIMIT 1
`, organizationID, groupName).Scan(
		&result.ID,
		&result.BranchID,
		&result.Name,
		&result.Branch,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return GroupLookup{}, ErrGroupNotFound
		}

		return GroupLookup{}, err
	}

	return result, nil
}

func (r *Repository) ImportStudents(ctx context.Context, organizationID uuid.UUID, rows []ValidStudentImportRow) (ImportConfirmResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return ImportConfirmResult{}, err
	}
	defer tx.Rollback(ctx)

	var result ImportConfirmResult

	for _, row := range rows {
		studentID, studentCreated, err := r.findOrCreateStudent(
			ctx,
			tx,
			organizationID,
			row.Group.BranchID,
			row.Row.StudentFullName,
			row.Row.StudentPhone,
		)
		if err != nil {
			return ImportConfirmResult{}, err
		}

		if studentCreated {
			result.CreatedStudents++
		} else {
			result.ReusedStudents++
		}

		parentID, parentCreated, err := r.findOrCreateParent(
			ctx,
			tx,
			organizationID,
			row.Row.ParentFullName,
			row.Row.ParentPhone,
			row.Row.ParentEmail,
		)
		if err != nil {
			return ImportConfirmResult{}, err
		}

		if parentCreated {
			result.CreatedParents++
		} else {
			result.ReusedParents++
		}

		if err := r.linkParentToStudent(ctx, tx, studentID, parentID, row.Row.Relation); err != nil {
			return ImportConfirmResult{}, err
		}
		result.LinkedParentsToStudents++

		if err := r.linkStudentToGroup(ctx, tx, row.Group.ID, studentID); err != nil {
			return ImportConfirmResult{}, err
		}
		result.LinkedStudentsToGroups++
	}

	if err := tx.Commit(ctx); err != nil {
		return ImportConfirmResult{}, err
	}

	return result, nil
}

func (r *Repository) findOrCreateStudent(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	branchID uuid.UUID,
	fullName string,
	phone string,
) (uuid.UUID, bool, error) {
	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)

	var existingID uuid.UUID
	var err error

	if phone != "" {
		err = tx.QueryRow(ctx, `
SELECT id
FROM students
WHERE organization_id = $1
  AND LOWER(TRIM(full_name)) = LOWER(TRIM($2))
  AND COALESCE(phone, '') = $3
LIMIT 1
`, organizationID, fullName, phone).Scan(&existingID)
	} else {
		err = tx.QueryRow(ctx, `
SELECT id
FROM students
WHERE organization_id = $1
  AND branch_id = $2
  AND LOWER(TRIM(full_name)) = LOWER(TRIM($3))
LIMIT 1
`, organizationID, branchID, fullName).Scan(&existingID)
	}

	if err == nil {
		return existingID, false, nil
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, false, err
	}

	newID := uuid.New()

	_, err = tx.Exec(ctx, `
INSERT INTO students (
id,
organization_id,
branch_id,
full_name,
phone,
status,
source,
notes
)
VALUES ($1, $2, $3, $4, $5, 'active', 'google_form_import', 'Imported from Excel/CSV')
`, newID, organizationID, branchID, fullName, nullableString(phone))
	if err != nil {
		return uuid.Nil, false, err
	}

	return newID, true, nil
}

func (r *Repository) findOrCreateParent(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	fullName string,
	phone string,
	email string,
) (uuid.UUID, bool, error) {
	fullName = strings.TrimSpace(fullName)
	phone = strings.TrimSpace(phone)
	email = strings.ToLower(strings.TrimSpace(email))

	var existingID uuid.UUID
	var err error

	switch {
	case phone != "":
		err = tx.QueryRow(ctx, `
SELECT id
FROM parents
WHERE organization_id = $1
  AND COALESCE(phone, '') = $2
LIMIT 1
`, organizationID, phone).Scan(&existingID)
	case email != "":
		err = tx.QueryRow(ctx, `
SELECT id
FROM parents
WHERE organization_id = $1
  AND LOWER(COALESCE(email, '')) = LOWER($2)
LIMIT 1
`, organizationID, email).Scan(&existingID)
	default:
		err = tx.QueryRow(ctx, `
SELECT id
FROM parents
WHERE organization_id = $1
  AND LOWER(TRIM(full_name)) = LOWER(TRIM($2))
LIMIT 1
`, organizationID, fullName).Scan(&existingID)
	}

	if err == nil {
		return existingID, false, nil
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, false, err
	}

	newID := uuid.New()

	_, err = tx.Exec(ctx, `
INSERT INTO parents (
id,
organization_id,
full_name,
phone,
email
)
VALUES ($1, $2, $3, $4, $5)
`, newID, organizationID, fullName, nullableString(phone), nullableString(email))
	if err != nil {
		return uuid.Nil, false, err
	}

	return newID, true, nil
}

func (r *Repository) linkParentToStudent(ctx context.Context, tx pgx.Tx, studentID uuid.UUID, parentID uuid.UUID, relation string) error {
	_, err := tx.Exec(ctx, `
INSERT INTO student_parents (
student_id,
parent_id,
relation
)
VALUES ($1, $2, $3)
ON CONFLICT (student_id, parent_id)
DO UPDATE SET relation = EXCLUDED.relation
`, studentID, parentID, nullableString(relation))

	return err
}

func (r *Repository) linkStudentToGroup(ctx context.Context, tx pgx.Tx, groupID uuid.UUID, studentID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
INSERT INTO group_students (
group_id,
student_id,
joined_at,
status
)
VALUES ($1, $2, CURRENT_DATE, 'active')
ON CONFLICT (group_id, student_id)
DO UPDATE SET
status = 'active',
left_at = NULL
`, groupID, studentID)

	return err
}

func nullableString(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	return value
}
