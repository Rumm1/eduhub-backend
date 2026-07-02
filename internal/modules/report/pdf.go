package report

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/go-pdf/fpdf"
)

const pdfContentType = "application/pdf"

type pdfColumn struct {
	Title string
	Width float64
}

var pdfTranslations = map[string]map[string]string{
	"en": {
		"teacher_schedule_report": "Teacher Schedule Report",
		"payments_report":         "Payments Report",
		"student_balances_report": "Student Balances Report",
		"payroll_report":          "Payroll Report",

		"teacher":             "Teacher",
		"teacher_id":          "Teacher ID",
		"period":              "Period",
		"total_lessons":       "Total lessons",
		"actual_lessons":      "Actual lessons",
		"planned_only":        "Planned only lessons",
		"substitutions":       "Substitutions",
		"total_actual_hours":  "Total actual hours",
		"date":                "Date",
		"time":                "Time",
		"group":               "Group",
		"subject":             "Subject",
		"topic":               "Topic",
		"status":              "Status",
		"role":                "Role",
		"subst":               "Subst.",
		"payment_period":      "Payment period",
		"total_payments":      "Total payments",
		"total_amount":        "Total amount",
		"paid_amount":         "Paid amount",
		"pending_amount":      "Pending amount",
		"refunded_amount":     "Refunded amount",
		"cancelled_amount":    "Cancelled amount",
		"student":             "Student",
		"branch":              "Branch",
		"amount":              "Amount",
		"method":              "Method",
		"total_students":      "Total students",
		"paid_count":          "Paid count",
		"partial_count":       "Partial count",
		"unpaid_count":        "Unpaid count",
		"expected_amount":     "Expected amount",
		"debt_amount":         "Debt amount",
		"monthly":             "Monthly",
		"paid":                "Paid",
		"debt":                "Debt",
		"total_entries":       "Total entries",
		"total_substitutions": "Total substitutions",
		"total_hours":         "Total hours",
		"base_amount":         "Base amount",
		"bonus_amount":        "Bonus amount",
		"penalty_amount":      "Penalty amount",
		"correction_amount":   "Correction amount",
		"final_amount":        "Final amount",
		"lessons":             "Lessons",
		"hours":               "Hours",
		"base":                "Base",
		"bonus":               "Bonus",
		"penalty":             "Penalty",
		"final":               "Final",
	},
	"ru": {
		"teacher_schedule_report": "Отчёт по расписанию преподавателя",
		"payments_report":         "Отчёт по оплатам",
		"student_balances_report": "Отчёт по балансам студентов",
		"payroll_report":          "Отчёт по зарплате",

		"teacher":             "Преподаватель",
		"teacher_id":          "ID преподавателя",
		"period":              "Период",
		"total_lessons":       "Всего уроков",
		"actual_lessons":      "Проведено уроков",
		"planned_only":        "Только запланированные",
		"substitutions":       "Замены",
		"total_actual_hours":  "Всего часов",
		"date":                "Дата",
		"time":                "Время",
		"group":               "Группа",
		"subject":             "Предмет",
		"topic":               "Тема",
		"status":              "Статус",
		"role":                "Роль",
		"subst":               "Замена",
		"payment_period":      "Период оплаты",
		"total_payments":      "Всего платежей",
		"total_amount":        "Общая сумма",
		"paid_amount":         "Оплачено",
		"pending_amount":      "Ожидает оплаты",
		"refunded_amount":     "Возврат",
		"cancelled_amount":    "Отменено",
		"student":             "Студент",
		"branch":              "Филиал",
		"amount":              "Сумма",
		"method":              "Метод",
		"total_students":      "Всего студентов",
		"paid_count":          "Оплатили",
		"partial_count":       "Частично",
		"unpaid_count":        "Не оплатили",
		"expected_amount":     "Ожидаемая сумма",
		"debt_amount":         "Сумма долга",
		"monthly":             "Месячная цена",
		"paid":                "Оплачено",
		"debt":                "Долг",
		"total_entries":       "Всего записей",
		"total_substitutions": "Всего замен",
		"total_hours":         "Всего часов",
		"base_amount":         "Базовая сумма",
		"bonus_amount":        "Бонус",
		"penalty_amount":      "Штраф",
		"correction_amount":   "Коррекция",
		"final_amount":        "Итоговая сумма",
		"lessons":             "Уроки",
		"hours":               "Часы",
		"base":                "База",
		"bonus":               "Бонус",
		"penalty":             "Штраф",
		"final":               "Итог",
	},
	"kk": {
		"teacher_schedule_report": "Мұғалім кестесі бойынша есеп",
		"payments_report":         "Төлемдер есебі",
		"student_balances_report": "Студенттер балансы бойынша есеп",
		"payroll_report":          "Жалақы есебі",

		"teacher":             "Мұғалім",
		"teacher_id":          "Мұғалім ID",
		"period":              "Кезең",
		"total_lessons":       "Барлық сабақтар",
		"actual_lessons":      "Өткізілген сабақтар",
		"planned_only":        "Тек жоспарланған сабақтар",
		"substitutions":       "Ауыстырулар",
		"total_actual_hours":  "Жалпы сағат",
		"date":                "Күні",
		"time":                "Уақыты",
		"group":               "Топ",
		"subject":             "Пән",
		"topic":               "Тақырып",
		"status":              "Статус",
		"role":                "Рөл",
		"subst":               "Ауыст.",
		"payment_period":      "Төлем кезеңі",
		"total_payments":      "Барлық төлемдер",
		"total_amount":        "Жалпы сома",
		"paid_amount":         "Төленді",
		"pending_amount":      "Күтілуде",
		"refunded_amount":     "Қайтарылды",
		"cancelled_amount":    "Болдырылмады",
		"student":             "Студент",
		"branch":              "Филиал",
		"amount":              "Сома",
		"method":              "Әдіс",
		"total_students":      "Барлық студенттер",
		"paid_count":          "Төледі",
		"partial_count":       "Ішінара",
		"unpaid_count":        "Төлемеген",
		"expected_amount":     "Күтілетін сома",
		"debt_amount":         "Қарыз сомасы",
		"monthly":             "Айлық баға",
		"paid":                "Төленді",
		"debt":                "Қарыз",
		"total_entries":       "Барлық жазбалар",
		"total_substitutions": "Барлық ауыстырулар",
		"total_hours":         "Барлық сағат",
		"base_amount":         "Негізгі сома",
		"bonus_amount":        "Бонус",
		"penalty_amount":      "Айыппұл",
		"correction_amount":   "Түзету",
		"final_amount":        "Қорытынды сома",
		"lessons":             "Сабақтар",
		"hours":               "Сағат",
		"base":                "Негізгі",
		"bonus":               "Бонус",
		"penalty":             "Айыппұл",
		"final":               "Қорытынды",
	},
}

