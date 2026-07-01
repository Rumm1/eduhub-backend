package auth

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (User, error) {
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	var user User

	err := r.db.QueryRow(ctx, `
SELECT
id,
organization_id,
email,
password_hash,
full_name,
status
FROM users
WHERE LOWER(email) = $1
LIMIT 1
`, normalizedEmail).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (User, error) {
	var user User

	err := r.db.QueryRow(ctx, `
SELECT
id,
organization_id,
email,
password_hash,
full_name,
status
FROM users
WHERE id = $1
LIMIT 1
`, userID).Scan(
		&user.ID,
		&user.OrganizationID,
		&user.Email,
		&user.PasswordHash,
		&user.FullName,
		&user.Status,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
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

	var roles []string

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

func (r *Repository) GetPermissionCodesByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT DISTINCT p.code
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_roles ur ON ur.role_id = rp.role_id
WHERE ur.user_id = $1
ORDER BY p.code
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string

	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
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

	var branchIDs []uuid.UUID

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

func (r *Repository) GetUserAccessData(ctx context.Context, email string) (UserAccessData, error) {
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get user by email: %w", err)
	}

	return r.buildUserAccessData(ctx, user)
}

func (r *Repository) GetUserAccessDataByID(ctx context.Context, userID uuid.UUID) (UserAccessData, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get user by id: %w", err)
	}

	return r.buildUserAccessData(ctx, user)
}

func (r *Repository) buildUserAccessData(ctx context.Context, user User) (UserAccessData, error) {
	roles, err := r.GetRoleCodesByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get roles: %w", err)
	}

	permissions, err := r.GetPermissionCodesByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get permissions: %w", err)
	}

	branchIDs, err := r.GetBranchIDsByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get branch ids: %w", err)
	}

	return UserAccessData{
		User:        user,
		Roles:       roles,
		Permissions: permissions,
		BranchIDs:   branchIDs,
	}, nil
}
