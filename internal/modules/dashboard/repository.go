package dashboard

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

func (r *Repository) GetOverview(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (Overview, error) {
	var overview Overview

	err := r.db.QueryRow(ctx, `
SELECT
(SELECT COUNT(*) FROM students WHERE organization_id = $1 AND status = 'active') AS students_count,

(SELECT COUNT(*) FROM teacher_profiles WHERE organization_id = $1) AS teachers_count,

(SELECT COUNT(*) FROM groups WHERE organization_id = $1 AND status = 'active') AS groups_count,

(
SELECT COUNT(*)
FROM lessons
WHERE organization_id = $1
  AND lesson_date = CURRENT_DATE
  AND status <> 'cancelled'
) AS lessons_today,

(
SELECT COUNT(*)
FROM payments
WHERE organization_id = $1
  AND status = 'paid'
  AND payment_date >= date_trunc('month', CURRENT_DATE)::date
  AND payment_date < (date_trunc('month', CURRENT_DATE) + INTERVAL '1 month')::date
) AS payments_this_month,

(
SELECT COALESCE(SUM(amount), 0)::double precision
FROM payments
WHERE organization_id = $1
  AND status = 'paid'
  AND payment_date >= date_trunc('month', CURRENT_DATE)::date
  AND payment_date < (date_trunc('month', CURRENT_DATE) + INTERVAL '1 month')::date
) AS payments_amount_this_month,

(
WITH expected AS (
SELECT COALESCE(SUM(g.monthly_price), 0) AS total
FROM group_students gs
JOIN groups g ON g.id = gs.group_id
JOIN students s ON s.id = gs.student_id
WHERE g.organization_id = $1
  AND g.status = 'active'
  AND s.status = 'active'
),
paid AS (
SELECT COALESCE(SUM(amount), 0) AS total
FROM payments
WHERE organization_id = $1
  AND status = 'paid'
  AND payment_date >= date_trunc('month', CURRENT_DATE)::date
  AND payment_date < (date_trunc('month', CURRENT_DATE) + INTERVAL '1 month')::date
)
SELECT GREATEST((SELECT total FROM expected) - (SELECT total FROM paid), 0)::double precision
) AS student_debt_total,

(
SELECT COUNT(*)
FROM payroll_entries
WHERE organization_id = $1
  AND status <> 'paid'
) AS pending_payroll_entries,

(
SELECT COUNT(*)
FROM notifications
WHERE user_id = $2
  AND is_read = false
  AND (organization_id = $1 OR organization_id IS NULL)
) AS unread_notifications
`, organizationID, userID).Scan(
		&overview.StudentsCount,
		&overview.TeachersCount,
		&overview.GroupsCount,
		&overview.LessonsToday,
		&overview.PaymentsThisMonth,
		&overview.PaymentsAmountThisMonth,
		&overview.StudentDebtTotal,
		&overview.PendingPayrollEntries,
		&overview.UnreadNotifications,
	)
	if err != nil {
		return Overview{}, err
	}

	recentAuditLogs, err := r.ListRecentAuditLogs(ctx, organizationID)
	if err != nil {
		return Overview{}, err
	}

	overview.RecentAuditLogs = recentAuditLogs

	return overview, nil
}

func (r *Repository) ListRecentAuditLogs(ctx context.Context, organizationID uuid.UUID) ([]RecentAuditLog, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
COALESCE(user_id, '00000000-0000-0000-0000-000000000000'::uuid),
action,
COALESCE(entity_type, ''),
COALESCE(entity_id, '00000000-0000-0000-0000-000000000000'::uuid),
COALESCE(description, ''),
created_at::text
FROM audit_logs
WHERE organization_id = $1
ORDER BY created_at DESC
LIMIT 5
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]RecentAuditLog, 0)

	for rows.Next() {
		var item RecentAuditLog

		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.Action,
			&item.EntityType,
			&item.EntityID,
			&item.Description,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
