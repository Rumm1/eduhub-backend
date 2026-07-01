package organization

import "github.com/google/uuid"

type Organization struct {
	ID     uuid.UUID
	Name   string
	BIN    string
	Phone  string
	Email  string
	Status string
}

type User struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	PasswordHash   string
	FullName       string
	Phone          string
	Status         string
}

type Role struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Code           string
	Description    string
	IsSystem       bool
}
