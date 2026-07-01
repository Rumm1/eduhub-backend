package schedule

type CreateScheduleRequest struct {
	GroupID   string `json:"group_id"`
	Weekday   int    `json:"weekday"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Room      string `json:"room"`
}

type GenerateLessonsRequest struct {
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
	Topic    string `json:"topic"`
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

type GeneratedLessonResponse struct {
	ID         string `json:"id"`
	ScheduleID string `json:"schedule_id"`
	GroupID    string `json:"group_id"`
	LessonDate string `json:"lesson_date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Topic      string `json:"topic,omitempty"`
	Status     string `json:"status"`
}

type GenerateLessonsResponse struct {
	Created      []GeneratedLessonResponse `json:"created"`
	SkippedDates []string                  `json:"skipped_dates"`
	CreatedCount int                       `json:"created_count"`
	SkippedCount int                       `json:"skipped_count"`
}

type ListSchedulesResponse struct {
	Items []ScheduleResponse `json:"items"`
	Total int                `json:"total"`
}
