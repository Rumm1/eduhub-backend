package importer

type StudentImportSummary struct {
	TotalRows   int `json:"total_rows"`
	ValidRows   int `json:"valid_rows"`
	InvalidRows int `json:"invalid_rows"`
	WarningRows int `json:"warning_rows"`
}

type StudentImportPreviewResponse struct {
	Summary StudentImportSummary      `json:"summary"`
	Rows    []StudentImportPreviewRow `json:"rows"`
}

type StudentImportConfirmSummary struct {
	TotalRows               int `json:"total_rows"`
	ValidRows               int `json:"valid_rows"`
	InvalidRows             int `json:"invalid_rows"`
	WarningRows             int `json:"warning_rows"`
	CreatedStudents         int `json:"created_students"`
	ReusedStudents          int `json:"reused_students"`
	CreatedParents          int `json:"created_parents"`
	ReusedParents           int `json:"reused_parents"`
	LinkedParentsToStudents int `json:"linked_parents_to_students"`
	LinkedStudentsToGroups  int `json:"linked_students_to_groups"`
}

type StudentImportConfirmResponse struct {
	Summary StudentImportConfirmSummary `json:"summary"`
	Rows    []StudentImportPreviewRow   `json:"rows"`
}
