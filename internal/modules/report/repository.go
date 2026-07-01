package report

import (
	"context"
	"strconv"
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

func (r *Repository) GetTeacherSchedule(
	ctx context.Context,
	organizationID uuid.UUID,
	teacherID uuid.UUID,
	fromDate string,
	toDate string,
) (TeacherScheduleReport, error) {
	var teacherName string

	err := r.db.QueryRow(ctx, `
SELECT full_name
FROM users
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, teacherID, organizationID).Scan(&teacherName)
	if err != nil {
		return TeacherScheduleReport{}, ErrTeacherNotFound
	}

	rows, err := r.db.Query(ctx, `
SELECT
l.id::text,
l.lesson_date::text,
l.start_time::text,
l.end_time::text,
ROUND((EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600)::numeric, 2)::text AS hours,
COALESCE(l.topic, ''),
l.status,
g.id::text,
g.name,
b.id::text,
b.name,
s.id::text,
s.name,
COALESCE(l.planned_teacher_id::text, ''),
COALESCE(planned_user.full_name, ''),
COALESCE(l.actual_teacher_id::text, ''),
COALESCE(actual_user.full_name, ''),
(
COALESCE(l.actual_teacher_id, l.teacher_id) <> COALESCE(l.planned_teacher_id, l.teacher_id)
) AS is_substitution,
CASE
WHEN COALESCE(l.actual_teacher_id, l.teacher_id) = $2 THEN 'actual'
WHEN l.planned_teacher_id = $2 THEN 'planned_only'
ELSE 'unknown'
END AS teacher_role_in_lesson,
COALESCE(l.substitution_reason, '')
FROM lessons l
JOIN groups g ON g.id = l.group_id
JOIN branches b ON b.id = l.branch_id
JOIN subjects s ON s.id = l.subject_id
LEFT JOIN users planned_user ON planned_user.id = l.planned_teacher_id
LEFT JOIN users actual_user ON actual_user.id = l.actual_teacher_id
WHERE l.organization_id = $1
  AND l.lesson_date >= $3::date
  AND l.lesson_date <= $4::date
  AND (
COALESCE(l.actual_teacher_id, l.teacher_id) = $2
OR l.planned_teacher_id = $2
  )
ORDER BY l.lesson_date ASC, l.start_time ASC
`, organizationID, teacherID, fromDate, toDate)
	if err != nil {
		return TeacherScheduleReport{}, err
	}
	defer rows.Close()

	items := make([]TeacherScheduleItem, 0)
	actualLessons := 0
	plannedOnlyLessons := 0
	substitutions := 0
	totalActualHours := "0"

	for rows.Next() {
		var item TeacherScheduleItem

		if err := rows.Scan(
			&item.LessonID,
			&item.LessonDate,
			&item.StartTime,
			&item.EndTime,
			&item.Hours,
			&item.Topic,
			&item.Status,
			&item.GroupID,
			&item.GroupName,
			&item.BranchID,
			&item.BranchName,
			&item.SubjectID,
			&item.SubjectName,
			&item.PlannedTeacherID,
			&item.PlannedTeacherName,
			&item.ActualTeacherID,
			&item.ActualTeacherName,
			&item.IsSubstitution,
			&item.TeacherRoleInLesson,
			&item.SubstitutionReason,
		); err != nil {
			return TeacherScheduleReport{}, err
		}

		if item.TeacherRoleInLesson == "actual" {
			actualLessons++
		}

		if item.TeacherRoleInLesson == "planned_only" {
			plannedOnlyLessons++
		}

		if item.IsSubstitution {
			substitutions++
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return TeacherScheduleReport{}, err
	}

	err = r.db.QueryRow(ctx, `
SELECT
ROUND(COALESCE(SUM(EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600), 0)::numeric, 2)::text
FROM lessons l
WHERE l.organization_id = $1
  AND l.lesson_date >= $3::date
  AND l.lesson_date <= $4::date
  AND COALESCE(l.actual_teacher_id, l.teacher_id) = $2
  AND l.status <> 'cancelled'
`, organizationID, teacherID, fromDate, toDate).Scan(&totalActualHours)
	if err != nil {
		return TeacherScheduleReport{}, err
	}

	return TeacherScheduleReport{
		TeacherID:          teacherID.String(),
		TeacherName:        teacherName,
		FromDate:           fromDate,
		ToDate:             toDate,
		TotalLessons:       len(items),
		ActualLessons:      actualLessons,
		PlannedOnlyLessons: plannedOnlyLessons,
		Substitutions:      substitutions,
		TotalActualHours:   totalActualHours,
		Items:              items,
	}, nil
}

func (r *Repository) GetPaymentsReport(
	ctx context.Context,
	organizationID uuid.UUID,
	fromDate string,
	toDate string,
	branchID string,
	groupID string,
	studentID string,
	status string,
) (PaymentsReport, error) {
	whereParts := []string{
		"p.organization_id = $1",
		"p.payment_date >= $2::date",
		"p.payment_date <= $3::date",
	}

	args := []interface{}{organizationID, fromDate, toDate}
	argIndex := 4

	if branchID != "" {
		whereParts = append(whereParts, "p.branch_id = $"+itoa(argIndex))
		args = append(args, branchID)
		argIndex++
	}

	if groupID != "" {
		whereParts = append(whereParts, "p.group_id = $"+itoa(argIndex)+"::uuid")
		args = append(args, groupID)
		argIndex++
	}

	if studentID != "" {
		whereParts = append(whereParts, "p.student_id = $"+itoa(argIndex)+"::uuid")
		args = append(args, studentID)
		argIndex++
	}

	if status != "" {
		whereParts = append(whereParts, "p.status = $"+itoa(argIndex))
		args = append(args, status)
		argIndex++
	}

	whereSQL := strings.Join(whereParts, " AND ")

	query := `
SELECT
p.id::text,
p.payment_date::text,
COALESCE(p.payment_period::text, ''),
s.id::text,
s.full_name,
COALESCE(g.id::text, ''),
COALESCE(g.name, ''),
b.id::text,
b.name,
p.amount::text,
COALESCE(p.payment_method, ''),
p.status,
COALESCE(p.comment, '')
FROM payments p
JOIN students s ON s.id = p.student_id
JOIN branches b ON b.id = p.branch_id
LEFT JOIN groups g ON g.id = p.group_id
WHERE ` + whereSQL + `
ORDER BY p.payment_date ASC, p.created_at ASC
`

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return PaymentsReport{}, err
	}
	defer rows.Close()

	items := make([]PaymentReportItem, 0)

	for rows.Next() {
		var item PaymentReportItem

		if err := rows.Scan(
			&item.PaymentID,
			&item.PaymentDate,
			&item.PaymentPeriod,
			&item.StudentID,
			&item.StudentName,
			&item.GroupID,
			&item.GroupName,
			&item.BranchID,
			&item.BranchName,
			&item.Amount,
			&item.PaymentMethod,
			&item.Status,
			&item.Comment,
		); err != nil {
			return PaymentsReport{}, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return PaymentsReport{}, err
	}

	totalAmount := "0"
	paidAmount := "0"
	pendingAmount := "0"
	refundedAmount := "0"
	cancelledAmount := "0"

	summaryQuery := `
SELECT
COALESCE(SUM(p.amount), 0)::text AS total_amount,
COALESCE(SUM(p.amount) FILTER (WHERE p.status = 'paid'), 0)::text AS paid_amount,
COALESCE(SUM(p.amount) FILTER (WHERE p.status = 'pending'), 0)::text AS pending_amount,
COALESCE(SUM(p.amount) FILTER (WHERE p.status = 'refunded'), 0)::text AS refunded_amount,
COALESCE(SUM(p.amount) FILTER (WHERE p.status = 'cancelled'), 0)::text AS cancelled_amount
FROM payments p
WHERE ` + whereSQL

	err = r.db.QueryRow(ctx, summaryQuery, args...).Scan(
		&totalAmount,
		&paidAmount,
		&pendingAmount,
		&refundedAmount,
		&cancelledAmount,
	)
	if err != nil {
		return PaymentsReport{}, err
	}

	return PaymentsReport{
		FromDate:        fromDate,
		ToDate:          toDate,
		TotalPayments:   len(items),
		TotalAmount:     totalAmount,
		PaidAmount:      paidAmount,
		PendingAmount:   pendingAmount,
		RefundedAmount:  refundedAmount,
		CancelledAmount: cancelledAmount,
		Items:           items,
	}, nil
}

func itoa(value int) string {
	return strconv.Itoa(value)
}
