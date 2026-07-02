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

type UserProfile struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID uuid.UUID
	BranchID       *uuid.UUID
	DisplayName    string
	Position       string
	ProfileType    string
	Status         string
	IsDefault      bool
	RoleCodes      []string
	BranchIDs      []uuid.UUID
}
