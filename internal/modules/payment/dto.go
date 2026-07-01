package payment

type CreatePaymentRequest struct {
StudentID     string `json:"student_id"`
GroupID       string `json:"group_id"`
Amount        string `json:"amount"`
PaymentDate   string `json:"payment_date"`
PaymentMethod string `json:"payment_method"`
Status        string `json:"status"`
Comment       string `json:"comment"`
}

type PaymentResponse struct {
ID             string `json:"id"`
OrganizationID string `json:"organization_id"`
BranchID       string `json:"branch_id"`
StudentID      string `json:"student_id"`
GroupID        string `json:"group_id,omitempty"`
Amount         string `json:"amount"`
PaymentDate    string `json:"payment_date"`
PaymentMethod  string `json:"payment_method,omitempty"`
Status         string `json:"status"`
Comment        string `json:"comment,omitempty"`
}

type ListPaymentsResponse struct {
Items []PaymentResponse `json:"items"`
Total int               `json:"total"`
}
