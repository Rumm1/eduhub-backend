package report

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"
)

type docxSection struct {
	Label string
	Value interface{}
}

func BuildTeacherScheduleDOCX(report TeacherScheduleReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)

	sections := []docxSection{
		{pdfText(lang, "teacher"), report.TeacherName},
		{pdfText(lang, "teacher_id"), report.TeacherID},
		{pdfText(lang, "period"), report.FromDate + " - " + report.ToDate},
		{pdfText(lang, "total_lessons"), report.TotalLessons},
		{pdfText(lang, "actual_lessons"), report.ActualLessons},
		{pdfText(lang, "planned_only"), report.PlannedOnlyLessons},
		{pdfText(lang, "substitutions"), report.Substitutions},
		{pdfText(lang, "total_actual_hours"), report.TotalActualHours},
	}

	headers := []string{
		pdfText(lang, "date"),
		pdfText(lang, "time"),
		pdfText(lang, "group"),
		pdfText(lang, "subject"),
		pdfText(lang, "topic"),
		pdfText(lang, "status"),
		pdfText(lang, "role"),
		pdfText(lang, "subst"),
	}

	rows := make([][]string, 0, len(report.Items))
	for _, item := range report.Items {
		rows = append(rows, []string{
			item.LessonDate,
			item.StartTime + "-" + item.EndTime,
			item.GroupName,
			item.SubjectName,
			item.Topic,
			item.Status,
			item.TeacherRoleInLesson,
			fmt.Sprint(item.IsSubstitution),
		})
	}

	filename := fmt.Sprintf("teacher_schedule_%s_%s_%s_%s.docx", lang, safePDFName(report.TeacherID), report.FromDate, report.ToDate)

	return buildDOCXReport(pdfText(lang, "teacher_schedule_report"), sections, headers, rows, filename)
}

func BuildPaymentsReportDOCX(report PaymentsReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)

	sections := []docxSection{
		{pdfText(lang, "period"), report.FromDate + " - " + report.ToDate},
		{pdfText(lang, "total_payments"), report.TotalPayments},
		{pdfText(lang, "total_amount"), report.TotalAmount},
		{pdfText(lang, "paid_amount"), report.PaidAmount},
		{pdfText(lang, "pending_amount"), report.PendingAmount},
		{pdfText(lang, "refunded_amount"), report.RefundedAmount},
		{pdfText(lang, "cancelled_amount"), report.CancelledAmount},
	}

	headers := []string{
		pdfText(lang, "date"),
		pdfText(lang, "payment_period"),
		pdfText(lang, "student"),
		pdfText(lang, "group"),
		pdfText(lang, "branch"),
		pdfText(lang, "amount"),
		pdfText(lang, "method"),
		pdfText(lang, "status"),
	}

	rows := make([][]string, 0, len(report.Items))
	for _, item := range report.Items {
		rows = append(rows, []string{
			item.PaymentDate,
			item.PaymentPeriod,
			item.StudentName,
			item.GroupName,
			item.BranchName,
			fmt.Sprint(item.Amount),
			item.PaymentMethod,
			item.Status,
		})
	}

	filename := fmt.Sprintf("payments_report_%s_%s_%s.docx", lang, report.FromDate, report.ToDate)

	return buildDOCXReport(pdfText(lang, "payments_report"), sections, headers, rows, filename)
}

