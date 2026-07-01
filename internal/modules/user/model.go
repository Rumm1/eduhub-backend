package user

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	PasswordHash   string
	FullName       string
	Phone          string
	AvatarPath     string
	Status         string
}

type RoleTemplate struct {
	Code        string
	Name        string
	Description string
	Permissions []string
}
