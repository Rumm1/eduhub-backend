package permission

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) List(ctx context.Context) ([]Permission, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
code,
code AS name,
COALESCE(description, '')
FROM permissions
ORDER BY code
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Permission, 0)

	for rows.Next() {
		var item Permission

		if err := rows.Scan(
			&item.ID,
			&item.Code,
			&item.Name,
			&item.Description,
		); err != nil {
			return nil, err
		}

		item.Group = permissionGroupFromCode(item.Code)

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
