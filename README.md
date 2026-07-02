@'
# EduHub Backend

EduHub Backend is a Go-based CRM/SaaS backend for educational centers.  
The system is designed for managing organizations, branches, users, roles, students, parents, teachers, groups, lessons, attendance, homework, payments, payroll, reports, files, imports, notifications, and AI dashboard insights.

## Tech Stack

- Go 1.26+
- PostgreSQL
- Redis
- Chi Router
- JWT Authentication
- RBAC Permissions
- Docker Compose
- Excelize for XLSX reports
- FPDF for PDF reports
- Custom DOCX export via OOXML

## Main Features

- Multi-tenant organization structure
- Authentication with JWT access and refresh tokens
- User profiles with multiple roles
- Role and permission management
- Branch management
- Subject management
- Teacher management
- Student management
- Parent management
- Group management
- Lesson scheduling
- Attendance tracking
- Homework management
- Payments and student balance tracking
- Payroll workflow for teachers
- Audit logs
- File upload and storage
- Notifications
- Dashboard overview
- AI dashboard insights
- Import students and parents from CSV/XLSX
- Reports export in JSON, XLSX, PDF, and DOCX
- Report language support: English, Russian, Kazakh

## Project Structure

```txt
cmd/api                 Application entry point
internal/app            Router setup
internal/config         Environment configuration
internal/middleware     Auth, tenant, role and permission middleware
internal/modules        Business modules
internal/platform       Database, JWT, Redis, logger, storage
internal/shared         Shared response, context, constants and helpers
migrations              Database migrations
scripts                 Utility scripts
Main Modules
auth
organization
branch
user
profile
role
permission
subject
teacher
student
parent
group
lesson
attendance
homework
schedule
payment
payroll
report
audit
file
notification
dashboard
importer
ai
API Overview
Auth
POST /api/v1/auth/login
GET  /api/v1/auth/me
POST /api/v1/auth/switch-profile
Platform
POST /api/v1/platform/organizations
Branches
POST /api/v1/branches
GET  /api/v1/branches
Users and Profiles
POST   /api/v1/users
GET    /api/v1/users
GET    /api/v1/users/{userID}/profiles
POST   /api/v1/users/{userID}/profiles
GET    /api/v1/profiles/{profileID}
PATCH  /api/v1/profiles/{profileID}
DELETE /api/v1/profiles/{profileID}
POST   /api/v1/profiles/{profileID}/set-default
POST   /api/v1/profiles/{profileID}/roles
DELETE /api/v1/profiles/{profileID}/roles/{roleCode}
POST   /api/v1/profiles/{profileID}/branches
DELETE /api/v1/profiles/{profileID}/branches/{branchID}
Roles and Permissions
GET    /api/v1/roles
POST   /api/v1/roles
GET    /api/v1/roles/{roleID}
PATCH  /api/v1/roles/{roleID}
DELETE /api/v1/roles/{roleID}

GET    /api/v1/permissions
GET    /api/v1/permissions/groups
POST   /api/v1/roles/{roleID}/permissions
DELETE /api/v1/roles/{roleID}/permissions/{permissionCode}
Students, Parents, Teachers and Groups
POST /api/v1/students
GET  /api/v1/students

GET    /api/v1/parents
POST   /api/v1/parents
GET    /api/v1/parents/{parentID}
PATCH  /api/v1/parents/{parentID}
DELETE /api/v1/parents/{parentID}
GET    /api/v1/parents/{parentID}/students
POST   /api/v1/parents/{parentID}/students/{studentID}
DELETE /api/v1/parents/{parentID}/students/{studentID}

POST /api/v1/teachers
GET  /api/v1/teachers

POST /api/v1/groups
GET  /api/v1/groups
POST /api/v1/groups/{groupID}/students
GET  /api/v1/groups/{groupID}/students
Lessons, Attendance, Homework and Schedule
POST  /api/v1/lessons
GET   /api/v1/lessons
PATCH /api/v1/lessons/{lessonID}/teacher

GET  /api/v1/attendance/lessons/{lessonID}
POST /api/v1/attendance/lessons/{lessonID}/mark

POST /api/v1/homeworks
GET  /api/v1/homeworks
GET  /api/v1/homeworks/lessons/{lessonID}

POST /api/v1/schedules
GET  /api/v1/schedules
POST /api/v1/schedules/{scheduleID}/generate-lessons
Payments and Payroll
POST  /api/v1/payments
GET   /api/v1/payments
GET   /api/v1/payments/students/{studentID}
PATCH /api/v1/payments/groups/{groupID}/price
GET   /api/v1/payments/students/{studentID}/balance

POST /api/v1/payroll/periods
POST /api/v1/payroll/periods/{periodID}/generate
GET  /api/v1/payroll/entries
GET  /api/v1/payroll/entries/my
POST /api/v1/payroll/entries/{entryID}/adjustments
POST /api/v1/payroll/entries/{entryID}/send-to-teacher
POST /api/v1/payroll/entries/{entryID}/confirm
POST /api/v1/payroll/entries/{entryID}/dispute
POST /api/v1/payroll/entries/{entryID}/approve
POST /api/v1/payroll/entries/{entryID}/mark-paid
Reports

Reports support JSON, XLSX, PDF and DOCX formats.

GET /api/v1/reports/teacher-schedule
GET /api/v1/reports/payments
GET /api/v1/reports/student-balances
GET /api/v1/reports/payroll

Example:

GET /api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=pdf&lang=ru

Supported formats:

json
xlsx
pdf
docx

Supported languages:

en
ru
kk
Files
POST   /api/v1/files/upload
GET    /api/v1/files
GET    /api/v1/files/{fileID}
DELETE /api/v1/files/{fileID}
Notifications
GET    /api/v1/notifications
POST   /api/v1/notifications
GET    /api/v1/notifications/types
PATCH  /api/v1/notifications/{notificationID}/read
PATCH  /api/v1/notifications/read-all
DELETE /api/v1/notifications/{notificationID}
Import
POST /api/v1/imports/students/preview
POST /api/v1/imports/students/confirm

Supported import formats:

.csv
.xlsx

Expected columns:

student_full_name
student_phone
parent_full_name
parent_phone
parent_email
group_name
relation
Dashboard and AI Insights
GET /api/v1/dashboard/overview
GET /api/v1/ai/insights/dashboard
Environment Variables

Create a .env file:

APP_ENV=local
APP_PORT=8080

DATABASE_URL=postgres://eduhub:eduhub_password@localhost:5432/eduhub?sslmode=disable

REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_ACCESS_SECRET=change_me_access_secret
JWT_REFRESH_SECRET=change_me_refresh_secret
JWT_ACCESS_TTL_MINUTES=15
JWT_REFRESH_TTL_DAYS=30

CORS_ALLOWED_ORIGINS=http://localhost:5173

SUPER_ADMIN_EMAIL=superadmin@eduhub.kz
SUPER_ADMIN_PASSWORD=SuperAdmin123!
SUPER_ADMIN_FULL_NAME=EduHub Super Admin
Run Locally

Start PostgreSQL and Redis:

docker compose up -d

Run migrations:

migrate -path migrations -database "postgres://eduhub:eduhub_password@localhost:5432/eduhub?sslmode=disable" up

Run backend:

go run ./cmd/api

Health check:

GET /health
GET /api/v1/health
GET /health/db
GET /api/v1/health/db
Tests
go test ./...
Default Accounts

Super admin:

email: superadmin@eduhub.kz
password: SuperAdmin123!

Organization admin example:

email: admin@smarthub.kz
password: Admin123!

Teacher example:

email: teacher@smarthub.kz
password: Teacher123!
Report Export Examples

Payments PDF:

GET /api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=pdf&lang=ru

Payments DOCX:

GET /api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=docx&lang=kk

Payroll XLSX:

GET /api/v1/reports/payroll?period=2026-07&format=xlsx&lang=en
Import Example

CSV file:

student_full_name,student_phone,parent_full_name,parent_phone,parent_email,group_name,relation
Import Test Student,+77001112233,Import Test Parent,+77001112244,import.parent@example.com,Programming Group A,father

Preview:

POST /api/v1/imports/students/preview

Confirm:

POST /api/v1/imports/students/confirm
Status
Core CRM backend: completed
Reports JSON/XLSX/PDF/DOCX: completed
File module: completed
Dashboard overview: completed
Notifications: completed
Parents: completed
Import students and parents: completed
AI dashboard insights: completed
Messaging campaigns: future feature
License

This project is currently private and intended for internal development.
'@ | Set-Content -Encoding UTF8 README.md