func BuildTeacherSchedulePDF(report TeacherScheduleReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)
	pdf, fontFamily := newReportPDF(pdfText(lang, "teacher_schedule_report"))

	writeKeyValue(pdf, fontFamily, pdfText(lang, "teacher"), report.TeacherName)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "teacher_id"), report.TeacherID)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "period"), report.FromDate+" - "+report.ToDate)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_lessons"), report.TotalLessons)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "actual_lessons"), report.ActualLessons)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "planned_only"), report.PlannedOnlyLessons)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "substitutions"), report.Substitutions)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_actual_hours"), report.TotalActualHours)
	pdf.Ln(4)

	columns := []pdfColumn{
		{pdfText(lang, "date"), 24},
		{pdfText(lang, "time"), 30},
		{pdfText(lang, "group"), 42},
		{pdfText(lang, "subject"), 36},
		{pdfText(lang, "topic"), 45},
		{pdfText(lang, "status"), 25},
		{pdfText(lang, "role"), 28},
		{pdfText(lang, "subst"), 18},
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

	writeTable(pdf, fontFamily, columns, rows)

	filename := fmt.Sprintf("teacher_schedule_%s_%s_%s_%s.pdf", lang, safePDFName(report.TeacherID), report.FromDate, report.ToDate)

	return outputPDF(pdf, filename)
}

