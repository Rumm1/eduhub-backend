DELETE FROM role_permissions rp
USING permissions p
WHERE rp.permission_id = p.id
  AND p.code IN (
'profiles.read',
'profiles.manage',
'profiles.switch'
  );

DELETE FROM permissions
WHERE code IN (
'profiles.read',
'profiles.manage',
'profiles.switch'
);

DROP TABLE IF EXISTS user_profile_branches;
DROP TABLE IF EXISTS user_profile_roles;
DROP TABLE IF EXISTS user_profiles;
