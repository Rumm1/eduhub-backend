package parent

import "github.com/google/uuid"

type Parent struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	FullName       string
	Phone          string
	Email          string
	CreatedAt      string
	UpdatedAt      string
}

type Student struct {
	ID       uuid.UUID
	BranchID uuid.UUID
	FullName string
	Phone    string
	Status   string
	Relation string
}
