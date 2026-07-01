package schedule

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

func (r *Repository) Create(ctx context.Context, schedule Schedule) (Schedule, error) {
	err := r.db.QueryRow(ctx, `
SELECT branch_id
FROM groups
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, schedule.GroupID, schedule.OrganizationID).Scan(&schedule.BranchID)
	if err != nil {
		return Schedule{}, ErrGroupNotFound
	}

	err = r.db.QueryRow(ctx, `
INSERT INTO schedules (
id,
organization_id,
branch_id,
group_id,
weekday,
start_time,
end_time,
room
)
VALUES ($1, $2, $3, $4, $5, $6::time, $7::time, $8)
RETURNING
id,
organization_id,
branch_id,
group_id,
weekday,
start_time::text,
end_time::text,
COALESCE(room, '')
`,
		schedule.ID,
		schedule.OrganizationID,
		schedule.BranchID,
		schedule.GroupID,
		schedule.Weekday,
		schedule.StartTime,
		schedule.EndTime,
		schedule.Room,
	).Scan(
		&schedule.ID,
		&schedule.OrganizationID,
		&schedule.BranchID,
		&schedule.GroupID,
		&schedule.Weekday,
		&schedule.StartTime,
		&schedule.EndTime,
		&schedule.Room,
	)
	if err != nil {
		return Schedule{}, err
	}

	return schedule, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Schedule, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
group_id,
weekday,
start_time::text,
end_time::text,
COALESCE(room, '')
FROM schedules
WHERE organization_id = $1
ORDER BY weekday ASC, start_time ASC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	schedules := make([]Schedule, 0)

	for rows.Next() {
		var item Schedule

		if err := rows.Scan(
			&item.ID,
			&item.OrganizationID,
			&item.BranchID,
			&item.GroupID,
			&item.Weekday,
			&item.StartTime,
			&item.EndTime,
			&item.Room,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return schedules, nil
}
