DELETE FROM role_permissions rp
USING permissions p
WHERE rp.permission_id = p.id
  AND p.code = 'dashboard.overview.read';

DELETE FROM permissions
WHERE code = 'dashboard.overview.read';
