package payroll

import "github.com/google/uuid"

type PayrollPeriod struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Month          int
	Year           int
	Status         string
}

type PayrollEntry struct {
	ID                        uuid.UUID
	OrganizationID            uuid.UUID
	PeriodID                  uuid.UUID
	TeacherID                 uuid.UUID
	LessonsCount              int
	SubstitutionCount         int
	HoursWorked               string
	HourlyRate                string
	BaseAmount                string
	BonusAmount               string
	PenaltyAmount             string
	CorrectionAmount          string
	TotalAmount               string
	FinalAmount               string
	Status                    string
	TeacherConfirmationStatus string
	TeacherDisputeReason      string
	Comment                   string
}

type PayrollAdjustment struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	PeriodID       uuid.UUID
	PayrollEntryID uuid.UUID
	EmployeeID     uuid.UUID
	AdjustmentType string
	Amount         string
	Reason         string
	Status         string
	CreatedBy      string
}
