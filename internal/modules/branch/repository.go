package branch

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

func (r *Repository) Create(ctx context.Context, branch Branch) (Branch, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO branches (
id,
organization_id,
name,
address,
phone,
status
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, organization_id, name, address, phone, status
`,
		branch.ID,
		branch.OrganizationID,
		branch.Name,
		branch.Address,
		branch.Phone,
		branch.Status,
	).Scan(
		&branch.ID,
		&branch.OrganizationID,
		&branch.Name,
		&branch.Address,
		&branch.Phone,
		&branch.Status,
	)

	if err != nil {
		return Branch{}, err
	}

	return branch, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Branch, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
name,
address,
phone,
status
FROM branches
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	branches := make([]Branch, 0)

	for rows.Next() {
		var branch Branch

		if err := rows.Scan(
			&branch.ID,
			&branch.OrganizationID,
			&branch.Name,
			&branch.Address,
			&branch.Phone,
			&branch.Status,
		); err != nil {
			return nil, err
		}

		branches = append(branches, branch)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}
