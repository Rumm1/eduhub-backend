package payroll

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreatePeriod(ctx context.Context, period PayrollPeriod) (PayrollPeriod, error) {
	err := r.db.QueryRow(ctx, `
INSERT INTO payroll_periods (
id,
organization_id,
month,
year,
status
)
VALUES ($1, $2, $3, $4, 'draft')
ON CONFLICT (organization_id, month, year)
DO UPDATE SET updated_at = now()
RETURNING
id,
organization_id,
month,
year,
status
`,
		period.ID,
		period.OrganizationID,
		period.Month,
		period.Year,
	).Scan(
		&period.ID,
		&period.OrganizationID,
		&period.Month,
		&period.Year,
		&period.Status,
	)
	if err != nil {
		return PayrollPeriod{}, err
	}

	return period, nil
}

func (r *Repository) GetPeriod(ctx context.Context, organizationID uuid.UUID, periodID uuid.UUID) (PayrollPeriod, error) {
	var period PayrollPeriod

	err := r.db.QueryRow(ctx, `
SELECT
id,
organization_id,
month,
year,
status
FROM payroll_periods
WHERE id = $1
  AND organization_id = $2
`, periodID, organizationID).Scan(
		&period.ID,
		&period.OrganizationID,
		&period.Month,
		&period.Year,
		&period.Status,
	)
	if err != nil {
		return PayrollPeriod{}, ErrPeriodNotFound
	}

	return period, nil
}

