package payment

import (
"context"
"errors"
"strconv"
"strings"
"time"

usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
"github.com/google/uuid"
)

var (
ErrTenantRequired       = errors.New("tenant organization is required")
ErrStudentIDRequired    = errors.New("student id is required")
ErrStudentIDInvalid     = errors.New("student id is invalid")
ErrStudentNotFound      = errors.New("student not found in organization")
ErrGroupIDInvalid       = errors.New("group id is invalid")
ErrStudentNotInGroup    = errors.New("student is not in group")
ErrAmountRequired       = errors.New("amount is required")
ErrAmountInvalid        = errors.New("amount is invalid")
ErrPaymentDateRequired  = errors.New("payment date is required")
ErrPaymentDateInvalid   = errors.New("payment date is invalid")
ErrPaymentStatusInvalid = errors.New("payment status is invalid")
)

type Service struct {
repo *Repository
}

func NewService(repo *Repository) *Service {
return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, req CreatePaymentRequest) (PaymentResponse, error) {
currentUser, ok := usercontext.GetUser(ctx)
if !ok || currentUser.OrganizationID == nil {
return PaymentResponse{}, ErrTenantRequired
}

studentIDRaw := strings.TrimSpace(req.StudentID)
if studentIDRaw == "" {
return PaymentResponse{}, ErrStudentIDRequired
}

studentID, err := uuid.Parse(studentIDRaw)
if err != nil {
return PaymentResponse{}, ErrStudentIDInvalid
}

groupID := strings.TrimSpace(req.GroupID)
if groupID != "" {
if _, err := uuid.Parse(groupID); err != nil {
return PaymentResponse{}, ErrGroupIDInvalid
}
}

amount := strings.TrimSpace(req.Amount)
if amount == "" {
return PaymentResponse{}, ErrAmountRequired
}

parsedAmount, err := strconv.ParseFloat(amount, 64)
if err != nil || parsedAmount <= 0 {
return PaymentResponse{}, ErrAmountInvalid
}

paymentDate := strings.TrimSpace(req.PaymentDate)
if paymentDate == "" {
return PaymentResponse{}, ErrPaymentDateRequired
}

if _, err := time.Parse("2006-01-02", paymentDate); err != nil {
return PaymentResponse{}, ErrPaymentDateInvalid
}

status := strings.ToLower(strings.TrimSpace(req.Status))
if status == "" {
status = "paid"
}

if !isValidStatus(status) {
return PaymentResponse{}, ErrPaymentStatusInvalid
}

newPayment := Payment{
ID:             uuid.New(),
OrganizationID: *currentUser.OrganizationID,
StudentID:      studentID,
GroupID:        groupID,
Amount:         amount,
PaymentDate:    paymentDate,
PaymentMethod:  strings.ToLower(strings.TrimSpace(req.PaymentMethod)),
Status:         status,
Comment:        strings.TrimSpace(req.Comment),
}

createdPayment, err := s.repo.Create(ctx, newPayment)
if err != nil {
return PaymentResponse{}, err
}

return buildPaymentResponse(createdPayment), nil
}

func (s *Service) List(ctx context.Context) (ListPaymentsResponse, error) {
currentUser, ok := usercontext.GetUser(ctx)
if !ok || currentUser.OrganizationID == nil {
return ListPaymentsResponse{}, ErrTenantRequired
}

payments, err := s.repo.ListByOrganizationID(ctx, *currentUser.OrganizationID)
if err != nil {
return ListPaymentsResponse{}, err
}

items := make([]PaymentResponse, 0, len(payments))

for _, item := range payments {
items = append(items, buildPaymentResponse(item))
}

return ListPaymentsResponse{
Items: items,
Total: len(items),
}, nil
}

func (s *Service) ListByStudentID(ctx context.Context, studentIDRaw string) (ListPaymentsResponse, error) {
currentUser, ok := usercontext.GetUser(ctx)
if !ok || currentUser.OrganizationID == nil {
return ListPaymentsResponse{}, ErrTenantRequired
}

studentID, err := uuid.Parse(strings.TrimSpace(studentIDRaw))
if err != nil {
return ListPaymentsResponse{}, ErrStudentIDInvalid
}

payments, err := s.repo.ListByStudentID(ctx, *currentUser.OrganizationID, studentID)
if err != nil {
return ListPaymentsResponse{}, err
}

items := make([]PaymentResponse, 0, len(payments))

for _, item := range payments {
items = append(items, buildPaymentResponse(item))
}

return ListPaymentsResponse{
Items: items,
Total: len(items),
}, nil
}

func buildPaymentResponse(payment Payment) PaymentResponse {
return PaymentResponse{
ID:             payment.ID.String(),
OrganizationID: payment.OrganizationID.String(),
BranchID:       payment.BranchID.String(),
StudentID:      payment.StudentID.String(),
GroupID:        payment.GroupID,
Amount:         payment.Amount,
PaymentDate:    payment.PaymentDate,
PaymentMethod:  payment.PaymentMethod,
Status:         payment.Status,
Comment:        payment.Comment,
}
}

func isValidStatus(status string) bool {
switch status {
case "paid", "pending", "cancelled", "refunded":
return true
default:
return false
}
}
