package report

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func BuildTeacherScheduleXLSX(report TeacherScheduleReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizeExportLang(langRaw)
	t := exportTranslations(lang)

	file := excelize.NewFile()

	sheetName := t["teacher_schedule_sheet"]
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)

	title := fmt.Sprintf("%s: %s", t["teacher_schedule_report"], report.TeacherName)
	file.SetCellValue(sheetName, "A1", title)
	file.SetCellValue(sheetName, "A2", t["from"])
	file.SetCellValue(sheetName, "B2", report.FromDate)
	file.SetCellValue(sheetName, "C2", t["to"])
	file.SetCellValue(sheetName, "D2", report.ToDate)

	file.SetCellValue(sheetName, "A4", t["total_lessons"])
	file.SetCellValue(sheetName, "B4", report.TotalLessons)
	file.SetCellValue(sheetName, "C4", t["actual_lessons"])
	file.SetCellValue(sheetName, "D4", report.ActualLessons)
	file.SetCellValue(sheetName, "E4", t["planned_only"])
	file.SetCellValue(sheetName, "F4", report.PlannedOnlyLessons)
	file.SetCellValue(sheetName, "G4", t["substitutions"])
	file.SetCellValue(sheetName, "H4", report.Substitutions)
	file.SetCellValue(sheetName, "I4", t["actual_hours"])
	file.SetCellValue(sheetName, "J4", report.TotalActualHours)

	headers := []string{
		t["date"],
		t["start"],
		t["end"],
		t["hours"],
		t["group"],
		t["branch"],
		t["subject"],
		t["topic"],
		t["status"],
		t["planned_teacher"],
		t["actual_teacher"],
		t["substitution"],
		t["role"],
		t["reason"],
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
			translateValue(lang, item.Status),
			item.PlannedTeacherName,
			item.ActualTeacherName,
			translateBool(lang, item.IsSubstitution),
			translateValue(lang, item.TeacherRoleInLesson),
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

	filename := fmt.Sprintf("teacher_schedule_%s_%s_%s_%s.xlsx", lang, report.TeacherID, report.FromDate, report.ToDate)

	return writeWorkbook(file, filename)
}

func BuildPaymentsReportXLSX(report PaymentsReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizeExportLang(langRaw)
	t := exportTranslations(lang)

	file := excelize.NewFile()

	sheetName := t["payments_sheet"]
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)

	file.SetCellValue(sheetName, "A1", t["payments_report"])
	file.SetCellValue(sheetName, "A2", t["from"])
	file.SetCellValue(sheetName, "B2", report.FromDate)
	file.SetCellValue(sheetName, "C2", t["to"])
	file.SetCellValue(sheetName, "D2", report.ToDate)

	file.SetCellValue(sheetName, "A4", t["total_payments"])
	file.SetCellValue(sheetName, "B4", report.TotalPayments)
	file.SetCellValue(sheetName, "C4", t["total_amount"])
	file.SetCellValue(sheetName, "D4", report.TotalAmount)
	file.SetCellValue(sheetName, "E4", t["paid"])
	file.SetCellValue(sheetName, "F4", report.PaidAmount)
	file.SetCellValue(sheetName, "G4", t["pending"])
	file.SetCellValue(sheetName, "H4", report.PendingAmount)
	file.SetCellValue(sheetName, "I4", t["refunded"])
	file.SetCellValue(sheetName, "J4", report.RefundedAmount)
	file.SetCellValue(sheetName, "K4", t["cancelled"])
	file.SetCellValue(sheetName, "L4", report.CancelledAmount)

	headers := []string{
		t["payment_date"],
		t["payment_period"],
		t["student"],
		t["group"],
		t["branch"],
		t["amount"],
		t["method"],
		t["status"],
		t["comment"],
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
			translateValue(lang, item.PaymentMethod),
			translateValue(lang, item.Status),
			item.Comment,
		}

		writeRow(file, sheetName, row, values)
	}

	autoWidth(file, sheetName, len(headers))

	file.SetColWidth(sheetName, "C", "E", 24)
	file.SetColWidth(sheetName, "I", "I", 42)

	styleHeader(file, sheetName, "A1", "A1")
	styleHeader(file, sheetName, "A6", "I6")

	filename := fmt.Sprintf("payments_report_%s_%s_%s.xlsx", lang, report.FromDate, report.ToDate)

	return writeWorkbook(file, filename)
}

func normalizeExportLang(langRaw string) string {
	lang := strings.ToLower(strings.TrimSpace(langRaw))

	switch lang {
	case "ru", "ru-ru", "russian", "русский":
		return "ru"
	case "kk", "kz", "kk-kz", "kazakh", "?аза?ша", "казахский":
		return "kk"
	case "en", "en-us", "en-gb", "english":
		return "en"
	default:
		return "ru"
	}
}

