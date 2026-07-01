package lesson

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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateLessonRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeLessonError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		writeLessonError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) UpdateActualTeacher(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonID")

	var req UpdateLessonTeacherRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.UpdateActualTeacher(r.Context(), lessonID, req)
	if err != nil {
		writeLessonError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeLessonError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrGroupIDRequired):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_REQUIRED", "Group id is required")
	case errors.Is(err, ErrGroupIDInvalid):
		response.Error(w, http.StatusBadRequest, "GROUP_ID_INVALID", "Group id is invalid")
	case errors.Is(err, ErrGroupNotFound):
		response.Error(w, http.StatusBadRequest, "GROUP_NOT_FOUND", "Group not found in organization")
	case errors.Is(err, ErrLessonIDInvalid):
		response.Error(w, http.StatusBadRequest, "LESSON_ID_INVALID", "Lesson id is invalid")
	case errors.Is(err, ErrLessonNotFound):
		response.Error(w, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson not found in organization")
	case errors.Is(err, ErrLessonDateRequired):
		response.Error(w, http.StatusBadRequest, "LESSON_DATE_REQUIRED", "Lesson date is required")
	case errors.Is(err, ErrLessonDateInvalid):
		response.Error(w, http.StatusBadRequest, "LESSON_DATE_INVALID", "Lesson date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrStartTimeRequired):
		response.Error(w, http.StatusBadRequest, "START_TIME_REQUIRED", "Start time is required")
	case errors.Is(err, ErrStartTimeInvalid):
		response.Error(w, http.StatusBadRequest, "START_TIME_INVALID", "Start time must be in HH:MM format")
	case errors.Is(err, ErrEndTimeRequired):
		response.Error(w, http.StatusBadRequest, "END_TIME_REQUIRED", "End time is required")
	case errors.Is(err, ErrEndTimeInvalid):
		response.Error(w, http.StatusBadRequest, "END_TIME_INVALID", "End time must be in HH:MM format")
	case errors.Is(err, ErrLessonTimeInvalid):
		response.Error(w, http.StatusBadRequest, "LESSON_TIME_INVALID", "End time must be after start time")
	case errors.Is(err, ErrActualTeacherIDRequired):
		response.Error(w, http.StatusBadRequest, "ACTUAL_TEACHER_ID_REQUIRED", "Actual teacher id is required")
	case errors.Is(err, ErrActualTeacherIDInvalid):
		response.Error(w, http.StatusBadRequest, "ACTUAL_TEACHER_ID_INVALID", "Actual teacher id is invalid")
	case errors.Is(err, ErrActualTeacherNotFound):
		response.Error(w, http.StatusBadRequest, "ACTUAL_TEACHER_NOT_FOUND", "Actual teacher not found in organization")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
