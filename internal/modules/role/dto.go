package role

type CreateRoleRequest struct {
	Name            string   `json:"name"`
	Code            string   `json:"code"`
	Description     string   `json:"description"`
	PermissionCodes []string `json:"permission_codes"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
}

type AddPermissionRequest struct {
	PermissionCode string `json:"permission_code"`
}

type RoleResponse struct {
	ID              string   `json:"id"`
	OrganizationID  string   `json:"organization_id,omitempty"`
	Name            string   `json:"name"`
	Code            string   `json:"code"`
	Description     string   `json:"description"`
	IsSystem        bool     `json:"is_system"`
	PermissionCodes []string `json:"permission_codes"`
}

type ListRolesResponse struct {
	Items []RoleResponse `json:"items"`
	Total int            `json:"total"`
}
