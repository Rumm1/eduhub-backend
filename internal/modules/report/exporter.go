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

func BuildStudentBalancesReportXLSX(report StudentBalancesReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizeExportLang(langRaw)
	t := exportTranslations(lang)

	file := excelize.NewFile()

	sheetName := t["student_balances_sheet"]
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)

	file.SetCellValue(sheetName, "A1", t["student_balances_report"])
	file.SetCellValue(sheetName, "A2", t["period"])
	file.SetCellValue(sheetName, "B2", report.Period)

	file.SetCellValue(sheetName, "A4", t["total_students"])
	file.SetCellValue(sheetName, "B4", report.TotalStudents)
	file.SetCellValue(sheetName, "C4", t["paid"])
	file.SetCellValue(sheetName, "D4", report.PaidCount)
	file.SetCellValue(sheetName, "E4", t["partial"])
	file.SetCellValue(sheetName, "F4", report.PartialCount)
	file.SetCellValue(sheetName, "G4", t["unpaid"])
	file.SetCellValue(sheetName, "H4", report.UnpaidCount)

	file.SetCellValue(sheetName, "A5", t["total_expected_amount"])
	file.SetCellValue(sheetName, "B5", report.TotalExpectedAmount)
	file.SetCellValue(sheetName, "C5", t["total_paid_amount"])
	file.SetCellValue(sheetName, "D5", report.TotalPaidAmount)
	file.SetCellValue(sheetName, "E5", t["total_debt_amount"])
	file.SetCellValue(sheetName, "F5", report.TotalDebtAmount)

	headers := []string{
		t["student"],
		t["group"],
		t["branch"],
		t["monthly_price"],
		t["paid_amount"],
		t["debt_amount"],
		t["status"],
	}

	headerRow := 7
	writeHeaders(file, sheetName, headers, headerRow)

	for rowIndex, item := range report.Items {
		row := headerRow + rowIndex + 1

		values := []interface{}{
			item.StudentName,
			item.GroupName,
			item.BranchName,
			item.MonthlyPrice,
			item.PaidAmount,
			item.DebtAmount,
			translateValue(lang, item.PaymentStatus),
		}

		writeRow(file, sheetName, row, values)
	}

	autoWidth(file, sheetName, len(headers))

	file.SetColWidth(sheetName, "A", "C", 26)
	file.SetColWidth(sheetName, "D", "F", 18)

	styleHeader(file, sheetName, "A1", "A1")
	styleHeader(file, sheetName, "A7", "G7")

	filename := fmt.Sprintf("student_balances_%s_%s.xlsx", lang, report.Period)

	return writeWorkbook(file, filename)
}