func (r *Repository) GenerateForPeriod(ctx context.Context, organizationID uuid.UUID, periodID uuid.UUID) (PayrollPeriod, []PayrollEntry, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return PayrollPeriod{}, nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var period PayrollPeriod

	err = tx.QueryRow(ctx, `
		SELECT
			id,
			organization_id,
			month,
			year,
			status
		FROM payroll_periods
		WHERE id = $1
		  AND organization_id = $2
	`, periodID, organizationID).Scan(
		&period.ID,
		&period.OrganizationID,
		&period.Month,
		&period.Year,
		&period.Status,
	)
	if err != nil {
		return PayrollPeriod{}, nil, ErrPeriodNotFound
	}

	fromDate := time.Date(period.Year, time.Month(period.Month), 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)

	type payrollGenerationRow struct {
		TeacherID         uuid.UUID
		LessonsCount      int
		SubstitutionCount int
		HoursWorked       string
		HourlyRate        string
		BaseAmount        string
	}

	rows, err := tx.Query(ctx, `
		SELECT
			COALESCE(l.actual_teacher_id, l.teacher_id) AS teacher_id,
			COUNT(*)::int AS lessons_count,
			COUNT(*) FILTER (
				WHERE COALESCE(l.actual_teacher_id, l.teacher_id) <> COALESCE(l.planned_teacher_id, l.teacher_id)
			)::int AS substitution_count,
			ROUND(COALESCE(SUM(EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600), 0)::numeric, 2)::text AS hours_worked,
			COALESCE(tp.hourly_rate, 0)::text AS hourly_rate,
			ROUND((
				COALESCE(SUM(EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600), 0)
				* COALESCE(tp.hourly_rate, 0)
			)::numeric, 2)::text AS base_amount
		FROM lessons l
		JOIN teacher_profiles tp
		  ON tp.user_id = COALESCE(l.actual_teacher_id, l.teacher_id)
		 AND tp.organization_id = l.organization_id
		WHERE l.organization_id = $1
		  AND l.lesson_date >= $2::date
		  AND l.lesson_date < $3::date
		  AND COALESCE(l.actual_teacher_id, l.teacher_id) IS NOT NULL
		  AND l.status <> 'cancelled'
		GROUP BY COALESCE(l.actual_teacher_id, l.teacher_id), tp.hourly_rate
		ORDER BY teacher_id
	`, organizationID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"))
	if err != nil {
		return PayrollPeriod{}, nil, err
	}

	generationRows := make([]payrollGenerationRow, 0)

	for rows.Next() {
		var item payrollGenerationRow

		if err := rows.Scan(
			&item.TeacherID,
			&item.LessonsCount,
			&item.SubstitutionCount,
			&item.HoursWorked,
			&item.HourlyRate,
			&item.BaseAmount,
		); err != nil {
			rows.Close()
			return PayrollPeriod{}, nil, err
		}

		generationRows = append(generationRows, item)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return PayrollPeriod{}, nil, err
	}

	rows.Close()

	entries := make([]PayrollEntry, 0, len(generationRows))

	for _, item := range generationRows {
		entryID := uuid.New()

		var entry PayrollEntry

		err = tx.QueryRow(ctx, `
			INSERT INTO payroll_entries (
				id,
				organization_id,
				period_id,
				teacher_id,
				lessons_count,
				substitution_count,
				hours_worked,
				hourly_rate,
				base_amount,
				bonus_amount,
				penalty_amount,
				correction_amount,
				total_amount,
				final_amount,
				status,
				teacher_confirmation_status,
				worked_minutes,
				comment
			)
			VALUES (
				$1, $2, $3, $4,
				$5, $6, $7::numeric, $8::numeric, $9::numeric,
				0, 0, 0,
				$9::numeric,
				$9::numeric,
				'draft',
				'not_sent',
				ROUND(($7::numeric * 60))::int,
				'Generated from lessons'
			)
			ON CONFLICT (period_id, teacher_id)
			DO UPDATE SET
				lessons_count = EXCLUDED.lessons_count,
				substitution_count = EXCLUDED.substitution_count,
				hours_worked = EXCLUDED.hours_worked,
				hourly_rate = EXCLUDED.hourly_rate,
				base_amount = EXCLUDED.base_amount,
				total_amount = EXCLUDED.base_amount + payroll_entries.bonus_amount - payroll_entries.penalty_amount + payroll_entries.correction_amount,
				final_amount = EXCLUDED.base_amount + payroll_entries.bonus_amount - payroll_entries.penalty_amount + payroll_entries.correction_amount,
				worked_minutes = EXCLUDED.worked_minutes,
				updated_at = now()
			RETURNING
				id,
				organization_id,
				period_id,
				teacher_id,
				lessons_count,
				substitution_count,
				hours_worked::text,
				hourly_rate::text,
				base_amount::text,
				bonus_amount::text,
				penalty_amount::text,
				correction_amount::text,
				total_amount::text,
				final_amount::text,
				status,
				teacher_confirmation_status,
				COALESCE(teacher_dispute_reason, ''),
				COALESCE(comment, '')
		`,
			entryID,
			organizationID,
			periodID,
			item.TeacherID,
			item.LessonsCount,
			item.SubstitutionCount,
			item.HoursWorked,
			item.HourlyRate,
			item.BaseAmount,
		).Scan(
			&entry.ID,
			&entry.OrganizationID,
			&entry.PeriodID,
			&entry.TeacherID,
			&entry.LessonsCount,
			&entry.SubstitutionCount,
			&entry.HoursWorked,
			&entry.HourlyRate,
			&entry.BaseAmount,
			&entry.BonusAmount,
			&entry.PenaltyAmount,
			&entry.CorrectionAmount,
			&entry.TotalAmount,
			&entry.FinalAmount,
			&entry.Status,
			&entry.TeacherConfirmationStatus,
			&entry.TeacherDisputeReason,
			&entry.Comment,
		)
		if err != nil {
			return PayrollPeriod{}, nil, err
		}

		if err := r.rebuildEntryLessons(ctx, tx, organizationID, entry.ID, item.TeacherID, fromDate, toDate, item.HourlyRate); err != nil {
			return PayrollPeriod{}, nil, err
		}

		entries = append(entries, entry)
	}

	if err := tx.Commit(ctx); err != nil {
		return PayrollPeriod{}, nil, err
	}

	return period, entries, nil
}

func (r *Repository) rebuildEntryLessons(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	entryID uuid.UUID,
	teacherID uuid.UUID,
	fromDate time.Time,
	toDate time.Time,
	hourlyRate string,
) error {
	_, err := tx.Exec(ctx, `
		DELETE FROM payroll_entry_lessons
		WHERE payroll_entry_id = $1
	`, entryID)
	if err != nil {
		return err
	}

	type lessonPayrollRow struct {
		LessonID       uuid.UUID
		LessonDate     string
		StartTime      string
		EndTime        string
		Hours          string
		Amount         string
		IsSubstitution bool
	}

	rows, err := tx.Query(ctx, `
		SELECT
			l.id,
			l.lesson_date::text,
			l.start_time::text,
			l.end_time::text,
			ROUND((EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600)::numeric, 2)::text AS hours,
			ROUND((
				(EXTRACT(EPOCH FROM (l.end_time::time - l.start_time::time)) / 3600) * $5::numeric
			)::numeric, 2)::text AS amount,
			(
				COALESCE(l.actual_teacher_id, l.teacher_id) <> COALESCE(l.planned_teacher_id, l.teacher_id)
			) AS is_substitution
		FROM lessons l
		WHERE l.organization_id = $1
		  AND COALESCE(l.actual_teacher_id, l.teacher_id) = $2
		  AND l.lesson_date >= $3::date
		  AND l.lesson_date < $4::date
		  AND l.status <> 'cancelled'
		ORDER BY l.lesson_date ASC, l.start_time ASC
	`, organizationID, teacherID, fromDate.Format("2006-01-02"), toDate.Format("2006-01-02"), hourlyRate)
	if err != nil {
		return err
	}

	lessonRows := make([]lessonPayrollRow, 0)

	for rows.Next() {
		var item lessonPayrollRow

		if err := rows.Scan(
			&item.LessonID,
			&item.LessonDate,
			&item.StartTime,
			&item.EndTime,
			&item.Hours,
			&item.Amount,
			&item.IsSubstitution,
		); err != nil {
			rows.Close()
			return err
		}

		lessonRows = append(lessonRows, item)
	}

	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}

	rows.Close()

	for _, item := range lessonRows {
		_, err = tx.Exec(ctx, `
			INSERT INTO payroll_entry_lessons (
				id,
				organization_id,
				payroll_entry_id,
				lesson_id,
				teacher_id,
				lesson_date,
				start_time,
				end_time,
				hours,
				hourly_rate,
				amount,
				is_substitution
			)
			VALUES ($1, $2, $3, $4, $5, $6::date, $7::time, $8::time, $9::numeric, $10::numeric, $11::numeric, $12)
		`,
			uuid.New(),
			organizationID,
			entryID,
			item.LessonID,
			teacherID,
			item.LessonDate,
			item.StartTime,
			item.EndTime,
			item.Hours,
			hourlyRate,
			item.Amount,
			item.IsSubstitution,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) ListEntries(ctx context.Context, organizationID uuid.UUID) ([]PayrollEntry, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
period_id,
teacher_id,
lessons_count,
substitution_count,
hours_worked::text,
hourly_rate::text,
base_amount::text,
bonus_amount::text,
penalty_amount::text,
correction_amount::text,
total_amount::text,
final_amount::text,
status,
teacher_confirmation_status,
COALESCE(teacher_dispute_reason, ''),
COALESCE(comment, '')
FROM payroll_entries
WHERE organization_id = $1
ORDER BY created_at DESC
`, organizationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEntries(rows)
}

func (r *Repository) ListMyEntries(ctx context.Context, organizationID uuid.UUID, teacherID uuid.UUID) ([]PayrollEntry, error) {
	rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
period_id,
teacher_id,
lessons_count,
substitution_count,
hours_worked::text,
hourly_rate::text,
base_amount::text,
bonus_amount::text,
penalty_amount::text,
correction_amount::text,
total_amount::text,
final_amount::text,
status,
teacher_confirmation_status,
COALESCE(teacher_dispute_reason, ''),
COALESCE(comment, '')
FROM payroll_entries
WHERE organization_id = $1
  AND teacher_id = $2
ORDER BY created_at DESC
`, organizationID, teacherID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanEntries(rows)
}

func (r *Repository) CreateAdjustment(
	ctx context.Context,
	organizationID uuid.UUID,
	entryID uuid.UUID,
	adjustment PayrollAdjustment,
) (PayrollAdjustment, PayrollEntry, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return PayrollAdjustment{}, PayrollEntry{}, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var periodID uuid.UUID
	var teacherID uuid.UUID

	err = tx.QueryRow(ctx, `
SELECT period_id, teacher_id
FROM payroll_entries
WHERE id = $1
  AND organization_id = $2
`, entryID, organizationID).Scan(&periodID, &teacherID)
	if err != nil {
		return PayrollAdjustment{}, PayrollEntry{}, ErrEntryNotFound
	}

	var created PayrollAdjustment

	err = tx.QueryRow(ctx, `
INSERT INTO payroll_adjustments (
id,
organization_id,
period_id,
payroll_entry_id,
employee_id,
adjustment_type,
amount,
reason,
status,
created_by
)
VALUES ($1, $2, $3, $4, $5, $6, $7::numeric, $8, 'active', $9)
RETURNING
id,
organization_id,
period_id,
payroll_entry_id,
employee_id,
adjustment_type,
amount::text,
COALESCE(reason, ''),
status,
COALESCE(created_by::text, '')
`,
		adjustment.ID,
		organizationID,
		periodID,
		entryID,
		teacherID,
		adjustment.AdjustmentType,
		adjustment.Amount,
		adjustment.Reason,
		adjustment.CreatedBy,
	).Scan(
		&created.ID,
		&created.OrganizationID,
		&created.PeriodID,
		&created.PayrollEntryID,
		&created.EmployeeID,
		&created.AdjustmentType,
		&created.Amount,
		&created.Reason,
		&created.Status,
		&created.CreatedBy,
	)
	if err != nil {
		return PayrollAdjustment{}, PayrollEntry{}, err
	}

	entry, err := r.recalculateEntry(ctx, tx, organizationID, entryID)
	if err != nil {
		return PayrollAdjustment{}, PayrollEntry{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return PayrollAdjustment{}, PayrollEntry{}, err
	}

	return created, entry, nil
}

func (r *Repository) SendToTeacher(ctx context.Context, organizationID uuid.UUID, entryID uuid.UUID) (PayrollEntry, error) {
	var entry PayrollEntry

	err := r.db.QueryRow(ctx, `
		UPDATE payroll_entries
		SET
			status = 'sent_to_teacher',
			teacher_confirmation_status = 'pending',
			sent_to_teacher_at = now(),
			teacher_dispute_reason = NULL,
			teacher_confirmed_at = NULL,
			updated_at = now()
		WHERE id = $1
		  AND organization_id = $2
		  AND status IN ('draft', 'teacher_disputed')
		RETURNING
			id,
			organization_id,
			period_id,
			teacher_id,
			lessons_count,
			substitution_count,
			hours_worked::text,
			hourly_rate::text,
			base_amount::text,
			bonus_amount::text,
			penalty_amount::text,
			correction_amount::text,
			total_amount::text,
			final_amount::text,
			status,
			teacher_confirmation_status,
			COALESCE(teacher_dispute_reason, ''),
			COALESCE(comment, '')
	`, entryID, organizationID).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrInvalidPayrollStatus
	}

	return entry, nil
}

func (r *Repository) ConfirmByTeacher(ctx context.Context, organizationID uuid.UUID, entryID uuid.UUID, teacherID uuid.UUID) (PayrollEntry, error) {
	var entry PayrollEntry

	err := r.db.QueryRow(ctx, `
		UPDATE payroll_entries
		SET
			status = 'teacher_approved',
			teacher_confirmation_status = 'approved',
			teacher_confirmed_at = now(),
			teacher_dispute_reason = NULL,
			updated_at = now()
		WHERE id = $1
		  AND organization_id = $2
		  AND teacher_id = $3
		  AND status = 'sent_to_teacher'
		  AND teacher_confirmation_status = 'pending'
		RETURNING
			id,
			organization_id,
			period_id,
			teacher_id,
			lessons_count,
			substitution_count,
			hours_worked::text,
			hourly_rate::text,
			base_amount::text,
			bonus_amount::text,
			penalty_amount::text,
			correction_amount::text,
			total_amount::text,
			final_amount::text,
			status,
			teacher_confirmation_status,
			COALESCE(teacher_dispute_reason, ''),
			COALESCE(comment, '')
	`, entryID, organizationID, teacherID).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrInvalidPayrollStatus
	}

	return entry, nil
}
func (r *Repository) DisputeByTeacher(ctx context.Context, organizationID uuid.UUID, entryID uuid.UUID, teacherID uuid.UUID, reason string) (PayrollEntry, error) {
	var entry PayrollEntry

	err := r.db.QueryRow(ctx, `
UPDATE payroll_entries
SET
status = 'teacher_disputed',
teacher_confirmation_status = 'disputed',
teacher_dispute_reason = $4,
teacher_confirmed_at = now(),
updated_at = now()
WHERE id = $1
  AND organization_id = $2
  AND teacher_id = $3
  AND status = 'sent_to_teacher'
  AND teacher_confirmation_status = 'pending'
RETURNING
id,
organization_id,
period_id,
teacher_id,
lessons_count,
substitution_count,
hours_worked::text,
hourly_rate::text,
base_amount::text,
bonus_amount::text,
penalty_amount::text,
correction_amount::text,
total_amount::text,
final_amount::text,
status,
teacher_confirmation_status,
COALESCE(teacher_dispute_reason, ''),
COALESCE(comment, '')
`, entryID, organizationID, teacherID, reason).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrInvalidPayrollStatus
	}

	return entry, nil
}

