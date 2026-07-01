ALTER TABLE payments
ADD COLUMN IF NOT EXISTS group_id UUID REFERENCES groups(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_payments_group_id ON payments(group_id);
CREATE INDEX IF NOT EXISTS idx_payments_student_id ON payments(student_id);
CREATE INDEX IF NOT EXISTS idx_payments_branch_id ON payments(branch_id);
CREATE INDEX IF NOT EXISTS idx_payments_payment_date ON payments(payment_date);
