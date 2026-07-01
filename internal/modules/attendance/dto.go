package attendance

type MarkLessonAttendanceRequest struct {
	Items []MarkAttendanceItemRequest `json:"items"`
}

type MarkAttendanceItemRequest struct {
	StudentID string `json:"student_id"`
	Status    string `json:"status"`
	Reason    string `json:"reason"`
	Comment   string `json:"comment"`
}

type AttendanceResponse struct {
	ID              string `json:"id,omitempty"`
	LessonID        string `json:"lesson_id"`
	StudentID       string `json:"student_id"`
	StudentFullName string `json:"student_full_name"`
	Status          string `json:"status"`
	Reason          string `json:"reason,omitempty"`
	Comment         string `json:"comment,omitempty"`
	MarkedBy        string `json:"marked_by,omitempty"`
	MarkedAt        string `json:"marked_at,omitempty"`
}

type ListLessonAttendanceResponse struct {
	Items []AttendanceResponse `json:"items"`
	Total int                  `json:"total"`
}
