DELETE FROM role_permissions
WHERE permission_id IN (
SELECT id
FROM permissions
WHERE code IN (
'clients.read',
'clients.create',
'clients.update',
'clients.delete',

'parents.read',
'parents.create',
'parents.update',
'parents.delete',

'messages.read',
'messages.send',
'messages.send.group',
'messages.templates.manage',
'messages.channels.manage',

'payroll.read_limited',
'payroll.preview',

'notifications.send.organization',
'notifications.send.platform'
)
);

DELETE FROM permissions
WHERE code IN (
'clients.read',
'clients.create',
'clients.update',
'clients.delete',

'parents.read',
'parents.create',
'parents.update',
'parents.delete',

'messages.read',
'messages.send',
'messages.send.group',
'messages.templates.manage',
'messages.channels.manage',

'payroll.read_limited',
'payroll.preview',

'notifications.send.organization',
'notifications.send.platform'
);