func BuildPayrollReportXLSX(report PayrollReportResponse, langRaw string) ([]byte, string, error) {
	lang := normalizeExportLang(langRaw)
	t := exportTranslations(lang)
	file := excelize.NewFile()
	sheetName := tr(t, "payroll_sheet", "Payroll")
	defaultSheet := file.GetSheetName(0)
	file.SetSheetName(defaultSheet, sheetName)
	file.SetCellValue(sheetName, "A1", tr(t, "payroll_report", "Payroll Report"))
	file.SetCellValue(sheetName, "A2", tr(t, "period", "Period"))
	file.SetCellValue(sheetName, "B2", report.Period)
	file.SetCellValue(sheetName, "A4", tr(t, "total_entries", "Total entries"))
	file.SetCellValue(sheetName, "B4", report.TotalEntries)
	file.SetCellValue(sheetName, "C4", tr(t, "total_lessons", "Total lessons"))
	file.SetCellValue(sheetName, "D4", report.TotalLessons)
	file.SetCellValue(sheetName, "E4", tr(t, "substitutions", "Substitutions"))
	file.SetCellValue(sheetName, "F4", report.TotalSubstitutions)
	file.SetCellValue(sheetName, "G4", tr(t, "hours", "Hours"))
	file.SetCellValue(sheetName, "H4", report.TotalHours)
	file.SetCellValue(sheetName, "I4", tr(t, "total_final_amount", "Total final amount"))
	file.SetCellValue(sheetName, "J4", report.TotalFinalAmount)
	file.SetCellValue(sheetName, "A5", tr(t, "draft", "Draft"))
	file.SetCellValue(sheetName, "B5", report.DraftCount)
	file.SetCellValue(sheetName, "C5", tr(t, "sent_to_teacher", "Sent to teacher"))
	file.SetCellValue(sheetName, "D5", report.SentToTeacherCount)
	file.SetCellValue(sheetName, "E5", tr(t, "teacher_approved", "Teacher approved"))
	file.SetCellValue(sheetName, "F5", report.TeacherApprovedCount)
	file.SetCellValue(sheetName, "G5", tr(t, "teacher_disputed", "Teacher disputed"))
	file.SetCellValue(sheetName, "H5", report.TeacherDisputedCount)
	file.SetCellValue(sheetName, "I5", tr(t, "approved_by_finance", "Approved by finance"))
	file.SetCellValue(sheetName, "J5", report.ApprovedByFinanceCount)
	file.SetCellValue(sheetName, "K5", tr(t, "paid", "Paid"))
	file.SetCellValue(sheetName, "L5", report.PaidCount)

	headers := []string{
		tr(t, "teacher", "Teacher"),
		tr(t, "lessons", "Lessons"),
		tr(t, "substitutions", "Substitutions"),
		tr(t, "hours_worked", "Hours worked"),
		tr(t, "hourly_rate", "Hourly rate"),
		tr(t, "base_amount", "Base amount"),
		tr(t, "bonus_amount", "Bonus amount"),
		tr(t, "penalty_amount", "Penalty amount"),
		tr(t, "correction_amount", "Correction amount"),
		tr(t, "final_amount", "Final amount"),
		tr(t, "status", "Status"),
		tr(t, "teacher_confirmation", "Teacher confirmation"),
		tr(t, "dispute_reason", "Dispute reason"),
		tr(t, "comment", "Comment"),
	}

	headerRow := 7
	writeHeaders(file, sheetName, headers, headerRow)

	for rowIndex, item := range report.Items {
		row := headerRow + rowIndex + 1

		values := []interface{}{
			item.TeacherName,
			item.LessonsCount,
			item.SubstitutionCount,
			item.HoursWorked,
			item.HourlyRate,
			item.BaseAmount,
			item.BonusAmount,
			item.PenaltyAmount,
			item.CorrectionAmount,
			item.FinalAmount,
			translateValue(lang, item.Status),
			translateValue(lang, item.TeacherConfirmationStatus),
			item.TeacherDisputeReason,
			item.Comment,
		}

		writeRow(file, sheetName, row, values)
	}

	autoWidth(file, sheetName, len(headers))

	file.SetColWidth(sheetName, "A", "A", 28)
	file.SetColWidth(sheetName, "M", "N", 42)

	styleHeader(file, sheetName, "A1", "A1")
	styleHeader(file, sheetName, "A7", "N7")

	filename := fmt.Sprintf("payroll_report_%s_%s.xlsx", lang, report.Period)

	return writeWorkbook(file, filename)
}

