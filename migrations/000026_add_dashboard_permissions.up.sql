INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'dashboard.overview.read', 'Read dashboard overview')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code = 'dashboard.overview.read'
WHERE r.code IN (
'ORG_ADMIN',
'OWNER',
'DIRECTOR',
'FINANCE_MANAGER',
'ACCOUNTANT',
'SALES_MANAGER'
)
ON CONFLICT DO NOTHING;
