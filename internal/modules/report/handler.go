package report

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetTeacherSchedule(w http.ResponseWriter, r *http.Request) {
	teacherID := r.URL.Query().Get("teacher_id")
	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")
	format := normalizeReportFormat(r.URL.Query().Get("format"))
	lang := r.URL.Query().Get("lang")

	result, err := h.service.GetTeacherSchedule(r.Context(), teacherID, fromDate, toDate)
	if err != nil {
		writeReportError(w, err)
		return
	}

	switch format {
	case "xlsx":
		fileBytes, filename, err := BuildTeacherScheduleXLSX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, xlsxContentType)
		return
	case "pdf":
		fileBytes, filename, err := BuildTeacherSchedulePDF(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, pdfContentType)
		return
	case "docx":
		fileBytes, filename, err := BuildTeacherScheduleDOCX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, docxContentType)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetPaymentsReport(w http.ResponseWriter, r *http.Request) {
	fromDate := r.URL.Query().Get("from_date")
	toDate := r.URL.Query().Get("to_date")
	branchID := r.URL.Query().Get("branch_id")
	groupID := r.URL.Query().Get("group_id")
	studentID := r.URL.Query().Get("student_id")
	status := r.URL.Query().Get("status")
	format := normalizeReportFormat(r.URL.Query().Get("format"))
	lang := r.URL.Query().Get("lang")

	result, err := h.service.GetPaymentsReport(
		r.Context(),
		fromDate,
		toDate,
		branchID,
		groupID,
		studentID,
		status,
	)
	if err != nil {
		writeReportError(w, err)
		return
	}

	switch format {
	case "xlsx":
		fileBytes, filename, err := BuildPaymentsReportXLSX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, xlsxContentType)
		return
	case "pdf":
		fileBytes, filename, err := BuildPaymentsReportPDF(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, pdfContentType)
		return
	case "docx":
		fileBytes, filename, err := BuildPaymentsReportDOCX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, docxContentType)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetStudentBalancesReport(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	branchID := r.URL.Query().Get("branch_id")
	groupID := r.URL.Query().Get("group_id")
	studentID := r.URL.Query().Get("student_id")
	status := r.URL.Query().Get("status")
	format := normalizeReportFormat(r.URL.Query().Get("format"))
	lang := r.URL.Query().Get("lang")

	result, err := h.service.GetStudentBalancesReport(
		r.Context(),
		period,
		branchID,
		groupID,
		studentID,
		status,
	)
	if err != nil {
		writeReportError(w, err)
		return
	}

	switch format {
	case "xlsx":
		fileBytes, filename, err := BuildStudentBalancesReportXLSX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, xlsxContentType)
		return
	case "pdf":
		fileBytes, filename, err := BuildStudentBalancesReportPDF(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, pdfContentType)
		return
	case "docx":
		fileBytes, filename, err := BuildStudentBalancesReportDOCX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, docxContentType)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetPayrollReport(w http.ResponseWriter, r *http.Request) {
	period := r.URL.Query().Get("period")
	teacherID := r.URL.Query().Get("teacher_id")
	status := r.URL.Query().Get("status")
	teacherConfirmationStatus := r.URL.Query().Get("teacher_confirmation_status")
	format := normalizeReportFormat(r.URL.Query().Get("format"))
	lang := r.URL.Query().Get("lang")

	result, err := h.service.GetPayrollReport(
		r.Context(),
		period,
		teacherID,
		status,
		teacherConfirmationStatus,
	)
	if err != nil {
		writeReportError(w, err)
		return
	}

	switch format {
	case "xlsx":
		fileBytes, filename, err := BuildPayrollReportXLSX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, xlsxContentType)
		return
	case "pdf":
		fileBytes, filename, err := BuildPayrollReportPDF(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, pdfContentType)
		return
	case "docx":
		fileBytes, filename, err := BuildPayrollReportDOCX(result, lang)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		writeBinaryFile(w, fileBytes, filename, docxContentType)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func normalizeReportFormat(formatRaw string) string {
	return strings.ToLower(strings.TrimSpace(formatRaw))
}

func writeBinaryFile(w http.ResponseWriter, fileBytes []byte, filename string, contentType string) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(fileBytes)
}

func writeReportError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrTeacherIDRequired):
		response.Error(w, http.StatusBadRequest, "TEACHER_ID_REQUIRED", "Teacher id is required")
	case errors.Is(err, ErrTeacherIDInvalid):
		response.Error(w, http.StatusBadRequest, "TEACHER_ID_INVALID", "Teacher id is invalid")
	case errors.Is(err, ErrTeacherNotFound):
		response.Error(w, http.StatusNotFound, "TEACHER_NOT_FOUND", "Teacher not found in organization")
	case errors.Is(err, ErrFromDateRequired):
		response.Error(w, http.StatusBadRequest, "FROM_DATE_REQUIRED", "From date is required")
	case errors.Is(err, ErrToDateRequired):
		response.Error(w, http.StatusBadRequest, "TO_DATE_REQUIRED", "To date is required")
	case errors.Is(err, ErrFromDateInvalid):
		response.Error(w, http.StatusBadRequest, "FROM_DATE_INVALID", "From date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrToDateInvalid):
		response.Error(w, http.StatusBadRequest, "TO_DATE_INVALID", "To date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrDateRangeInvalid):
		response.Error(w, http.StatusBadRequest, "DATE_RANGE_INVALID", "To date must be after or equal from date")
	case errors.Is(err, ErrForbiddenReport):
		response.Error(w, http.StatusForbidden, "FORBIDDEN_REPORT", "You are not allowed to view this report")
	case errors.Is(err, ErrBranchIDInvalid):
		response.Error(w, http.StatusBadRequest, "BRANCH_ID_INVALID", "Branch id is invalid")
	case errors.Is(err, ErrGroupIDInvalid):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_INVALID", "Group id is invalid")
	case errors.Is(err, ErrStudentIDInvalid):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_INVALID", "Student id is invalid")
	case errors.Is(err, ErrPaymentStatusInvalid):
		response.Error(w, http.StatusBadRequest, "PAYMENT_STATUS_INVALID", "Payment status is invalid")
	case errors.Is(err, ErrPeriodRequired):
		response.Error(w, http.StatusBadRequest, "PERIOD_REQUIRED", "Period is required")
	case errors.Is(err, ErrPeriodInvalid):
		response.Error(w, http.StatusBadRequest, "PERIOD_INVALID", "Period must be in YYYY-MM format")
	case errors.Is(err, ErrBalanceStatusInvalid):
		response.Error(w, http.StatusBadRequest, "BALANCE_STATUS_INVALID", "Balance status is invalid")
	case errors.Is(err, ErrPayrollStatusInvalid):
		response.Error(w, http.StatusBadRequest, "PAYROLL_STATUS_INVALID", "Payroll status is invalid")
	case errors.Is(err, ErrTeacherConfirmationStatusInvalid):
		response.Error(w, http.StatusBadRequest, "TEACHER_CONFIRMATION_STATUS_INVALID", "Teacher confirmation status is invalid")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
