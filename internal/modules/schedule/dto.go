package schedule

type CreateScheduleRequest struct {
	GroupID   string `json:"group_id"`
	Weekday   int    `json:"weekday"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Room      string `json:"room"`
}

type ScheduleResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	BranchID       string `json:"branch_id"`
	GroupID        string `json:"group_id"`
	Weekday        int    `json:"weekday"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	Room           string `json:"room,omitempty"`
}

type ListSchedulesResponse struct {
	Items []ScheduleResponse `json:"items"`
	Total int                `json:"total"`
}
