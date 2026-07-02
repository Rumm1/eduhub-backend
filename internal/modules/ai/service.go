package ai

import (
	"context"
	"fmt"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
)

type Service struct {
	repository *Repository
	provider   ChatProvider
}

func NewService(repository *Repository) *Service {
	return &Service{
		repository: repository,
		provider:   NewRuleBasedProvider(),
	}
}

func (s *Service) GetDashboardInsights(ctx context.Context) (DashboardInsightsResponse, error) {
	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return DashboardInsightsResponse{}, ErrTenantRequired
	}

	metrics, err := s.repository.GetDashboardMetrics(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return DashboardInsightsResponse{}, err
	}

	insights := buildInsights(metrics)
	recommendations := buildRecommendations(metrics)
	riskLevel := calculateRiskLevel(metrics)

	return DashboardInsightsResponse{
		Summary:         buildSummary(metrics, riskLevel),
		RiskLevel:       riskLevel,
		Metrics:         mapMetrics(metrics),
		Insights:        insights,
		Recommendations: recommendations,
	}, nil
}

func (s *Service) Chat(ctx context.Context, request ChatRequest) (ChatResponse, error) {
	message := strings.TrimSpace(request.Message)
	if message == "" {
		return ChatResponse{}, ErrChatMessageEmpty
	}

	currentUser, ok := usercontext.GetUser(ctx)
	if !ok || currentUser.OrganizationID == nil {
		return ChatResponse{}, ErrTenantRequired
	}

	metrics, err := s.repository.GetDashboardMetrics(ctx, *currentUser.OrganizationID, currentUser.UserID)
	if err != nil {
		return ChatResponse{}, err
	}

	riskLevel := calculateRiskLevel(metrics)

	output, err := s.provider.GenerateChatReply(ctx, ChatInput{
		Message:   message,
		Metrics:   metrics,
		RiskLevel: riskLevel,
	})
	if err != nil {
		return ChatResponse{}, err
	}

	return ChatResponse{
		Provider:         s.provider.Name(),
		Intent:           output.Intent,
		RiskLevel:        riskLevel,
		Reply:            output.Reply,
		SuggestedActions: output.SuggestedActions,
		Metrics:          mapMetrics(metrics),
	}, nil
}

func buildSummary(metrics DashboardMetrics, riskLevel string) string {
	return fmt.Sprintf(
		"CRM overview: %d active students, %d active groups, %d teachers, %d lessons today. Current operational risk level is %s.",
		metrics.StudentsCount,
		metrics.GroupsCount,
		metrics.TeachersCount,
		metrics.LessonsToday,
		riskLevel,
	)
}

func buildInsights(metrics DashboardMetrics) []Insight {
	insights := make([]Insight, 0)

	if metrics.StudentsCount == 0 {
		insights = append(insights, Insight{
			Type:     "students",
			Severity: "high",
			Title:    "No active students",
			Message:  "There are no active students in the organization. The sales or import flow should be checked.",
		})
	}

	if metrics.GroupsCount == 0 {
		insights = append(insights, Insight{
			Type:     "groups",
			Severity: "high",
			Title:    "No active groups",
			Message:  "There are no active groups. Students may not be assigned to learning groups yet.",
		})
	}

	if metrics.TeachersCount == 0 {
		insights = append(insights, Insight{
			Type:     "teachers",
			Severity: "high",
			Title:    "No teachers",
			Message:  "No teacher profiles were found. Lessons and schedules may not work correctly.",
		})
	}

	if metrics.StudentDebtTotal > 0 {
		insights = append(insights, Insight{
			Type:     "payments",
			Severity: debtSeverity(metrics.StudentDebtTotal),
			Title:    "Student debt detected",
			Message:  fmt.Sprintf("Current estimated student debt is %.2f. Finance or managers should review unpaid balances.", metrics.StudentDebtTotal),
		})
	}

	if metrics.PaymentsThisMonth == 0 && metrics.StudentsCount > 0 {
		insights = append(insights, Insight{
			Type:     "payments",
			Severity: "medium",
			Title:    "No payments this month",
			Message:  "There are active students, but no payments were recorded for the current month.",
		})
	}

	if metrics.PendingPayrollEntries > 0 {
		insights = append(insights, Insight{
			Type:     "payroll",
			Severity: "medium",
			Title:    "Pending payroll entries",
			Message:  fmt.Sprintf("There are %d payroll entries that are not paid yet.", metrics.PendingPayrollEntries),
		})
	}

	if metrics.UnreadNotifications > 0 {
		insights = append(insights, Insight{
			Type:     "notifications",
			Severity: "low",
			Title:    "Unread notifications",
			Message:  fmt.Sprintf("You have %d unread notifications.", metrics.UnreadNotifications),
		})
	}

	if metrics.LessonsToday == 0 && metrics.GroupsCount > 0 {
		insights = append(insights, Insight{
			Type:     "schedule",
			Severity: "low",
			Title:    "No lessons today",
			Message:  "There are active groups, but no lessons are scheduled for today.",
		})
	}

	if len(insights) == 0 {
		insights = append(insights, Insight{
			Type:     "general",
			Severity: "low",
			Title:    "No major issues detected",
			Message:  "The main CRM metrics look stable at the moment.",
		})
	}

	return insights
}

