package attendance

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) MarkLessonAttendance(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonID")

	var req MarkLessonAttendanceRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.MarkLessonAttendance(r.Context(), lessonID, req)
	if err != nil {
		writeAttendanceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ListByLessonID(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonID")

	result, err := h.service.ListByLessonID(r.Context(), lessonID)
	if err != nil {
		writeAttendanceError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeAttendanceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrLessonIDInvalid):
		response.Error(w, http.StatusBadRequest, "LESSON_ID_INVALID", "Lesson id is invalid")
	case errors.Is(err, ErrLessonNotFound):
		response.Error(w, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson not found in organization")
	case errors.Is(err, ErrAttendanceItemsRequired):
		response.Error(w, http.StatusBadRequest, "ATTENDANCE_ITEMS_REQUIRED", "Attendance items are required")
	case errors.Is(err, ErrStudentIDRequired):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_REQUIRED", "Student id is required")
	case errors.Is(err, ErrStudentIDInvalid):
		response.Error(w, http.StatusBadRequest, "STUDENT_ID_INVALID", "Student id is invalid")
	case errors.Is(err, ErrAttendanceStatusRequired):
		response.Error(w, http.StatusBadRequest, "ATTENDANCE_STATUS_REQUIRED", "Attendance status is required")
	case errors.Is(err, ErrAttendanceStatusInvalid):
		response.Error(w, http.StatusBadRequest, "ATTENDANCE_STATUS_INVALID", "Attendance status is invalid")
	case errors.Is(err, ErrStudentNotInLessonGroup):
		response.Error(w, http.StatusBadRequest, "STUDENT_NOT_IN_LESSON_GROUP", "Student is not in lesson group")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
