package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SwitchProfileRequest struct {
	ProfileID string `json:"profile_id"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	User        UserResponse `json:"user"`
}

type UserResponse struct {
	ID                string            `json:"id"`
	ProfileID         *string           `json:"profile_id,omitempty"`
	OrganizationID    *string           `json:"organization_id,omitempty"`
	Email             string            `json:"email"`
	FullName          string            `json:"full_name"`
	Roles             []string          `json:"roles"`
	Permissions       []string          `json:"permissions"`
	BranchIDs         []string          `json:"branch_ids"`
	CurrentProfile    *ProfileResponse  `json:"current_profile,omitempty"`
	AvailableProfiles []ProfileResponse `json:"available_profiles"`
}

type ProfileResponse struct {
	ID             string   `json:"id"`
	OrganizationID *string  `json:"organization_id,omitempty"`
	BranchID       *string  `json:"branch_id,omitempty"`
	DisplayName    string   `json:"display_name"`
	Position       string   `json:"position"`
	ProfileType    string   `json:"profile_type"`
	Status         string   `json:"status"`
	IsDefault      bool     `json:"is_default"`
	Roles          []string `json:"roles"`
	BranchIDs      []string `json:"branch_ids"`
}
