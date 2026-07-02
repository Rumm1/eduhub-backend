package parent

type CreateParentRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type UpdateParentRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
}

type AttachStudentRequest struct {
	Relation string `json:"relation"`
}

type ParentResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	FullName       string `json:"full_name"`
	Phone          string `json:"phone,omitempty"`
	Email          string `json:"email,omitempty"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}

type StudentResponse struct {
	ID       string `json:"id"`
	BranchID string `json:"branch_id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Status   string `json:"status"`
	Relation string `json:"relation,omitempty"`
}

type ListParentsResponse struct {
	Items []ParentResponse `json:"items"`
	Total int              `json:"total"`
}

type ListParentStudentsResponse struct {
	Items []StudentResponse `json:"items"`
	Total int               `json:"total"`
}
