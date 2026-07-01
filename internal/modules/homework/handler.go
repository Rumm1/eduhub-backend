package homework

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
	var req CreateHomeworkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeHomeworkError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		writeHomeworkError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ListByLessonID(w http.ResponseWriter, r *http.Request) {
	lessonID := chi.URLParam(r, "lessonID")

	result, err := h.service.ListByLessonID(r.Context(), lessonID)
	if err != nil {
		writeHomeworkError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeHomeworkError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrLessonIDRequired):
		response.Error(w, http.StatusBadRequest, "LESSON_ID_REQUIRED", "Lesson id is required")
	case errors.Is(err, ErrLessonIDInvalid):
		response.Error(w, http.StatusBadRequest, "LESSON_ID_INVALID", "Lesson id is invalid")
	case errors.Is(err, ErrLessonNotFound):
		response.Error(w, http.StatusNotFound, "LESSON_NOT_FOUND", "Lesson not found in organization")
	case errors.Is(err, ErrTeacherNotFound):
		response.Error(w, http.StatusBadRequest, "TEACHER_NOT_FOUND", "Teacher not found for lesson")
	case errors.Is(err, ErrTitleRequired):
		response.Error(w, http.StatusBadRequest, "HOMEWORK_TITLE_REQUIRED", "Homework title is required")
	case errors.Is(err, ErrDueDateInvalid):
		response.Error(w, http.StatusBadRequest, "DUE_DATE_INVALID", "Due date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrHomeworkDisabled):
		response.Error(w, http.StatusBadRequest, "HOMEWORK_DISABLED", "Homework is disabled for this subject or group")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
