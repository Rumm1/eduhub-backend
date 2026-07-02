package dashboard

type OverviewResponse struct {
	StudentsCount           int64                    `json:"students_count"`
	TeachersCount           int64                    `json:"teachers_count"`
	GroupsCount             int64                    `json:"groups_count"`
	LessonsToday            int64                    `json:"lessons_today"`
	PaymentsThisMonth       int64                    `json:"payments_this_month"`
	PaymentsAmountThisMonth float64                  `json:"payments_amount_this_month"`
	StudentDebtTotal        float64                  `json:"student_debt_total"`
	PendingPayrollEntries   int64                    `json:"pending_payroll_entries"`
	UnreadNotifications     int64                    `json:"unread_notifications"`
	RecentAuditLogs         []RecentAuditLogResponse `json:"recent_audit_logs"`
}

type RecentAuditLogResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id,omitempty"`
	Action      string `json:"action"`
	EntityType  string `json:"entity_type,omitempty"`
	EntityID    string `json:"entity_id,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}