func normalizeExportLang(langRaw string) string {
	lang := strings.ToLower(strings.TrimSpace(langRaw))

	switch lang {
	case "ru", "ru-ru", "russian", "русский":
		return "ru"
	case "kk", "kz", "kk-kz", "kazakh", "қазақша", "казахский":
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
			"payroll_sheet":           "Жалақы",
			"payroll_report":          "Жалақы есебі",
			"total_entries":           "Барлық жазбалар",
			"total_final_amount":      "Қорытынды сома",
			"teacher":                 "Оқытушы",
			"lessons":                 "Сабақтар",
			"hours_worked":            "Жұмыс істеген сағаттар",
			"hourly_rate":             "Сағаттық ақы",
			"base_amount":             "Негізгі сома",
			"bonus_amount":            "Бонус сомасы",
			"penalty_amount":          "Айыппұл сомасы",
			"correction_amount":       "Түзету сомасы",
			"final_amount":            "Қорытынды сома",
			"teacher_confirmation":    "Оқытушының растауы",
			"dispute_reason":          "Дау себебі",
			"draft":                   "Жоба",
			"sent_to_teacher":         "Оқытушыға жіберілді",
			"teacher_approved":        "Оқытушы растады",
			"teacher_disputed":        "Оқытушы даулаған",
			"approved_by_finance":     "Қаржы бөлімі мақұлдады",
			"teacher_schedule_sheet":  "Кесте",
			"teacher_schedule_report": "Оқытушы кестесі есебі",
			"payments_sheet":          "Төлемдер",
			"payments_report":         "Төлемдер есебі",

			"from": "Бастап",
			"to":   "Дейін",

			"total_lessons":  "Барлық сабақтар",
			"actual_lessons": "Нақты сабақтар",
			"planned_only":   "Тек жоспарланған",
			"substitutions":  "Ауыстырулар",
			"actual_hours":   "Нақты сағаттар",

			"date":            "Күні",
			"start":           "Басталуы",
			"end":             "Аяқталуы",
			"hours":           "Сағат",
			"group":           "Топ",
			"branch":          "Филиал",
			"subject":         "Пән",
			"topic":           "Тақырып",
			"status":          "Статус",
			"planned_teacher": "Жоспарланған мұғалім",
			"actual_teacher":  "Нақты мұғалім",
			"substitution":    "Ауыстыру",
			"role":            "Рөлі",
			"reason":          "Себебі",

			"total_payments": "Барлық төлемдер",
			"total_amount":   "Жалпы сома",
			"paid":           "Төленді",
			"pending":        "Күтуде",
			"refunded":       "Қайтарылды",
			"cancelled":      "Болдырылмады",

			"payment_date":            "Төлем күні",
			"payment_period":          "Төлем кезеңі",
			"student":                 "Оқушы",
			"amount":                  "Сомасы",
			"method":                  "Әдіс",
			"comment":                 "Пікір",
			"student_balances_sheet":  "Қарыздар",
			"student_balances_report": "Оқушылар балансы есебі",
			"period":                  "Кезең",
			"total_students":          "Барлық оқушылар",
			"partial":                 "Жартылай",
			"unpaid":                  "Төленбеген",
			"total_expected_amount":   "Күтілетін жалпы сома",
			"total_paid_amount":       "Төленген жалпы сома",
			"total_debt_amount":       "Жалпы қарыз",
			"monthly_price":           "Айлық баға",
			"paid_amount":             "Төленген сома",
			"debt_amount":             "Қарыз сомасы",
		}
	case "en":
		return map[string]string{
			"payroll_sheet":           "Payroll",
			"payroll_report":          "Payroll Report",
			"total_entries":           "Total entries",
			"total_final_amount":      "Total final amount",
			"teacher":                 "Teacher",
			"lessons":                 "Lessons",
			"hours_worked":            "Hours worked",
			"hourly_rate":             "Hourly rate",
			"base_amount":             "Base amount",
			"bonus_amount":            "Bonus amount",
			"penalty_amount":          "Penalty amount",
			"correction_amount":       "Correction amount",
			"final_amount":            "Final amount",
			"teacher_confirmation":    "Teacher confirmation",
			"dispute_reason":          "Dispute reason",
			"draft":                   "Draft",
			"sent_to_teacher":         "Sent to teacher",
			"teacher_approved":        "Teacher approved",
			"teacher_disputed":        "Teacher disputed",
			"approved_by_finance":     "Approved by finance",
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

			"payment_date":            "Payment Date",
			"payment_period":          "Payment Period",
			"student":                 "Student",
			"amount":                  "Amount",
			"method":                  "Method",
			"comment":                 "Comment",
			"student_balances_sheet":  "Balances",
			"student_balances_report": "Student Balances Report",
			"period":                  "Period",
			"total_students":          "Total students",
			"partial":                 "Partial",
			"unpaid":                  "Unpaid",
			"total_expected_amount":   "Total expected amount",
			"total_paid_amount":       "Total paid amount",
			"total_debt_amount":       "Total debt amount",
			"monthly_price":           "Monthly price",
			"paid_amount":             "Paid amount",
			"debt_amount":             "Debt amount",
		}
	default:
		return map[string]string{
			"payroll_sheet":           "Зарплата",
			"payroll_report":          "Отчёт по зарплате",
			"total_entries":           "Всего записей",
			"total_final_amount":      "Итоговая сумма",
			"teacher":                 "Преподаватель",
			"lessons":                 "Уроки",
			"hours_worked":            "Отработанные часы",
			"hourly_rate":             "Ставка за час",
			"base_amount":             "Базовая сумма",
			"bonus_amount":            "Бонусная сумма",
			"penalty_amount":          "Сумма штрафа",
			"correction_amount":       "Сумма корректировки",
			"final_amount":            "Итоговая сумма",
			"teacher_confirmation":    "Подтверждение преподавателя",
			"dispute_reason":          "Причина спора",
			"draft":                   "Черновик",
			"sent_to_teacher":         "Отправлено преподавателю",
			"teacher_approved":        "Подтверждено преподавателем",
			"teacher_disputed":        "Оспорено преподавателем",
			"approved_by_finance":     "Одобрено финансами",
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

			"payment_date":            "Дата оплаты",
			"payment_period":          "Период оплаты",
			"student":                 "Ученик",
			"amount":                  "Сумма",
			"method":                  "Метод",
			"comment":                 "Комментарий",
			"student_balances_sheet":  "Долги",
			"student_balances_report": "Отчёт по балансам учеников",
			"period":                  "Период",
			"total_students":          "Всего учеников",
			"partial":                 "Частично",
			"unpaid":                  "Не оплачено",
			"total_expected_amount":   "Ожидаемая сумма",
			"total_paid_amount":       "Оплаченная сумма",
			"total_debt_amount":       "Общий долг",
			"monthly_price":           "Месячная цена",
			"paid_amount":             "Оплачено",
			"debt_amount":             "Долг",
		}
	}
}

