CREATE TABLE IF NOT EXISTS user_profiles (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
organization_id UUID REFERENCES organizations(id) ON DELETE CASCADE,
branch_id UUID REFERENCES branches(id) ON DELETE SET NULL,
display_name VARCHAR(255),
position VARCHAR(120),
profile_type VARCHAR(80) NOT NULL DEFAULT 'staff',
status VARCHAR(32) NOT NULL DEFAULT 'active',
is_default BOOLEAN NOT NULL DEFAULT false,
created_at TIMESTAMP NOT NULL DEFAULT now(),
updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS user_profile_roles (
profile_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
PRIMARY KEY (profile_id, role_id)
);

CREATE TABLE IF NOT EXISTS user_profile_branches (
profile_id UUID NOT NULL REFERENCES user_profiles(id) ON DELETE CASCADE,
branch_id UUID NOT NULL REFERENCES branches(id) ON DELETE CASCADE,
PRIMARY KEY (profile_id, branch_id)
);

CREATE INDEX IF NOT EXISTS idx_user_profiles_user_id ON user_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_organization_id ON user_profiles(organization_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_branch_id ON user_profiles(branch_id);
CREATE INDEX IF NOT EXISTS idx_user_profiles_status ON user_profiles(status);
CREATE INDEX IF NOT EXISTS idx_user_profile_roles_role_id ON user_profile_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_user_profile_branches_branch_id ON user_profile_branches(branch_id);

CREATE UNIQUE INDEX IF NOT EXISTS idx_user_profiles_one_default
ON user_profiles(user_id)
WHERE is_default = true;

-- Create one default profile for every existing user.
INSERT INTO user_profiles (
id,
user_id,
organization_id,
branch_id,
display_name,
position,
profile_type,
status,
is_default
)
SELECT
gen_random_uuid(),
u.id,
u.organization_id,
(
SELECT ub.branch_id
FROM user_branches ub
WHERE ub.user_id = u.id
ORDER BY ub.branch_id::text
LIMIT 1
),
u.full_name,
COALESCE(
(
SELECT r.code
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = u.id
ORDER BY
CASE r.code
WHEN 'SUPER_ADMIN' THEN 1
WHEN 'OWNER' THEN 2
WHEN 'DIRECTOR' THEN 3
WHEN 'ORG_ADMIN' THEN 4
WHEN 'FINANCE_MANAGER' THEN 5
WHEN 'ACCOUNTANT' THEN 6
WHEN 'SALES_MANAGER' THEN 7
WHEN 'TEACHER' THEN 8
WHEN 'STUDENT' THEN 9
WHEN 'PARENT' THEN 10
ELSE 99
END,
r.code
LIMIT 1
),
'STAFF'
),
COALESCE(
(
SELECT lower(r.code)
FROM user_roles ur
JOIN roles r ON r.id = ur.role_id
WHERE ur.user_id = u.id
ORDER BY
CASE r.code
WHEN 'SUPER_ADMIN' THEN 1
WHEN 'OWNER' THEN 2
WHEN 'DIRECTOR' THEN 3
WHEN 'ORG_ADMIN' THEN 4
WHEN 'FINANCE_MANAGER' THEN 5
WHEN 'ACCOUNTANT' THEN 6
WHEN 'SALES_MANAGER' THEN 7
WHEN 'TEACHER' THEN 8
WHEN 'STUDENT' THEN 9
WHEN 'PARENT' THEN 10
ELSE 99
END,
r.code
LIMIT 1
),
'staff'
),
u.status,
true
FROM users u
WHERE NOT EXISTS (
SELECT 1
FROM user_profiles up
WHERE up.user_id = u.id
);

-- Copy current user_roles into default profile roles.
INSERT INTO user_profile_roles (profile_id, role_id)
SELECT
up.id,
ur.role_id
FROM user_profiles up
JOIN user_roles ur ON ur.user_id = up.user_id
WHERE up.is_default = true
ON CONFLICT DO NOTHING;

-- Copy current user_branches into default profile branches.
INSERT INTO user_profile_branches (profile_id, branch_id)
SELECT
up.id,
ub.branch_id
FROM user_profiles up
JOIN user_branches ub ON ub.user_id = up.user_id
WHERE up.is_default = true
ON CONFLICT DO NOTHING;

-- Profile permissions.
INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'profiles.read', 'Read user profiles'),
(gen_random_uuid(), 'profiles.manage', 'Manage user profiles'),
(gen_random_uuid(), 'profiles.switch', 'Switch active user profile')
ON CONFLICT (code) DO NOTHING;

-- ORG_ADMIN / OWNER / DIRECTOR can manage profiles.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'profiles.read',
'profiles.manage',
'profiles.switch'
)
WHERE r.code IN (
'ORG_ADMIN',
'OWNER',
'DIRECTOR'
)
ON CONFLICT DO NOTHING;

-- Other staff can read/switch their own profiles later through auth endpoints.
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'profiles.read',
'profiles.switch'
)
WHERE r.code IN (
'FINANCE_MANAGER',
'ACCOUNTANT',
'SALES_MANAGER',
'TEACHER',
'STUDENT',
'PARENT'
)
ON CONFLICT DO NOTHING;
