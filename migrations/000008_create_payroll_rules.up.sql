CREATE TABLE payroll_rules (
    id UUID PRIMARY KEY,
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    branch_id UUID REFERENCES branches(id) ON DELETE CASCADE,
    subject_id UUID REFERENCES subjects(id) ON DELETE SET NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT,

    lesson_rate NUMERIC(12,2) NOT NULL DEFAULT 0,
    lesson_duration_minutes INT NOT NULL DEFAULT 60,

    formula TEXT NOT NULL DEFAULT 'lessons_count * lesson_rate + bonus - penalty',

    currency VARCHAR(10) NOT NULL DEFAULT 'KZT',
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    is_default BOOLEAN NOT NULL DEFAULT false,

    effective_from DATE,
    effective_to DATE,

    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),

    CHECK (lesson_rate >= 0),
    CHECK (lesson_duration_minutes > 0)
);

CREATE TABLE teacher_payroll_rules (
    teacher_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payroll_rule_id UUID NOT NULL REFERENCES payroll_rules(id) ON DELETE CASCADE,

    assigned_at TIMESTAMP NOT NULL DEFAULT now(),

    PRIMARY KEY (teacher_id, payroll_rule_id)
);

ALTER TABLE payroll_entries
ADD COLUMN payroll_rule_id UUID REFERENCES payroll_rules(id) ON DELETE SET NULL,
ADD COLUMN lesson_rate NUMERIC(12,2) NOT NULL DEFAULT 0,
ADD COLUMN lesson_duration_minutes INT NOT NULL DEFAULT 60,
ADD COLUMN worked_minutes INT NOT NULL DEFAULT 0,
ADD COLUMN formula_snapshot TEXT;

CREATE INDEX idx_payroll_rules_organization_id ON payroll_rules(organization_id);
CREATE INDEX idx_payroll_rules_branch_id ON payroll_rules(branch_id);
CREATE INDEX idx_teacher_payroll_rules_teacher_id ON teacher_payroll_rules(teacher_id);
