package group

type CreateGroupRequest struct {
	BranchID    string `json:"branch_id"`
	SubjectID   string `json:"subject_id"`
	TeacherID   string `json:"teacher_id"`
	Name        string `json:"name"`
	Level       string `json:"level"`
	MaxStudents int    `json:"max_students"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

type AddStudentToGroupRequest struct {
	StudentID string `json:"student_id"`
}

type GroupResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	BranchID       string `json:"branch_id"`
	SubjectID      string `json:"subject_id"`
	TeacherID      string `json:"teacher_id,omitempty"`
	Name           string `json:"name"`
	Level          string `json:"level,omitempty"`
	Status         string `json:"status"`
	MaxStudents    int    `json:"max_students"`
	StartDate      string `json:"start_date,omitempty"`
	EndDate        string `json:"end_date,omitempty"`
	StudentsCount  int    `json:"students_count"`
}

type ListGroupsResponse struct {
	Items []GroupResponse `json:"items"`
	Total int             `json:"total"`
}

type GroupStudentResponse struct {
	StudentID string `json:"student_id"`
	FullName  string `json:"full_name"`
	Phone     string `json:"phone,omitempty"`
	Status    string `json:"status"`
	JoinedAt  string `json:"joined_at"`
}

type ListGroupStudentsResponse struct {
	Items []GroupStudentResponse `json:"items"`
	Total int                    `json:"total"`
}
