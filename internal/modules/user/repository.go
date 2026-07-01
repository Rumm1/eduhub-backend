package user

import (
	"context"
	"fmt"

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

func (r *Repository) CreateUserWithRoleAndBranches(
	ctx context.Context,
	newUser User,
	role RoleTemplate,
	branchIDs []uuid.UUID,
) (User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	for _, branchID := range branchIDs {
		var exists bool

		err := tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM branches
WHERE id = $1 AND organization_id = $2
)
`, branchID, newUser.OrganizationID).Scan(&exists)
		if err != nil {
			return User{}, err
		}

		if !exists {
			return User{}, ErrBranchNotFound
		}
	}

	_, err = tx.Exec(ctx, `
INSERT INTO users (
id,
organization_id,
email,
password_hash,
full_name,
phone,
avatar_path,
status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`,
		newUser.ID,
		newUser.OrganizationID,
		newUser.Email,
		newUser.PasswordHash,
		newUser.FullName,
		newUser.Phone,
		newUser.AvatarPath,
		newUser.Status,
	)
	if err != nil {
		return User{}, err
	}

	roleID, err := r.ensureRole(ctx, tx, newUser.OrganizationID, role)
	if err != nil {
		return User{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
`, newUser.ID, roleID)
	if err != nil {
		return User{}, err
	}

	for _, branchID := range branchIDs {
		_, err = tx.Exec(ctx, `
INSERT INTO user_branches (user_id, branch_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, newUser.ID, branchID)
		if err != nil {
			return User{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (r *Repository) ensureRole(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	role RoleTemplate,
) (uuid.UUID, error) {
	var roleID uuid.UUID

	err := tx.QueryRow(ctx, `
SELECT id
FROM roles
WHERE organization_id = $1 AND code = $2
LIMIT 1
`, organizationID, role.Code).Scan(&roleID)

	if err == nil {
		return roleID, nil
	}

	if err != pgx.ErrNoRows {
		return uuid.Nil, err
	}

	roleID = uuid.New()

	_, err = tx.Exec(ctx, `
INSERT INTO roles (
id,
organization_id,
name,
code,
description,
is_system
)
VALUES ($1, $2, $3, $4, $5, true)
`, roleID, organizationID, role.Name, role.Code, role.Description)
	if err != nil {
		return uuid.Nil, err
	}

	for _, permissionCode := range role.Permissions {
		_, err = tx.Exec(ctx, `
INSERT INTO role_permissions (role_id, permission_id)
SELECT $1, id
FROM permissions
WHERE code = $2
ON CONFLICT DO NOTHING
`, roleID, permissionCode)
		if err != nil {
			return uuid.Nil, err
		}
	}

	return roleID, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]User, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
email,
full_name,
COALESCE(phone, ''),
COALESCE(avatar_path, ''),
status
FROM users
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var item User

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.Email,
			&item.FullName,
			&item.Phone,
			&item.AvatarPath,
			&item.Status,
		); err != nil {
			return nil, err
		}

		users = append(users, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *Repository) GetRoleCodesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT r.code
FROM roles r
JOIN user_roles ur ON ur.role_id = r.id
WHERE ur.user_id = $1
ORDER BY r.code
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	roles := make([]string, 0)

	for rows.Next() {
		var role string
		if err := rows.Scan(&role); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *Repository) GetBranchIDsByUserID(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
SELECT branch_id
FROM user_branches
WHERE user_id = $1
ORDER BY branch_id
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branchIDs := make([]uuid.UUID, 0)

	for rows.Next() {
		var branchID uuid.UUID
		if err := rows.Scan(&branchID); err != nil {
			return nil, err
		}

		branchIDs = append(branchIDs, branchID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branchIDs, nil
}

func (r *Repository) BuildUserResponseData(ctx context.Context, userID uuid.UUID) ([]string, []uuid.UUID, error) {
	roles, err := r.GetRoleCodesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get roles: %w", err)
	}

	branchIDs, err := r.GetBranchIDsByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("get branch ids: %w", err)
	}

	return roles, branchIDs, nil
}
