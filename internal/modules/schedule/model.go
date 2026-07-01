package schedule

import "github.com/google/uuid"

type Schedule struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	GroupID        uuid.UUID
	SubjectID      uuid.UUID
	TeacherID      string
	Weekday        int
	StartTime      string
	EndTime        string
	Room           string
}

type GeneratedLesson struct {
	ID             uuid.UUID
	ScheduleID     uuid.UUID
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
