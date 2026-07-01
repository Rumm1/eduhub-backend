package teacher

import "github.com/google/uuid"

type Teacher struct {
	UserID          uuid.UUID
	OrganizationID  uuid.UUID
	Email           string
	FullName        string
	Phone           string
	Bio             string
	ExperienceYears int
	EmploymentType  string
	HourlyRate      float64
	FixedSalary     float64
}

type Subject struct {
	ID   uuid.UUID
	Name string
}
