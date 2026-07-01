package group

import "github.com/google/uuid"

type Group struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	SubjectID      uuid.UUID
	TeacherID      string
	Name           string
	Level          string
	Status         string
	MaxStudents    int
	StartDate      string
	EndDate        string
	StudentsCount  int
}

type GroupStudent struct {
	StudentID uuid.UUID
	FullName  string
	Phone     string
	Status    string
	JoinedAt  string
}
