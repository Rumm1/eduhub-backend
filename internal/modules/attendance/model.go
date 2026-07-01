package attendance

import "github.com/google/uuid"

type Attendance struct {
	ID              string
	LessonID        uuid.UUID
	StudentID       uuid.UUID
	StudentFullName string
	Status          string
	Reason          string
	Comment         string
	MarkedBy        string
	MarkedAt        string
}
