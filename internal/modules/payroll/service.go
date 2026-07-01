package payroll

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
	ErrTenantRequired          = errors.New("tenant organization is required")
	ErrPeriodNotFound          = errors.New("payroll period not found")
	ErrPeriodIDInvalid         = errors.New("payroll period id is invalid")
	ErrEntryNotFound           = errors.New("payroll entry not found")
	ErrEntryIDInvalid          = errors.New("payroll entry id is invalid")
	ErrMonthInvalid            = errors.New("month is invalid")
	ErrYearInvalid             = errors.New("year is invalid")
	ErrAdjustmentTypeInvalid   = errors.New("adjustment type is invalid")
	ErrAdjustmentAmountInvalid = errors.New("adjustment amount is invalid")
	ErrDisputeReasonRequired   = errors.New("dispute reason is required")
	ErrInvalidPayrollStatus    = errors.New("invalid payroll status transition")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreatePeriod(ctx context.Context, req CreatePeriodRequest) (PayrollPeriodResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollPeriodResponse{}, ErrTenantRequired
	}

	if req.Month < 1 || req.Month > 12 {
		return PayrollPeriodResponse{}, ErrMonthInvalid
	}

	if req.Year < 2020 || req.Year > 2100 {
		return PayrollPeriodResponse{}, ErrYearInvalid
	}

	period := PayrollPeriod{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		Month:          req.Month,
		Year:           req.Year,
		Status:         "draft",
	}

	createdPeriod, err := s.repo.CreatePeriod(ctx, period)
	if err != nil {
		return PayrollPeriodResponse{}, err
	}

	return buildPeriodResponse(createdPeriod), nil
}

func (s *Service) GenerateForPeriod(ctx context.Context, periodIDRaw string) (GeneratePayrollResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return GeneratePayrollResponse{}, ErrTenantRequired
	}

	periodID, err := uuid.Parse(strings.TrimSpace(periodIDRaw))
	if err != nil {
		return GeneratePayrollResponse{}, ErrPeriodIDInvalid
	}

	period, entries, err := s.repo.GenerateForPeriod(ctx, *currentUser.OrganizationID, periodID)
	if err != nil {
		return GeneratePayrollResponse{}, err
	}

	items := make([]PayrollEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, buildEntryResponse(entry))
	}

	return GeneratePayrollResponse{
		Period: buildPeriodResponse(period),
		Items:  items,
		Total:  len(items),
	}, nil
}

func (s *Service) ListEntries(ctx context.Context) (ListPayrollEntriesResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListPayrollEntriesResponse{}, ErrTenantRequired
	}

	entries, err := s.repo.ListEntries(ctx, *currentUser.OrganizationID)
	if err != nil {
		return ListPayrollEntriesResponse{}, err
	}

	items := make([]PayrollEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, buildEntryResponse(entry))
	}

	return ListPayrollEntriesResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) ListMyEntries(ctx context.Context) (ListPayrollEntriesResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ListPayrollEntriesResponse{}, ErrTenantRequired
	}

	entries, err := s.repo.ListMyEntries(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return ListPayrollEntriesResponse{}, err
	}

	items := make([]PayrollEntryResponse, 0, len(entries))
	for _, entry := range entries {
		items = append(items, buildEntryResponse(entry))
	}

	return ListPayrollEntriesResponse{
		Items: items,
		Total: len(items),
	}, nil
}

func (s *Service) CreateAdjustment(
	ctx context.Context,
	entryIDRaw string,
	req CreateAdjustmentRequest,
) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	adjustmentType := strings.ToLower(strings.TrimSpace(req.AdjustmentType))
	if !isValidAdjustmentType(adjustmentType) {
		return PayrollEntryResponse{}, ErrAdjustmentTypeInvalid
	}

	amount := strings.TrimSpace(req.Amount)
	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil || parsedAmount <= 0 {
		return PayrollEntryResponse{}, ErrAdjustmentAmountInvalid
	}

	adjustment := PayrollAdjustment{
		ID:             uuid.New(),
		OrganizationID: *currentUser.OrganizationID,
		AdjustmentType: adjustmentType,
		Amount:         amount,
		Reason:         strings.TrimSpace(req.Reason),
		CreatedBy:      currentUser.UserID.String(),
	}

	_, updatedEntry, err := s.repo.CreateAdjustment(ctx, *currentUser.OrganizationID, entryID, adjustment)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(updatedEntry), nil
}