func BuildStudentBalancesReportDOCX(report StudentBalancesReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)

	sections := []docxSection{
		{pdfText(lang, "period"), report.Period},
		{pdfText(lang, "total_students"), report.TotalStudents},
		{pdfText(lang, "paid_count"), report.PaidCount},
		{pdfText(lang, "partial_count"), report.PartialCount},
		{pdfText(lang, "unpaid_count"), report.UnpaidCount},
		{pdfText(lang, "expected_amount"), report.TotalExpectedAmount},
		{pdfText(lang, "paid_amount"), report.TotalPaidAmount},
		{pdfText(lang, "debt_amount"), report.TotalDebtAmount},
	}

	headers := []string{
		pdfText(lang, "student"),
		pdfText(lang, "group"),
		pdfText(lang, "branch"),
		pdfText(lang, "monthly"),
		pdfText(lang, "paid"),
		pdfText(lang, "debt"),
		pdfText(lang, "status"),
	}

	rows := make([][]string, 0, len(report.Items))
	for _, item := range report.Items {
		rows = append(rows, []string{
			item.StudentName,
			item.GroupName,
			item.BranchName,
			fmt.Sprint(item.MonthlyPrice),
			fmt.Sprint(item.PaidAmount),
			fmt.Sprint(item.DebtAmount),
			item.PaymentStatus,
		})
	}

	filename := fmt.Sprintf("student_balances_%s_%s.docx", lang, report.Period)

	return buildDOCXReport(pdfText(lang, "student_balances_report"), sections, headers, rows, filename)
}

func BuildPayrollReportDOCX(report PayrollReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)

	sections := []docxSection{
		{pdfText(lang, "period"), report.Period},
		{pdfText(lang, "total_entries"), report.TotalEntries},
		{pdfText(lang, "total_lessons"), report.TotalLessons},
		{pdfText(lang, "total_substitutions"), report.TotalSubstitutions},
		{pdfText(lang, "total_hours"), report.TotalHours},
		{pdfText(lang, "base_amount"), report.TotalBaseAmount},
		{pdfText(lang, "bonus_amount"), report.TotalBonusAmount},
		{pdfText(lang, "penalty_amount"), report.TotalPenaltyAmount},
		{pdfText(lang, "correction_amount"), report.TotalCorrectionAmount},
		{pdfText(lang, "final_amount"), report.TotalFinalAmount},
	}

	headers := []string{
		pdfText(lang, "teacher"),
		pdfText(lang, "lessons"),
		pdfText(lang, "subst"),
		pdfText(lang, "hours"),
		pdfText(lang, "base"),
		pdfText(lang, "bonus"),
		pdfText(lang, "penalty"),
		pdfText(lang, "final"),
		pdfText(lang, "status"),
	}

	rows := make([][]string, 0, len(report.Items))
	for _, item := range report.Items {
		rows = append(rows, []string{
			item.TeacherName,
			fmt.Sprint(item.LessonsCount),
			fmt.Sprint(item.SubstitutionCount),
			fmt.Sprint(item.HoursWorked),
			fmt.Sprint(item.BaseAmount),
			fmt.Sprint(item.BonusAmount),
			fmt.Sprint(item.PenaltyAmount),
			fmt.Sprint(item.FinalAmount),
			item.Status,
		})
	}

	filename := fmt.Sprintf("payroll_report_%s_%s.docx", lang, report.Period)

	return buildDOCXReport(pdfText(lang, "payroll_report"), sections, headers, rows, filename)
}

