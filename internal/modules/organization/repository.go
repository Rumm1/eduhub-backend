package organization

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

func (r *Repository) CreateOrganizationWithAdmin(
	ctx context.Context,
	org Organization,
	admin User,
	role Role,
	permissionCodes []string,
) (Organization, User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return Organization{}, User{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(ctx, `
INSERT INTO organizations (
id,
name,
name_ru,
name_kk,
name_en,
default_language,
bin,
phone,
email,
status
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`,
		org.ID,
		org.Name,
		org.NameRU,
		org.NameKK,
		org.NameEN,
		org.DefaultLanguage,
		org.BIN,
		org.Phone,
		org.Email,
		org.Status,
	)
	if err != nil {
		return Organization{}, User{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO users (
id,
organization_id,
email,
password_hash,
full_name,
phone,
status
)
VALUES ($1, $2, $3, $4, $5, $6, $7)
`, admin.ID, admin.OrganizationID, admin.Email, admin.PasswordHash, admin.FullName, admin.Phone, admin.Status)
	if err != nil {
		return Organization{}, User{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO roles (
id,
organization_id,
name,
code,
description,
is_system
)
VALUES ($1, $2, $3, $4, $5, $6)
`, role.ID, role.OrganizationID, role.Name, role.Code, role.Description, role.IsSystem)
	if err != nil {
		return Organization{}, User{}, err
	}

	_, err = tx.Exec(ctx, `
INSERT INTO user_roles (user_id, role_id)
VALUES ($1, $2)
`, admin.ID, role.ID)
	if err != nil {
		return Organization{}, User{}, err
	}

	for _, code := range permissionCodes {
		_, err = tx.Exec(ctx, `
INSERT INTO role_permissions (role_id, permission_id)
SELECT $1, id
FROM permissions
WHERE code = $2
ON CONFLICT DO NOTHING
`, role.ID, code)
		if err != nil {
			return Organization{}, User{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return Organization{}, User{}, err
	}

	return org, admin, nil
}

func NewOrganization(
	name string,
	nameRU string,
	nameKK string,
	nameEN string,
	defaultLanguage string,
	bin string,
	phone string,
	email string,
) Organization {
	return Organization{
		ID:              uuid.New(),
		Name:            name,
		NameRU:          nameRU,
		NameKK:          nameKK,
		NameEN:          nameEN,
		DefaultLanguage: defaultLanguage,
		BIN:             bin,
		Phone:           phone,
		Email:           email,
		Status:          "active",
	}
}
