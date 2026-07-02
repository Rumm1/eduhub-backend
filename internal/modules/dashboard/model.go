package dashboard

import "github.com/google/uuid"

type Overview struct {
	StudentsCount           int64
	TeachersCount           int64
	GroupsCount             int64
	LessonsToday            int64
	PaymentsThisMonth       int64
	PaymentsAmountThisMonth float64
	StudentDebtTotal        float64
	PendingPayrollEntries   int64
	UnreadNotifications     int64
	RecentAuditLogs         []RecentAuditLog
}

type RecentAuditLog struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Action      string
	EntityType  string
	EntityID    uuid.UUID
	Description string
	CreatedAt   string
}
