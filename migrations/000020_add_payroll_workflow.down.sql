DELETE FROM role_permissions
WHERE permission_id IN (
SELECT id
FROM permissions
WHERE code IN (
'payroll.read_all',
'payroll.read_own',
'payroll.generate',
'payroll.adjustments.manage',
'payroll.send_to_teacher',
'payroll.confirm',
'payroll.dispute',
'payroll.approve',
'payroll.mark_paid',
'payroll.export'
)
);

DELETE FROM permissions
WHERE code IN (
'payroll.read_all',
'payroll.read_own',
'payroll.generate',
'payroll.adjustments.manage',
'payroll.send_to_teacher',
'payroll.confirm',
'payroll.dispute',
'payroll.approve',
'payroll.mark_paid',
'payroll.export'
);

DROP INDEX IF EXISTS idx_payroll_entry_lessons_lesson_id;
DROP INDEX IF EXISTS idx_payroll_entry_lessons_teacher_id;
DROP INDEX IF EXISTS idx_payroll_entry_lessons_payroll_entry_id;
DROP INDEX IF EXISTS idx_payroll_adjustments_payroll_entry_id;
DROP INDEX IF EXISTS idx_payroll_adjustments_employee_id;
DROP INDEX IF EXISTS idx_payroll_adjustments_period_id;
DROP INDEX IF EXISTS idx_payroll_adjustments_organization_id;
DROP INDEX IF EXISTS idx_payroll_entries_teacher_confirmation_status;
DROP INDEX IF EXISTS idx_payroll_entries_teacher_status;

DROP TABLE IF EXISTS payroll_entry_lessons;
DROP TABLE IF EXISTS payroll_adjustments;

ALTER TABLE payroll_entries
DROP COLUMN IF EXISTS paid_at,
DROP COLUMN IF EXISTS paid_by,
DROP COLUMN IF EXISTS finance_approved_at,
DROP COLUMN IF EXISTS finance_approved_by,
DROP COLUMN IF EXISTS teacher_dispute_reason,
DROP COLUMN IF EXISTS teacher_confirmed_at,
DROP COLUMN IF EXISTS sent_to_teacher_at,
DROP COLUMN IF EXISTS teacher_confirmation_status,
DROP COLUMN IF EXISTS status,
DROP COLUMN IF EXISTS final_amount,
DROP COLUMN IF EXISTS correction_amount,
DROP COLUMN IF EXISTS penalty_amount,
DROP COLUMN IF EXISTS bonus_amount,
DROP COLUMN IF EXISTS base_amount,
DROP COLUMN IF EXISTS hourly_rate,
DROP COLUMN IF EXISTS hours_worked,
DROP COLUMN IF EXISTS substitution_count,
DROP COLUMN IF EXISTS lessons_count;
