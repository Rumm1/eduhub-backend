DELETE FROM permissions
WHERE code IN (
    'organizations.read',
    'organizations.manage',

    'branches.read',
    'branches.create',
    'branches.update',
    'branches.delete',

    'users.read',
    'users.create',
    'users.update',
    'users.delete',

    'roles.read',
    'roles.manage',

    'subjects.read',
    'subjects.create',
    'subjects.update',
    'subjects.delete',

    'teachers.read',
    'teachers.create',
    'teachers.update',
    'teachers.delete',

    'students.read',
    'students.create',
    'students.update',
    'students.delete',

    'groups.read',
    'groups.create',
    'groups.update',
    'groups.delete',

    'lessons.read',
    'lessons.create',
    'lessons.update',
    'lessons.delete',

    'attendance.read',
    'attendance.manage',

    'homeworks.read',
    'homeworks.manage',

    'payments.read',
    'payments.manage',

    'payroll.read',
    'payroll.manage',
    'payroll.approve',
    'payroll.rules.manage',

    'files.upload',
    'files.read',
    'files.delete',

    'notifications.read',
    'notifications.manage',

    'audit.read'
);
