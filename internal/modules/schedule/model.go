package schedule

import "github.com/google/uuid"

type Schedule struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	GroupID        uuid.UUID
	Weekday        int
	StartTime      string
	EndTime        string
	Room           string
}
