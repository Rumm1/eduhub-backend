package organization

type CreateOrganizationRequest struct {
	Name          string `json:"name"`
	BIN           string `json:"bin"`
	Phone         string `json:"phone"`
	Email         string `json:"email"`
	AdminEmail    string `json:"admin_email"`
	AdminPassword string `json:"admin_password"`
	AdminFullName string `json:"admin_full_name"`
	AdminPhone    string `json:"admin_phone"`
}

type CreateOrganizationResponse struct {
	Organization OrganizationResponse `json:"organization"`
	Admin        AdminResponse        `json:"admin"`
}

type OrganizationResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	BIN    string `json:"bin,omitempty"`
	Phone  string `json:"phone,omitempty"`
	Email  string `json:"email,omitempty"`
	Status string `json:"status"`
}

type AdminResponse struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	FullName string   `json:"full_name"`
	Phone    string   `json:"phone,omitempty"`
	Roles    []string `json:"roles"`
}
