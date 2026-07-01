DELETE FROM role_permissions
WHERE permission_id IN (
SELECT id
FROM permissions
WHERE code IN (
'schedules.read',
'schedules.create',
'schedules.update',
'schedules.delete'
)
);

DELETE FROM permissions
WHERE code IN (
'schedules.read',
'schedules.create',
'schedules.update',
'schedules.delete'
);
