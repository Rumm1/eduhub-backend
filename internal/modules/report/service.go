package report

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/google/uuid"
)

var (
	ErrTenantRequired                   = errors.New("tenant organization is required")
	ErrTeacherIDRequired                = errors.New("teacher id is required")
	ErrTeacherIDInvalid                 = errors.New("teacher id is invalid")
	ErrTeacherNotFound                  = errors.New("teacher not found in organization")
	ErrFromDateRequired                 = errors.New("from date is required")
	ErrToDateRequired                   = errors.New("to date is required")
	ErrFromDateInvalid                  = errors.New("from date is invalid")
	ErrToDateInvalid                    = errors.New("to date is invalid")
	ErrDateRangeInvalid                 = errors.New("to date must be after or equal from date")
	ErrForbiddenReport                  = errors.New("forbidden report access")
	ErrBranchIDInvalid                  = errors.New("branch id is invalid")
	ErrGroupIDInvalid                   = errors.New("group id is invalid")
	ErrStudentIDInvalid                 = errors.New("student id is invalid")
	ErrPaymentStatusInvalid             = errors.New("payment status is invalid")
	ErrPeriodRequired                   = errors.New("period is required")
	ErrPeriodInvalid                    = errors.New("period is invalid")
	ErrBalanceStatusInvalid             = errors.New("balance status is invalid")
	ErrPayrollStatusInvalid             = errors.New("payroll status is invalid")
	ErrTeacherConfirmationStatusInvalid = errors.New("teacher confirmation status is invalid")
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetTeacherSchedule(
	ctx context.Context,
	teacherIDRaw string,
	fromDateRaw string,
	toDateRaw string,
) (TeacherScheduleReportResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return TeacherScheduleReportResponse{}, ErrTenantRequired
	}

	fromDate, toDate, err := validateDateRange(fromDateRaw, toDateRaw)
	if err != nil {
		return TeacherScheduleReportResponse{}, err
	}

	teacherIDText := strings.TrimSpace(teacherIDRaw)

	if teacherIDText == "" {
		if hasRole(currentUser.Roles, "TEACHER") {
			teacherIDText = currentUser.UserID.String()
		} else {
			return TeacherScheduleReportResponse{}, ErrTeacherIDRequired
		}
	}

	teacherID, err := uuid.Parse(teacherIDText)
	if err != nil {
		return TeacherScheduleReportResponse{}, ErrTeacherIDInvalid
	}

	if hasRole(currentUser.Roles, "TEACHER") && teacherID != currentUser.UserID && !hasRole(currentUser.Roles, "ORG_ADMIN") {
		return TeacherScheduleReportResponse{}, ErrForbiddenReport
	}

	report, err := s.repo.GetTeacherSchedule(ctx, *currentUser.OrganizationID, teacherID, fromDate, toDate)
	if err != nil {
		return TeacherScheduleReportResponse{}, err
	}

	return buildTeacherScheduleResponse(report), nil
}

func (s *Service) GetPaymentsReport(
	ctx context.Context,
	fromDateRaw string,
	toDateRaw string,
	branchIDRaw string,
	groupIDRaw string,
	studentIDRaw string,
	statusRaw string,
) (PaymentsReportResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PaymentsReportResponse{}, ErrTenantRequired
	}

	fromDate, toDate, err := validateDateRange(fromDateRaw, toDateRaw)
	if err != nil {
		return PaymentsReportResponse{}, err
	}

	branchID := strings.TrimSpace(branchIDRaw)
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return PaymentsReportResponse{}, ErrBranchIDInvalid
		}
	}

	groupID := strings.TrimSpace(groupIDRaw)
	if groupID != "" {
		if _, err := uuid.Parse(groupID); err != nil {
			return PaymentsReportResponse{}, ErrGroupIDInvalid
		}
	}

	studentID := strings.TrimSpace(studentIDRaw)
	if studentID != "" {
		if _, err := uuid.Parse(studentID); err != nil {
			return PaymentsReportResponse{}, ErrStudentIDInvalid
		}
	}

	status := strings.ToLower(strings.TrimSpace(statusRaw))
	if status != "" && !isValidPaymentStatus(status) {
		return PaymentsReportResponse{}, ErrPaymentStatusInvalid
	}

	report, err := s.repo.GetPaymentsReport(
		ctx,
		*currentUser.OrganizationID,
		fromDate,
		toDate,
		branchID,
		groupID,
		studentID,
		status,
	)
	if err != nil {
		return PaymentsReportResponse{}, err
	}

	return buildPaymentsReportResponse(report), nil
}

