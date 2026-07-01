DELETE FROM role_permissions
WHERE permission_id IN (
SELECT id
FROM permissions
WHERE code IN (
'audit_logs.read',
'audit_logs.export'
)
);

DELETE FROM permissions
WHERE code IN (
'audit_logs.read',
'audit_logs.export'
);

DROP TABLE IF EXISTS audit_logs;
