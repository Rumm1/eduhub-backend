package report

type TeacherScheduleReport struct {
	TeacherID          string
	TeacherName        string
	FromDate           string
	ToDate             string
	TotalLessons       int
	ActualLessons      int
	PlannedOnlyLessons int
	Substitutions      int
	TotalActualHours   string
	Items              []TeacherScheduleItem
}

type TeacherScheduleItem struct {
	LessonID            string
	LessonDate          string
	StartTime           string
	EndTime             string
	Hours               string
	Topic               string
	Status              string
	GroupID             string
	GroupName           string
	BranchID            string
	BranchName          string
	SubjectID           string
	SubjectName         string
	PlannedTeacherID    string
	PlannedTeacherName  string
	ActualTeacherID     string
	ActualTeacherName   string
	IsSubstitution      bool
	TeacherRoleInLesson string
	SubstitutionReason  string
}

type PaymentsReport struct {
	FromDate        string
	ToDate          string
	TotalPayments   int
	TotalAmount     string
	PaidAmount      string
	PendingAmount   string
	RefundedAmount  string
	CancelledAmount string
	Items           []PaymentReportItem
}

type PaymentReportItem struct {
	PaymentID     string
	PaymentDate   string
	PaymentPeriod string
	StudentID     string
	StudentName   string
	GroupID       string
	GroupName     string
	BranchID      string
	BranchName    string
	Amount        string
	PaymentMethod string
	Status        string
	Comment       string
}

type StudentBalancesReport struct {
	Period              string
	TotalStudents       int
	PaidCount           int
	PartialCount        int
	UnpaidCount         int
	TotalExpectedAmount string
	TotalPaidAmount     string
	TotalDebtAmount     string
	Items               []StudentBalanceItem
}

type StudentBalanceItem struct {
	StudentID     string
	StudentName   string
	GroupID       string
	GroupName     string
	BranchID      string
	BranchName    string
	MonthlyPrice  string
	PaidAmount    string
	DebtAmount    string
	PaymentStatus string
}