func (s *Service) SendToTeacher(ctx context.Context, entryIDRaw string) (PayrollEntryResponse, error) {
	return s.updateEntryByFinance(ctx, entryIDRaw, s.repo.SendToTeacher)
}

func (s *Service) ConfirmByTeacher(ctx context.Context, entryIDRaw string) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	entry, err := s.repo.ConfirmByTeacher(ctx, *currentUser.OrganizationID, entryID, currentUser.UserID)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(entry), nil
}

func (s *Service) DisputeByTeacher(ctx context.Context, entryIDRaw string, req DisputePayrollRequest) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	reason := strings.TrimSpace(req.Reason)
	if reason == "" {
		return PayrollEntryResponse{}, ErrDisputeReasonRequired
	}

	entry, err := s.repo.DisputeByTeacher(ctx, *currentUser.OrganizationID, entryID, currentUser.UserID, reason)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(entry), nil
}

func (s *Service) ApproveByFinance(ctx context.Context, entryIDRaw string) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	entry, err := s.repo.ApproveByFinance(ctx, *currentUser.OrganizationID, entryID, currentUser.UserID)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(entry), nil
}

func (s *Service) MarkPaid(ctx context.Context, entryIDRaw string) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	entry, err := s.repo.MarkPaid(ctx, *currentUser.OrganizationID, entryID, currentUser.UserID)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(entry), nil
}

func (s *Service) updateEntryByFinance(
	ctx context.Context,
	entryIDRaw string,
	updateFn func(context.Context, uuid.UUID, uuid.UUID) (PayrollEntry, error),
) (PayrollEntryResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollEntryResponse{}, ErrTenantRequired
	}

	entryID, err := uuid.Parse(strings.TrimSpace(entryIDRaw))
	if err != nil {
		return PayrollEntryResponse{}, ErrEntryIDInvalid
	}

	entry, err := updateFn(ctx, *currentUser.OrganizationID, entryID)
	if err != nil {
		return PayrollEntryResponse{}, err
	}

	return buildEntryResponse(entry), nil
}

func buildPeriodResponse(period PayrollPeriod) PayrollPeriodResponse {
	return PayrollPeriodResponse{
		ID:             period.ID.String(),
		OrganizationID: period.OrganizationID.String(),
		Month:          period.Month,
		Year:           period.Year,
		Status:         period.Status,
	}
}

func buildEntryResponse(entry PayrollEntry) PayrollEntryResponse {
	return PayrollEntryResponse{
		ID:                        entry.ID.String(),
		OrganizationID:            entry.OrganizationID.String(),
		PeriodID:                  entry.PeriodID.String(),
		TeacherID:                 entry.TeacherID.String(),
		LessonsCount:              entry.LessonsCount,
		SubstitutionCount:         entry.SubstitutionCount,
		HoursWorked:               entry.HoursWorked,
		HourlyRate:                entry.HourlyRate,
		BaseAmount:                entry.BaseAmount,
		BonusAmount:               entry.BonusAmount,
		PenaltyAmount:             entry.PenaltyAmount,
		CorrectionAmount:          entry.CorrectionAmount,
		TotalAmount:               entry.TotalAmount,
		FinalAmount:               entry.FinalAmount,
		Status:                    entry.Status,
		TeacherConfirmationStatus: entry.TeacherConfirmationStatus,
		TeacherDisputeReason:      entry.TeacherDisputeReason,
		Comment:                   entry.Comment,
	}
}

func isValidAdjustmentType(adjustmentType string) bool {
	switch adjustmentType {
	case "bonus", "premium", "extra_work", "penalty", "deduction", "correction":
		return true
	default:
		return false
	}
}

var _ = time.Now
