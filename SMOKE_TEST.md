# EduHub Backend Smoke Test Checklist

This file contains the final smoke test checklist for EduHub Backend.

Base URL:

```txt
http://localhost:8080/api/v1
1. Health Check
Invoke-RestMethod -Uri "http://localhost:8080/health"
Invoke-RestMethod -Uri "http://localhost:8080/health/db"
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health"
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/health/db"

Expected result:

EduHub backend is running
PostgreSQL is connected
EduHub API is running
PostgreSQL is connected
2. Login
$login = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body '{
    "email": "admin@smarthub.kz",
    "password": "Admin123!"
  }'

$token = $login.data.access_token

Expected result:

Access token received
3. Current User
$me = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/auth/me" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$me | ConvertTo-Json -Depth 20

Expected result:

Current user, organization, profile, roles and permissions are returned.
4. Dashboard
$dashboard = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/dashboard/overview" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$dashboard | ConvertTo-Json -Depth 20

Expected result:

Dashboard metrics are returned.
5. AI Insights
$ai = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/ai/insights/dashboard" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$ai | ConvertTo-Json -Depth 20

Expected result:

AI summary, risk level, metrics, insights and recommendations are returned.
6. Branches
$branches = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/branches" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$branches | ConvertTo-Json -Depth 20

Expected result:

Branches are returned.
7. Subjects
$subjects = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/subjects" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$subjects | ConvertTo-Json -Depth 20

Expected result:

Subjects are returned.
8. Students
$students = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/students" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$students | ConvertTo-Json -Depth 20

Expected result:

Students are returned.
9. Parents
$parents = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/parents" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$parents | ConvertTo-Json -Depth 20

Expected result:

Parents are returned.
10. Teachers
$teachers = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/teachers" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$teachers | ConvertTo-Json -Depth 20

Expected result:

Teachers are returned.
11. Groups
$groups = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/groups" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$groups | ConvertTo-Json -Depth 20

Expected result:

Groups are returned.
12. Lessons
$lessons = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/lessons" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$lessons | ConvertTo-Json -Depth 20

Expected result:

Lessons are returned.
13. Schedules
$schedules = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/schedules" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$schedules | ConvertTo-Json -Depth 20

Expected result:

Schedules are returned.
14. Payments
$payments = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/payments" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$payments | ConvertTo-Json -Depth 20

Expected result:

Payments are returned.
15. Payroll
$payroll = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/payroll/entries" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$payroll | ConvertTo-Json -Depth 20

Expected result:

Payroll entries are returned.
16. Notifications
$notifications = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/notifications" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$notifications | ConvertTo-Json -Depth 20

Expected result:

Notifications are returned.
17. Files
$files = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/files" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$files | ConvertTo-Json -Depth 20

Expected result:

Files are returned.
18. Reports JSON
$paymentsReport = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$paymentsReport | ConvertTo-Json -Depth 20

Expected result:

Payments report JSON is returned.
19. Reports XLSX
Invoke-WebRequest `
  -Uri "http://localhost:8080/api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=xlsx&lang=ru" `
  -Headers @{ Authorization = "Bearer $token" } `
  -OutFile "payments-report.xlsx"

Expected result:

payments-report.xlsx file is downloaded.
20. Reports PDF
Invoke-WebRequest `
  -Uri "http://localhost:8080/api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=pdf&lang=ru" `
  -Headers @{ Authorization = "Bearer $token" } `
  -OutFile "payments-report.pdf"

Expected result:

payments-report.pdf file is downloaded.
21. Reports DOCX
Invoke-WebRequest `
  -Uri "http://localhost:8080/api/v1/reports/payments?from_date=2026-07-01&to_date=2026-07-31&format=docx&lang=ru" `
  -Headers @{ Authorization = "Bearer $token" } `
  -OutFile "payments-report.docx"

Expected result:

payments-report.docx file is downloaded.
22. Audit Logs
$auditLogs = Invoke-RestMethod `
  -Uri "http://localhost:8080/api/v1/audit-logs?limit=100&offset=0" `
  -Method GET `
  -Headers @{ Authorization = "Bearer $token" }

$auditLogs | ConvertTo-Json -Depth 20

Expected result:

Audit logs are returned.
23. Import Preview

Use a CSV or XLSX file with columns:

student_full_name
student_phone
parent_full_name
parent_phone
parent_email
group_name
relation

Endpoint:

POST /api/v1/imports/students/preview

Expected result:

Import validation result is returned.
24. Import Confirm

Endpoint:

POST /api/v1/imports/students/confirm

Expected result:

Students, parents and group links are created or reused.
25. Go Tests
go test ./...

Expected result:

All packages passed.
26. Git Status
git status

Expected result:

nothing to commit, working tree clean
Final Project Status
Auth: completed
Users: completed
Profiles: completed
Roles: completed
Permissions: completed
Branches: completed
Subjects: completed
Teachers: completed
Students: completed
Parents: completed
Groups: completed
Lessons: completed
Attendance: completed
Homework: completed
Schedules: completed
Payments: completed
Payroll: completed
Reports JSON/XLSX/PDF/DOCX: completed
Files: completed
Notifications: completed
Dashboard: completed
AI Insights: completed
Imports CSV/XLSX: completed
Audit logs: completed
Messaging campaigns: future feature

