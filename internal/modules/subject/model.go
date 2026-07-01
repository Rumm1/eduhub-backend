package subject

import "github.com/google/uuid"

type Subject struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Name           string
	Description    string
	Status         string
}