func (s *Service) GetStudentBalancesReport(
	ctx context.Context,
	periodRaw string,
	branchIDRaw string,
	groupIDRaw string,
	studentIDRaw string,
	statusRaw string,
) (StudentBalancesReportResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return StudentBalancesReportResponse{}, ErrTenantRequired
	}

	period := strings.TrimSpace(periodRaw)
	if period == "" {
		return StudentBalancesReportResponse{}, ErrPeriodRequired
	}

	periodDate, err := time.Parse("2006-01-02", period+"-01")
	if err != nil {
		return StudentBalancesReportResponse{}, ErrPeriodInvalid
	}

	periodStart := periodDate.Format("2006-01-02")

	branchID := strings.TrimSpace(branchIDRaw)
	if branchID != "" {
		if _, err := uuid.Parse(branchID); err != nil {
			return StudentBalancesReportResponse{}, ErrBranchIDInvalid
		}
	}

	groupID := strings.TrimSpace(groupIDRaw)
	if groupID != "" {
		if _, err := uuid.Parse(groupID); err != nil {
			return StudentBalancesReportResponse{}, ErrGroupIDInvalid
		}
	}

	studentID := strings.TrimSpace(studentIDRaw)
	if studentID != "" {
		if _, err := uuid.Parse(studentID); err != nil {
			return StudentBalancesReportResponse{}, ErrStudentIDInvalid
		}
	}

	status := strings.ToLower(strings.TrimSpace(statusRaw))
	if status != "" && !isValidBalanceStatus(status) {
		return StudentBalancesReportResponse{}, ErrBalanceStatusInvalid
	}

	report, err := s.repo.GetStudentBalancesReport(
		ctx,
		*currentUser.OrganizationID,
		periodStart,
		branchID,
		groupID,
		studentID,
	)
	if err != nil {
		return StudentBalancesReportResponse{}, err
	}

	return buildStudentBalancesReportResponse(report, status), nil
}

func (s *Service) GetPayrollReport(
	ctx context.Context,
	periodRaw string,
	teacherIDRaw string,
	statusRaw string,
	teacherConfirmationStatusRaw string,
) (PayrollReportResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return PayrollReportResponse{}, ErrTenantRequired
	}

	year, month, period, err := parseReportPeriod(periodRaw)
	if err != nil {
		return PayrollReportResponse{}, err
	}

	teacherID := strings.TrimSpace(teacherIDRaw)
	if teacherID != "" {
		if _, err := uuid.Parse(teacherID); err != nil {
			return PayrollReportResponse{}, ErrTeacherIDInvalid
		}
	}

	status := strings.ToLower(strings.TrimSpace(statusRaw))
	if status != "" && !isValidPayrollReportStatus(status) {
		return PayrollReportResponse{}, ErrPayrollStatusInvalid
	}

	teacherConfirmationStatus := strings.ToLower(strings.TrimSpace(teacherConfirmationStatusRaw))
	if teacherConfirmationStatus != "" && !isValidTeacherConfirmationStatus(teacherConfirmationStatus) {
		return PayrollReportResponse{}, ErrTeacherConfirmationStatusInvalid
	}

	report, err := s.repo.GetPayrollReport(
		ctx,
		*currentUser.OrganizationID,
		year,
		month,
		teacherID,
		status,
		teacherConfirmationStatus,
	)
	if err != nil {
		return PayrollReportResponse{}, err
	}

	report.Period = period

	return buildPayrollReportResponse(report), nil
}

func validateDateRange(fromDateRaw string, toDateRaw string) (string, string, error) {
	fromDate := strings.TrimSpace(fromDateRaw)
	if fromDate == "" {
		return "", "", ErrFromDateRequired
	}

	parsedFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return "", "", ErrFromDateInvalid
	}

	toDate := strings.TrimSpace(toDateRaw)
	if toDate == "" {
		return "", "", ErrToDateRequired
	}

	parsedToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return "", "", ErrToDateInvalid
	}

	if parsedToDate.Before(parsedFromDate) {
		return "", "", ErrDateRangeInvalid
	}

	return fromDate, toDate, nil
}

