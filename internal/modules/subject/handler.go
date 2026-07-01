package subject

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
	var req CreateSubjectRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrTenantRequired):
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
		case errors.Is(err, ErrSubjectNameRequired):
			response.Error(w, http.StatusBadRequest, "SUBJECT_NAME_REQUIRED", "Subject name is required")
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
