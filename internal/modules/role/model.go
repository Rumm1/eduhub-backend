package role

import "github.com/google/uuid"

type Role struct {
	ID              uuid.UUID
	OrganizationID  *uuid.UUID
	Name            string
	Code            string
	Description     string
	IsSystem        bool
	PermissionCodes []string
}