func BuildPaymentsReportPDF(report PaymentsReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)
	pdf, fontFamily := newReportPDF(pdfText(lang, "payments_report"))

	writeKeyValue(pdf, fontFamily, pdfText(lang, "period"), report.FromDate+" - "+report.ToDate)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_payments"), report.TotalPayments)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_amount"), report.TotalAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "paid_amount"), report.PaidAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "pending_amount"), report.PendingAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "refunded_amount"), report.RefundedAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "cancelled_amount"), report.CancelledAmount)
	pdf.Ln(4)

	columns := []pdfColumn{
		{pdfText(lang, "date"), 24},
		{pdfText(lang, "payment_period"), 30},
		{pdfText(lang, "student"), 45},
		{pdfText(lang, "group"), 42},
		{pdfText(lang, "branch"), 38},
		{pdfText(lang, "amount"), 28},
		{pdfText(lang, "method"), 25},
		{pdfText(lang, "status"), 25},
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

	writeTable(pdf, fontFamily, columns, rows)

	filename := fmt.Sprintf("payments_report_%s_%s_%s.pdf", lang, report.FromDate, report.ToDate)

	return outputPDF(pdf, filename)
}

func BuildStudentBalancesReportPDF(report StudentBalancesReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)
	pdf, fontFamily := newReportPDF(pdfText(lang, "student_balances_report"))

	writeKeyValue(pdf, fontFamily, pdfText(lang, "period"), report.Period)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_students"), report.TotalStudents)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "paid_count"), report.PaidCount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "partial_count"), report.PartialCount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "unpaid_count"), report.UnpaidCount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "expected_amount"), report.TotalExpectedAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "paid_amount"), report.TotalPaidAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "debt_amount"), report.TotalDebtAmount)
	pdf.Ln(4)

	columns := []pdfColumn{
		{pdfText(lang, "student"), 55},
		{pdfText(lang, "group"), 50},
		{pdfText(lang, "branch"), 45},
		{pdfText(lang, "monthly"), 30},
		{pdfText(lang, "paid"), 30},
		{pdfText(lang, "debt"), 30},
		{pdfText(lang, "status"), 28},
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

	writeTable(pdf, fontFamily, columns, rows)

	filename := fmt.Sprintf("student_balances_%s_%s.pdf", lang, report.Period)

	return outputPDF(pdf, filename)
}

func BuildPayrollReportPDF(report PayrollReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizePDFLang(langRaw)
	pdf, fontFamily := newReportPDF(pdfText(lang, "payroll_report"))

	writeKeyValue(pdf, fontFamily, pdfText(lang, "period"), report.Period)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_entries"), report.TotalEntries)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_lessons"), report.TotalLessons)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_substitutions"), report.TotalSubstitutions)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "total_hours"), report.TotalHours)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "base_amount"), report.TotalBaseAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "bonus_amount"), report.TotalBonusAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "penalty_amount"), report.TotalPenaltyAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "correction_amount"), report.TotalCorrectionAmount)
	writeKeyValue(pdf, fontFamily, pdfText(lang, "final_amount"), report.TotalFinalAmount)
	pdf.Ln(4)

	columns := []pdfColumn{
		{pdfText(lang, "teacher"), 50},
		{pdfText(lang, "lessons"), 22},
		{pdfText(lang, "subst"), 22},
		{pdfText(lang, "hours"), 25},
		{pdfText(lang, "base"), 30},
		{pdfText(lang, "bonus"), 30},
		{pdfText(lang, "penalty"), 30},
		{pdfText(lang, "final"), 30},
		{pdfText(lang, "status"), 33},
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

	writeTable(pdf, fontFamily, columns, rows)

	filename := fmt.Sprintf("payroll_report_%s_%s.pdf", lang, report.Period)

	return outputPDF(pdf, filename)
}

