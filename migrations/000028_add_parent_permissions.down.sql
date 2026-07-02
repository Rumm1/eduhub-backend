DELETE FROM role_permissions
WHERE permission_id IN (
  SELECT id FROM permissions
  WHERE code IN ('parents.read', 'parents.manage')
);

DELETE FROM permissions
WHERE code IN ('parents.read', 'parents.manage');