func buildRecommendations(metrics DashboardMetrics) []string {
	recommendations := make([]string, 0)

	if metrics.StudentDebtTotal > 0 {
		recommendations = append(recommendations, "Review student balances and contact parents with unpaid invoices.")
	}

	if metrics.PendingPayrollEntries > 0 {
		recommendations = append(recommendations, "Check payroll entries and finish the approval or payment workflow.")
	}

	if metrics.PaymentsThisMonth == 0 && metrics.StudentsCount > 0 {
		recommendations = append(recommendations, "Check whether payments for the current month were recorded correctly.")
	}

	if metrics.LessonsToday == 0 && metrics.GroupsCount > 0 {
		recommendations = append(recommendations, "Review today's schedule and generate lessons if schedules exist.")
	}

	if metrics.StudentsCount == 0 {
		recommendations = append(recommendations, "Import students from Excel/CSV or create students manually.")
	}

	if metrics.GroupsCount == 0 {
		recommendations = append(recommendations, "Create groups and assign students to them.")
	}

	if metrics.TeachersCount == 0 {
		recommendations = append(recommendations, "Create teacher profiles before generating schedules.")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Continue monitoring payments, attendance, lessons, and payroll.")
	}

	return recommendations
}

func calculateRiskLevel(metrics DashboardMetrics) string {
	score := 0

	if metrics.StudentsCount == 0 {
		score += 3
	}

	if metrics.GroupsCount == 0 {
		score += 3
	}

	if metrics.TeachersCount == 0 {
		score += 3
	}

	if metrics.StudentDebtTotal > 0 {
		score += 2
	}

	if metrics.StudentDebtTotal >= 50000 {
		score += 2
	}

	if metrics.PendingPayrollEntries > 0 {
		score += 1
	}

	if metrics.PaymentsThisMonth == 0 && metrics.StudentsCount > 0 {
		score += 2
	}

	switch {
	case score >= 6:
		return "high"
	case score >= 3:
		return "medium"
	default:
		return "low"
	}
}

func debtSeverity(debt float64) string {
	switch {
	case debt >= 100000:
		return "high"
	case debt >= 30000:
		return "medium"
	default:
		return "low"
	}
}

func mapMetrics(metrics DashboardMetrics) DashboardMetricsResponse {
	return DashboardMetricsResponse{
		StudentsCount:           metrics.StudentsCount,
		TeachersCount:           metrics.TeachersCount,
		GroupsCount:             metrics.GroupsCount,
		LessonsToday:            metrics.LessonsToday,
		PaymentsThisMonth:       metrics.PaymentsThisMonth,
		PaymentsAmountThisMonth: metrics.PaymentsAmountThisMonth,
		StudentDebtTotal:        metrics.StudentDebtTotal,
		PendingPayrollEntries:   metrics.PendingPayrollEntries,
		UnreadNotifications:     metrics.UnreadNotifications,
		RecentAuditLogsCount:    metrics.RecentAuditLogsCount,
	}
}
