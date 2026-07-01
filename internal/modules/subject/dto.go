package subject

type CreateSubjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type SubjectResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Description    string `json:"description,omitempty"`
	Status         string `json:"status"`
}

type ListSubjectsResponse struct {
	Items []SubjectResponse `json:"items"`
	Total int               `json:"total"`
}
