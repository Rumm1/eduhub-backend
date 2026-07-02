package notification

import "github.com/google/uuid"

const (
	NotificationTypeNormal    = "normal"
	NotificationTypeImportant = "important"
	NotificationTypeSystem    = "system"
	NotificationTypeWarning   = "warning"
	NotificationTypePayment   = "payment"
	NotificationTypeSchedule  = "schedule"
	NotificationTypeHomework  = "homework"
	NotificationTypePayroll   = "payroll"
	NotificationTypeLesson    = "lesson"
	NotificationTypeMessage   = "message"
)

type Notification struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Title          string
	Message        string
	Type           string
	IsRead         bool
	CreatedAt      string
}
