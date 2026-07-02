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

func (r *Repository) GetUserAccessData(ctx context.Context, email string) (UserAccessData, error) {
	user, err := r.GetUserByEmail(ctx, email)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get user by email: %w", err)
	}

	profile, err := r.GetDefaultProfileByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get default profile: %w", err)
	}

	return r.buildUserAccessData(ctx, user, profile)
}

func (r *Repository) GetUserAccessDataByID(ctx context.Context, userID uuid.UUID) (UserAccessData, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get user by id: %w", err)
	}

	profile, err := r.GetDefaultProfileByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get default profile: %w", err)
	}

	return r.buildUserAccessData(ctx, user, profile)
}

func (r *Repository) GetUserAccessDataByProfileID(ctx context.Context, userID uuid.UUID, profileID uuid.UUID) (UserAccessData, error) {
	user, err := r.GetUserByID(ctx, userID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get user by id: %w", err)
	}

	profile, err := r.GetProfileByIDForUser(ctx, userID, profileID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get profile by id: %w", err)
	}

	return r.buildUserAccessData(ctx, user, profile)
}

func (r *Repository) GetDefaultProfileByUserID(ctx context.Context, userID uuid.UUID) (UserProfile, error) {
	var profile UserProfile

	err := r.db.QueryRow(ctx, `
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
  AND is_default = true
LIMIT 1
`, userID).Scan(
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

	if err != nil {
		return UserProfile{}, err
	}

	return r.hydrateProfile(ctx, profile)
}

func (r *Repository) GetProfileByIDForUser(ctx context.Context, userID uuid.UUID, profileID uuid.UUID) (UserProfile, error) {
	var profile UserProfile

	err := r.db.QueryRow(ctx, `
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
  AND user_id = $2
LIMIT 1
`, profileID, userID).Scan(
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

	if err != nil {
		return UserProfile{}, err
	}

	return r.hydrateProfile(ctx, profile)
}

func (r *Repository) GetAvailableProfilesByUserID(ctx context.Context, userID uuid.UUID) ([]UserProfile, error) {
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

		hydratedProfile, err := r.hydrateProfile(ctx, profile)
		if err != nil {
			return nil, err
		}

		profiles = append(profiles, hydratedProfile)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *Repository) hydrateProfile(ctx context.Context, profile UserProfile) (UserProfile, error) {
	roles, err := r.GetRoleCodesByProfileID(ctx, profile.ID)
	if err != nil {
		return UserProfile{}, fmt.Errorf("get profile roles: %w", err)
	}

	permissions, err := r.GetPermissionCodesByProfileID(ctx, profile.ID)
	if err != nil {
		return UserProfile{}, fmt.Errorf("get profile permissions: %w", err)
	}

	branchIDs, err := r.GetBranchIDsByProfileID(ctx, profile.ID)
	if err != nil {
		return UserProfile{}, fmt.Errorf("get profile branch ids: %w", err)
	}

	profile.Roles = roles
	profile.Permissions = permissions
	profile.BranchIDs = branchIDs

	return profile, nil
}

func (r *Repository) buildUserAccessData(ctx context.Context, user User, profile UserProfile) (UserAccessData, error) {
	availableProfiles, err := r.GetAvailableProfilesByUserID(ctx, user.ID)
	if err != nil {
		return UserAccessData{}, fmt.Errorf("get available profiles: %w", err)
	}

	return UserAccessData{
		User:              user,
		Profile:           profile,
		AvailableProfiles: availableProfiles,
		Roles:             profile.Roles,
		Permissions:       profile.Permissions,
		BranchIDs:         profile.BranchIDs,
	}, nil
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

func (r *Repository) GetPermissionCodesByProfileID(ctx context.Context, profileID uuid.UUID) ([]string, error) {
	rows, err := r.db.Query(ctx, `
SELECT DISTINCT p.code
FROM permissions p
JOIN role_permissions rp ON rp.permission_id = p.id
JOIN user_profile_roles upr ON upr.role_id = rp.role_id
WHERE upr.profile_id = $1
ORDER BY p.code
`, profileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	permissions := make([]string, 0)

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
