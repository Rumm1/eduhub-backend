package permission

type PermissionResponse struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Group       string `json:"group"`
}

type ListPermissionsResponse struct {
	Items []PermissionResponse `json:"items"`
	Total int                  `json:"total"`
}

type PermissionGroupResponse struct {
	Name        string               `json:"name"`
	Title       string               `json:"title"`
	Permissions []PermissionResponse `json:"permissions"`
}

type ListPermissionGroupsResponse struct {
	Items []PermissionGroupResponse `json:"items"`
	Total int                       `json:"total"`
}
