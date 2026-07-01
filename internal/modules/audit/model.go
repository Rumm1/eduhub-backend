package audit

import "time"

type AuditLog struct {
	ID             string
	OrganizationID string
	UserID         string
	UserName       string
	Action         string
	EntityType     string
	EntityID       string
	Description    string
	Metadata       string
	IPAddress      string
	UserAgent      string
	CreatedAt      time.Time
}

type CreateAuditLogInput struct {
	OrganizationID string
	UserID         string
	Action         string
	EntityType     string
	EntityID       string
	Description    string
	Metadata       map[string]interface{}
	IPAddress      string
	UserAgent      string
}

type AuditLogFilter struct {
	OrganizationID string
	UserID         string
	Action         string
	EntityType     string
	EntityID       string
	FromDate       string
	ToDate         string
	Limit          int
	Offset         int
}