func buildTeacherScheduleResponse(report TeacherScheduleReport) TeacherScheduleReportResponse {
	items := make([]TeacherScheduleReportItem, 0, len(report.Items))

	for _, item := range report.Items {
		items = append(items, TeacherScheduleReportItem{
			LessonID:            item.LessonID,
			LessonDate:          item.LessonDate,
			StartTime:           item.StartTime,
			EndTime:             item.EndTime,
			Hours:               item.Hours,
			Topic:               item.Topic,
			Status:              item.Status,
			GroupID:             item.GroupID,
			GroupName:           item.GroupName,
			BranchID:            item.BranchID,
			BranchName:          item.BranchName,
			SubjectID:           item.SubjectID,
			SubjectName:         item.SubjectName,
			PlannedTeacherID:    item.PlannedTeacherID,
			PlannedTeacherName:  item.PlannedTeacherName,
			ActualTeacherID:     item.ActualTeacherID,
			ActualTeacherName:   item.ActualTeacherName,
			IsSubstitution:      item.IsSubstitution,
			TeacherRoleInLesson: item.TeacherRoleInLesson,
			SubstitutionReason:  item.SubstitutionReason,
		})
	}

	return TeacherScheduleReportResponse{
		TeacherID:          report.TeacherID,
		TeacherName:        report.TeacherName,
		FromDate:           report.FromDate,
		ToDate:             report.ToDate,
		TotalLessons:       report.TotalLessons,
		ActualLessons:      report.ActualLessons,
		PlannedOnlyLessons: report.PlannedOnlyLessons,
		Substitutions:      report.Substitutions,
		TotalActualHours:   report.TotalActualHours,
		Items:              items,
	}
}

func buildPaymentsReportResponse(report PaymentsReport) PaymentsReportResponse {
	items := make([]PaymentsReportItem, 0, len(report.Items))

	for _, item := range report.Items {
		items = append(items, PaymentsReportItem{
			PaymentID:     item.PaymentID,
			PaymentDate:   item.PaymentDate,
			PaymentPeriod: item.PaymentPeriod,
			StudentID:     item.StudentID,
			StudentName:   item.StudentName,
			GroupID:       item.GroupID,
			GroupName:     item.GroupName,
			BranchID:      item.BranchID,
			BranchName:    item.BranchName,
			Amount:        item.Amount,
			PaymentMethod: item.PaymentMethod,
			Status:        item.Status,
			Comment:       item.Comment,
		})
	}

	return PaymentsReportResponse{
		FromDate:        report.FromDate,
		ToDate:          report.ToDate,
		TotalPayments:   report.TotalPayments,
		TotalAmount:     report.TotalAmount,
		PaidAmount:      report.PaidAmount,
		PendingAmount:   report.PendingAmount,
		RefundedAmount:  report.RefundedAmount,
		CancelledAmount: report.CancelledAmount,
		Items:           items,
	}
}

func buildStudentBalancesReportResponse(report StudentBalancesReport, statusFilter string) StudentBalancesReportResponse {
	items := make([]StudentBalanceReportItem, 0, len(report.Items))

	totalExpected := 0.0
	totalPaid := 0.0
	totalDebt := 0.0
	paidCount := 0
	partialCount := 0
	unpaidCount := 0

	for _, item := range report.Items {
		if statusFilter != "" && item.PaymentStatus != statusFilter {
			continue
		}

		totalExpected += moneyToFloat(item.MonthlyPrice)
		totalPaid += moneyToFloat(item.PaidAmount)
		totalDebt += moneyToFloat(item.DebtAmount)

		switch item.PaymentStatus {
		case "paid":
			paidCount++
		case "partial":
			partialCount++
		case "unpaid":
			unpaidCount++
		}

		items = append(items, StudentBalanceReportItem{
			StudentID:     item.StudentID,
			StudentName:   item.StudentName,
			GroupID:       item.GroupID,
			GroupName:     item.GroupName,
			BranchID:      item.BranchID,
			BranchName:    item.BranchName,
			MonthlyPrice:  item.MonthlyPrice,
			PaidAmount:    item.PaidAmount,
			DebtAmount:    item.DebtAmount,
			PaymentStatus: item.PaymentStatus,
		})
	}

	return StudentBalancesReportResponse{
		Period:              report.Period,
		TotalStudents:       len(items),
		PaidCount:           paidCount,
		PartialCount:        partialCount,
		UnpaidCount:         unpaidCount,
		TotalExpectedAmount: formatMoney(totalExpected),
		TotalPaidAmount:     formatMoney(totalPaid),
		TotalDebtAmount:     formatMoney(totalDebt),
		Items:               items,
	}
}

