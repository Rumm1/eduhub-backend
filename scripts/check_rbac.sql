-- 1. Migration version
SELECT 'schema_migrations' AS check_name, version::text AS value, dirty::text AS extra
FROM schema_migrations;

-- 2. All roles
SELECT
	'roles' AS check_name,
	COALESCE(organization_id::text, 'SYSTEM') AS scope,
	code,
	name,
	COALESCE(description, '') AS description
FROM roles
ORDER BY code, scope;

-- 3. Permissions count by module
SELECT
'permissions_by_module' AS check_name,
split_part(code, '.', 1) AS module,
COUNT(*) AS permissions_count
FROM permissions
GROUP BY split_part(code, '.', 1)
ORDER BY module;

-- 4. Role permissions count
SELECT
	'role_permissions_count' AS check_name,
	COALESCE(r.organization_id::text, 'SYSTEM') AS scope,
	r.code AS role,
	COUNT(p.id) AS permissions_count
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
GROUP BY r.organization_id, r.code
ORDER BY r.code, scope;

-- 5. Full role-permission matrix
SELECT
r.code AS role,
p.code AS permission
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
ORDER BY r.code, p.code;

-- 6. Missing expected base roles
WITH expected_roles(code) AS (
VALUES
('SUPER_ADMIN'),
('ORG_ADMIN'),
('TEACHER'),
('STUDENT'),
('PARENT'),
('OWNER'),
('DIRECTOR'),
('FINANCE_MANAGER'),
('ACCOUNTANT'),
('SALES_MANAGER')
)
SELECT
'missing_roles' AS check_name,
er.code
FROM expected_roles er
LEFT JOIN roles r ON r.code = er.code
WHERE r.id IS NULL
ORDER BY er.code;

-- 7. Missing critical permissions
WITH expected_permissions(code) AS (
VALUES
('reports.teacher_schedule.read'),
('reports.payments.read'),
('reports.payroll.read'),
('reports.student_balance.read'),
('reports.export'),
('audit_logs.read'),
('audit_logs.export'),
('payroll.read_all'),
('payroll.read_own'),
('payroll.generate'),
('payroll.adjustments.manage'),
('payroll.send_to_teacher'),
('payroll.confirm'),
('payroll.dispute'),
('payroll.approve'),
('payroll.mark_paid'),
('payments.create'),
('payments.read'),
('payments.update_group_price')
)
SELECT
'missing_permissions' AS check_name,
ep.code
FROM expected_permissions ep
LEFT JOIN permissions p ON p.code = ep.code
WHERE p.id IS NULL
ORDER BY ep.code;

-- 8. Dangerous SUPER_ADMIN financial access
SELECT
'super_admin_sensitive_permissions' AS check_name,
p.code
FROM roles r
JOIN role_permissions rp ON rp.role_id = r.id
JOIN permissions p ON p.id = rp.permission_id
WHERE r.code = 'SUPER_ADMIN'
  AND (
p.code LIKE 'payroll.%'
OR p.code LIKE 'payments.%'
OR p.code LIKE 'reports.payments.%'
OR p.code LIKE 'reports.payroll.%'
OR p.code LIKE 'reports.student_balance.%'
  )
ORDER BY p.code;

-- 9. ORG_ADMIN missing important permissions
WITH required_org_admin_permissions(code) AS (
VALUES
('reports.teacher_schedule.read'),
('reports.payments.read'),
('reports.payroll.read'),
('reports.student_balance.read'),
('reports.export'),
('audit_logs.read'),
('audit_logs.export'),
('payroll.read_all'),
('payroll.generate'),
('payroll.adjustments.manage'),
('payroll.send_to_teacher'),
('payroll.approve'),
('payroll.mark_paid')
)
SELECT
'org_admin_missing_permissions' AS check_name,
req.code
FROM required_org_admin_permissions req
LEFT JOIN permissions p ON p.code = req.code
LEFT JOIN roles r ON r.code = 'ORG_ADMIN'
LEFT JOIN role_permissions rp ON rp.role_id = r.id AND rp.permission_id = p.id
WHERE rp.role_id IS NULL
ORDER BY req.code;

-- 10. TEACHER missing own permissions
WITH required_teacher_permissions(code) AS (
VALUES
('reports.teacher_schedule.read'),
('payroll.read_own'),
('payroll.confirm'),
('payroll.dispute')
)
SELECT
'teacher_missing_permissions' AS check_name,
req.code
FROM required_teacher_permissions req
LEFT JOIN permissions p ON p.code = req.code
LEFT JOIN roles r ON r.code = 'TEACHER'
LEFT JOIN role_permissions rp ON rp.role_id = r.id AND rp.permission_id = p.id
WHERE rp.role_id IS NULL
ORDER BY req.code;
