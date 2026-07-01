DROP INDEX IF EXISTS idx_groups_monthly_price;
DROP INDEX IF EXISTS idx_payments_student_group_period;
DROP INDEX IF EXISTS idx_payments_payment_period;

ALTER TABLE payments
DROP COLUMN IF EXISTS payment_period;

ALTER TABLE groups
DROP COLUMN IF EXISTS monthly_price;
