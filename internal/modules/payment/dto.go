package payment

type CreatePaymentRequest struct {
	StudentID     string `json:"student_id"`
	GroupID       string `json:"group_id"`
	Amount        string `json:"amount"`
	PaymentDate   string `json:"payment_date"`
	PaymentPeriod string `json:"payment_period"`
	PaymentMethod string `json:"payment_method"`
	Status        string `json:"status"`
	Comment       string `json:"comment"`
}

type UpdateGroupPriceRequest struct {
	MonthlyPrice string `json:"monthly_price"`
}

type PaymentResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	BranchID       string `json:"branch_id"`
	StudentID      string `json:"student_id"`
	GroupID        string `json:"group_id,omitempty"`
	Amount         string `json:"amount"`
	PaymentDate    string `json:"payment_date"`
	PaymentPeriod  string `json:"payment_period,omitempty"`
	PaymentMethod  string `json:"payment_method,omitempty"`
	Status         string `json:"status"`
	Comment        string `json:"comment,omitempty"`
}

type GroupPriceResponse struct {
	GroupID      string `json:"group_id"`
	MonthlyPrice string `json:"monthly_price"`
}

type StudentBalanceResponse struct {
	StudentID      string `json:"student_id"`
	GroupID        string `json:"group_id"`
	BranchID       string `json:"branch_id"`
	PaymentPeriod  string `json:"payment_period"`
	ExpectedAmount string `json:"expected_amount"`
	PaidAmount     string `json:"paid_amount"`
	DebtAmount     string `json:"debt_amount"`
	OverpaidAmount string `json:"overpaid_amount"`
	Status         string `json:"status"`
}

type ListPaymentsResponse struct {
	Items []PaymentResponse `json:"items"`
	Total int               `json:"total"`
}
