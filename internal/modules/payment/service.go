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
	ErrGroupIDRequired      = errors.New("group id is required")
	ErrGroupIDInvalid       = errors.New("group id is invalid")
	ErrGroupNotFound        = errors.New("group not found in organization")
	ErrStudentNotInGroup    = errors.New("student is not in group")
	ErrAmountRequired       = errors.New("amount is required")
	ErrAmountInvalid        = errors.New("amount is invalid")
	ErrMonthlyPriceRequired = errors.New("monthly price is required")
	ErrMonthlyPriceInvalid  = errors.New("monthly price is invalid")
	ErrPaymentDateRequired  = errors.New("payment date is required")
	ErrPaymentDateInvalid   = errors.New("payment date is invalid")
	ErrPaymentPeriodInvalid = errors.New("payment period is invalid")
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

	parsedPaymentDate, err := time.Parse("2006-01-02", paymentDate)
	if err != nil {
		return PaymentResponse{}, ErrPaymentDateInvalid
	}

	paymentPeriod, err := normalizePaymentPeriod(req.PaymentPeriod, parsedPaymentDate)
	if err != nil {
		return PaymentResponse{}, err
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
		PaymentPeriod:  paymentPeriod,
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

func (s *Service) UpdateGroupPrice(
	ctx context.Context,
	groupIDRaw string,
	req UpdateGroupPriceRequest,
) (GroupPriceResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return GroupPriceResponse{}, ErrTenantRequired
	}

	groupID, err := uuid.Parse(strings.TrimSpace(groupIDRaw))
	if err != nil {
		return GroupPriceResponse{}, ErrGroupIDInvalid
	}

	monthlyPrice := strings.TrimSpace(req.MonthlyPrice)
	if monthlyPrice == "" {
		return GroupPriceResponse{}, ErrMonthlyPriceRequired
	}

	parsedPrice, err := strconv.ParseFloat(monthlyPrice, 64)
	if err != nil || parsedPrice < 0 {
		return GroupPriceResponse{}, ErrMonthlyPriceInvalid
	}

	updatedPrice, err := s.repo.UpdateGroupPrice(ctx, *currentUser.OrganizationID, groupID, monthlyPrice)
	if err != nil {
		return GroupPriceResponse{}, err
	}

	return GroupPriceResponse{
		GroupID:      groupID.String(),
		MonthlyPrice: updatedPrice,
	}, nil
}

func (s *Service) GetStudentBalance(
	ctx context.Context,
	studentIDRaw string,
	groupIDRaw string,
	periodRaw string,
) (StudentBalanceResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return StudentBalanceResponse{}, ErrTenantRequired
	}

	studentID, err := uuid.Parse(strings.TrimSpace(studentIDRaw))
	if err != nil {
		return StudentBalanceResponse{}, ErrStudentIDInvalid
	}

	groupID, err := uuid.Parse(strings.TrimSpace(groupIDRaw))
	if err != nil {
		return StudentBalanceResponse{}, ErrGroupIDInvalid
	}

	paymentPeriod, err := normalizePaymentPeriodFromString(periodRaw)
	if err != nil {
		return StudentBalanceResponse{}, err
	}

	balance, err := s.repo.GetStudentBalance(ctx, *currentUser.OrganizationID, studentID, groupID, paymentPeriod)
	if err != nil {
		return StudentBalanceResponse{}, err
	}

	return StudentBalanceResponse{
		StudentID:      balance.StudentID.String(),
		GroupID:        balance.GroupID.String(),
		BranchID:       balance.BranchID.String(),
		PaymentPeriod:  balance.PaymentPeriod,
		ExpectedAmount: balance.ExpectedAmount,
		PaidAmount:     balance.PaidAmount,
		DebtAmount:     balance.DebtAmount,
		OverpaidAmount: balance.OverpaidAmount,
		Status:         balance.Status,
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
		PaymentPeriod:  payment.PaymentPeriod,
		PaymentMethod:  payment.PaymentMethod,
		Status:         payment.Status,
		Comment:        payment.Comment,
	}
}

func normalizePaymentPeriod(paymentPeriodRaw string, paymentDate time.Time) (string, error) {
	paymentPeriodRaw = strings.TrimSpace(paymentPeriodRaw)
	if paymentPeriodRaw == "" {
		return time.Date(paymentDate.Year(), paymentDate.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"), nil
	}

	return normalizePaymentPeriodFromString(paymentPeriodRaw)
}

func normalizePaymentPeriodFromString(paymentPeriodRaw string) (string, error) {
	paymentPeriodRaw = strings.TrimSpace(paymentPeriodRaw)
	if paymentPeriodRaw == "" {
		return "", ErrPaymentPeriodInvalid
	}

	parsedPeriod, err := time.Parse("2006-01", paymentPeriodRaw)
	if err != nil {
		return "", ErrPaymentPeriodInvalid
	}

	return time.Date(parsedPeriod.Year(), parsedPeriod.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"), nil
}

func isValidStatus(status string) bool {
	switch status {
	case "paid", "pending", "cancelled", "refunded":
		return true
	default:
		return false
	}
}