func buildDOCXReport(title string, sections []docxSection, headers []string, rows [][]string, filename string) ([]byte, string, error) {
	documentXML := buildDocumentXML(title, sections, headers, rows)

	var buffer bytes.Buffer
	zipWriter := zip.NewWriter(&buffer)

	files := map[string]string{
		"[Content_Types].xml": docxContentTypesXML(),
		"_rels/.rels":         docxRootRelsXML(),
		"word/document.xml":   documentXML,
	}

	for path, content := range files {
		writer, err := zipWriter.Create(path)
		if err != nil {
			_ = zipWriter.Close()
			return nil, "", err
		}

		if _, err := writer.Write([]byte(content)); err != nil {
			_ = zipWriter.Close()
			return nil, "", err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, "", err
	}

	return buffer.Bytes(), filename, nil
}

func buildDocumentXML(title string, sections []docxSection, headers []string, rows [][]string) string {
	var builder strings.Builder

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">`)
	builder.WriteString(`<w:body>`)

	builder.WriteString(docxParagraph(title, true, 32))

	for _, section := range sections {
		builder.WriteString(docxParagraph(section.Label+": "+fmt.Sprint(section.Value), false, 20))
	}

	builder.WriteString(docxParagraph("", false, 20))
	builder.WriteString(docxTable(headers, rows))

	builder.WriteString(`<w:sectPr>`)
	builder.WriteString(`<w:pgSz w:w="16838" w:h="11906" w:orient="landscape"/>`)
	builder.WriteString(`<w:pgMar w:top="720" w:right="720" w:bottom="720" w:left="720" w:header="360" w:footer="360" w:gutter="0"/>`)
	builder.WriteString(`</w:sectPr>`)

	builder.WriteString(`</w:body>`)
	builder.WriteString(`</w:document>`)

	return builder.String()
}

func docxTable(headers []string, rows [][]string) string {
	var builder strings.Builder

	builder.WriteString(`<w:tbl>`)
	builder.WriteString(`<w:tblPr>`)
	builder.WriteString(`<w:tblW w:w="0" w:type="auto"/>`)
	builder.WriteString(`<w:tblBorders>`)
	builder.WriteString(`<w:top w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`<w:left w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`<w:bottom w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`<w:right w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`<w:insideH w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`<w:insideV w:val="single" w:sz="4" w:space="0" w:color="999999"/>`)
	builder.WriteString(`</w:tblBorders>`)
	builder.WriteString(`</w:tblPr>`)

	builder.WriteString(`<w:tr>`)
	for _, header := range headers {
		builder.WriteString(docxCell(header, true))
	}
	builder.WriteString(`</w:tr>`)

	for _, row := range rows {
		builder.WriteString(`<w:tr>`)

		for index := range headers {
			value := ""
			if index < len(row) {
				value = row[index]
			}

			builder.WriteString(docxCell(value, false))
		}

		builder.WriteString(`</w:tr>`)
	}

	builder.WriteString(`</w:tbl>`)

	return builder.String()
}

func docxCell(text string, bold bool) string {
	var builder strings.Builder

	builder.WriteString(`<w:tc>`)
	builder.WriteString(`<w:tcPr><w:tcW w:w="2400" w:type="dxa"/>`)

	if bold {
		builder.WriteString(`<w:shd w:fill="EAF2F8"/>`)
	}

	builder.WriteString(`</w:tcPr>`)
	builder.WriteString(docxParagraph(text, bold, 18))
	builder.WriteString(`</w:tc>`)

	return builder.String()
}

func docxParagraph(text string, bold bool, size int) string {
	var builder strings.Builder

	builder.WriteString(`<w:p>`)
	builder.WriteString(`<w:r>`)
	builder.WriteString(`<w:rPr>`)
	builder.WriteString(`<w:rFonts w:ascii="Arial" w:hAnsi="Arial" w:eastAsia="Arial" w:cs="Arial"/>`)

	if bold {
		builder.WriteString(`<w:b/>`)
		builder.WriteString(`<w:bCs/>`)
	}

	builder.WriteString(fmt.Sprintf(`<w:sz w:val="%d"/>`, size))
	builder.WriteString(fmt.Sprintf(`<w:szCs w:val="%d"/>`, size))
	builder.WriteString(`</w:rPr>`)
	builder.WriteString(`<w:t xml:space="preserve">`)
	builder.WriteString(escapeXML(sanitizeXMLText(text)))
	builder.WriteString(`</w:t>`)
	builder.WriteString(`</w:r>`)
	builder.WriteString(`</w:p>`)

	return builder.String()
}

func docxContentTypesXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
<Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
<Default Extension="xml" ContentType="application/xml"/>
<Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`
}

func docxRootRelsXML() string {
	return `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`
}

func escapeXML(value string) string {
	var buffer bytes.Buffer

	if err := xml.EscapeText(&buffer, []byte(value)); err != nil {
		return ""
	}

	return buffer.String()
}

func sanitizeXMLText(value string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == 0x09:
			return r
		case r == 0x0A:
			return ' '
		case r == 0x0D:
			return ' '
		case r >= 0x20 && r <= 0xD7FF:
			return r
		case r >= 0xE000 && r <= 0xFFFD:
			return r
		case r >= 0x10000 && r <= 0x10FFFF:
			return r
		default:
			return -1
		}
	}, value)
}