func exportTranslations(lang string) map[string]string {
	switch lang {
	case "kk":
		return map[string]string{
			"teacher_schedule_sheet":  "Кесте",
			"teacher_schedule_report": "О?ытушы кестесі есебі",
			"payments_sheet":          "Т?лемдер",
			"payments_report":         "Т?лемдер есебі",

			"from": "Бастап",
			"to":   "Дейін",

			"total_lessons":  "Барлы? саба?тар",
			"actual_lessons": "На?ты саба?тар",
			"planned_only":   "Тек жоспарлан?ан",
			"substitutions":  "Ауыстырулар",
			"actual_hours":   "На?ты са?аттар",

			"date":            "К?ні",
			"start":           "Басталуы",
			"end":             "Ая?талуы",
			"hours":           "Са?ат",
			"group":           "Топ",
			"branch":          "Филиал",
			"subject":         "П?н",
			"topic":           "Та?ырып",
			"status":          "Статус",
			"planned_teacher": "Жоспарлан?ан о?ытушы",
			"actual_teacher":  "На?ты о?ытушы",
			"substitution":    "Ауыстыру",
			"role":            "Р?лі",
			"reason":          "Себебі",

			"total_payments": "Барлы? т?лемдер",
			"total_amount":   "Жалпы сома",
			"paid":           "Т?ленді",
			"pending":        "К?туде",
			"refunded":       "?айтарылды",
			"cancelled":      "Болдырылмады",

			"payment_date":   "Т?лем к?ні",
			"payment_period": "Т?лем кезе?і",
			"student":        "О?ушы",
			"amount":         "Сома",
			"method":         "?діс",
			"comment":        "Пікір",
		}
	case "en":
		return map[string]string{
			"teacher_schedule_sheet":  "Teacher Schedule",
			"teacher_schedule_report": "Teacher Schedule Report",
			"payments_sheet":          "Payments",
			"payments_report":         "Payments Report",

			"from": "From",
			"to":   "To",

			"total_lessons":  "Total lessons",
			"actual_lessons": "Actual lessons",
			"planned_only":   "Planned only",
			"substitutions":  "Substitutions",
			"actual_hours":   "Actual hours",

			"date":            "Date",
			"start":           "Start",
			"end":             "End",
			"hours":           "Hours",
			"group":           "Group",
			"branch":          "Branch",
			"subject":         "Subject",
			"topic":           "Topic",
			"status":          "Status",
			"planned_teacher": "Planned Teacher",
			"actual_teacher":  "Actual Teacher",
			"substitution":    "Substitution",
			"role":            "Role",
			"reason":          "Reason",

			"total_payments": "Total payments",
			"total_amount":   "Total amount",
			"paid":           "Paid",
			"pending":        "Pending",
			"refunded":       "Refunded",
			"cancelled":      "Cancelled",

			"payment_date":   "Payment Date",
			"payment_period": "Payment Period",
			"student":        "Student",
			"amount":         "Amount",
			"method":         "Method",
			"comment":        "Comment",
		}
	default:
		return map[string]string{
			"teacher_schedule_sheet":  "Расписание",
			"teacher_schedule_report": "Отчёт по расписанию преподавателя",
			"payments_sheet":          "Платежи",
			"payments_report":         "Отчёт по платежам",

			"from": "С",
			"to":   "По",

			"total_lessons":  "Всего уроков",
			"actual_lessons": "Фактические уроки",
			"planned_only":   "Только запланированные",
			"substitutions":  "Замены",
			"actual_hours":   "Фактические часы",

			"date":            "Дата",
			"start":           "Начало",
			"end":             "Конец",
			"hours":           "Часы",
			"group":           "Группа",
			"branch":          "Филиал",
			"subject":         "Предмет",
			"topic":           "Тема",
			"status":          "Статус",
			"planned_teacher": "Плановый преподаватель",
			"actual_teacher":  "Фактический преподаватель",
			"substitution":    "Замена",
			"role":            "Роль",
			"reason":          "Причина",

			"total_payments": "Всего платежей",
			"total_amount":   "Общая сумма",
			"paid":           "Оплачено",
			"pending":        "Ожидает",
			"refunded":       "Возвращено",
			"cancelled":      "Отменено",

			"payment_date":   "Дата оплаты",
			"payment_period": "Период оплаты",
			"student":        "Ученик",
			"amount":         "Сумма",
			"method":         "Метод",
			"comment":        "Комментарий",
		}
	}
}

func translateBool(lang string, value bool) string {
	if value {
		switch lang {
		case "kk":
			return "И?"
		case "en":
			return "Yes"
		default:
			return "Да"
		}
	}

	switch lang {
	case "kk":
		return "Жо?"
	case "en":
		return "No"
	default:
		return "Нет"
	}
}

func translateValue(lang string, valueRaw string) string {
	value := strings.ToLower(strings.TrimSpace(valueRaw))

	if value == "" {
		return ""
	}

	translations := map[string]map[string]string{
		"ru": {
			"paid":          "Оплачено",
			"pending":       "Ожидает",
			"cancelled":     "Отменено",
			"refunded":      "Возвращено",
			"scheduled":     "Запланирован",
			"completed":     "Проведён",
			"missed":        "Пропущен",
			"actual":        "Фактический преподаватель",
			"planned_only":  "Только плановый преподаватель",
			"cash":          "Наличные",
			"card":          "Карта",
			"kaspi":         "Kaspi",
			"bank_transfer": "Банковский перевод",
		},
		"kk": {
			"paid":          "Т?ленді",
			"pending":       "К?туде",
			"cancelled":     "Болдырылмады",
			"refunded":      "?айтарылды",
			"scheduled":     "Жоспарлан?ан",
			"completed":     "?ткізілді",
			"missed":        "?ткізілмеді",
			"actual":        "На?ты о?ытушы",
			"planned_only":  "Тек жоспарлан?ан о?ытушы",
			"cash":          "?олма-?ол",
			"card":          "Карта",
			"kaspi":         "Kaspi",
			"bank_transfer": "Банк аударымы",
		},
		"en": {
			"paid":          "Paid",
			"pending":       "Pending",
			"cancelled":     "Cancelled",
			"refunded":      "Refunded",
			"scheduled":     "Scheduled",
			"completed":     "Completed",
			"missed":        "Missed",
			"actual":        "Actual teacher",
			"planned_only":  "Planned teacher only",
			"cash":          "Cash",
			"card":          "Card",
			"kaspi":         "Kaspi",
			"bank_transfer": "Bank transfer",
		},
	}

	if langTranslations, ok := translations[lang]; ok {
		if translated, ok := langTranslations[value]; ok {
			return translated
		}
	}

	return valueRaw
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
