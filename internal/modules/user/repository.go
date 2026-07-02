package user

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

func (r *Repository) CreateUserWithProfiles(
	ctx context.Context,
	newUser User,
	profiles []UserProfile,
) (User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return User{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
INSERT INTO users (
id,
organization_id,
email,
password_hash,
full_name,
phone,
avatar_path,
status,
must_change_password
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, TRUE)
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

	for _, profile := range profiles {
		if profile.BranchID != nil {
			if err := r.ensureBranchBelongsToOrganization(ctx, tx, *profile.BranchID, newUser.OrganizationID); err != nil {
				return User{}, err
			}
		}

		for _, branchID := range profile.BranchIDs {
			if err := r.ensureBranchBelongsToOrganization(ctx, tx, branchID, newUser.OrganizationID); err != nil {
				return User{}, err
			}
		}

		_, err = tx.Exec(ctx, `
INSERT INTO user_profiles (
id,
user_id,
organization_id,
branch_id,
display_name,
position,
profile_type,
status,
is_default
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`,
			profile.ID,
			newUser.ID,
			newUser.OrganizationID,
			profile.BranchID,
			profile.DisplayName,
			profile.Position,
			profile.ProfileType,
			profile.Status,
			profile.IsDefault,
		)
		if err != nil {
			return User{}, err
		}

		roleIDs := make([]uuid.UUID, 0, len(profile.RoleCodes))

		for _, roleCode := range profile.RoleCodes {
			roleID, err := r.getRoleIDByCode(ctx, tx, newUser.OrganizationID, roleCode)
			if err != nil {
				return User{}, err
			}

			roleIDs = append(roleIDs, roleID)

			_, err = tx.Exec(ctx, `
INSERT INTO user_profile_roles (profile_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, profile.ID, roleID)
			if err != nil {
				return User{}, err
			}
		}

		for _, branchID := range profile.BranchIDs {
			_, err = tx.Exec(ctx, `
INSERT INTO user_profile_branches (profile_id, branch_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, profile.ID, branchID)
			if err != nil {
				return User{}, err
			}
		}

		if profile.IsDefault {
			if err := r.syncDefaultProfileToLegacyTables(ctx, tx, newUser.ID, roleIDs, profile.BranchIDs); err != nil {
				return User{}, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return User{}, err
	}

	return newUser, nil
}

func (r *Repository) ensureBranchBelongsToOrganization(
	ctx context.Context,
	tx pgx.Tx,
	branchID uuid.UUID,
	organizationID uuid.UUID,
) error {
	var exists bool

	err := tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM branches
WHERE id = $1
  AND organization_id = $2
)
`, branchID, organizationID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrBranchNotFound
	}

	return nil
}

func (r *Repository) getRoleIDByCode(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	roleCode string,
) (uuid.UUID, error) {
	var roleID uuid.UUID

	err := tx.QueryRow(ctx, `
SELECT id
FROM roles
WHERE code = $1
  AND (organization_id = $2 OR organization_id IS NULL)
ORDER BY
CASE
WHEN organization_id = $2 THEN 0
ELSE 1
END
LIMIT 1
`, roleCode, organizationID).Scan(&roleID)

	if err == nil {
		return roleID, nil
	}

	if err == pgx.ErrNoRows {
		return uuid.Nil, ErrRoleInvalid
	}

	return uuid.Nil, err
}

func (r *Repository) syncDefaultProfileToLegacyTables(
	ctx context.Context,
	tx pgx.Tx,
	userID uuid.UUID,
	roleIDs []uuid.UUID,
	branchIDs []uuid.UUID,
) error {
	for _, roleID := range roleIDs {
		_, err := tx.Exec(ctx, `
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, userID, roleID)
		if err != nil {
			return err
		}
	}

	for _, branchID := range branchIDs {
		_, err := tx.Exec(ctx, `
INSERT INTO user_branches (user_id, branch_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, userID, branchID)
		if err != nil {
			return err
		}
	}

	return nil
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

func (r *Repository) BuildUserResponseData(ctx context.Context, userID uuid.UUID) ([]string, []uuid.UUID, []UserProfile, error) {
	profiles, err := r.GetProfilesByUserID(ctx, userID)
	if err != nil {
		return nil, nil, nil, err
	}

	roleSet := make(map[string]struct{})
	branchSet := make(map[uuid.UUID]struct{})

	roles := make([]string, 0)
	branchIDs := make([]uuid.UUID, 0)

	for _, profile := range profiles {
		for _, roleCode := range profile.RoleCodes {
			if _, exists := roleSet[roleCode]; exists {
				continue
			}

			roleSet[roleCode] = struct{}{}
			roles = append(roles, roleCode)
		}

		for _, branchID := range profile.BranchIDs {
			if _, exists := branchSet[branchID]; exists {
				continue
			}

			branchSet[branchID] = struct{}{}
			branchIDs = append(branchIDs, branchID)
		}
	}

	return roles, branchIDs, profiles, nil
}

func (r *Repository) GetProfilesByUserID(ctx context.Context, userID uuid.UUID) ([]UserProfile, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
user_id,
organization_id,
branch_id,
COALESCE(display_name, ''),
COALESCE(position, ''),
profile_type,
status,
is_default
FROM user_profiles
WHERE user_id = $1
ORDER BY is_default DESC, created_at ASC
`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	profiles := make([]UserProfile, 0)

	for rows.Next() {
		var profile UserProfile

		if err := rows.Scan(
			&profile.ID,
			&profile.UserID,
			&profile.OrganizationID,
			&profile.BranchID,
			&profile.DisplayName,
			&profile.Position,
			&profile.ProfileType,
			&profile.Status,
			&profile.IsDefault,
		); err != nil {
			return nil, err
		}

		roleCodes, err := r.GetRoleCodesByProfileID(ctx, profile.ID)
		if err != nil {
			return nil, err
		}

		branchIDs, err := r.GetBranchIDsByProfileID(ctx, profile.ID)
		if err != nil {
			return nil, err
		}

		profile.RoleCodes = roleCodes
		profile.BranchIDs = branchIDs

		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *Repository) GetRoleCodesByProfileID(ctx context.Context, profileID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT r.code
FROM roles r
JOIN user_profile_roles upr ON upr.role_id = r.id
WHERE upr.profile_id = $1
ORDER BY r.code
`, profileID)
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

func (r *Repository) GetBranchIDsByProfileID(ctx context.Context, profileID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `
SELECT branch_id
FROM user_profile_branches
WHERE profile_id = $1
ORDER BY branch_id
`, profileID)
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
