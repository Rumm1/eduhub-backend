package report

import (
	"bytes"
	"fmt"

	"github.com/xuri/excelize/v2"
)

func BuildTeacherScheduleXLSX(report TeacherScheduleReportResponse) ([]byte, string, error) {
	file := excelize.NewFile()

	sheetName := "Teacher Schedule"
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)

	title := fmt.Sprintf("Teacher Schedule Report: %s", report.TeacherName)
	file.SetCellValue(sheetName, "A1", title)
	file.SetCellValue(sheetName, "A2", "From")
	file.SetCellValue(sheetName, "B2", report.FromDate)
	file.SetCellValue(sheetName, "C2", "To")
	file.SetCellValue(sheetName, "D2", report.ToDate)

	file.SetCellValue(sheetName, "A4", "Total lessons")
	file.SetCellValue(sheetName, "B4", report.TotalLessons)
	file.SetCellValue(sheetName, "C4", "Actual lessons")
	file.SetCellValue(sheetName, "D4", report.ActualLessons)
	file.SetCellValue(sheetName, "E4", "Planned only")
	file.SetCellValue(sheetName, "F4", report.PlannedOnlyLessons)
	file.SetCellValue(sheetName, "G4", "Substitutions")
	file.SetCellValue(sheetName, "H4", report.Substitutions)
	file.SetCellValue(sheetName, "I4", "Actual hours")
	file.SetCellValue(sheetName, "J4", report.TotalActualHours)

	headers := []string{
		"Date", "Start", "End", "Hours", "Group", "Branch", "Subject",
		"Topic", "Status", "Planned Teacher", "Actual Teacher", "Substitution", "Role", "Reason",
	}

	headerRow := 6
	writeHeaders(file, sheetName, headers, headerRow)

	for rowIndex, item := range report.Items {
		row := headerRow + rowIndex + 1

		values := []interface{}{
			item.LessonDate,
			item.StartTime,
			item.EndTime,
			item.Hours,
			item.GroupName,
			item.BranchName,
			item.SubjectName,
			item.Topic,
			item.Status,
			item.PlannedTeacherName,
			item.ActualTeacherName,
			item.IsSubstitution,
			item.TeacherRoleInLesson,
			item.SubstitutionReason,
		}

		writeRow(file, sheetName, row, values)
	}

	autoWidth(file, sheetName, len(headers))

	file.SetColWidth(sheetName, "H", "H", 32)
	file.SetColWidth(sheetName, "J", "K", 26)
	file.SetColWidth(sheetName, "N", "N", 40)

	styleHeader(file, sheetName, "A1", "A1")
	styleHeader(file, sheetName, "A6", "N6")

	return writeWorkbook(file, fmt.Sprintf("teacher_schedule_%s_%s_%s.xlsx", report.TeacherID, report.FromDate, report.ToDate))
}

func BuildPaymentsReportXLSX(report PaymentsReportResponse) ([]byte, string, error) {
	file := excelize.NewFile()

	sheetName := "Payments"
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)

	file.SetCellValue(sheetName, "A1", "Payments Report")
	file.SetCellValue(sheetName, "A2", "From")
	file.SetCellValue(sheetName, "B2", report.FromDate)
	file.SetCellValue(sheetName, "C2", "To")
	file.SetCellValue(sheetName, "D2", report.ToDate)

	file.SetCellValue(sheetName, "A4", "Total payments")
	file.SetCellValue(sheetName, "B4", report.TotalPayments)
	file.SetCellValue(sheetName, "C4", "Total amount")
	file.SetCellValue(sheetName, "D4", report.TotalAmount)
	file.SetCellValue(sheetName, "E4", "Paid")
	file.SetCellValue(sheetName, "F4", report.PaidAmount)
	file.SetCellValue(sheetName, "G4", "Pending")
	file.SetCellValue(sheetName, "H4", report.PendingAmount)
	file.SetCellValue(sheetName, "I4", "Refunded")
	file.SetCellValue(sheetName, "J4", report.RefundedAmount)
	file.SetCellValue(sheetName, "K4", "Cancelled")
	file.SetCellValue(sheetName, "L4", report.CancelledAmount)

	headers := []string{
		"Payment Date",
		"Payment Period",
		"Student",
		"Group",
		"Branch",
		"Amount",
		"Method",
		"Status",
		"Comment",
	}

	headerRow := 6
	writeHeaders(file, sheetName, headers, headerRow)

	for rowIndex, item := range report.Items {
		row := headerRow + rowIndex + 1

		values := []interface{}{
			item.PaymentDate,
			item.PaymentPeriod,
			item.StudentName,
			item.GroupName,
			item.BranchName,
			item.Amount,
			item.PaymentMethod,
			item.Status,
			item.Comment,
		}

		writeRow(file, sheetName, row, values)
	}

	autoWidth(file, sheetName, len(headers))

	file.SetColWidth(sheetName, "C", "E", 24)
	file.SetColWidth(sheetName, "I", "I", 42)

	styleHeader(file, sheetName, "A1", "A1")
	styleHeader(file, sheetName, "A6", "I6")

	return writeWorkbook(file, fmt.Sprintf("payments_report_%s_%s.xlsx", report.FromDate, report.ToDate))
}

func writeHeaders(file *excelize.File, sheetName string, headers []string, headerRow int) {
	for index, header := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, headerRow)
		if err != nil {
			continue
		}

		file.SetCellValue(sheetName, cell, header)
	}
}

func writeRow(file *excelize.File, sheetName string, row int, values []interface{}) {
	for colIndex, value := range values {
		cell, err := excelize.CoordinatesToCellName(colIndex+1, row)
		if err != nil {
			continue
		}

		file.SetCellValue(sheetName, cell, value)
	}
}

func autoWidth(file *excelize.File, sheetName string, columns int) {
	for col := 1; col <= columns; col++ {
		columnName, err := excelize.ColumnNumberToName(col)
		if err != nil {
			continue
		}

		file.SetColWidth(sheetName, columnName, columnName, 18)
	}
}

func styleHeader(file *excelize.File, sheetName string, startCell string, endCell string) {
	headerStyle, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})
	if err != nil {
		return
	}

	file.SetCellStyle(sheetName, startCell, endCell, headerStyle)
}

func writeWorkbook(file *excelize.File, filename string) ([]byte, string, error) {
	var buffer bytes.Buffer
	if err := file.Write(&buffer); err != nil {
		return nil, "", err
	}

	return buffer.Bytes(), filename, nil
}
