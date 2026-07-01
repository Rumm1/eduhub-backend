package branch

type CreateBranchRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

type BranchResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Name           string `json:"name"`
	Address        string `json:"address,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Status         string `json:"status"`
}

type ListBranchesResponse struct {
	Items []BranchResponse `json:"items"`
	Total int              `json:"total"`
}
