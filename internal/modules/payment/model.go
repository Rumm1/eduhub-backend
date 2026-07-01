package payment

import "github.com/google/uuid"

type Payment struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	BranchID       uuid.UUID
	StudentID      uuid.UUID
	GroupID        string
	Amount         string
	PaymentDate    string
	PaymentPeriod  string
	PaymentMethod  string
	Status         string
	Comment        string
}

type StudentBalance struct {
	StudentID      uuid.UUID
	GroupID        uuid.UUID
	BranchID       uuid.UUID
	PaymentPeriod  string
	ExpectedAmount string
	PaidAmount     string
	DebtAmount     string
	OverpaidAmount string
	Status         string
}
