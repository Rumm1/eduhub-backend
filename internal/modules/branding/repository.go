package branding

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

type BrandingData struct {
	UserAvatarPath       string
	OrganizationLogoPath string
}

func (r *Repository) GetCurrentBranding(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID) (BrandingData, error) {
	var data BrandingData

	if organizationID == nil {
		err := r.db.QueryRow(ctx, `
SELECT COALESCE(avatar_path, '')
FROM users
WHERE id = $1
`, userID).Scan(&data.UserAvatarPath)

		return data, err
	}

	err := r.db.QueryRow(ctx, `
SELECT
COALESCE(u.avatar_path, '') AS avatar_path,
COALESCE(o.logo_path, '') AS logo_path
FROM users u
JOIN organizations o ON o.id = $2
WHERE u.id = $1
`, userID, *organizationID).Scan(
		&data.UserAvatarPath,
		&data.OrganizationLogoPath,
	)
	if err != nil {
		return BrandingData{}, err
	}

	return data, nil
}

func (r *Repository) UpdateUserAvatar(ctx context.Context, userID uuid.UUID, avatarPath string) error {
	_, err := r.db.Exec(ctx, `
UPDATE users
SET avatar_path = $2
WHERE id = $1
`, userID, avatarPath)

	return err
}

func (r *Repository) ClearUserAvatar(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
UPDATE users
SET avatar_path = NULL
WHERE id = $1
`, userID)

	return err
}

func (r *Repository) UpdateOrganizationLogo(ctx context.Context, organizationID uuid.UUID, logoPath string) error {
	_, err := r.db.Exec(ctx, `
UPDATE organizations
SET logo_path = $2
WHERE id = $1
`, organizationID, logoPath)

	return err
}

func (r *Repository) ClearOrganizationLogo(ctx context.Context, organizationID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `
UPDATE organizations
SET logo_path = NULL
WHERE id = $1
`, organizationID)

	return err
}
