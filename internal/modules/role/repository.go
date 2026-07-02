package role

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

func (r *Repository) List(ctx context.Context, organizationID uuid.UUID) ([]Role, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
name,
code,
COALESCE(description, ''),
is_system
FROM roles
WHERE organization_id = $1
   OR organization_id IS NULL
ORDER BY organization_id NULLS FIRST, code
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Role, 0)

	for rows.Next() {
		var item Role

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.Name,
			&item.Code,
			&item.Description,
			&item.IsSystem,
		); err != nil {
			return nil, err
		}

		permissionCodes, err := r.GetPermissionCodesByRoleID(ctx, item.ID)
		if err != nil {
			return nil, err
		}

		item.PermissionCodes = permissionCodes
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, organizationID uuid.UUID, roleID uuid.UUID) (Role, error) {
	item, err := r.getByID(ctx, r.db, organizationID, roleID)
	if err != nil {
		return Role{}, err
	}

	permissionCodes, err := r.GetPermissionCodesByRoleID(ctx, item.ID)
	if err != nil {
		return Role{}, err
	}

	item.PermissionCodes = permissionCodes

	return item, nil
}

func (r *Repository) Create(ctx context.Context, organizationID uuid.UUID, input Role) (Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Role{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
INSERT INTO roles (
id,
organization_id,
name,
code,
description,
is_system
)
VALUES ($1, $2, $3, $4, $5, false)
`,
		input.ID,
		organizationID,
		input.Name,
		input.Code,
		input.Description,
	)
	if err != nil {
		return Role{}, err
	}

	for _, permissionCode := range input.PermissionCodes {
		permissionID, err := r.getPermissionIDByCode(ctx, tx, permissionCode)
		if err != nil {
			return Role{}, err
		}

		_, err = tx.Exec(ctx, `
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, input.ID, permissionID)
		if err != nil {
			return Role{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Role{}, err
	}

	return r.GetByID(ctx, organizationID, input.ID)
}

type UpdateRoleInput struct {
	Name        *string
	Code        *string
	Description *string
}

func (r *Repository) Update(ctx context.Context, organizationID uuid.UUID, roleID uuid.UUID, input UpdateRoleInput) (Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Role{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentRole, err := r.getByID(ctx, tx, organizationID, roleID)
	if err != nil {
		return Role{}, err
	}

	if currentRole.OrganizationID == nil {
		return Role{}, ErrSystemRoleReadonly
	}

	_, err = tx.Exec(ctx, `
UPDATE roles
SET
name = CASE WHEN $3 THEN $4 ELSE name END,
code = CASE WHEN $5 THEN $6 ELSE code END,
description = CASE WHEN $7 THEN $8 ELSE description END
WHERE id = $1
  AND organization_id = $2
`,
		roleID,
		organizationID,
		input.Name != nil,
		stringPointerValue(input.Name),
		input.Code != nil,
		stringPointerValue(input.Code),
		input.Description != nil,
		stringPointerValue(input.Description),
	)
	if err != nil {
		return Role{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Role{}, err
	}

	return r.GetByID(ctx, organizationID, roleID)
}

func (r *Repository) Delete(ctx context.Context, organizationID uuid.UUID, roleID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentRole, err := r.getByID(ctx, tx, organizationID, roleID)
	if err != nil {
		return err
	}

	if currentRole.OrganizationID == nil || currentRole.IsSystem {
		return ErrSystemRoleReadonly
	}

	var usageCount int

	err = tx.QueryRow(ctx, `
SELECT
(
SELECT COUNT(*)
FROM user_roles
WHERE role_id = $1
)
+
(
SELECT COUNT(*)
FROM user_profile_roles
WHERE role_id = $1
)
`, roleID).Scan(&usageCount)
	if err != nil {
		return err
	}

	if usageCount > 0 {
		return ErrRoleInUse
	}

	_, err = tx.Exec(ctx, `
DELETE FROM role_permissions
WHERE role_id = $1
`, roleID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM roles
WHERE id = $1
  AND organization_id = $2
`, roleID, organizationID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) AddPermission(ctx context.Context, organizationID uuid.UUID, roleID uuid.UUID, permissionCode string) (Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Role{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentRole, err := r.getByID(ctx, tx, organizationID, roleID)
	if err != nil {
		return Role{}, err
	}

	if currentRole.OrganizationID == nil {
		return Role{}, ErrSystemRoleReadonly
	}

	permissionID, err := r.getPermissionIDByCode(ctx, tx, permissionCode)
	if err != nil {
		return Role{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO role_permissions (role_id, permission_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, roleID, permissionID)
	if err != nil {
		return Role{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Role{}, err
	}

	return r.GetByID(ctx, organizationID, roleID)
}

func (r *Repository) RemovePermission(ctx context.Context, organizationID uuid.UUID, roleID uuid.UUID, permissionCode string) (Role, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Role{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentRole, err := r.getByID(ctx, tx, organizationID, roleID)
	if err != nil {
		return Role{}, err
	}

	if currentRole.OrganizationID == nil {
		return Role{}, ErrSystemRoleReadonly
	}

	permissionID, err := r.getPermissionIDByCode(ctx, tx, permissionCode)
	if err != nil {
		return Role{}, err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM role_permissions
WHERE role_id = $1
  AND permission_id = $2
`, roleID, permissionID)
	if err != nil {
		return Role{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Role{}, err
	}

	return r.GetByID(ctx, organizationID, roleID)
}

type queryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *Repository) getByID(ctx context.Context, q queryer, organizationID uuid.UUID, roleID uuid.UUID) (Role, error) {
	var item Role

	err := q.QueryRow(ctx, `
SELECT
id,
organization_id,
name,
code,
COALESCE(description, ''),
is_system
FROM roles
WHERE id = $1
  AND (organization_id = $2 OR organization_id IS NULL)
LIMIT 1
`, roleID, organizationID).Scan(
		&item.ID,
		&item.OrganizationID,
		&item.Name,
		&item.Code,
		&item.Description,
		&item.IsSystem,
	)
	if err == pgx.ErrNoRows {
		return Role{}, ErrRoleNotFound
	}

	if err != nil {
		return Role{}, err
	}

	return item, nil
}

func (r *Repository) getPermissionIDByCode(ctx context.Context, q queryer, permissionCode string) (uuid.UUID, error) {
	var permissionID uuid.UUID

	err := q.QueryRow(ctx, `
SELECT id
FROM permissions
WHERE code = $1
LIMIT 1
`, permissionCode).Scan(&permissionID)

	if err == pgx.ErrNoRows {
		return uuid.Nil, ErrPermissionNotFound
	}

	if err != nil {
		return uuid.Nil, err
	}

	return permissionID, nil
}

func (r *Repository) GetPermissionCodesByRoleID(ctx context.Context, roleID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT p.code
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
WHERE rp.role_id = $1
ORDER BY p.code
`, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]string, 0)

	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func stringPointerValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
