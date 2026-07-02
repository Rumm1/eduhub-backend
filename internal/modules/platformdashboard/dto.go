package platformdashboard

type DashboardResponse struct {
	OrganizationsCount int64                 `json:"organizations_count"`
	UsersCount         int64                 `json:"users_count"`
	BranchesCount      int64                 `json:"branches_count"`
	StudentsCount      int64                 `json:"students_count"`
	TeachersCount      int64                 `json:"teachers_count"`
	GroupsCount        int64                 `json:"groups_count"`
	Organizations      []OrganizationSummary `json:"organizations"`
}

type OrganizationSummary struct {
	ID            string          `json:"id"`
	Name          string          `json:"name"`
	BIN           string          `json:"bin"`
	Status        string          `json:"status"`
	LogoPath      string          `json:"logo_path"`
	LogoURL       string          `json:"logo_url"`
	UsersCount    int64           `json:"users_count"`
	BranchesCount int64           `json:"branches_count"`
	StudentsCount int64           `json:"students_count"`
	TeachersCount int64           `json:"teachers_count"`
	GroupsCount   int64           `json:"groups_count"`
	Branches      []BranchSummary `json:"branches"`
}

type BranchSummary struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Status  string `json:"status"`
}
