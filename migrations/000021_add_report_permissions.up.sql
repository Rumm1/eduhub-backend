INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'reports.teacher_schedule.read', 'Read teacher schedule reports'),
(gen_random_uuid(), 'reports.payments.read', 'Read payments reports'),
(gen_random_uuid(), 'reports.payroll.read', 'Read payroll reports'),
(gen_random_uuid(), 'reports.student_balance.read', 'Read student balance reports'),
(gen_random_uuid(), 'reports.export', 'Export reports to files')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'reports.teacher_schedule.read',
'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export'
)
WHERE r.code IN ('ORG_ADMIN')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'reports.teacher_schedule.read'
)
WHERE r.code IN ('TEACHER')
ON CONFLICT DO NOTHING;
