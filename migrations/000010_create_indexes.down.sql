DROP INDEX IF EXISTS idx_audit_logs_created_at;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_user_id;
DROP INDEX IF EXISTS idx_audit_logs_organization_id;

DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_notifications_organization_id;

DROP INDEX IF EXISTS idx_files_uploaded_by;
DROP INDEX IF EXISTS idx_files_organization_id;

DROP INDEX IF EXISTS idx_payroll_entries_period_id;
DROP INDEX IF EXISTS idx_payroll_entries_teacher_id;
DROP INDEX IF EXISTS idx_payroll_entries_organization_id;
DROP INDEX IF EXISTS idx_payroll_periods_organization_id;

DROP INDEX IF EXISTS idx_payments_payment_date;
DROP INDEX IF EXISTS idx_payments_student_id;
DROP INDEX IF EXISTS idx_payments_branch_id;
DROP INDEX IF EXISTS idx_payments_organization_id;

DROP INDEX IF EXISTS idx_homeworks_group_id;
DROP INDEX IF EXISTS idx_homeworks_organization_id;

DROP INDEX IF EXISTS idx_attendance_student_id;
DROP INDEX IF EXISTS idx_attendance_lesson_id;

DROP INDEX IF EXISTS idx_lessons_lesson_date;
DROP INDEX IF EXISTS idx_lessons_teacher_id;
DROP INDEX IF EXISTS idx_lessons_group_id;
DROP INDEX IF EXISTS idx_lessons_branch_id;
DROP INDEX IF EXISTS idx_lessons_organization_id;

DROP INDEX IF EXISTS idx_schedules_group_id;
DROP INDEX IF EXISTS idx_schedules_branch_id;
DROP INDEX IF EXISTS idx_schedules_organization_id;

DROP INDEX IF EXISTS idx_group_students_student_id;

DROP INDEX IF EXISTS idx_groups_subject_id;
DROP INDEX IF EXISTS idx_groups_teacher_id;
DROP INDEX IF EXISTS idx_groups_branch_id;
DROP INDEX IF EXISTS idx_groups_organization_id;

DROP INDEX IF EXISTS idx_parents_phone;
DROP INDEX IF EXISTS idx_parents_organization_id;

DROP INDEX IF EXISTS idx_students_full_name;
DROP INDEX IF EXISTS idx_students_branch_id;
DROP INDEX IF EXISTS idx_students_organization_id;

DROP INDEX IF EXISTS idx_teacher_profiles_organization_id;

DROP INDEX IF EXISTS idx_subjects_organization_id;

DROP INDEX IF EXISTS idx_roles_organization_id;

DROP INDEX IF EXISTS idx_branches_organization_id;

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_organization_id;
