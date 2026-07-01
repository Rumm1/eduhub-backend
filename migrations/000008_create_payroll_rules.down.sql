DROP INDEX IF EXISTS idx_teacher_payroll_rules_teacher_id;
DROP INDEX IF EXISTS idx_payroll_rules_branch_id;
DROP INDEX IF EXISTS idx_payroll_rules_organization_id;

ALTER TABLE payroll_entries
DROP COLUMN IF EXISTS formula_snapshot,
DROP COLUMN IF EXISTS worked_minutes,
DROP COLUMN IF EXISTS lesson_duration_minutes,
DROP COLUMN IF EXISTS lesson_rate,
DROP COLUMN IF EXISTS payroll_rule_id;

DROP TABLE IF EXISTS teacher_payroll_rules;
DROP TABLE IF EXISTS payroll_rules;
