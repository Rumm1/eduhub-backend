package ai

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

func (r *Repository) GetDashboardMetrics(ctx context.Context, organizationID uuid.UUID, userID uuid.UUID) (DashboardMetrics, error) {
	var metrics DashboardMetrics

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
  AND gs.status = 'active'
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
) AS unread_notifications,

(
SELECT COUNT(*)
FROM audit_logs
WHERE organization_id = $1
  AND created_at >= NOW() - INTERVAL '7 days'
) AS recent_audit_logs_count
`, organizationID, userID).Scan(
		&metrics.StudentsCount,
		&metrics.TeachersCount,
		&metrics.GroupsCount,
		&metrics.LessonsToday,
		&metrics.PaymentsThisMonth,
		&metrics.PaymentsAmountThisMonth,
		&metrics.StudentDebtTotal,
		&metrics.PendingPayrollEntries,
		&metrics.UnreadNotifications,
		&metrics.RecentAuditLogsCount,
	)
	if err != nil {
		return DashboardMetrics{}, err
	}

	return metrics, nil
}