func translateBool(lang string, value bool) string {
	if value {
		switch lang {
		case "kk":
			return "Иә"
		case "en":
			return "Yes"
		default:
			return "Да"
		}
	}

	switch lang {
	case "kk":
		return "Жоқ"
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
			"draft":               "Черновик",
			"sent_to_teacher":     "Отправлено преподавателю",
			"teacher_approved":    "Подтверждено преподавателем",
			"teacher_disputed":    "Оспорено преподавателем",
			"approved_by_finance": "Одобрено финансами",
			"not_sent":            "Не отправлено",
			"approved":            "Подтверждено",
			"disputed":            "Оспорено",
			"partial":             "Частично",
			"unpaid":              "Не оплачено",
			"paid":                "Оплачено",
			"pending":             "Ожидает",
			"cancelled":           "Отменено",
			"refunded":            "Возвращено",
			"planned":             "Запланирован",
			"scheduled":           "Запланирован",
			"completed":           "Проведён",
			"missed":              "Пропущен",
			"actual":              "Фактический преподаватель",
			"planned_only":        "Только плановый преподаватель",
			"true":                "Да",
			"false":               "Нет",
			"cash":                "Наличные",
			"card":                "Карта",
			"kaspi":               "Kaspi",
			"bank_transfer":       "Банковский перевод",
		},
		"kk": {
			"draft":               "Жоба",
			"sent_to_teacher":     "Оқытушыға жіберілді",
			"teacher_approved":    "Оқытушы растады",
			"teacher_disputed":    "Оқытушы даулаған",
			"approved_by_finance": "Қаржы бөлімі мақұлдады",
			"not_sent":            "Жіберілмеген",
			"approved":            "Расталды",
			"disputed":            "Даулы",
			"partial":             "Жартылай",
			"unpaid":              "Төленбеген",
			"paid":                "Төленді",
			"pending":             "Күтуде",
			"cancelled":           "Болдырылмады",
			"refunded":            "Қайтарылды",
			"planned":             "Жоспарланған",
			"scheduled":           "Жоспарланған",
			"completed":           "Өткізілді",
			"missed":              "Өткізілмеді",
			"actual":              "Нақты оқытушы",
			"planned_only":        "Тек жоспарланған оқытушы",
			"true":                "Иә",
			"false":               "Жоқ",
			"cash":                "Қолма-қол",
			"card":                "Карта",
			"kaspi":               "Kaspi",
			"bank_transfer":       "Банк аударымы",
		},
		"en": {
			"draft":               "Draft",
			"sent_to_teacher":     "Sent to teacher",
			"teacher_approved":    "Teacher approved",
			"teacher_disputed":    "Teacher disputed",
			"approved_by_finance": "Approved by finance",
			"not_sent":            "Not sent",
			"approved":            "Approved",
			"disputed":            "Disputed",
			"partial":             "Partial",
			"unpaid":              "Unpaid",
			"paid":                "Paid",
			"pending":             "Pending",
			"cancelled":           "Cancelled",
			"refunded":            "Refunded",
			"planned":             "Planned",
			"scheduled":           "Scheduled",
			"completed":           "Completed",
			"missed":              "Missed",
			"actual":              "Actual teacher",
			"planned_only":        "Planned teacher only",
			"true":                "Yes",
			"false":               "No",
			"cash":                "Cash",
			"card":                "Card",
			"kaspi":               "Kaspi",
			"bank_transfer":       "Bank transfer",
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

func tr(t map[string]string, key string, fallback string) string {
	value := strings.TrimSpace(t[key])
	if value == "" {
		return fallback
	}

	return value
}
