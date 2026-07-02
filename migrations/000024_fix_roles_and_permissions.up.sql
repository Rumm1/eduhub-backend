-- 1. Merge duplicate system roles by code where organization_id IS NULL.
-- PostgreSQL UNIQUE (organization_id, code) allows duplicates when organization_id is NULL,
-- so we normalize them manually first.

WITH ranked_roles AS (
SELECT
id,
code,
ROW_NUMBER() OVER (
PARTITION BY code
ORDER BY created_at ASC, id::text ASC
) AS rn
FROM roles
WHERE organization_id IS NULL
),
canonical_roles AS (
SELECT id, code
FROM ranked_roles
WHERE rn = 1
),
duplicate_roles AS (
SELECT id, code
FROM ranked_roles
WHERE rn > 1
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT
canonical_roles.id,
role_permissions.permission_id
FROM duplicate_roles
JOIN canonical_roles ON canonical_roles.code = duplicate_roles.code
JOIN role_permissions ON role_permissions.role_id = duplicate_roles.id
ON CONFLICT DO NOTHING;

WITH ranked_roles AS (
SELECT
id,
code,
ROW_NUMBER() OVER (
PARTITION BY code
ORDER BY created_at ASC, id::text ASC
) AS rn
FROM roles
WHERE organization_id IS NULL
),
canonical_roles AS (
SELECT id, code
FROM ranked_roles
WHERE rn = 1
),
duplicate_roles AS (
SELECT id, code
FROM ranked_roles
WHERE rn > 1
)
INSERT INTO user_roles (user_id, role_id)
SELECT
user_roles.user_id,
canonical_roles.id
FROM duplicate_roles
JOIN canonical_roles ON canonical_roles.code = duplicate_roles.code
JOIN user_roles ON user_roles.role_id = duplicate_roles.id
ON CONFLICT DO NOTHING;

WITH ranked_roles AS (
SELECT
id,
code,
ROW_NUMBER() OVER (
PARTITION BY code
ORDER BY created_at ASC, id::text ASC
) AS rn
FROM roles
WHERE organization_id IS NULL
)
DELETE FROM roles
WHERE id IN (
SELECT id
FROM ranked_roles
WHERE rn > 1
);

-- 2. Add partial unique index for system roles.
CREATE UNIQUE INDEX IF NOT EXISTS roles_system_code_unique
ON roles(code)
WHERE organization_id IS NULL;

-- 3. Add missing business roles.
INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Owner', 'OWNER', 'Business owner role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'OWNER'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Director', 'DIRECTOR', 'Organization director role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'DIRECTOR'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Finance Manager', 'FINANCE_MANAGER', 'Finance manager role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'FINANCE_MANAGER'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Accountant', 'ACCOUNTANT', 'Accountant role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'ACCOUNTANT'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Sales Manager', 'SALES_MANAGER', 'Sales manager role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'SALES_MANAGER'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Student', 'STUDENT', 'Student role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'STUDENT'
);

INSERT INTO roles (id, organization_id, name, code, description, is_system)
SELECT gen_random_uuid(), NULL, 'Parent', 'PARENT', 'Parent role', true
WHERE NOT EXISTS (
SELECT 1 FROM roles WHERE organization_id IS NULL AND code = 'PARENT'
);

-- 4. Add missing granular payment permissions.
INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'payments.create', 'Create student payments'),
(gen_random_uuid(), 'payments.update_group_price', 'Update group monthly price')
ON CONFLICT (code) DO NOTHING;

-- 5. Remove sensitive tenant finance permissions from SUPER_ADMIN.
DELETE FROM role_permissions rp
USING roles r, permissions p
WHERE rp.role_id = r.id
  AND rp.permission_id = p.id
  AND r.organization_id IS NULL
  AND r.code = 'SUPER_ADMIN'
  AND (
p.code LIKE 'payments.%'
OR p.code LIKE 'payroll.%'
OR p.code LIKE 'reports.payments.%'
OR p.code LIKE 'reports.payroll.%'
OR p.code LIKE 'reports.student_balance.%'
  );

-- 6. ORG_ADMIN keeps full organization management access.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'payments.read',
'payments.manage',
'payments.create',
'payments.update_group_price',

'payroll.read_all',
'payroll.generate',
'payroll.adjustments.manage',
'payroll.send_to_teacher',
'payroll.approve',
'payroll.mark_paid',
'payroll.export',

'reports.teacher_schedule.read',
'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export',

'audit_logs.read',
'audit_logs.export'
)
WHERE r.organization_id IS NULL
  AND r.code = 'ORG_ADMIN'
ON CONFLICT DO NOTHING;

-- 7. OWNER and DIRECTOR can see business reports and audit.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'branches.read',
'users.read',
'teachers.read',
'students.read',
'groups.read',
'schedules.read',
'lessons.read',
'attendance.read',

'payments.read',
'payroll.read_all',

'reports.teacher_schedule.read',
'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export',

'audit_logs.read',
'audit_logs.export'
)
WHERE r.organization_id IS NULL
  AND r.code IN ('OWNER', 'DIRECTOR')
ON CONFLICT DO NOTHING;

-- 8. Finance roles.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'students.read',
'groups.read',
'branches.read',

'payments.read',
'payments.manage',
'payments.create',
'payments.update_group_price',

'payroll.read_all',
'payroll.generate',
'payroll.adjustments.manage',
'payroll.send_to_teacher',
'payroll.approve',
'payroll.mark_paid',
'payroll.export',

'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export',

'audit_logs.read'
)
WHERE r.organization_id IS NULL
  AND r.code IN ('FINANCE_MANAGER', 'ACCOUNTANT')
ON CONFLICT DO NOTHING;

-- 9. Sales manager can work with students/groups/payments, but not payroll.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'branches.read',
'students.read',
'students.create',
'students.update',
'parents.read',
'parents.create',
'parents.update',
'groups.read',
'teachers.read',
'schedules.read',

'payments.read',
'payments.create',

'reports.payments.read',
'reports.student_balance.read'
)
WHERE r.organization_id IS NULL
  AND r.code = 'SALES_MANAGER'
ON CONFLICT DO NOTHING;

-- 10. Teacher own permissions.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'groups.read',
'students.read',
'schedules.read',
'lessons.read',
'lessons.create',
'lessons.update',
'attendance.read',
'attendance.manage',
'homeworks.read',
'homeworks.manage',
'files.read',
'files.upload',
'notifications.read',
'payroll.read_own',
'payroll.confirm',
'payroll.dispute',
'reports.teacher_schedule.read'
)
WHERE r.organization_id IS NULL
  AND r.code = 'TEACHER'
ON CONFLICT DO NOTHING;

-- 11. Student and parent minimal read permissions.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'notifications.read',
'homeworks.read',
'files.read',
'schedules.read',
'lessons.read'
)
WHERE r.organization_id IS NULL
  AND r.code IN ('STUDENT', 'PARENT')
ON CONFLICT DO NOTHING;
