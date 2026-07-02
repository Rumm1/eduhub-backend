package permission

import "github.com/google/uuid"

type Permission struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
	Group       string
}

type PermissionGroup struct {
	Name        string
	Title       string
	Permissions []Permission
}