func (r *Repository) ApproveByFinance(ctx context.Context, organizationID uuid.UUID, entryID uuid.UUID, approvedBy uuid.UUID) (PayrollEntry, error) {
	var entry PayrollEntry

	err := r.db.QueryRow(ctx, `
UPDATE payroll_entries
SET
status = 'approved_by_finance',
finance_approved_by = $3,
finance_approved_at = now(),
updated_at = now()
WHERE id = $1
  AND organization_id = $2
  AND status = 'teacher_approved'
  AND teacher_confirmation_status = 'approved'
RETURNING
id,
organization_id,
period_id,
teacher_id,
lessons_count,
substitution_count,
hours_worked::text,
hourly_rate::text,
base_amount::text,
bonus_amount::text,
penalty_amount::text,
correction_amount::text,
total_amount::text,
final_amount::text,
status,
teacher_confirmation_status,
COALESCE(teacher_dispute_reason, ''),
COALESCE(comment, '')
`, entryID, organizationID, approvedBy).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrInvalidPayrollStatus
	}

	return entry, nil
}

func (r *Repository) MarkPaid(ctx context.Context, organizationID uuid.UUID, entryID uuid.UUID, paidBy uuid.UUID) (PayrollEntry, error) {
	var entry PayrollEntry

	err := r.db.QueryRow(ctx, `
UPDATE payroll_entries
SET
status = 'paid',
paid_by = $3,
paid_at = now(),
updated_at = now()
WHERE id = $1
  AND organization_id = $2
  AND status = 'approved_by_finance'
RETURNING
id,
organization_id,
period_id,
teacher_id,
lessons_count,
substitution_count,
hours_worked::text,
hourly_rate::text,
base_amount::text,
bonus_amount::text,
penalty_amount::text,
correction_amount::text,
total_amount::text,
final_amount::text,
status,
teacher_confirmation_status,
COALESCE(teacher_dispute_reason, ''),
COALESCE(comment, '')
`, entryID, organizationID, paidBy).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrInvalidPayrollStatus
	}

	return entry, nil
}

