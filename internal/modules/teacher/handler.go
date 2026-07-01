package teacher

import (
	"encoding/json"
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

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTeacherRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTenantRequired):
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
		case errors.Is(err, ErrUserIDRequired):
			response.Error(w, http.StatusBadRequest, "USER_ID_REQUIRED", "User id is required")
		case errors.Is(err, ErrUserIDInvalid):
			response.Error(w, http.StatusBadRequest, "USER_ID_INVALID", "User id is invalid")
		case errors.Is(err, ErrUserNotFound):
			response.Error(w, http.StatusBadRequest, "USER_NOT_FOUND", "User not found in organization")
		case errors.Is(err, ErrUserIsNotTeacher):
			response.Error(w, http.StatusBadRequest, "USER_IS_NOT_TEACHER", "User does not have TEACHER role")
		case errors.Is(err, ErrSubjectIDInvalid):
			response.Error(w, http.StatusBadRequest, "SUBJECT_ID_INVALID", "Subject id is invalid")
		case errors.Is(err, ErrSubjectNotFound):
			response.Error(w, http.StatusBadRequest, "SUBJECT_NOT_FOUND", "Subject not found in organization")
		case errors.Is(err, ErrExperienceInvalid):
			response.Error(w, http.StatusBadRequest, "EXPERIENCE_INVALID", "Experience years must be greater than or equal to zero")
		case errors.Is(err, ErrHourlyRateInvalid):
			response.Error(w, http.StatusBadRequest, "HOURLY_RATE_INVALID", "Hourly rate must be greater than or equal to zero")
		case errors.Is(err, ErrFixedSalaryInvalid):
			response.Error(w, http.StatusBadRequest, "FIXED_SALARY_INVALID", "Fixed salary must be greater than or equal to zero")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		switch {
		case errors.Is(err, ErrTenantRequired):
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	response.Success(w, http.StatusOK, result)
}
