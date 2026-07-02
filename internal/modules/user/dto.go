package user

type CreateUserRequest struct {
	Email      string                     `json:"email"`
	Password   string                     `json:"password"`
	FullName   string                     `json:"full_name"`
	Phone      string                     `json:"phone"`
	AvatarPath string                     `json:"avatar_path"`
	Profiles   []CreateUserProfileRequest `json:"profiles"`
}

type CreateUserProfileRequest struct {
	BranchID    string   `json:"branch_id"`
	DisplayName string   `json:"display_name"`
	Position    string   `json:"position"`
	ProfileType string   `json:"profile_type"`
	IsDefault   bool     `json:"is_default"`
	RoleCodes   []string `json:"role_codes"`
	BranchIDs   []string `json:"branch_ids"`
}

type CreateUserResponse struct {
	User                 UserResponse                 `json:"user"`
	TemporaryCredentials TemporaryCredentialsResponse `json:"temporary_credentials"`
}

type TemporaryCredentialsResponse struct {
	Login              string `json:"login"`
	Password           string `json:"password"`
	MustChangePassword bool   `json:"must_change_password"`
}

type UserResponse struct {
	ID             string                `json:"id"`
	OrganizationID string                `json:"organization_id"`
	Email          string                `json:"email"`
	FullName       string                `json:"full_name"`
	Phone          string                `json:"phone,omitempty"`
	AvatarPath     string                `json:"avatar_path,omitempty"`
	Status         string                `json:"status"`
	Roles          []string              `json:"roles"`
	BranchIDs      []string              `json:"branch_ids"`
	Profiles       []UserProfileResponse `json:"profiles"`
}

type UserProfileResponse struct {
	ID             string   `json:"id"`
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

type ListUsersResponse struct {
	Items []UserResponse `json:"items"`
	Total int            `json:"total"`
}
