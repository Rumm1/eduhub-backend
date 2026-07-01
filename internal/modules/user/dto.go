package user

type CreateUserRequest struct {
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	FullName   string   `json:"full_name"`
	Phone      string   `json:"phone"`
	AvatarPath string   `json:"avatar_path"`
	RoleCode   string   `json:"role_code"`
	BranchIDs  []string `json:"branch_ids"`
}

type UserResponse struct {
	ID             string   `json:"id"`
	OrganizationID string   `json:"organization_id"`
	Email          string   `json:"email"`
	FullName       string   `json:"full_name"`
	Phone          string   `json:"phone,omitempty"`
	AvatarPath     string   `json:"avatar_path,omitempty"`
	Status         string   `json:"status"`
	Roles          []string `json:"roles"`
	BranchIDs      []string `json:"branch_ids"`
}

type ListUsersResponse struct {
	Items []UserResponse `json:"items"`
	Total int            `json:"total"`
}
