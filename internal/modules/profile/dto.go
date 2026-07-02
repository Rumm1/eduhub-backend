package profile

type CreateProfileRequest struct {
	BranchID    string   `json:"branch_id"`
	DisplayName string   `json:"display_name"`
	Position    string   `json:"position"`
	ProfileType string   `json:"profile_type"`
	IsDefault   bool     `json:"is_default"`
	RoleCodes   []string `json:"role_codes"`
	BranchIDs   []string `json:"branch_ids"`
}

type UpdateProfileRequest struct {
	BranchID    *string `json:"branch_id"`
	DisplayName *string `json:"display_name"`
	Position    *string `json:"position"`
	ProfileType *string `json:"profile_type"`
	Status      *string `json:"status"`
	IsDefault   *bool   `json:"is_default"`
}

type AddRoleRequest struct {
	RoleCode string `json:"role_code"`
}

type AddBranchRequest struct {
	BranchID string `json:"branch_id"`
}

type ProfileResponse struct {
	ID             string   `json:"id"`
	UserID         string   `json:"user_id"`
	OrganizationID string   `json:"organization_id"`
	BranchID       string   `json:"branch_id,omitempty"`
	DisplayName    string   `json:"display_name"`
	Position       string   `json:"position"`
	ProfileType    string   `json:"profile_type"`
	Status         string   `json:"status"`
	IsDefault      bool     `json:"is_default"`
	RoleCodes      []string `json:"role_codes"`
	BranchIDs      []string `json:"branch_ids"`
}

type ListProfilesResponse struct {
	Items []ProfileResponse `json:"items"`
	Total int               `json:"total"`
}
