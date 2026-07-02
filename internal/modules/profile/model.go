package profile

import "github.com/google/uuid"

type Profile struct {
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
