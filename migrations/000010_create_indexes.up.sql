CREATE INDEX idx_users_organization_id ON users(organization_id);
CREATE INDEX idx_users_email ON users(email);

CREATE INDEX idx_branches_organization_id ON branches(organization_id);

CREATE INDEX idx_roles_organization_id ON roles(organization_id);

CREATE INDEX idx_subjects_organization_id ON subjects(organization_id);

CREATE INDEX idx_teacher_profiles_organization_id ON teacher_profiles(organization_id);

CREATE INDEX idx_students_organization_id ON students(organization_id);
CREATE INDEX idx_students_branch_id ON students(branch_id);
CREATE INDEX idx_students_full_name ON students(full_name);

CREATE INDEX idx_parents_organization_id ON parents(organization_id);
CREATE INDEX idx_parents_phone ON parents(phone);

CREATE INDEX idx_groups_organization_id ON groups(organization_id);
CREATE INDEX idx_groups_branch_id ON groups(branch_id);
CREATE INDEX idx_groups_teacher_id ON groups(teacher_id);
CREATE INDEX idx_groups_subject_id ON groups(subject_id);

CREATE INDEX idx_group_students_student_id ON group_students(student_id);

CREATE INDEX idx_schedules_organization_id ON schedules(organization_id);
CREATE INDEX idx_schedules_branch_id ON schedules(branch_id);
CREATE INDEX idx_schedules_group_id ON schedules(group_id);

CREATE INDEX idx_lessons_organization_id ON lessons(organization_id);
CREATE INDEX idx_lessons_branch_id ON lessons(branch_id);
CREATE INDEX idx_lessons_group_id ON lessons(group_id);
CREATE INDEX idx_lessons_teacher_id ON lessons(teacher_id);
CREATE INDEX idx_lessons_lesson_date ON lessons(lesson_date);

CREATE INDEX idx_attendance_lesson_id ON attendance(lesson_id);
CREATE INDEX idx_attendance_student_id ON attendance(student_id);

CREATE INDEX idx_homeworks_organization_id ON homeworks(organization_id);
CREATE INDEX idx_homeworks_group_id ON homeworks(group_id);

CREATE INDEX idx_payments_organization_id ON payments(organization_id);
CREATE INDEX idx_payments_branch_id ON payments(branch_id);
CREATE INDEX idx_payments_student_id ON payments(student_id);
CREATE INDEX idx_payments_payment_date ON payments(payment_date);

CREATE INDEX idx_payroll_periods_organization_id ON payroll_periods(organization_id);
CREATE INDEX idx_payroll_entries_organization_id ON payroll_entries(organization_id);
CREATE INDEX idx_payroll_entries_teacher_id ON payroll_entries(teacher_id);
CREATE INDEX idx_payroll_entries_period_id ON payroll_entries(period_id);

CREATE INDEX idx_files_organization_id ON files(organization_id);
CREATE INDEX idx_files_uploaded_by ON files(uploaded_by);

CREATE INDEX idx_notifications_organization_id ON notifications(organization_id);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);

CREATE INDEX idx_audit_logs_organization_id ON audit_logs(organization_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
