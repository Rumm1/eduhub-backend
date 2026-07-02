package ai

type DashboardMetricsResponse struct {
	StudentsCount           int64   `json:"students_count"`
	TeachersCount           int64   `json:"teachers_count"`
	GroupsCount             int64   `json:"groups_count"`
	LessonsToday            int64   `json:"lessons_today"`
	PaymentsThisMonth       int64   `json:"payments_this_month"`
	PaymentsAmountThisMonth float64 `json:"payments_amount_this_month"`
	StudentDebtTotal        float64 `json:"student_debt_total"`
	PendingPayrollEntries   int64   `json:"pending_payroll_entries"`
	UnreadNotifications     int64   `json:"unread_notifications"`
	RecentAuditLogsCount    int64   `json:"recent_audit_logs_count"`
}

type DashboardInsightsResponse struct {
	Summary         string                   `json:"summary"`
	RiskLevel       string                   `json:"risk_level"`
	Metrics         DashboardMetricsResponse `json:"metrics"`
	Insights        []Insight                `json:"insights"`
	Recommendations []string                 `json:"recommendations"`
}

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Provider         string                   `json:"provider"`
	Intent           string                   `json:"intent"`
	RiskLevel        string                   `json:"risk_level"`
	Reply            string                   `json:"reply"`
	SuggestedActions []string                 `json:"suggested_actions"`
	Metrics          DashboardMetricsResponse `json:"metrics"`
}
