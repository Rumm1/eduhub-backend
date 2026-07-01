package payroll

type CreatePeriodRequest struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type CreateAdjustmentRequest struct {
	AdjustmentType string `json:"adjustment_type"`
	Amount         string `json:"amount"`
	Reason         string `json:"reason"`
}

type DisputePayrollRequest struct {
	Reason string `json:"reason"`
}

type PayrollPeriodResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Month          int    `json:"month"`
	Year           int    `json:"year"`
	Status         string `json:"status"`
}

type PayrollEntryResponse struct {
	ID                        string `json:"id"`
	OrganizationID            string `json:"organization_id"`
	PeriodID                  string `json:"period_id"`
	TeacherID                 string `json:"teacher_id"`
	LessonsCount              int    `json:"lessons_count"`
	SubstitutionCount         int    `json:"substitution_count"`
	HoursWorked               string `json:"hours_worked"`
	HourlyRate                string `json:"hourly_rate"`
	BaseAmount                string `json:"base_amount"`
	BonusAmount               string `json:"bonus_amount"`
	PenaltyAmount             string `json:"penalty_amount"`
	CorrectionAmount          string `json:"correction_amount"`
	TotalAmount               string `json:"total_amount"`
	FinalAmount               string `json:"final_amount"`
	Status                    string `json:"status"`
	TeacherConfirmationStatus string `json:"teacher_confirmation_status"`
	TeacherDisputeReason      string `json:"teacher_dispute_reason,omitempty"`
	Comment                   string `json:"comment,omitempty"`
}

type PayrollAdjustmentResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	PeriodID       string `json:"period_id"`
	PayrollEntryID string `json:"payroll_entry_id"`
	EmployeeID     string `json:"employee_id"`
	AdjustmentType string `json:"adjustment_type"`
	Amount         string `json:"amount"`
	Reason         string `json:"reason,omitempty"`
	Status         string `json:"status"`
	CreatedBy      string `json:"created_by,omitempty"`
}

type GeneratePayrollResponse struct {
	Period PayrollPeriodResponse  `json:"period"`
	Items  []PayrollEntryResponse `json:"items"`
	Total  int                    `json:"total"`
}

type ListPayrollEntriesResponse struct {
	Items []PayrollEntryResponse `json:"items"`
	Total int                    `json:"total"`
}
