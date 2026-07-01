INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'clients.read', 'View clients and leads'),
(gen_random_uuid(), 'clients.create', 'Create clients and leads'),
(gen_random_uuid(), 'clients.update', 'Update clients and leads'),
(gen_random_uuid(), 'clients.delete', 'Delete clients and leads'),

(gen_random_uuid(), 'parents.read', 'View parent records'),
(gen_random_uuid(), 'parents.create', 'Create parent records'),
(gen_random_uuid(), 'parents.update', 'Update parent records'),
(gen_random_uuid(), 'parents.delete', 'Delete parent records'),

(gen_random_uuid(), 'messages.read', 'View message history'),
(gen_random_uuid(), 'messages.send', 'Send messages to parents and groups'),
(gen_random_uuid(), 'messages.send.group', 'Send messages to group parent chats'),
(gen_random_uuid(), 'messages.templates.manage', 'Create and update message templates'),
(gen_random_uuid(), 'messages.channels.manage', 'Manage Telegram, WhatsApp and VK integrations'),

(gen_random_uuid(), 'payroll.read_limited', 'View limited teacher payroll preview'),
(gen_random_uuid(), 'payroll.preview', 'Preview teacher salary calculations'),

(gen_random_uuid(), 'notifications.send.organization', 'Send notifications inside organization'),
(gen_random_uuid(), 'notifications.send.platform', 'Send notifications to all organizations')
ON CONFLICT (code) DO NOTHING;

WITH role_permission_map(role_code, permission_code) AS (
VALUES
('SUPER_ADMIN', 'clients.read'),
('SUPER_ADMIN', 'clients.create'),
('SUPER_ADMIN', 'clients.update'),
('SUPER_ADMIN', 'clients.delete'),
('SUPER_ADMIN', 'parents.read'),
('SUPER_ADMIN', 'parents.create'),
('SUPER_ADMIN', 'parents.update'),
('SUPER_ADMIN', 'parents.delete'),
('SUPER_ADMIN', 'messages.read'),
('SUPER_ADMIN', 'messages.send'),
('SUPER_ADMIN', 'messages.send.group'),
('SUPER_ADMIN', 'messages.templates.manage'),
('SUPER_ADMIN', 'messages.channels.manage'),
('SUPER_ADMIN', 'payroll.read_limited'),
('SUPER_ADMIN', 'payroll.preview'),
('SUPER_ADMIN', 'notifications.send.organization'),
('SUPER_ADMIN', 'notifications.send.platform'),

('ORG_ADMIN', 'clients.read'),
('ORG_ADMIN', 'clients.create'),
('ORG_ADMIN', 'clients.update'),
('ORG_ADMIN', 'clients.delete'),
('ORG_ADMIN', 'parents.read'),
('ORG_ADMIN', 'parents.create'),
('ORG_ADMIN', 'parents.update'),
('ORG_ADMIN', 'parents.delete'),
('ORG_ADMIN', 'messages.read'),
('ORG_ADMIN', 'messages.send'),
('ORG_ADMIN', 'messages.send.group'),
('ORG_ADMIN', 'messages.templates.manage'),
('ORG_ADMIN', 'messages.channels.manage'),
('ORG_ADMIN', 'payroll.read_limited'),
('ORG_ADMIN', 'payroll.preview'),
('ORG_ADMIN', 'notifications.send.organization'),

('BRANCH_ADMIN', 'clients.read'),
('BRANCH_ADMIN', 'clients.create'),
('BRANCH_ADMIN', 'clients.update'),
('BRANCH_ADMIN', 'parents.read'),
('BRANCH_ADMIN', 'parents.create'),
('BRANCH_ADMIN', 'parents.update'),
('BRANCH_ADMIN', 'messages.read'),
('BRANCH_ADMIN', 'messages.send'),
('BRANCH_ADMIN', 'messages.send.group'),
('BRANCH_ADMIN', 'messages.templates.manage'),
('BRANCH_ADMIN', 'payroll.read_limited'),
('BRANCH_ADMIN', 'payroll.preview'),

('MANAGER', 'clients.read'),
('MANAGER', 'clients.create'),
('MANAGER', 'clients.update'),
('MANAGER', 'parents.read'),
('MANAGER', 'parents.create'),
('MANAGER', 'parents.update'),
('MANAGER', 'messages.read'),
('MANAGER', 'messages.send'),
('MANAGER', 'messages.send.group'),
('MANAGER', 'messages.templates.manage'),
('MANAGER', 'payroll.read_limited'),
('MANAGER', 'payroll.preview'),

('ACCOUNTANT', 'clients.read'),
('ACCOUNTANT', 'parents.read'),
('ACCOUNTANT', 'messages.read'),
('ACCOUNTANT', 'payroll.read_limited'),
('ACCOUNTANT', 'payroll.preview'),

('RECEPTIONIST', 'clients.read'),
('RECEPTIONIST', 'clients.create'),
('RECEPTIONIST', 'clients.update'),
('RECEPTIONIST', 'parents.read'),
('RECEPTIONIST', 'parents.create'),
('RECEPTIONIST', 'parents.update'),
('RECEPTIONIST', 'messages.read'),
('RECEPTIONIST', 'messages.send')
)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN role_permission_map rpm ON rpm.role_code = r.code
JOIN permissions p ON p.code = rpm.permission_code
ON CONFLICT DO NOTHING;
