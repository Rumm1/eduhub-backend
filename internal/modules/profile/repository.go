package profile

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

func (r *Repository) ListByUserID(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) ([]Profile, error) {
	if err := r.ensureUserBelongsToOrganization(ctx, r.db, userID, organizationID); err != nil {
		return nil, err
	}

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
  AND organization_id = $2
ORDER BY is_default DESC, created_at ASC
`, userID, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Profile, 0)

	for rows.Next() {
		var item Profile

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.OrganizationID,
			&item.BranchID,
			&item.DisplayName,
			&item.Position,
			&item.ProfileType,
			&item.Status,
			&item.IsDefault,
		); err != nil {
			return nil, err
		}

		hydrated, err := r.hydrateProfile(ctx, item)
		if err != nil {
			return nil, err
		}

		items = append(items, hydrated)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID) (Profile, error) {
	profile, err := r.getByID(ctx, r.db, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	return r.hydrateProfile(ctx, profile)
}

func (r *Repository) Create(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID, input Profile) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := r.ensureUserBelongsToOrganization(ctx, tx, userID, organizationID); err != nil {
		return Profile{}, err
	}

	if input.BranchID != nil {
		if err := r.ensureBranchBelongsToOrganization(ctx, tx, *input.BranchID, organizationID); err != nil {
			return Profile{}, err
		}
	}

	for _, branchID := range input.BranchIDs {
		if err := r.ensureBranchBelongsToOrganization(ctx, tx, branchID, organizationID); err != nil {
			return Profile{}, err
		}
	}

	if input.IsDefault {
		if err := r.clearDefaultProfiles(ctx, tx, userID); err != nil {
			return Profile{}, err
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
		input.ID,
		userID,
		organizationID,
		input.BranchID,
		input.DisplayName,
		input.Position,
		input.ProfileType,
		input.Status,
		input.IsDefault,
	)
	if err != nil {
		return Profile{}, err
	}

	for _, roleCode := range input.RoleCodes {
		roleID, err := r.getRoleIDByCode(ctx, tx, organizationID, roleCode)
		if err != nil {
			return Profile{}, err
		}

		_, err = tx.Exec(ctx, `
INSERT INTO user_profile_roles (profile_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, input.ID, roleID)
		if err != nil {
			return Profile{}, err
		}
	}

	for _, branchID := range input.BranchIDs {
		_, err = tx.Exec(ctx, `
INSERT INTO user_profile_branches (profile_id, branch_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, input.ID, branchID)
		if err != nil {
			return Profile{}, err
		}
	}

	if input.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, input.ID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, input.ID)
}

type UpdateProfileInput struct {
	BranchID    *uuid.UUID
	HasBranchID bool
	DisplayName *string
	Position    *string
	ProfileType *string
	Status      *string
	IsDefault   *bool
}

func (r *Repository) Update(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID, input UpdateProfileInput) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	if input.HasBranchID && input.BranchID != nil {
		if err := r.ensureBranchBelongsToOrganization(ctx, tx, *input.BranchID, organizationID); err != nil {
			return Profile{}, err
		}
	}

	if input.Status != nil && *input.Status != "active" && currentProfile.IsDefault {
		return Profile{}, ErrCannotDisableDefaultProfile
	}

	if input.IsDefault != nil {
		if !*input.IsDefault && currentProfile.IsDefault {
			return Profile{}, ErrDefaultProfileRequired
		}

		if *input.IsDefault {
			if err := r.clearDefaultProfiles(ctx, tx, currentProfile.UserID); err != nil {
				return Profile{}, err
			}
		}
	}

	_, err = tx.Exec(ctx, `
UPDATE user_profiles
SET
branch_id = CASE WHEN $3 THEN $4 ELSE branch_id END,
display_name = CASE WHEN $5 THEN $6 ELSE display_name END,
position = CASE WHEN $7 THEN $8 ELSE position END,
profile_type = CASE WHEN $9 THEN $10 ELSE profile_type END,
status = CASE WHEN $11 THEN $12 ELSE status END,
is_default = CASE WHEN $13 THEN $14 ELSE is_default END,
updated_at = now()
WHERE id = $1
  AND organization_id = $2
`,
		profileID,
		organizationID,
		input.HasBranchID,
		input.BranchID,
		input.DisplayName != nil,
		stringPointerValue(input.DisplayName),
		input.Position != nil,
		stringPointerValue(input.Position),
		input.ProfileType != nil,
		stringPointerValue(input.ProfileType),
		input.Status != nil,
		stringPointerValue(input.Status),
		input.IsDefault != nil,
		boolPointerValue(input.IsDefault),
	)
	if err != nil {
		return Profile{}, err
	}

	updatedProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	if updatedProfile.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, updatedProfile.ID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

func (r *Repository) Disable(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return err
	}

	if currentProfile.IsDefault {
		return ErrCannotDisableDefaultProfile
	}

	_, err = tx.Exec(ctx, `
UPDATE user_profiles
SET status = 'inactive',
    updated_at = now()
WHERE id = $1
  AND organization_id = $2
`, profileID, organizationID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *Repository) SetDefault(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	if currentProfile.Status != "active" {
		return Profile{}, ErrProfileInactive
	}

	if err := r.clearDefaultProfiles(ctx, tx, currentProfile.UserID); err != nil {
		return Profile{}, err
	}

	_, err = tx.Exec(ctx, `
UPDATE user_profiles
SET is_default = true,
    updated_at = now()
WHERE id = $1
  AND organization_id = $2
`, profileID, organizationID)
	if err != nil {
		return Profile{}, err
	}

	if err := r.syncDefaultProfileToLegacyTables(ctx, tx, profileID); err != nil {
		return Profile{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

func (r *Repository) AddRole(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID, roleCode string) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	roleID, err := r.getRoleIDByCode(ctx, tx, organizationID, roleCode)
	if err != nil {
		return Profile{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_profile_roles (profile_id, role_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, profileID, roleID)
	if err != nil {
		return Profile{}, err
	}

	if currentProfile.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, profileID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

func (r *Repository) RemoveRole(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID, roleCode string) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	roleID, err := r.getRoleIDByCode(ctx, tx, organizationID, roleCode)
	if err != nil {
		return Profile{}, err
	}

	var roleCount int
	err = tx.QueryRow(ctx, `
SELECT COUNT(*)
FROM user_profile_roles
WHERE profile_id = $1
`, profileID).Scan(&roleCount)
	if err != nil {
		return Profile{}, err
	}

	if roleCount <= 1 {
		return Profile{}, ErrRoleRequired
	}

	_, err = tx.Exec(ctx, `
DELETE FROM user_profile_roles
WHERE profile_id = $1
  AND role_id = $2
`, profileID, roleID)
	if err != nil {
		return Profile{}, err
	}

	if currentProfile.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, profileID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

func (r *Repository) AddBranch(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID, branchID uuid.UUID) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	if err := r.ensureBranchBelongsToOrganization(ctx, tx, branchID, organizationID); err != nil {
		return Profile{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_profile_branches (profile_id, branch_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING
`, profileID, branchID)
	if err != nil {
		return Profile{}, err
	}

	if currentProfile.BranchID == nil {
		_, err = tx.Exec(ctx, `
UPDATE user_profiles
SET branch_id = $1,
    updated_at = now()
WHERE id = $2
  AND organization_id = $3
`, branchID, profileID, organizationID)
		if err != nil {
			return Profile{}, err
		}
	}

	if currentProfile.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, profileID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

func (r *Repository) RemoveBranch(ctx context.Context, organizationID uuid.UUID, profileID uuid.UUID, branchID uuid.UUID) (Profile, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Profile{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	currentProfile, err := r.getByID(ctx, tx, organizationID, profileID)
	if err != nil {
		return Profile{}, err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM user_profile_branches
WHERE profile_id = $1
  AND branch_id = $2
`, profileID, branchID)
	if err != nil {
		return Profile{}, err
	}

	if currentProfile.BranchID != nil && *currentProfile.BranchID == branchID {
		_, err = tx.Exec(ctx, `
UPDATE user_profiles
SET branch_id = NULL,
    updated_at = now()
WHERE id = $1
  AND organization_id = $2
`, profileID, organizationID)
		if err != nil {
			return Profile{}, err
		}
	}

	if currentProfile.IsDefault {
		if err := r.syncDefaultProfileToLegacyTables(ctx, tx, profileID); err != nil {
			return Profile{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Profile{}, err
	}

	return r.GetByID(ctx, organizationID, profileID)
}

type queryer interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func (r *Repository) ensureUserBelongsToOrganization(
	ctx context.Context,
	q queryer,
	userID uuid.UUID,
	organizationID uuid.UUID,
) error {
	var exists bool

	err := q.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM users
WHERE id = $1
  AND organization_id = $2
)
`, userID, organizationID).Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return ErrUserNotFound
	}

	return nil
}

func (r *Repository) ensureBranchBelongsToOrganization(
	ctx context.Context,
	q queryer,
	branchID uuid.UUID,
	organizationID uuid.UUID,
) error {
	var exists bool

	err := q.QueryRow(ctx, `
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
	q queryer,
	organizationID uuid.UUID,
	roleCode string,
) (uuid.UUID, error) {
	var roleID uuid.UUID

	err := q.QueryRow(ctx, `
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

func (r *Repository) getByID(ctx context.Context, q queryer, organizationID uuid.UUID, profileID uuid.UUID) (Profile, error) {
	var profile Profile

	err := q.QueryRow(ctx, `
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
WHERE id = $1
  AND organization_id = $2
LIMIT 1
`, profileID, organizationID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.OrganizationID,
		&profile.BranchID,
		&profile.DisplayName,
		&profile.Position,
		&profile.ProfileType,
		&profile.Status,
		&profile.IsDefault,
	)
	if err == pgx.ErrNoRows {
		return Profile{}, ErrProfileNotFound
	}

	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

func (r *Repository) hydrateProfile(ctx context.Context, profile Profile) (Profile, error) {
	roleCodes, err := r.GetRoleCodesByProfileID(ctx, profile.ID)
	if err != nil {
		return Profile{}, err
	}

	branchIDs, err := r.GetBranchIDsByProfileID(ctx, profile.ID)
	if err != nil {
		return Profile{}, err
	}

	profile.RoleCodes = roleCodes
	profile.BranchIDs = branchIDs

	return profile, nil
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

	items := make([]uuid.UUID, 0)

	for rows.Next() {
		var item uuid.UUID
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

func (r *Repository) clearDefaultProfiles(ctx context.Context, tx pgx.Tx, userID uuid.UUID) error {
	_, err := tx.Exec(ctx, `
UPDATE user_profiles
SET is_default = false,
    updated_at = now()
WHERE user_id = $1
`, userID)

	return err
}

func (r *Repository) syncDefaultProfileToLegacyTables(ctx context.Context, tx pgx.Tx, profileID uuid.UUID) error {
	var userID uuid.UUID

	err := tx.QueryRow(ctx, `
SELECT user_id
FROM user_profiles
WHERE id = $1
`, profileID).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM user_roles
WHERE user_id = $1
`, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
DELETE FROM user_branches
WHERE user_id = $1
`, userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_roles (user_id, role_id)
SELECT $1, role_id
FROM user_profile_roles
WHERE profile_id = $2
ON CONFLICT DO NOTHING
`, userID, profileID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_branches (user_id, branch_id)
SELECT $1, branch_id
FROM user_profile_branches
WHERE profile_id = $2
ON CONFLICT DO NOTHING
`, userID, profileID)

	return err
}

func stringPointerValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func boolPointerValue(value *bool) bool {
	if value == nil {
		return false
	}

	return *value
}
