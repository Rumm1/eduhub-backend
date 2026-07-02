package audit

import (
	"bytes"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func BuildAuditLogsXLSX(report AuditLogsListResponse, lang string) ([]byte, error) {
	lang = normalizeAuditExportLang(lang)
	t := auditExportTranslations(lang)

	file := excelize.NewFile()
	sheet := file.GetSheetName(0)
	if sheet == "" {
		sheet = "Sheet1"
	}

	file.SetSheetName(sheet, tr(t, "sheet.audit_logs", "Audit Logs"))
	sheet = tr(t, "sheet.audit_logs", "Audit Logs")

	headers := []string{
		tr(t, "header.created_at", "Created at"),
		tr(t, "header.user", "User"),
		tr(t, "header.action", "Action"),
		tr(t, "header.entity_type", "Entity type"),
		tr(t, "header.entity_id", "Entity ID"),
		tr(t, "header.description", "Description"),
		tr(t, "header.ip_address", "IP address"),
		tr(t, "header.metadata", "Metadata"),
	}

	for index, header := range headers {
		cell, err := excelize.CoordinatesToCellName(index+1, 1)
		if err != nil {
			return nil, err
		}

		if err := file.SetCellValue(sheet, cell, header); err != nil {
			return nil, err
		}
	}

	for rowIndex, item := range report.Items {
		row := rowIndex + 2

		values := []string{
			item.CreatedAt,
			displayAuditUser(item),
			item.Action,
			item.EntityType,
			item.EntityID,
			item.Description,
			item.IPAddress,
			item.Metadata,
		}

		for columnIndex, value := range values {
			cell, err := excelize.CoordinatesToCellName(columnIndex+1, row)
			if err != nil {
				return nil, err
			}

			if err := file.SetCellValue(sheet, cell, value); err != nil {
				return nil, err
			}
		}
	}

	if err := styleAuditHeader(file, sheet, len(headers)); err != nil {
		return nil, err
	}

	if err := autoWidthAuditColumns(file, sheet, len(headers)); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheet, "J1", tr(t, "summary.total", "Total")); err != nil {
		return nil, err
	}

	if err := file.SetCellValue(sheet, "K1", strconv.Itoa(report.Total)); err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	if err := file.Write(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func displayAuditUser(item AuditLogResponse) string {
	if item.UserName != "" {
		return item.UserName
	}

	return item.UserID
}

func normalizeAuditExportLang(lang string) string {
	switch lang {
	case "ru", "en", "kk":
		return lang
	default:
		return "ru"
	}
}

func auditExportTranslations(lang string) map[string]string {
	translations := map[string]map[string]string{
		"ru": {
			"sheet.audit_logs":   "Журнал аудита",
			"header.created_at":  "Дата и время",
			"header.user":        "Пользователь",
			"header.action":      "Действие",
			"header.entity_type": "Тип объекта",
			"header.entity_id":   "ID объекта",
			"header.description": "Описание",
			"header.ip_address":  "IP адрес",
			"header.metadata":    "Метаданные",
			"summary.total":      "Всего",
		},
		"en": {
			"sheet.audit_logs":   "Audit Logs",
			"header.created_at":  "Created at",
			"header.user":        "User",
			"header.action":      "Action",
			"header.entity_type": "Entity type",
			"header.entity_id":   "Entity ID",
			"header.description": "Description",
			"header.ip_address":  "IP address",
			"header.metadata":    "Metadata",
			"summary.total":      "Total",
		},
		"kk": {
			"sheet.audit_logs":   "Аудит журналы",
			"header.created_at":  "К?ні мен уа?ыты",
			"header.user":        "Пайдаланушы",
			"header.action":      "?рекет",
			"header.entity_type": "Объект т?рі",
			"header.entity_id":   "Объект ID",
			"header.description": "Сипаттама",
			"header.ip_address":  "IP мекенжай",
			"header.metadata":    "Метадеректер",
			"summary.total":      "Барлы?ы",
		},
	}

	return translations[normalizeAuditExportLang(lang)]
}

func tr(translations map[string]string, key string, fallback string) string {
	if value, ok := translations[key]; ok && value != "" {
		return value
	}

	return fallback
}

func styleAuditHeader(file *excelize.File, sheet string, headerCount int) error {
	styleID, err := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return err
	}

	startCell, err := excelize.CoordinatesToCellName(1, 1)
	if err != nil {
		return err
	}

	endCell, err := excelize.CoordinatesToCellName(headerCount, 1)
	if err != nil {
		return err
	}

	return file.SetCellStyle(sheet, startCell, endCell, styleID)
}

func autoWidthAuditColumns(file *excelize.File, sheet string, headerCount int) error {
	for columnIndex := 1; columnIndex <= headerCount; columnIndex++ {
		columnName, err := excelize.ColumnNumberToName(columnIndex)
		if err != nil {
			return err
		}

		width := 22.0
		if columnIndex == 8 {
			width = 60.0
		}

		if err := file.SetColWidth(sheet, columnName, columnName, width); err != nil {
			return err
		}
	}

	return nil
}
