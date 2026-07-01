package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID             string   `json:"id"`
	OrganizationID *string  `json:"organization_id,omitempty"`
	Email          string   `json:"email"`
	FullName       string   `json:"full_name"`
	Roles          []string `json:"roles"`
	Permissions    []string `json:"permissions"`
	BranchIDs      []string `json:"branch_ids"`
}
