package user

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
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTenantRequired):
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
		case errors.Is(err, ErrEmailRequired):
			response.Error(w, http.StatusBadRequest, "EMAIL_REQUIRED", "Email is required")
		case errors.Is(err, ErrPasswordRequired):
			response.Error(w, http.StatusBadRequest, "PASSWORD_REQUIRED", "Password is required")
		case errors.Is(err, ErrFullNameRequired):
			response.Error(w, http.StatusBadRequest, "FULL_NAME_REQUIRED", "Full name is required")
		case errors.Is(err, ErrRoleRequired):
			response.Error(w, http.StatusBadRequest, "ROLE_REQUIRED", "Role is required")
		case errors.Is(err, ErrRoleInvalid):
			response.Error(w, http.StatusBadRequest, "ROLE_INVALID", "Role is invalid")
		case errors.Is(err, ErrBranchIDInvalid):
			response.Error(w, http.StatusBadRequest, "BRANCH_ID_INVALID", "Branch id is invalid")
		case errors.Is(err, ErrBranchNotFound):
			response.Error(w, http.StatusBadRequest, "BRANCH_NOT_FOUND", "Branch not found in organization")
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
