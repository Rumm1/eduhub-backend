package branch

import "github.com/google/uuid"

type Branch struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Address        string
	Phone          string
	Status         string
}
