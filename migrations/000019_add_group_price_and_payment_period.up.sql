ALTER TABLE groups
ADD COLUMN IF NOT EXISTS monthly_price NUMERIC(12, 2) NOT NULL DEFAULT 0;

ALTER TABLE payments
ADD COLUMN IF NOT EXISTS payment_period DATE;

UPDATE payments
SET payment_period = date_trunc('month', payment_date)::date
WHERE payment_period IS NULL;

CREATE INDEX IF NOT EXISTS idx_payments_payment_period ON payments(payment_period);
CREATE INDEX IF NOT EXISTS idx_payments_student_group_period ON payments(student_id, group_id, payment_period);
CREATE INDEX IF NOT EXISTS idx_groups_monthly_price ON groups(monthly_price);
