CREATE EXTENSION IF NOT EXISTS pgcrypto;

INSERT INTO permissions (id, code, description) VALUES
(gen_random_uuid(), 'organizations.read', 'Read organizations'),
(gen_random_uuid(), 'organizations.manage', 'Manage organizations'),

(gen_random_uuid(), 'branches.read', 'Read branches'),
(gen_random_uuid(), 'branches.create', 'Create branches'),
(gen_random_uuid(), 'branches.update', 'Update branches'),
(gen_random_uuid(), 'branches.delete', 'Delete branches'),

(gen_random_uuid(), 'users.read', 'Read users'),
(gen_random_uuid(), 'users.create', 'Create users'),
(gen_random_uuid(), 'users.update', 'Update users'),
(gen_random_uuid(), 'users.delete', 'Delete users'),

(gen_random_uuid(), 'roles.read', 'Read roles'),
(gen_random_uuid(), 'roles.manage', 'Manage roles'),

(gen_random_uuid(), 'subjects.read', 'Read subjects'),
(gen_random_uuid(), 'subjects.create', 'Create subjects'),
(gen_random_uuid(), 'subjects.update', 'Update subjects'),
(gen_random_uuid(), 'subjects.delete', 'Delete subjects'),

(gen_random_uuid(), 'teachers.read', 'Read teachers'),
(gen_random_uuid(), 'teachers.create', 'Create teachers'),
(gen_random_uuid(), 'teachers.update', 'Update teachers'),
(gen_random_uuid(), 'teachers.delete', 'Delete teachers'),

(gen_random_uuid(), 'students.read', 'Read students'),
(gen_random_uuid(), 'students.create', 'Create students'),
(gen_random_uuid(), 'students.update', 'Update students'),
(gen_random_uuid(), 'students.delete', 'Delete students'),

(gen_random_uuid(), 'groups.read', 'Read groups'),
(gen_random_uuid(), 'groups.create', 'Create groups'),
(gen_random_uuid(), 'groups.update', 'Update groups'),
(gen_random_uuid(), 'groups.delete', 'Delete groups'),

(gen_random_uuid(), 'lessons.read', 'Read lessons'),
(gen_random_uuid(), 'lessons.create', 'Create lessons'),
(gen_random_uuid(), 'lessons.update', 'Update lessons'),
(gen_random_uuid(), 'lessons.delete', 'Delete lessons'),

(gen_random_uuid(), 'attendance.read', 'Read attendance'),
(gen_random_uuid(), 'attendance.manage', 'Manage attendance'),

(gen_random_uuid(), 'homeworks.read', 'Read homeworks'),
(gen_random_uuid(), 'homeworks.manage', 'Manage homeworks'),

(gen_random_uuid(), 'payments.read', 'Read payments'),
(gen_random_uuid(), 'payments.manage', 'Manage payments'),

(gen_random_uuid(), 'payroll.read', 'Read payroll'),
(gen_random_uuid(), 'payroll.manage', 'Manage payroll'),
(gen_random_uuid(), 'payroll.approve', 'Approve payroll'),
(gen_random_uuid(), 'payroll.rules.manage', 'Manage payroll rules'),

(gen_random_uuid(), 'files.upload', 'Upload files'),
(gen_random_uuid(), 'files.read', 'Read files'),
(gen_random_uuid(), 'files.delete', 'Delete files'),

(gen_random_uuid(), 'notifications.read', 'Read notifications'),
(gen_random_uuid(), 'notifications.manage', 'Manage notifications'),

(gen_random_uuid(), 'audit.read', 'Read audit logs')
ON CONFLICT (code) DO NOTHING;
