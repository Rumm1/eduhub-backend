DELETE FROM role_permissions
WHERE permission_id IN (
SELECT id
FROM permissions
WHERE code IN (
'reports.teacher_schedule.read',
'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export'
)
);

DELETE FROM permissions
WHERE code IN (
'reports.teacher_schedule.read',
'reports.payments.read',
'reports.payroll.read',
'reports.student_balance.read',
'reports.export'
);
