CREATE TABLE IF NOT EXISTS audit_logs (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
organization_id UUID REFERENCES organizations(id) ON DELETE SET NULL,
user_id UUID REFERENCES users(id) ON DELETE SET NULL,

action VARCHAR(120) NOT NULL,
entity_type VARCHAR(80) NOT NULL,
entity_id UUID,

description TEXT,
metadata JSONB NOT NULL DEFAULT '{}'::jsonb,

ip_address VARCHAR(80),
user_agent TEXT,

created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at);

INSERT INTO permissions (id, code, description)
VALUES
(gen_random_uuid(), 'audit_logs.read', 'Read audit logs'),
(gen_random_uuid(), 'audit_logs.export', 'Export audit logs')
ON CONFLICT (code) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code IN (
'audit_logs.read',
'audit_logs.export'
)
WHERE r.code IN ('ORG_ADMIN')
ON CONFLICT DO NOTHING;
