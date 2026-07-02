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

type UserProfile struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
	BranchID       *uuid.UUID
	DisplayName    string
	Position       string
	ProfileType    string
	Status         string
	IsDefault      bool
	Roles          []string
	Permissions    []string
	BranchIDs      []uuid.UUID
}

type UserAccessData struct {
	User              User
	Profile           UserProfile
	AvailableProfiles []UserProfile
	Roles             []string
	Permissions       []string
	BranchIDs         []uuid.UUID
}
