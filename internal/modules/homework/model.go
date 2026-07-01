package homework

import "github.com/google/uuid"

type Homework struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	GroupID        uuid.UUID
	LessonID       uuid.UUID
	TeacherID      string
	Title          string
	Description    string
	DueDate        string
}
