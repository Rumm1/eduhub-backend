-- This down migration only removes objects introduced by migration 24.
-- It does not recreate duplicate ORG_ADMIN roles.

DELETE FROM role_permissions rp
USING roles r
WHERE rp.role_id = r.id
  AND r.organization_id IS NULL
  AND r.code IN (
'OWNER',
'DIRECTOR',
'FINANCE_MANAGER',
'ACCOUNTANT',
'SALES_MANAGER',
'STUDENT',
'PARENT'
  );

DELETE FROM roles
WHERE organization_id IS NULL
  AND code IN (
'OWNER',
'DIRECTOR',
'FINANCE_MANAGER',
'ACCOUNTANT',
'SALES_MANAGER',
'STUDENT',
'PARENT'
  );

DELETE FROM role_permissions rp
USING permissions p
WHERE rp.permission_id = p.id
  AND p.code IN (
'payments.create',
'payments.update_group_price'
  );

DELETE FROM permissions
WHERE code IN (
'payments.create',
'payments.update_group_price'
);

DROP INDEX IF EXISTS roles_system_code_unique;
