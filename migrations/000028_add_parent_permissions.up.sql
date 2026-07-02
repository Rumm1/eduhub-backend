INSERT INTO permissions (id, code, description)
VALUES
  (gen_random_uuid(), 'parents.read', 'View parents'),
  (gen_random_uuid(), 'parents.manage', 'Manage parents')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN ('parents.read', 'parents.manage')
WHERE r.code IN ('SUPER_ADMIN', 'ORG_ADMIN', 'OWNER', 'DIRECTOR')
ON CONFLICT DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code = 'parents.read'
WHERE r.code IN ('SALES_MANAGER', 'TEACHER')
ON CONFLICT DO NOTHING;
