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
