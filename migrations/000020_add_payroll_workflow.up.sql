ALTER TABLE payroll_entries
ADD COLUMN IF NOT EXISTS lessons_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS substitution_count INTEGER NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS hours_worked NUMERIC(10, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS hourly_rate NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS base_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS bonus_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS penalty_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS correction_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS final_amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
ADD COLUMN IF NOT EXISTS status VARCHAR(40) NOT NULL DEFAULT 'draft',
ADD COLUMN IF NOT EXISTS teacher_confirmation_status VARCHAR(40) NOT NULL DEFAULT 'not_sent',
ADD COLUMN IF NOT EXISTS sent_to_teacher_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS teacher_confirmed_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS teacher_dispute_reason TEXT,
ADD COLUMN IF NOT EXISTS finance_approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS finance_approved_at TIMESTAMP,
ADD COLUMN IF NOT EXISTS paid_by UUID REFERENCES users(id) ON DELETE SET NULL,
ADD COLUMN IF NOT EXISTS paid_at TIMESTAMP;

CREATE TABLE IF NOT EXISTS payroll_adjustments (
id UUID PRIMARY KEY,
organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
period_id UUID NOT NULL REFERENCES payroll_periods(id) ON DELETE CASCADE,
payroll_entry_id UUID REFERENCES payroll_entries(id) ON DELETE CASCADE,
employee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
adjustment_type VARCHAR(40) NOT NULL,
amount NUMERIC(12, 2) NOT NULL,
reason TEXT,
status VARCHAR(40) NOT NULL DEFAULT 'active',
created_by UUID REFERENCES users(id) ON DELETE SET NULL,
approved_by UUID REFERENCES users(id) ON DELETE SET NULL,
approved_at TIMESTAMP,
created_at TIMESTAMP NOT NULL DEFAULT now(),
updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS payroll_entry_lessons (
id UUID PRIMARY KEY,
organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
payroll_entry_id UUID NOT NULL REFERENCES payroll_entries(id) ON DELETE CASCADE,
lesson_id UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
lesson_date DATE NOT NULL,
start_time TIME NOT NULL,
end_time TIME NOT NULL,
hours NUMERIC(8, 2) NOT NULL DEFAULT 0,
hourly_rate NUMERIC(12, 2) NOT NULL DEFAULT 0,
amount NUMERIC(12, 2) NOT NULL DEFAULT 0,
is_substitution BOOLEAN NOT NULL DEFAULT false,
created_at TIMESTAMP NOT NULL DEFAULT now(),
UNIQUE (payroll_entry_id, lesson_id)
);

CREATE INDEX IF NOT EXISTS idx_payroll_entries_teacher_status
ON payroll_entries(teacher_id, status);

CREATE INDEX IF NOT EXISTS idx_payroll_entries_teacher_confirmation_status
ON payroll_entries(teacher_confirmation_status);

CREATE INDEX IF NOT EXISTS idx_payroll_adjustments_organization_id
ON payroll_adjustments(organization_id);

CREATE INDEX IF NOT EXISTS idx_payroll_adjustments_period_id
ON payroll_adjustments(period_id);

CREATE INDEX IF NOT EXISTS idx_payroll_adjustments_employee_id
ON payroll_adjustments(employee_id);

CREATE INDEX IF NOT EXISTS idx_payroll_adjustments_payroll_entry_id
ON payroll_adjustments(payroll_entry_id);

CREATE INDEX IF NOT EXISTS idx_payroll_entry_lessons_payroll_entry_id
ON payroll_entry_lessons(payroll_entry_id);

CREATE INDEX IF NOT EXISTS idx_payroll_entry_lessons_teacher_id
ON payroll_entry_lessons(teacher_id);

CREATE INDEX IF NOT EXISTS idx_payroll_entry_lessons_lesson_id
ON payroll_entry_lessons(lesson_id);

INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'payroll.read_all', 'Read all payroll entries in organization'),
(gen_random_uuid(), 'payroll.read_own', 'Read own payroll entries'),
(gen_random_uuid(), 'payroll.generate', 'Generate payroll entries'),
(gen_random_uuid(), 'payroll.adjustments.manage', 'Manage payroll bonuses, premiums and corrections'),
(gen_random_uuid(), 'payroll.send_to_teacher', 'Send payroll to teacher for confirmation'),
(gen_random_uuid(), 'payroll.confirm', 'Confirm own payroll as employee'),
(gen_random_uuid(), 'payroll.dispute', 'Dispute own payroll as employee'),
(gen_random_uuid(), 'payroll.approve', 'Approve payroll by finance or director'),
(gen_random_uuid(), 'payroll.mark_paid', 'Mark payroll as paid'),
(gen_random_uuid(), 'payroll.export', 'Export payroll reports')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'payroll.read_all',
'payroll.generate',
'payroll.adjustments.manage',
'payroll.send_to_teacher',
'payroll.approve',
'payroll.mark_paid',
'payroll.export'
)
WHERE r.code IN ('ORG_ADMIN')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'payroll.read_own',
'payroll.confirm',
'payroll.dispute'
)
WHERE r.code IN ('TEACHER')
ON CONFLICT DO NOTHING;
