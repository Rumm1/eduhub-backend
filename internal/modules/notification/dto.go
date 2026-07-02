package notification

type CreateNotificationRequest struct {
	UserID  string   `json:"user_id"`
	UserIDs []string `json:"user_ids"`
	Title   string   `json:"title"`
	Message string   `json:"message"`
	Type    string   `json:"type"`
}

type NotificationResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id,omitempty"`
	UserID         string `json:"user_id"`
	Title          string `json:"title"`
	Message        string `json:"message,omitempty"`
	Type           string `json:"type,omitempty"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
}

type NotificationTypeResponse struct {
	Code        string `json:"code"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ListNotificationsResponse struct {
	Items  []NotificationResponse `json:"items"`
	Total  int                    `json:"total"`
	Unread int                    `json:"unread"`
}

type ListNotificationTypesResponse struct {
	Items []NotificationTypeResponse `json:"items"`
	Total int                        `json:"total"`
}

type CreateNotificationsResponse struct {
	Items []NotificationResponse `json:"items"`
	Total int                    `json:"total"`
}

type MarkAllReadResponse struct {
	Updated int64 `json:"updated"`
}
