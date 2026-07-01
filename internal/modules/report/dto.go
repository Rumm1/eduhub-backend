package report

type TeacherScheduleReportResponse struct {
	TeacherID          string                      `json:"teacher_id"`
	TeacherName        string                      `json:"teacher_name"`
	FromDate           string                      `json:"from_date"`
	ToDate             string                      `json:"to_date"`
	TotalLessons       int                         `json:"total_lessons"`
	ActualLessons      int                         `json:"actual_lessons"`
	PlannedOnlyLessons int                         `json:"planned_only_lessons"`
	Substitutions      int                         `json:"substitutions"`
	TotalActualHours   string                      `json:"total_actual_hours"`
	Items              []TeacherScheduleReportItem `json:"items"`
}

type TeacherScheduleReportItem struct {
	LessonID            string `json:"lesson_id"`
	LessonDate          string `json:"lesson_date"`
	StartTime           string `json:"start_time"`
	EndTime             string `json:"end_time"`
	Hours               string `json:"hours"`
	Topic               string `json:"topic,omitempty"`
	Status              string `json:"status"`
	GroupID             string `json:"group_id"`
	GroupName           string `json:"group_name"`
	BranchID            string `json:"branch_id"`
	BranchName          string `json:"branch_name"`
	SubjectID           string `json:"subject_id"`
	SubjectName         string `json:"subject_name"`
	PlannedTeacherID    string `json:"planned_teacher_id,omitempty"`
	PlannedTeacherName  string `json:"planned_teacher_name,omitempty"`
	ActualTeacherID     string `json:"actual_teacher_id,omitempty"`
	ActualTeacherName   string `json:"actual_teacher_name,omitempty"`
	IsSubstitution      bool   `json:"is_substitution"`
	TeacherRoleInLesson string `json:"teacher_role_in_lesson"`
	SubstitutionReason  string `json:"substitution_reason,omitempty"`
}

type PaymentsReportResponse struct {
	FromDate        string               `json:"from_date"`
	ToDate          string               `json:"to_date"`
	TotalPayments   int                  `json:"total_payments"`
	TotalAmount     string               `json:"total_amount"`
	PaidAmount      string               `json:"paid_amount"`
	PendingAmount   string               `json:"pending_amount"`
	RefundedAmount  string               `json:"refunded_amount"`
	CancelledAmount string               `json:"cancelled_amount"`
	Items           []PaymentsReportItem `json:"items"`
}

type PaymentsReportItem struct {
	PaymentID     string `json:"payment_id"`
	PaymentDate   string `json:"payment_date"`
	PaymentPeriod string `json:"payment_period,omitempty"`
	StudentID     string `json:"student_id"`
	StudentName   string `json:"student_name"`
	GroupID       string `json:"group_id,omitempty"`
	GroupName     string `json:"group_name,omitempty"`
	BranchID      string `json:"branch_id"`
	BranchName    string `json:"branch_name"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method,omitempty"`
	Status        string `json:"status"`
	Comment       string `json:"comment,omitempty"`
}

type StudentBalancesReportResponse struct {
	Period              string                     `json:"period"`
	TotalStudents       int                        `json:"total_students"`
	PaidCount           int                        `json:"paid_count"`
	PartialCount        int                        `json:"partial_count"`
	UnpaidCount         int                        `json:"unpaid_count"`
	TotalExpectedAmount string                     `json:"total_expected_amount"`
	TotalPaidAmount     string                     `json:"total_paid_amount"`
	TotalDebtAmount     string                     `json:"total_debt_amount"`
	Items               []StudentBalanceReportItem `json:"items"`
}

type StudentBalanceReportItem struct {
	StudentID     string `json:"student_id"`
	StudentName   string `json:"student_name"`
	GroupID       string `json:"group_id"`
	GroupName     string `json:"group_name"`
	BranchID      string `json:"branch_id"`
	BranchName    string `json:"branch_name"`
	MonthlyPrice  string `json:"monthly_price"`
	PaidAmount    string `json:"paid_amount"`
	DebtAmount    string `json:"debt_amount"`
	PaymentStatus string `json:"payment_status"`
}

type PayrollReportResponse struct {
	Period                 string              `json:"period"`
	TotalEntries           int                 `json:"total_entries"`
	TotalLessons           int                 `json:"total_lessons"`
	TotalSubstitutions     int                 `json:"total_substitutions"`
	TotalHours             string              `json:"total_hours"`
	TotalBaseAmount        string              `json:"total_base_amount"`
	TotalBonusAmount       string              `json:"total_bonus_amount"`
	TotalPenaltyAmount     string              `json:"total_penalty_amount"`
	TotalCorrectionAmount  string              `json:"total_correction_amount"`
	TotalFinalAmount       string              `json:"total_final_amount"`
	DraftCount             int                 `json:"draft_count"`
	SentToTeacherCount     int                 `json:"sent_to_teacher_count"`
	TeacherApprovedCount   int                 `json:"teacher_approved_count"`
	TeacherDisputedCount   int                 `json:"teacher_disputed_count"`
	ApprovedByFinanceCount int                 `json:"approved_by_finance_count"`
	PaidCount              int                 `json:"paid_count"`
	Items                  []PayrollReportItem `json:"items"`
}

type PayrollReportItem struct {
	EntryID                   string `json:"entry_id"`
	PeriodID                  string `json:"period_id"`
	TeacherID                 string `json:"teacher_id"`
	TeacherName               string `json:"teacher_name"`
	LessonsCount              int    `json:"lessons_count"`
	SubstitutionCount         int    `json:"substitution_count"`
	HoursWorked               string `json:"hours_worked"`
	HourlyRate                string `json:"hourly_rate"`
	BaseAmount                string `json:"base_amount"`
	BonusAmount               string `json:"bonus_amount"`
	PenaltyAmount             string `json:"penalty_amount"`
	CorrectionAmount          string `json:"correction_amount"`
	FinalAmount               string `json:"final_amount"`
	Status                    string `json:"status"`
	TeacherConfirmationStatus string `json:"teacher_confirmation_status"`
	TeacherDisputeReason      string `json:"teacher_dispute_reason,omitempty"`
	Comment                   string `json:"comment,omitempty"`
}
