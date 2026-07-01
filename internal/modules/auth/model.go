package auth

import "github.com/google/uuid"

type User struct {
	ID             uuid.UUID
	OrganizationID *uuid.UUID
	Email          string
	PasswordHash   string
	FullName       string
	Status         string
}

type UserAccessData struct {
	User        User
	Roles       []string
	Permissions []string
	BranchIDs   []uuid.UUID
}
