package subject

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

func (r *Repository) Create(ctx context.Context, subject Subject) (Subject, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO subjects (
id,
organization_id,
name,
description,
status
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, organization_id, name, COALESCE(description, ''), status
`,
		subject.ID,
		subject.OrganizationID,
		subject.Name,
		subject.Description,
		subject.Status,
	).Scan(
		&subject.ID,
		&subject.OrganizationID,
		&subject.Name,
		&subject.Description,
		&subject.Status,
	)

	if err != nil {
		return Subject{}, err
	}

	return subject, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Subject, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
name,
COALESCE(description, ''),
status
FROM subjects
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subjects := make([]Subject, 0)

	for rows.Next() {
		var item Subject

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.Name,
			&item.Description,
			&item.Status,
		); err != nil {
			return nil, err
		}

		subjects = append(subjects, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subjects, nil
}
