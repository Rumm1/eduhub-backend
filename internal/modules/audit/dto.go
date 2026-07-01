package audit

type AuditLogResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id,omitempty"`
	UserID         string `json:"user_id,omitempty"`
	UserName       string `json:"user_name,omitempty"`
	Action         string `json:"action"`
	EntityType     string `json:"entity_type"`
	EntityID       string `json:"entity_id,omitempty"`
	Description    string `json:"description,omitempty"`
	Metadata       string `json:"metadata"`
	IPAddress      string `json:"ip_address,omitempty"`
	UserAgent      string `json:"user_agent,omitempty"`
	CreatedAt      string `json:"created_at"`
}

type AuditLogsListResponse struct {
	Items  []AuditLogResponse `json:"items"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
	Total  int                `json:"total"`
}
