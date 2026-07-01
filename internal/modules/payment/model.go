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
PaymentMethod  string
Status         string
Comment        string
}
