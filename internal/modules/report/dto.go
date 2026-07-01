package report

type TeacherScheduleReportResponse struct {
	TeacherID          string                      `json:"teacher_id"`
	TeacherName        string                      `json:"teacher_name"`
	FromDate           string                      `json:"from_date"`
	ToDate             string                      `json:"to_date"`
	TotalLessons       int                         `json:"total_lessons"`
	ActualLessons      int                         `json:"actual_lessons"`
	PlannedOnlyLessons int                         `json:"planned_only_lessons"`
	Substitutions      int                         `json:"substitutions"`
	TotalActualHours   string                      `json:"total_actual_hours"`
	Items              []TeacherScheduleReportItem `json:"items"`
}

type TeacherScheduleReportItem struct {
	LessonID            string `json:"lesson_id"`
	LessonDate          string `json:"lesson_date"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Hours               string `json:"hours"`
	Topic               string `json:"topic,omitempty"`
	Status              string `json:"status"`
	GroupID             string `json:"group_id"`
	GroupName           string `json:"group_name"`
	BranchID            string `json:"branch_id"`
	BranchName          string `json:"branch_name"`
	SubjectID           string `json:"subject_id"`
	SubjectName         string `json:"subject_name"`
	PlannedTeacherID    string `json:"planned_teacher_id,omitempty"`
	PlannedTeacherName  string `json:"planned_teacher_name,omitempty"`
	ActualTeacherID     string `json:"actual_teacher_id,omitempty"`
	ActualTeacherName   string `json:"actual_teacher_name,omitempty"`
	IsSubstitution      bool   `json:"is_substitution"`
	TeacherRoleInLesson string `json:"teacher_role_in_lesson"`
	SubstitutionReason  string `json:"substitution_reason,omitempty"`
}