func newReportPDF(title string) (*fpdf.Fpdf, string) {
	pdf := fpdf.New("L", "mm", "A4", "")
	pdf.SetMargins(10, 10, 10)
	pdf.SetAutoPageBreak(true, 12)

	fontFamily := configurePDFFont(pdf)

	pdf.AddPage()

	pdf.SetFont(fontFamily, "B", 16)
	pdf.CellFormat(0, 8, title, "", 1, "L", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont(fontFamily, "", 9)

	return pdf, fontFamily
}

func configurePDFFont(pdf *fpdf.Fpdf) string {
	regularFontPath := firstExistingPath(
		os.Getenv("PDF_FONT_PATH"),
		"C:/Windows/Fonts/arial.ttf",
		"C:/Windows/Fonts/tahoma.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/dejavu/DejaVuSans.ttf",
	)

	boldFontPath := firstExistingPath(
		os.Getenv("PDF_FONT_BOLD_PATH"),
		"C:/Windows/Fonts/arialbd.ttf",
		"C:/Windows/Fonts/tahomabd.ttf",
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/dejavu/DejaVuSans-Bold.ttf",
	)

	if regularFontPath == "" {
		return "Arial"
	}

	if boldFontPath == "" {
		boldFontPath = regularFontPath
	}

	pdf.AddUTF8Font("EduHub", "", regularFontPath)
	pdf.AddUTF8Font("EduHub", "B", boldFontPath)

	return "EduHub"
}

func firstExistingPath(paths ...string) string {
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

func writeKeyValue(pdf *fpdf.Fpdf, fontFamily string, label string, value interface{}) {
	pdf.SetFont(fontFamily, "B", 9)
	pdf.CellFormat(55, 6, label+":", "", 0, "L", false, 0, "")
	pdf.SetFont(fontFamily, "", 9)
	pdf.CellFormat(0, 6, trimPDFText(fmt.Sprint(value), 140), "", 1, "L", false, 0, "")
}

func writeTable(pdf *fpdf.Fpdf, fontFamily string, columns []pdfColumn, rows [][]string) {
	writeTableHeader(pdf, fontFamily, columns)

	pdf.SetFont(fontFamily, "", 7)

	for _, row := range rows {
		if pdf.GetY() > 185 {
			pdf.AddPage()
			writeTableHeader(pdf, fontFamily, columns)
			pdf.SetFont(fontFamily, "", 7)
		}

		for index, column := range columns {
			value := ""
			if index < len(row) {
				value = row[index]
			}

			pdf.CellFormat(column.Width, 7, trimPDFText(value, int(column.Width/2.2)), "1", 0, "L", false, 0, "")
		}

		pdf.Ln(-1)
	}
}

func writeTableHeader(pdf *fpdf.Fpdf, fontFamily string, columns []pdfColumn) {
	pdf.SetFont(fontFamily, "B", 7)

	for _, column := range columns {
		pdf.CellFormat(column.Width, 7, column.Title, "1", 0, "L", false, 0, "")
	}

	pdf.Ln(-1)
}

func outputPDF(pdf *fpdf.Fpdf, filename string) ([]byte, string, error) {
	var buffer bytes.Buffer

	if err := pdf.Output(&buffer); err != nil {
		return nil, "", err
	}

	return buffer.Bytes(), filename, nil
}

func normalizePDFLang(langRaw string) string {
	lang := strings.ToLower(strings.TrimSpace(langRaw))

	switch lang {
	case "ru", "kk", "en":
		return lang
	default:
		return "en"
	}
}

func pdfText(lang string, key string) string {
	if values, ok := pdfTranslations[lang]; ok {
		if value, ok := values[key]; ok {
			return value
		}
	}

	if values, ok := pdfTranslations["en"]; ok {
		if value, ok := values[key]; ok {
			return value
		}
	}

	return key
}

func trimPDFText(value string, limit int) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.ReplaceAll(value, "\r", " ")

	if limit <= 0 || len([]rune(value)) <= limit {
		return value
	}

	runes := []rune(value)

	return string(runes[:limit]) + "..."
}

func safePDFName(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "report"
	}

	replacer := strings.NewReplacer(
		" ", "_",
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		`"`, "_",
		"<", "_",
		">", "_",
		"|", "_",
	)

	return replacer.Replace(value)
}