func buildPayrollReportResponse(report PayrollReport) PayrollReportResponse {
	items := make([]PayrollReportItem, 0, len(report.Items))

	totalLessons := 0
	totalSubstitutions := 0
	totalHours := 0.0
	totalBaseAmount := 0.0
	totalBonusAmount := 0.0
	totalPenaltyAmount := 0.0
	totalCorrectionAmount := 0.0
	totalFinalAmount := 0.0

	draftCount := 0
	sentToTeacherCount := 0
	teacherApprovedCount := 0
	teacherDisputedCount := 0
	approvedByFinanceCount := 0
	paidCount := 0

	for _, item := range report.Items {
		totalLessons += item.LessonsCount
		totalSubstitutions += item.SubstitutionCount
		totalHours += moneyToFloat(item.HoursWorked)
		totalBaseAmount += moneyToFloat(item.BaseAmount)
		totalBonusAmount += moneyToFloat(item.BonusAmount)
		totalPenaltyAmount += moneyToFloat(item.PenaltyAmount)
		totalCorrectionAmount += moneyToFloat(item.CorrectionAmount)
		totalFinalAmount += moneyToFloat(item.FinalAmount)

		switch item.Status {
		case "draft":
			draftCount++
		case "sent_to_teacher":
			sentToTeacherCount++
		case "teacher_approved":
			teacherApprovedCount++
		case "teacher_disputed":
			teacherDisputedCount++
		case "approved_by_finance":
			approvedByFinanceCount++
		case "paid":
			paidCount++
		}

		items = append(items, PayrollReportItem{
			EntryID:                   item.EntryID,
			PeriodID:                  item.PeriodID,
			TeacherID:                 item.TeacherID,
			TeacherName:               item.TeacherName,
			LessonsCount:              item.LessonsCount,
			SubstitutionCount:         item.SubstitutionCount,
			HoursWorked:               item.HoursWorked,
			HourlyRate:                item.HourlyRate,
			BaseAmount:                item.BaseAmount,
			BonusAmount:               item.BonusAmount,
			PenaltyAmount:             item.PenaltyAmount,
			CorrectionAmount:          item.CorrectionAmount,
			FinalAmount:               item.FinalAmount,
			Status:                    item.Status,
			TeacherConfirmationStatus: item.TeacherConfirmationStatus,
			TeacherDisputeReason:      item.TeacherDisputeReason,
			Comment:                   item.Comment,
		})
	}

	return PayrollReportResponse{
		Period:                 report.Period,
		TotalEntries:           len(items),
		TotalLessons:           totalLessons,
		TotalSubstitutions:     totalSubstitutions,
		TotalHours:             formatMoney(totalHours),
		TotalBaseAmount:        formatMoney(totalBaseAmount),
		TotalBonusAmount:       formatMoney(totalBonusAmount),
		TotalPenaltyAmount:     formatMoney(totalPenaltyAmount),
		TotalCorrectionAmount:  formatMoney(totalCorrectionAmount),
		TotalFinalAmount:       formatMoney(totalFinalAmount),
		DraftCount:             draftCount,
		SentToTeacherCount:     sentToTeacherCount,
		TeacherApprovedCount:   teacherApprovedCount,
		TeacherDisputedCount:   teacherDisputedCount,
		ApprovedByFinanceCount: approvedByFinanceCount,
		PaidCount:              paidCount,
		Items:                  items,
	}
}

func hasRole(roles []string, role string) bool {
	for _, item := range roles {
		if item == role {
			return true
		}
	}

	return false
}

func isValidPaymentStatus(status string) bool {
	switch status {
	case "paid", "pending", "cancelled", "refunded":
		return true
	default:
		return false
	}
}

func isValidBalanceStatus(status string) bool {
	switch status {
	case "paid", "partial", "unpaid":
		return true
	default:
		return false
	}
}

func moneyToFloat(value string) float64 {
	normalized := strings.ReplaceAll(strings.TrimSpace(value), ",", "")

	result, err := strconv.ParseFloat(normalized, 64)
	if err != nil {
		return 0
	}

	return result
}

func formatMoney(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

func parseReportPeriod(periodRaw string) (int, int, string, error) {
	period := strings.TrimSpace(periodRaw)
	if period == "" {
		return 0, 0, "", ErrPeriodRequired
	}

	periodDate, err := time.Parse("2006-01-02", period+"-01")
	if err != nil {
		return 0, 0, "", ErrPeriodInvalid
	}

	year := periodDate.Year()
	month := int(periodDate.Month())

	return year, month, periodDate.Format("2006-01"), nil
}

func isValidPayrollReportStatus(status string) bool {
	switch status {
	case "draft", "sent_to_teacher", "teacher_approved", "teacher_disputed", "approved_by_finance", "paid":
		return true
	default:
		return false
	}
}

func isValidTeacherConfirmationStatus(status string) bool {
	switch status {
	case "not_sent", "pending", "approved", "disputed":
		return true
	default:
		return false
	}
}
