package lesson

import "github.com/google/uuid"

type Lesson struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	GroupID        uuid.UUID
	TeacherID      string
	SubjectID      uuid.UUID
	LessonDate     string
	StartTime      string
	EndTime        string
	Topic          string
	Status         string
}
