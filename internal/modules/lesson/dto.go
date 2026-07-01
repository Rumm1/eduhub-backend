package lesson

type CreateLessonRequest struct {
	GroupID    string `json:"group_id"`
	LessonDate string `json:"lesson_date"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Topic      string `json:"topic"`
}

type UpdateLessonTeacherRequest struct {
	ActualTeacherID    string `json:"actual_teacher_id"`
	SubstitutionReason string `json:"substitution_reason"`
}

type LessonResponse struct {
	ID                 string `json:"id"`
	OrganizationID     string `json:"organization_id"`
	BranchID           string `json:"branch_id"`
	GroupID            string `json:"group_id"`
	TeacherID          string `json:"teacher_id,omitempty"`
	PlannedTeacherID   string `json:"planned_teacher_id,omitempty"`
	ActualTeacherID    string `json:"actual_teacher_id,omitempty"`
	SubjectID          string `json:"subject_id"`
	ScheduleID         string `json:"schedule_id,omitempty"`
	LessonDate         string `json:"lesson_date"`
	StartTime          string `json:"start_time"`
	EndTime            string `json:"end_time"`
	Topic              string `json:"topic,omitempty"`
	Status             string `json:"status"`
	SubstitutionReason string `json:"substitution_reason,omitempty"`
}

type ListLessonsResponse struct {
	Items []LessonResponse `json:"items"`
	Total int              `json:"total"`
}
