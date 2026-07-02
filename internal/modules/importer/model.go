package importer

import "github.com/google/uuid"

type StudentImportRow struct {
	RowNumber       int
	StudentFullName string
	StudentPhone    string
	ParentFullName  string
	ParentPhone     string
	ParentEmail     string
	GroupName       string
	Relation        string
}

type StudentImportPreviewRow struct {
	RowNumber       int      `json:"row_number"`
	Status          string   `json:"status"`
	Errors          []string `json:"errors"`
	Warnings        []string `json:"warnings"`
	StudentFullName string   `json:"student_full_name"`
	StudentPhone    string   `json:"student_phone,omitempty"`
	ParentFullName  string   `json:"parent_full_name"`
	ParentPhone     string   `json:"parent_phone,omitempty"`
	ParentEmail     string   `json:"parent_email,omitempty"`
	GroupName       string   `json:"group_name"`
	BranchName      string   `json:"branch_name,omitempty"`
	Relation        string   `json:"relation,omitempty"`
}

type GroupLookup struct {
	ID       uuid.UUID
	BranchID uuid.UUID
	Name     string
	Branch   string
}

type ValidStudentImportRow struct {
	Row   StudentImportRow
	Group GroupLookup
}

type ImportConfirmResult struct {
	CreatedStudents         int `json:"created_students"`
	ReusedStudents          int `json:"reused_students"`
	CreatedParents          int `json:"created_parents"`
	ReusedParents           int `json:"reused_parents"`
	LinkedParentsToStudents int `json:"linked_parents_to_students"`
	LinkedStudentsToGroups  int `json:"linked_students_to_groups"`
}