func (r *Repository) recalculateEntry(
	ctx context.Context,
	tx pgx.Tx,
	organizationID uuid.UUID,
	entryID uuid.UUID,
) (PayrollEntry, error) {
	var entry PayrollEntry

	err := tx.QueryRow(ctx, `
WITH adjustment_totals AS (
SELECT
COALESCE(SUM(amount) FILTER (
WHERE adjustment_type IN ('bonus', 'premium', 'extra_work')
  AND status = 'active'
), 0) AS bonus_amount,
COALESCE(SUM(amount) FILTER (
WHERE adjustment_type IN ('penalty', 'deduction')
  AND status = 'active'
), 0) AS penalty_amount,
COALESCE(SUM(amount) FILTER (
WHERE adjustment_type = 'correction'
  AND status = 'active'
), 0) AS correction_amount
FROM payroll_adjustments
WHERE payroll_entry_id = $1
  AND organization_id = $2
)
UPDATE payroll_entries pe
SET
bonus_amount = adjustment_totals.bonus_amount,
penalty_amount = adjustment_totals.penalty_amount,
correction_amount = adjustment_totals.correction_amount,
total_amount = pe.base_amount + adjustment_totals.bonus_amount - adjustment_totals.penalty_amount + adjustment_totals.correction_amount,
final_amount = pe.base_amount + adjustment_totals.bonus_amount - adjustment_totals.penalty_amount + adjustment_totals.correction_amount,
updated_at = now()
FROM adjustment_totals
WHERE pe.id = $1
  AND pe.organization_id = $2
RETURNING
pe.id,
pe.organization_id,
pe.period_id,
pe.teacher_id,
pe.lessons_count,
pe.substitution_count,
pe.hours_worked::text,
pe.hourly_rate::text,
pe.base_amount::text,
pe.bonus_amount::text,
pe.penalty_amount::text,
pe.correction_amount::text,
pe.total_amount::text,
pe.final_amount::text,
pe.status,
pe.teacher_confirmation_status,
COALESCE(pe.teacher_dispute_reason, ''),
COALESCE(pe.comment, '')
`, entryID, organizationID).Scan(
		&entry.ID,
		&entry.OrganizationID,
		&entry.PeriodID,
		&entry.TeacherID,
		&entry.LessonsCount,
		&entry.SubstitutionCount,
		&entry.HoursWorked,
		&entry.HourlyRate,
		&entry.BaseAmount,
		&entry.BonusAmount,
		&entry.PenaltyAmount,
		&entry.CorrectionAmount,
		&entry.TotalAmount,
		&entry.FinalAmount,
		&entry.Status,
		&entry.TeacherConfirmationStatus,
		&entry.TeacherDisputeReason,
		&entry.Comment,
	)
	if err != nil {
		return PayrollEntry{}, ErrEntryNotFound
	}

	return entry, nil
}

type entryRows interface {
	Close()
	Err() error
	Next() bool
	Scan(...interface{}) error
}

func scanEntries(rows entryRows) ([]PayrollEntry, error) {
	entries := make([]PayrollEntry, 0)

	for rows.Next() {
		var entry PayrollEntry

		if err := rows.Scan(
			&entry.ID,
			&entry.OrganizationID,
			&entry.PeriodID,
			&entry.TeacherID,
			&entry.LessonsCount,
			&entry.SubstitutionCount,
			&entry.HoursWorked,
			&entry.HourlyRate,
			&entry.BaseAmount,
			&entry.BonusAmount,
			&entry.PenaltyAmount,
			&entry.CorrectionAmount,
			&entry.TotalAmount,
			&entry.FinalAmount,
			&entry.Status,
			&entry.TeacherConfirmationStatus,
			&entry.TeacherDisputeReason,
			&entry.Comment,
		); err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
