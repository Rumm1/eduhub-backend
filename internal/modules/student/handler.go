package student

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
	var req CreateStudentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTenantRequired):
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
		case errors.Is(err, ErrBranchIDRequired):
			response.Error(w, http.StatusBadRequest, "BRANCH_ID_REQUIRED", "Branch id is required")
		case errors.Is(err, ErrBranchIDInvalid):
			response.Error(w, http.StatusBadRequest, "BRANCH_ID_INVALID", "Branch id is invalid")
		case errors.Is(err, ErrBranchNotFound):
			response.Error(w, http.StatusBadRequest, "BRANCH_NOT_FOUND", "Branch not found in organization")
		case errors.Is(err, ErrStudentNameRequired):
			response.Error(w, http.StatusBadRequest, "STUDENT_NAME_REQUIRED", "Student name is required")
		case errors.Is(err, ErrBirthDateInvalid):
			response.Error(w, http.StatusBadRequest, "BIRTH_DATE_INVALID", "Birth date must be in YYYY-MM-DD format")
		case errors.Is(err, ErrParentNameRequired):
			response.Error(w, http.StatusBadRequest, "PARENT_NAME_REQUIRED", "Parent name is required")
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
