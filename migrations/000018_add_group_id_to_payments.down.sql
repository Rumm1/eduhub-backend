DROP INDEX IF EXISTS idx_payments_payment_date;
DROP INDEX IF EXISTS idx_payments_branch_id;
DROP INDEX IF EXISTS idx_payments_student_id;
DROP INDEX IF EXISTS idx_payments_group_id;

ALTER TABLE payments
DROP COLUMN IF EXISTS group_id;
