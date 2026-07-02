package ai

import (
	"context"
	"fmt"
	"strings"
)

type ChatProvider interface {
	Name() string
	GenerateChatReply(ctx context.Context, input ChatInput) (ChatOutput, error)
}

type ChatInput struct {
	Message   string
	Metrics   DashboardMetrics
	RiskLevel string
}

type ChatOutput struct {
	Reply            string
	Intent           string
	SuggestedActions []string
}

type RuleBasedProvider struct{}

func NewRuleBasedProvider() *RuleBasedProvider {
	return &RuleBasedProvider{}
}

func (p *RuleBasedProvider) Name() string {
	return "rule_based"
}

func (p *RuleBasedProvider) GenerateChatReply(ctx context.Context, input ChatInput) (ChatOutput, error) {
	_ = ctx

	message := strings.ToLower(strings.TrimSpace(input.Message))

	switch {
	case containsAny(message, "долг", "задолж", "balance", "debt", "payment", "оплат", "платеж"):
		return ChatOutput{
			Intent: "payments",
			Reply: fmt.Sprintf(
				"Current estimated student debt is %.2f. Payments this month: %d, total paid amount: %.2f. Finance or managers should review unpaid balances.",
				input.Metrics.StudentDebtTotal,
				input.Metrics.PaymentsThisMonth,
				input.Metrics.PaymentsAmountThisMonth,
			),
			SuggestedActions: []string{
				"Open student balances report",
				"Check unpaid students",
				"Export payments report",
			},
		}, nil

	case containsAny(message, "зарплат", "salary", "payroll", "teacher payment", "учител", "преподав"):
		return ChatOutput{
			Intent: "payroll",
			Reply: fmt.Sprintf(
				"There are %d payroll entries that are not paid yet. Payroll should be reviewed before closing the period.",
				input.Metrics.PendingPayrollEntries,
			),
			SuggestedActions: []string{
				"Open payroll entries",
				"Check disputed payroll entries",
				"Approve or mark paid completed entries",
			},
		}, nil

	case containsAny(message, "урок", "lesson", "schedule", "распис", "занят"):
		return ChatOutput{
			Intent: "schedule",
			Reply: fmt.Sprintf(
				"Today there are %d lessons. Active groups: %d. If lessons are missing, check schedules and generate lessons.",
				input.Metrics.LessonsToday,
				input.Metrics.GroupsCount,
			),
			SuggestedActions: []string{
				"Open lessons page",
				"Check schedules",
				"Generate lessons from schedules",
			},
		}, nil

	case containsAny(message, "студент", "student", "ученик", "групп", "group"):
		return ChatOutput{
			Intent: "students",
			Reply: fmt.Sprintf(
				"The organization has %d active students and %d active groups. If students are missing, use manual creation or CSV/XLSX import.",
				input.Metrics.StudentsCount,
				input.Metrics.GroupsCount,
			),
			SuggestedActions: []string{
				"Open students page",
				"Open groups page",
				"Import students from CSV/XLSX",
			},
		}, nil

	case containsAny(message, "уведом", "notification", "notif"):
		return ChatOutput{
			Intent: "notifications",
			Reply: fmt.Sprintf(
				"You have %d unread notifications.",
				input.Metrics.UnreadNotifications,
			),
			SuggestedActions: []string{
				"Open notifications",
				"Mark notifications as read",
			},
		}, nil

	case containsAny(message, "отчет", "report", "excel", "xlsx", "pdf", "docx"):
		return ChatOutput{
			Intent: "reports",
			Reply:  "Reports are available in JSON, XLSX, PDF and DOCX formats. You can export payments, payroll, teacher schedule and student balances.",
			SuggestedActions: []string{
				"Export payments report",
				"Export payroll report",
				"Export student balances report",
			},
		}, nil

	case containsAny(message, "что делать", "совет", "рекоменд", "recommend", "next", "след"):
		return ChatOutput{
			Intent:           "recommendations",
			Reply:            buildRecommendationReply(input.Metrics),
			SuggestedActions: buildSuggestedActions(input.Metrics),
		}, nil

	default:
		return ChatOutput{
			Intent: "overview",
			Reply: fmt.Sprintf(
				"CRM overview: %d active students, %d active groups, %d teachers, %d lessons today. Current risk level is %s.",
				input.Metrics.StudentsCount,
				input.Metrics.GroupsCount,
				input.Metrics.TeachersCount,
				input.Metrics.LessonsToday,
				input.RiskLevel,
			),
			SuggestedActions: []string{
				"Ask about student debts",
				"Ask about payroll",
				"Ask about today's lessons",
				"Ask for recommendations",
			},
		}, nil
	}
}

func containsAny(value string, keywords ...string) bool {
	for _, keyword := range keywords {
		if strings.Contains(value, keyword) {
			return true
		}
	}

	return false
}

func buildRecommendationReply(metrics DashboardMetrics) string {
	recommendations := buildRecommendations(metrics)

	if len(recommendations) == 0 {
		return "No urgent recommendations. Continue monitoring payments, attendance, lessons and payroll."
	}

	return "Recommended next steps: " + strings.Join(recommendations, " ")
}

func buildSuggestedActions(metrics DashboardMetrics) []string {
	actions := make([]string, 0)

	if metrics.StudentDebtTotal > 0 {
		actions = append(actions, "Open student balances report")
	}

	if metrics.PendingPayrollEntries > 0 {
		actions = append(actions, "Open payroll entries")
	}

	if metrics.LessonsToday == 0 && metrics.GroupsCount > 0 {
		actions = append(actions, "Check schedules")
	}

	if metrics.StudentsCount == 0 {
		actions = append(actions, "Import students")
	}

	if len(actions) == 0 {
		actions = append(actions, "Open dashboard")
	}

	return actions
}
