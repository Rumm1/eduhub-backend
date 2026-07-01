package payment

import (
"context"

"github.com/google/uuid"
"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, payment Payment) (Payment, error) {
tx, err := r.db.Begin(ctx)
if err != nil {
return Payment{}, err
}
defer func() {
_ = tx.Rollback(ctx)
}()

err = tx.QueryRow(ctx, `
SELECT branch_id
FROM students
WHERE id = $1
  AND organization_id = $2
  AND status = 'active'
`, payment.StudentID, payment.OrganizationID).Scan(&payment.BranchID)
if err != nil {
return Payment{}, ErrStudentNotFound
}

if payment.GroupID != "" {
var groupExists bool

err = tx.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM group_students gs
JOIN groups g ON g.id = gs.group_id
WHERE gs.group_id = $1
  AND gs.student_id = $2
  AND g.organization_id = $3
  AND g.branch_id = $4
  AND gs.status = 'active'
  AND g.status = 'active'
)
`, payment.GroupID, payment.StudentID, payment.OrganizationID, payment.BranchID).Scan(&groupExists)
if err != nil {
return Payment{}, err
}

if !groupExists {
return Payment{}, ErrStudentNotInGroup
}
}

var groupID interface{}
if payment.GroupID != "" {
groupID = payment.GroupID
}

err = tx.QueryRow(ctx, `
INSERT INTO payments (
id,
organization_id,
branch_id,
student_id,
group_id,
amount,
payment_date,
payment_method,
status,
comment
)
VALUES ($1, $2, $3, $4, $5::uuid, $6::numeric, $7::date, $8, $9, $10)
RETURNING
id,
organization_id,
branch_id,
student_id,
COALESCE(group_id::text, ''),
amount::text,
payment_date::text,
COALESCE(payment_method, ''),
status,
COALESCE(comment, '')
`,
payment.ID,
payment.OrganizationID,
payment.BranchID,
payment.StudentID,
groupID,
payment.Amount,
payment.PaymentDate,
payment.PaymentMethod,
payment.Status,
payment.Comment,
).Scan(
&payment.ID,
&payment.OrganizationID,
&payment.BranchID,
&payment.StudentID,
&payment.GroupID,
&payment.Amount,
&payment.PaymentDate,
&payment.PaymentMethod,
&payment.Status,
&payment.Comment,
)
if err != nil {
return Payment{}, err
}

if err := tx.Commit(ctx); err != nil {
return Payment{}, err
}

return payment, nil
}

func (r *Repository) ListByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]Payment, error) {
rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
student_id,
COALESCE(group_id::text, ''),
amount::text,
payment_date::text,
COALESCE(payment_method, ''),
status,
COALESCE(comment, '')
FROM payments
WHERE organization_id = $1
ORDER BY payment_date DESC, created_at DESC
`, organizationID)
if err != nil {
return nil, err
}
defer rows.Close()

payments := make([]Payment, 0)

for rows.Next() {
var item Payment

if err := rows.Scan(
&item.ID,
&item.OrganizationID,
&item.BranchID,
&item.StudentID,
&item.GroupID,
&item.Amount,
&item.PaymentDate,
&item.PaymentMethod,
&item.Status,
&item.Comment,
); err != nil {
return nil, err
}

payments = append(payments, item)
}

if err := rows.Err(); err != nil {
return nil, err
}

return payments, nil
}

func (r *Repository) ListByStudentID(ctx context.Context, organizationID uuid.UUID, studentID uuid.UUID) ([]Payment, error) {
var studentExists bool

err := r.db.QueryRow(ctx, `
SELECT EXISTS (
SELECT 1
FROM students
WHERE id = $1
  AND organization_id = $2
)
`, studentID, organizationID).Scan(&studentExists)
if err != nil {
return nil, err
}

if !studentExists {
return nil, ErrStudentNotFound
}

rows, err := r.db.Query(ctx, `
SELECT
id,
organization_id,
branch_id,
student_id,
COALESCE(group_id::text, ''),
amount::text,
payment_date::text,
COALESCE(payment_method, ''),
status,
COALESCE(comment, '')
FROM payments
WHERE organization_id = $1
  AND student_id = $2
ORDER BY payment_date DESC, created_at DESC
`, organizationID, studentID)
if err != nil {
return nil, err
}
defer rows.Close()

payments := make([]Payment, 0)

for rows.Next() {
var item Payment

if err := rows.Scan(
&item.ID,
&item.OrganizationID,
&item.BranchID,
&item.StudentID,
&item.GroupID,
&item.Amount,
&item.PaymentDate,
&item.PaymentMethod,
&item.Status,
&item.Comment,
); err != nil {
return nil, err
}

payments = append(payments, item)
}

if err := rows.Err(); err != nil {
return nil, err
}

return payments, nil
}
