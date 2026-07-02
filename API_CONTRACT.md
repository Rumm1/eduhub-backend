# EduHub Backend API Contract

## Base URL

```txt
http://localhost:8080/api/v1
```

## Authentication

Most endpoints require JWT access token.

```txt
Authorization: Bearer <access_token>
```

## Standard Success Response

```json
{
  "success": true,
  "data": {}
}
```

## Standard Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message"
  }
}
```

---

# 1. Auth

## Login

```http
POST /auth/login
```

### Request

```json
{
  "email": "admin@smarthub.kz",
  "password": "Admin123!"
}
```

### Response

```json
{
  "success": true,
  "data": {
    "access_token": "...",
    "refresh_token": "...",
    "user": {}
  }
}
```

## Current User

```http
GET /auth/me
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "profile_id": "uuid",
    "organization_id": "uuid",
    "email": "admin@smarthub.kz",
    "full_name": "SmartHub Admin",
    "roles": ["ORG_ADMIN"],
    "permissions": [],
    "branch_ids": [],
    "current_profile": {},
    "available_profiles": []
  }
}
```

## Switch Profile

```http
POST /auth/switch-profile
```

### Request

```json
{
  "profile_id": "uuid"
}
```

---

# 2. Platform

Only `SUPER_ADMIN`.

## Create Organization

```http
POST /platform/organizations
```

### Request

```json
{
  "name": "SmartHub",
  "bin": "123456789012",
  "admin_email": "admin@smarthub.kz",
  "admin_password": "Admin123!",
  "admin_full_name": "SmartHub Admin"
}
```

---

# 3. Dashboard

## Dashboard Overview

```http
GET /dashboard/overview
```

### Response

```json
{
  "success": true,
  "data": {
    "students_count": 1,
    "teachers_count": 2,
    "groups_count": 1,
    "lessons_today": 1,
    "payments_this_month": 1,
    "payments_amount_this_month": 45000,
    "student_debt_total": 5000,
    "pending_payroll_entries": 1,
    "unread_notifications": 0,
    "recent_audit_logs": []
  }
}
```

---

# 4. AI Insights

## Dashboard AI Insights

```http
GET /ai/insights/dashboard
```

### Response

```json
{
  "success": true,
  "data": {
    "summary": "CRM overview: 2 active students, 1 active groups, 2 teachers, 1 lessons today. Current operational risk level is low.",
    "risk_level": "low",
    "metrics": {
      "students_count": 2,
      "teachers_count": 2,
      "groups_count": 1,
      "lessons_today": 1,
      "payments_this_month": 1,
      "payments_amount_this_month": 45000,
      "student_debt_total": 5000,
      "pending_payroll_entries": 1,
      "unread_notifications": 0,
      "recent_audit_logs_count": 5
    },
    "insights": [
      {
        "type": "payments",
        "severity": "low",
        "title": "Student debt detected",
        "message": "Current estimated student debt is 5000.00. Finance or managers should review unpaid balances."
      }
    ],
    "recommendations": [
      "Review student balances and contact parents with unpaid invoices."
    ]
  }
}
```

---

# 5. Branches

## Create Branch

```http
POST /branches
```

### Request

```json
{
  "name": "Skillset",
  "address": "Astana",
  "phone": "+77000000000"
}
```

## List Branches

```http
GET /branches
```

---

# 6. Subjects

## Create Subject

```http
POST /subjects
```

### Request

```json
{
  "name": "Programming",
  "description": "Programming course"
}
```

## List Subjects

```http
GET /subjects
```

---

# 7. Users

## Create User

```http
POST /users
```

### Request

```json
{
  "email": "teacher@smarthub.kz",
  "password": "Teacher123!",
  "full_name": "Teacher User",
  "phone": "+77000000001",
  "profiles": [
    {
      "organization_id": "uuid",
      "display_name": "Teacher User",
      "position": "Teacher",
      "profile_type": "teacher",
      "is_default": true,
      "roles": ["TEACHER"],
      "branch_ids": ["uuid"]
    }
  ]
}
```

## List Users

```http
GET /users
```

---

# 8. Profiles

## List User Profiles

```http
GET /users/{userID}/profiles
```

## Create User Profile

```http
POST /users/{userID}/profiles
```

### Request

```json
{
  "organization_id": "uuid",
  "display_name": "Dias Teacher",
  "position": "Teacher",
  "profile_type": "teacher",
  "is_default": false,
  "roles": ["TEACHER"],
  "branch_ids": ["uuid"]
}
```

## Get Profile

```http
GET /profiles/{profileID}
```

## Update Profile

```http
PATCH /profiles/{profileID}
```

### Request

```json
{
  "display_name": "Updated Name",
  "position": "Director",
  "profile_type": "director",
  "status": "active"
}
```

## Disable Profile

```http
DELETE /profiles/{profileID}
```

## Set Default Profile

```http
POST /profiles/{profileID}/set-default
```

## Add Role to Profile

```http
POST /profiles/{profileID}/roles
```

### Request

```json
{
  "role_code": "TEACHER"
}
```

## Remove Role from Profile

```http
DELETE /profiles/{profileID}/roles/{roleCode}
```

## Add Branch to Profile

```http
POST /profiles/{profileID}/branches
```

### Request

```json
{
  "branch_id": "uuid"
}
```

## Remove Branch from Profile

```http
DELETE /profiles/{profileID}/branches/{branchID}
```

---

# 9. Roles and Permissions

## List Roles

```http
GET /roles
```

## Create Role

```http
POST /roles
```

### Request

```json
{
  "name": "Manager",
  "code": "MANAGER",
  "description": "Manager role"
}
```

## Get Role

```http
GET /roles/{roleID}
```

## Update Role

```http
PATCH /roles/{roleID}
```

### Request

```json
{
  "name": "Updated Manager",
  "description": "Updated description"
}
```

## Delete Role

```http
DELETE /roles/{roleID}
```

## List Permissions

```http
GET /permissions
```

## List Permission Groups

```http
GET /permissions/groups
```

## Add Permission to Role

```http
POST /roles/{roleID}/permissions
```

### Request

```json
{
  "permission_code": "students.read"
}
```

## Remove Permission from Role

```http
DELETE /roles/{roleID}/permissions/{permissionCode}
```

---

# 10. Teachers

## Create Teacher Profile

```http
POST /teachers
```

### Request

```json
{
  "user_id": "uuid",
  "bio": "Teacher bio",
  "experience_years": 3,
  "employment_type": "hourly",
  "hourly_rate": 5000,
  "fixed_salary": 0,
  "subject_ids": ["uuid"]
}
```

## List Teachers

```http
GET /teachers
```

---

# 11. Students

## Create Student

```http
POST /students
```

### Request

```json
{
  "branch_id": "uuid",
  "full_name": "Student Name",
  "phone": "+77000000000",
  "birth_date": "2010-01-01",
  "gender": "male",
  "source": "instagram",
  "notes": "Test student",
  "parent": {
    "full_name": "Parent Name",
    "phone": "+77000000001",
    "email": "parent@example.com",
    "relation": "father"
  }
}
```

## List Students

```http
GET /students
```

---

# 12. Parents

## List Parents

```http
GET /parents
```

## Create Parent

```http
POST /parents
```

### Request

```json
{
  "full_name": "Parent Name",
  "phone": "+77000000001",
  "email": "parent@example.com"
}
```

## Get Parent

```http
GET /parents/{parentID}
```

## Update Parent

```http
PATCH /parents/{parentID}
```

### Request

```json
{
  "full_name": "Updated Parent Name",
  "phone": "+77000000002",
  "email": "updated.parent@example.com"
}
```

## Delete Parent

```http
DELETE /parents/{parentID}
```

## List Parent Students

```http
GET /parents/{parentID}/students
```

## Attach Student to Parent

```http
POST /parents/{parentID}/students/{studentID}
```

### Request

```json
{
  "relation": "father"
}
```

## Detach Student from Parent

```http
DELETE /parents/{parentID}/students/{studentID}
```

---

# 13. Groups

## Create Group

```http
POST /groups
```

### Request

```json
{
  "branch_id": "uuid",
  "subject_id": "uuid",
  "teacher_id": "uuid",
  "name": "Programming Group A",
  "level": "Beginner",
  "max_students": 15,
  "start_date": "2026-07-01",
  "end_date": "2026-12-31",
  "homework_enabled": true,
  "monthly_price": 50000
}
```

## List Groups

```http
GET /groups
```

## Add Student to Group

```http
POST /groups/{groupID}/students
```

### Request

```json
{
  "student_id": "uuid"
}
```

## List Group Students

```http
GET /groups/{groupID}/students
```

---

# 14. Lessons

## Create Lesson

```http
POST /lessons
```

### Request

```json
{
  "branch_id": "uuid",
  "group_id": "uuid",
  "teacher_id": "uuid",
  "subject_id": "uuid",
  "lesson_date": "2026-07-02",
  "start_time": "10:00",
  "end_time": "11:30",
  "topic": "Introduction",
  "status": "planned"
}
```

## List Lessons

```http
GET /lessons
```

### Query Params

```txt
branch_id
group_id
teacher_id
from_date
to_date
status
```

## Replace Lesson Teacher

```http
PATCH /lessons/{lessonID}/teacher
```

### Request

```json
{
  "actual_teacher_id": "uuid",
  "substitution_reason": "Teacher replacement"
}
```

---

# 15. Attendance

## Get Lesson Attendance

```http
GET /attendance/lessons/{lessonID}
```

## Mark Attendance

```http
POST /attendance/lessons/{lessonID}/mark
```

### Request

```json
{
  "items": [
    {
      "student_id": "uuid",
      "status": "present",
      "comment": "On time"
    }
  ]
}
```

### Attendance Statuses

```txt
present
absent
late
excused
```

---

# 16. Homework

## Create Homework

```http
POST /homeworks
```

### Request

```json
{
  "group_id": "uuid",
  "lesson_id": "uuid",
  "title": "Homework title",
  "description": "Homework description",
  "deadline": "2026-07-10"
}
```

## List Homeworks

```http
GET /homeworks
```

## Get Lesson Homeworks

```http
GET /homeworks/lessons/{lessonID}
```

---

# 17. Schedules

## Create Schedule

```http
POST /schedules
```

### Request

```json
{
  "branch_id": "uuid",
  "group_id": "uuid",
  "teacher_id": "uuid",
  "subject_id": "uuid",
  "day_of_week": 1,
  "start_time": "10:00",
  "end_time": "11:30",
  "start_date": "2026-07-01",
  "end_date": "2026-12-31"
}
```

## List Schedules

```http
GET /schedules
```

## Generate Lessons from Schedule

```http
POST /schedules/{scheduleID}/generate-lessons
```

### Request

```json
{
  "from_date": "2026-07-01",
  "to_date": "2026-07-31"
}
```

---

# 18. Payments

## Create Payment

```http
POST /payments
```

### Request

```json
{
  "branch_id": "uuid",
  "student_id": "uuid",
  "group_id": "uuid",
  "amount": 45000,
  "payment_date": "2026-07-01",
  "payment_period": "2026-07-01",
  "payment_method": "cash",
  "status": "paid",
  "comment": "July payment"
}
```

## List Payments

```http
GET /payments
```

### Query Params

```txt
branch_id
group_id
student_id
from_date
to_date
status
```

## Student Payments

```http
GET /payments/students/{studentID}
```

## Update Group Monthly Price

```http
PATCH /payments/groups/{groupID}/price
```

### Request

```json
{
  "monthly_price": 50000
}
```

## Student Balance

```http
GET /payments/students/{studentID}/balance
```

### Query Params

```txt
group_id
period=YYYY-MM
```

---

# 19. Payroll

## Create Payroll Period

```http
POST /payroll/periods
```

### Request

```json
{
  "period": "2026-07",
  "comment": "July payroll"
}
```

## Generate Payroll

```http
POST /payroll/periods/{periodID}/generate
```

## List Payroll Entries

```http
GET /payroll/entries
```

### Query Params

```txt
period_id
teacher_id
status
teacher_confirmation_status
```

## My Payroll Entries

```http
GET /payroll/entries/my
```

## Add Payroll Adjustment

```http
POST /payroll/entries/{entryID}/adjustments
```

### Request

```json
{
  "type": "bonus",
  "amount": 5000,
  "reason": "Extra lesson"
}
```

## Send Payroll Entry to Teacher

```http
POST /payroll/entries/{entryID}/send-to-teacher
```

## Teacher Confirm Payroll Entry

```http
POST /payroll/entries/{entryID}/confirm
```

## Teacher Dispute Payroll Entry

```http
POST /payroll/entries/{entryID}/dispute
```

### Request

```json
{
  "reason": "Lesson count is incorrect"
}
```

## Finance Approve Payroll Entry

```http
POST /payroll/entries/{entryID}/approve
```

## Mark Payroll Entry as Paid

```http
POST /payroll/entries/{entryID}/mark-paid
```

### Payroll Statuses

```txt
draft
sent_to_teacher
teacher_approved
teacher_disputed
approved_by_finance
paid
```

---

# 20. Reports

Reports support four formats:

```txt
json
xlsx
pdf
docx
```

Reports support three languages:

```txt
en
ru
kk
```

If `format` is empty, backend returns JSON.

## Teacher Schedule Report

```http
GET /reports/teacher-schedule
```

### Query Params

```txt
teacher_id
from_date=YYYY-MM-DD
to_date=YYYY-MM-DD
format=json|xlsx|pdf|docx
lang=en|ru|kk
```

### Example

```http
GET /reports/teacher-schedule?teacher_id=uuid&from_date=2026-07-01&to_date=2026-07-31&format=pdf&lang=ru
```

## Payments Report

```http
GET /reports/payments
```

### Query Params

```txt
from_date=YYYY-MM-DD
to_date=YYYY-MM-DD
branch_id
group_id
student_id
status
format=json|xlsx|pdf|docx
lang=en|ru|kk
```

### Example

```http
GET /reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=docx&lang=kk
```

## Student Balances Report

```http
GET /reports/student-balances
```

### Query Params

```txt
period=YYYY-MM
branch_id
group_id
student_id
status=paid|partial|unpaid
format=json|xlsx|pdf|docx
lang=en|ru|kk
```

### Example

```http
GET /reports/student-balances?period=2026-07&format=xlsx&lang=ru
```

## Payroll Report

```http
GET /reports/payroll
```

### Query Params

```txt
period=YYYY-MM
teacher_id
status
teacher_confirmation_status
format=json|xlsx|pdf|docx
lang=en|ru|kk
```

### Example

```http
GET /reports/payroll?period=2026-07&format=pdf&lang=en
```

---

# 21. Files

## Upload File

```http
POST /files/upload
```

### Content-Type

```txt
multipart/form-data
```

### Fields

```txt
file
folder
```

### Example Folders

```txt
avatars
homeworks
materials
reports
general
```

### Response

```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "organization_id": "uuid",
    "uploaded_by": "uuid",
    "folder": "avatars",
    "file_name": "avatar.jpg",
    "file_path": "uploads/avatars/uuid_avatar.jpg",
    "file_url": "/uploads/avatars/uuid_avatar.jpg",
    "mime_type": "image/jpeg",
    "size_bytes": 12345,
    "created_at": "2026-07-02 12:12:02"
  }
}
```

## List Files

```http
GET /files
```

## Get File

```http
GET /files/{fileID}
```

## Delete File

```http
DELETE /files/{fileID}
```

## Static File URL

```txt
/uploads/{folder}/{file}
```

---

# 22. Notifications

## List Notifications

```http
GET /notifications
```

## Create Notification

```http
POST /notifications
```

### Send to One User

```json
{
  "user_id": "uuid",
  "title": "Important notification",
  "message": "Message text",
  "type": "important"
}
```

### Send to Multiple Users

```json
{
  "user_ids": ["uuid", "uuid"],
  "title": "Group notification",
  "message": "Message text",
  "type": "message"
}
```

### Broadcast to Organization

```json
{
  "title": "System notification",
  "message": "Message text",
  "type": "system"
}
```

## Notification Types

```http
GET /notifications/types
```

### Types

```txt
normal
important
system
warning
payment
schedule
homework
payroll
lesson
message
```

## Mark Notification as Read

```http
PATCH /notifications/{notificationID}/read
```

## Mark All Notifications as Read

```http
PATCH /notifications/read-all
```

## Delete Notification

```http
DELETE /notifications/{notificationID}
```

---

# 23. Imports

## Students and Parents Import Preview

```http
POST /imports/students/preview
```

### Content-Type

```txt
multipart/form-data
```

### Fields

```txt
file
```

### Supported Formats

```txt
.csv
.xlsx
```

### Expected Columns

```txt
student_full_name
student_phone
parent_full_name
parent_phone
parent_email
group_name
relation
```

### Example CSV

```csv
student_full_name,student_phone,parent_full_name,parent_phone,parent_email,group_name,relation
Import Test Student,+77001112233,Import Test Parent,+77001112244,import.parent@example.com,Programming Group A,father
```

### Response

```json
{
  "success": true,
  "data": {
    "summary": {
      "total_rows": 1,
      "valid_rows": 1,
      "invalid_rows": 0,
      "warning_rows": 0
    },
    "rows": [
      {
        "row_number": 2,
        "status": "valid",
        "errors": [],
        "warnings": [],
        "student_full_name": "Import Test Student",
        "student_phone": "+77001112233",
        "parent_full_name": "Import Test Parent",
        "parent_phone": "+77001112244",
        "parent_email": "import.parent@example.com",
        "group_name": "Programming Group A",
        "branch_name": "Skillset",
        "relation": "father"
      }
    ]
  }
}
```

## Students and Parents Import Confirm

```http
POST /imports/students/confirm
```

### Content-Type

```txt
multipart/form-data
```

### Fields

```txt
file
```

### Response

```json
{
  "success": true,
  "data": {
    "summary": {
      "total_rows": 1,
      "valid_rows": 1,
      "invalid_rows": 0,
      "warning_rows": 0,
      "created_students": 1,
      "reused_students": 0,
      "created_parents": 1,
      "reused_parents": 0,
      "linked_parents_to_students": 1,
      "linked_students_to_groups": 1
    },
    "rows": []
  }
}
```

---

# 24. Audit Logs

## List Audit Logs

```http
GET /audit-logs
```

### Query Params

```txt
user_id
action
entity_type
entity_id
from_date
to_date
limit
offset
format=xlsx
lang=en|ru|kk
```

### JSON Example

```http
GET /audit-logs?limit=100&offset=0
```

### XLSX Export Example

```http
GET /audit-logs?format=xlsx&lang=ru&limit=100&offset=0
```

---

# 25. Frontend Notes

## Auth Flow

1. Login using `/auth/login`
2. Save `access_token`
3. Send token in `Authorization` header
4. Load current user using `/auth/me`
5. Use `permissions` array to show or hide UI actions
6. Use `available_profiles` for profile switching

## Permission-Based UI

Frontend should check permissions before showing buttons.

Examples:

```txt
students.create
students.update
students.delete
parents.create
payments.manage
payroll.approve
reports.export
files.upload
notifications.manage
dashboard.overview.read
```

## Report Download Flow

Frontend should build download URL using:

```txt
report type
filters
format
lang
```

Example:

```txt
/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=pdf&lang=ru
```

## File Upload Flow

Use `multipart/form-data`.

Fields:

```txt
file
folder
```

## Import Flow

Recommended frontend flow:

```txt
1. Upload file to /imports/students/preview
2. Show validation result to manager
3. If all rows are valid, upload the same file to /imports/students/confirm
4. Refresh students, parents and groups
```

## AI Insights Usage

Use this endpoint on dashboard page:

```http
GET /ai/insights/dashboard
```

Display:

```txt
summary
risk_level
insights
recommendations
metrics
```

---

# 26. Status

```txt
Core CRM backend: completed
Auth: completed
Users / Profiles / Roles / Permissions: completed
Branches / Subjects / Teachers / Students / Parents / Groups: completed
Lessons / Attendance / Homework / Schedule: completed
Payments / Payroll: completed
Reports JSON / XLSX / PDF / DOCX: completed
Report languages en / ru / kk: completed
Files: completed
Notifications: completed
Dashboard: completed
Imports: completed
AI Insights: completed
Messaging campaigns: future feature
```
