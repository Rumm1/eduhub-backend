INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'schedules.read', 'View schedules'),
(gen_random_uuid(), 'schedules.create', 'Create schedules'),
(gen_random_uuid(), 'schedules.update', 'Update schedules'),
(gen_random_uuid(), 'schedules.delete', 'Delete schedules')
ON CONFLICT (code) DO NOTHING;

WITH role_permission_map(role_code, permission_code) AS (
VALUES
('SUPER_ADMIN', 'schedules.read'),
('SUPER_ADMIN', 'schedules.create'),
('SUPER_ADMIN', 'schedules.update'),
('SUPER_ADMIN', 'schedules.delete'),

('ORG_ADMIN', 'schedules.read'),
('ORG_ADMIN', 'schedules.create'),
('ORG_ADMIN', 'schedules.update'),
('ORG_ADMIN', 'schedules.delete'),

('BRANCH_ADMIN', 'schedules.read'),
('BRANCH_ADMIN', 'schedules.create'),
('BRANCH_ADMIN', 'schedules.update'),

('MANAGER', 'schedules.read'),

('TEACHER', 'schedules.read')
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN role_permission_map rpm ON rpm.role_code = r.code
JOIN permissions p ON p.code = rpm.permission_code
ON CONFLICT DO NOTHING;
