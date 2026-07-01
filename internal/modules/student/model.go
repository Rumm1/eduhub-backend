package student

import "github.com/google/uuid"

type Student struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	FullName       string
	Phone          string
	BirthDate      string
	Gender         string
	Status         string
	Source         string
	Notes          string
}

type Parent struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	FullName       string
	Phone          string
	Email          string
	Relation       string
}
