package ai

type DashboardMetrics struct {
	StudentsCount           int64
	TeachersCount           int64
	GroupsCount             int64
	LessonsToday            int64
	PaymentsThisMonth       int64
	PaymentsAmountThisMonth float64
	StudentDebtTotal        float64
	PendingPayrollEntries   int64
	UnreadNotifications     int64
	RecentAuditLogsCount    int64
}

type Insight struct {
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Title    string `json:"title"`
	Message  string `json:"message"`
}
