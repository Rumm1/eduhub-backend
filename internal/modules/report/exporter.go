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
"Date",
"Start",
"End",
"Hours",
"Group",
"Branch",
"Subject",
"Topic",
"Status",
"Planned Teacher",
"Actual Teacher",
"Substitution",
"Role",
"Reason",
}

headerRow := 6
for index, header := range headers {
cell, err := excelize.CoordinatesToCellName(index+1, headerRow)
if err != nil {
return nil, "", err
}

file.SetCellValue(sheetName, cell, header)
}

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

for colIndex, value := range values {
cell, err := excelize.CoordinatesToCellName(colIndex+1, row)
if err != nil {
return nil, "", err
}

file.SetCellValue(sheetName, cell, value)
}
}

for col := 1; col <= len(headers); col++ {
columnName, err := excelize.ColumnNumberToName(col)
if err != nil {
return nil, "", err
}

file.SetColWidth(sheetName, columnName, columnName, 18)
}

file.SetColWidth(sheetName, "H", "H", 32)
file.SetColWidth(sheetName, "J", "K", 26)
file.SetColWidth(sheetName, "N", "N", 40)

headerStyle, err := file.NewStyle(&excelize.Style{
Font: &excelize.Font{
Bold: true,
},
})
if err != nil {
return nil, "", err
}

file.SetCellStyle(sheetName, "A1", "A1", headerStyle)
file.SetCellStyle(sheetName, "A6", "N6", headerStyle)

var buffer bytes.Buffer
if err := file.Write(&buffer); err != nil {
return nil, "", err
}

filename := fmt.Sprintf("teacher_schedule_%s_%s_%s.xlsx", report.TeacherID, report.FromDate, report.ToDate)

return buffer.Bytes(), filename, nil
}
