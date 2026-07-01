package report

import (
	"errors"
	"net/http"

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
	format := r.URL.Query().Get("format")

	result, err := h.service.GetTeacherSchedule(r.Context(), teacherID, fromDate, toDate)
	if err != nil {
		writeReportError(w, err)
		return
	}

	if format == "xlsx" {
		fileBytes, filename, err := BuildTeacherScheduleXLSX(result)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "EXPORT_ERROR", "Failed to export report")
			return
		}

		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
		w.WriteHeader(http.StatusOK)

		_, _ = w.Write(fileBytes)
		return
	}

	response.Success(w, http.StatusOK, result)
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
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
