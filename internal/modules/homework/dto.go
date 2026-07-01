package homework

type CreateHomeworkRequest struct {
	LessonID    string `json:"lesson_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
}

type HomeworkResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	GroupID        string `json:"group_id"`
	LessonID       string `json:"lesson_id,omitempty"`
	TeacherID      string `json:"teacher_id"`
	Title          string `json:"title"`
	Description    string `json:"description,omitempty"`
	DueDate        string `json:"due_date,omitempty"`
}

type ListHomeworksResponse struct {
	Items []HomeworkResponse `json:"items"`
	Total int                `json:"total"`
}
