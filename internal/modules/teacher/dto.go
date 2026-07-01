package teacher

type CreateTeacherRequest struct {
	UserID          string   `json:"user_id"`
	Bio             string   `json:"bio"`
	ExperienceYears int      `json:"experience_years"`
	EmploymentType  string   `json:"employment_type"`
	HourlyRate      float64  `json:"hourly_rate"`
	FixedSalary     float64  `json:"fixed_salary"`
	SubjectIDs      []string `json:"subject_ids"`
}

type TeacherResponse struct {
	UserID          string            `json:"user_id"`
	OrganizationID  string            `json:"organization_id"`
	Email           string            `json:"email"`
	FullName        string            `json:"full_name"`
	Phone           string            `json:"phone,omitempty"`
	Bio             string            `json:"bio,omitempty"`
	ExperienceYears int               `json:"experience_years"`
	EmploymentType  string            `json:"employment_type,omitempty"`
	HourlyRate      float64           `json:"hourly_rate"`
	FixedSalary     float64           `json:"fixed_salary"`
	Subjects        []SubjectResponse `json:"subjects"`
}

type SubjectResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListTeachersResponse struct {
	Items []TeacherResponse `json:"items"`
	Total int               `json:"total"`
}
