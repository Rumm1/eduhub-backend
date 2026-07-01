package student

type CreateStudentRequest struct {
	BranchID  string               `json:"branch_id"`
	FullName  string               `json:"full_name"`
	Phone     string               `json:"phone"`
	BirthDate string               `json:"birth_date"`
	Gender    string               `json:"gender"`
	Source    string               `json:"source"`
	Notes     string               `json:"notes"`
	Parent    *CreateParentRequest `json:"parent"`
}

type CreateParentRequest struct {
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Relation string `json:"relation"`
}

type StudentResponse struct {
	ID             string           `json:"id"`
	OrganizationID string           `json:"organization_id"`
	BranchID       string           `json:"branch_id"`
	FullName       string           `json:"full_name"`
	Phone          string           `json:"phone,omitempty"`
	BirthDate      string           `json:"birth_date,omitempty"`
	Gender         string           `json:"gender,omitempty"`
	Status         string           `json:"status"`
	Source         string           `json:"source,omitempty"`
	Notes          string           `json:"notes,omitempty"`
	Parents        []ParentResponse `json:"parents"`
}

type ParentResponse struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Relation string `json:"relation,omitempty"`
}

type ListStudentsResponse struct {
	Items []StudentResponse `json:"items"`
	Total int               `json:"total"`
}
