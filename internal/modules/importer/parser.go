package importer

import (
	"encoding/csv"
	"io"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ParseStudentImportFile(filename string, reader io.Reader) ([]StudentImportRow, error) {
	extension := strings.ToLower(filepath.Ext(filename))

	switch extension {
	case ".csv":
		return parseStudentImportCSV(reader)
	case ".xlsx":
		return parseStudentImportXLSX(reader)
	default:
		return nil, ErrFileTypeUnsupported
	}
}

func parseStudentImportCSV(reader io.Reader) ([]StudentImportRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return buildRowsFromRawTable(records)
}

func parseStudentImportXLSX(reader io.Reader) ([]StudentImportRow, error) {
	file, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sheets := file.GetSheetList()
	if len(sheets) == 0 {
		return nil, ErrEmptyImportFile
	}

	rows, err := file.GetRows(sheets[0])
	if err != nil {
		return nil, err
	}

	return buildRowsFromRawTable(rows)
}

func buildRowsFromRawTable(rawRows [][]string) ([]StudentImportRow, error) {
	if len(rawRows) < 2 {
		return nil, ErrEmptyImportFile
	}

	headerMap := buildHeaderMap(rawRows[0])

	result := make([]StudentImportRow, 0, len(rawRows)-1)

	for index, rawRow := range rawRows[1:] {
		row := StudentImportRow{
			RowNumber:       index + 2,
			StudentFullName: getCell(rawRow, headerMap["student_full_name"]),
			StudentPhone:    getCell(rawRow, headerMap["student_phone"]),
			ParentFullName:  getCell(rawRow, headerMap["parent_full_name"]),
			ParentPhone:     getCell(rawRow, headerMap["parent_phone"]),
			ParentEmail:     getCell(rawRow, headerMap["parent_email"]),
			GroupName:       getCell(rawRow, headerMap["group_name"]),
			Relation:        getCell(rawRow, headerMap["relation"]),
		}

		if isEmptyImportRow(row) {
			continue
		}

		if row.Relation == "" {
			row.Relation = "parent"
		}

		result = append(result, row)
	}

	if len(result) == 0 {
		return nil, ErrEmptyImportFile
	}

	return result, nil
}

func buildHeaderMap(headers []string) map[string]int {
	result := make(map[string]int)

	for index, header := range headers {
		normalized := normalizeHeader(header)

		switch normalized {
		case "student_full_name", "student_name", "student", "fio_student", "fio_uchenika", "ученик", "фио_ученика", "имя_ученика", "оқушы", "оқушы_аты":
			result["student_full_name"] = index
		case "student_phone", "phone_student", "student_mobile", "телефон_ученика", "номер_ученика", "оқушы_телефоны":
			result["student_phone"] = index
		case "parent_full_name", "parent_name", "parent", "fio_parent", "fio_roditelya", "родитель", "фио_родителя", "имя_родителя", "ата_ана", "ата_ана_аты":
			result["parent_full_name"] = index
		case "parent_phone", "phone_parent", "parent_mobile", "телефон_родителя", "номер_родителя", "ата_ана_телефоны":
			result["parent_phone"] = index
		case "parent_email", "email_parent", "email", "почта_родителя", "email_родителя", "ата_ана_email":
			result["parent_email"] = index
		case "group_name", "group", "группа", "название_группы", "топ", "топ_атауы":
			result["group_name"] = index
		case "relation", "relationship", "связь", "родство", "кім":
			result["relation"] = index
		}
	}

	return result
}

func normalizeHeader(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "\ufeff")
	value = strings.ToLower(value)

	replacer := strings.NewReplacer(
		" ", "_",
		"-", "_",
		".", "_",
		"/", "_",
		"\\", "_",
		"(", "",
		")", "",
		":", "",
	)

	value = replacer.Replace(value)

	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}

	return strings.Trim(value, "_")
}

func getCell(row []string, index int) string {
	if index < 0 || index >= len(row) {
		return ""
	}

	return strings.TrimSpace(row[index])
}

func isEmptyImportRow(row StudentImportRow) bool {
	return strings.TrimSpace(row.StudentFullName) == "" &&
		strings.TrimSpace(row.StudentPhone) == "" &&
		strings.TrimSpace(row.ParentFullName) == "" &&
		strings.TrimSpace(row.ParentPhone) == "" &&
		strings.TrimSpace(row.ParentEmail) == "" &&
		strings.TrimSpace(row.GroupName) == ""
}
