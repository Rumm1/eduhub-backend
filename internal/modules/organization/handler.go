package organization

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
	var req CreateOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.CreateOrganization(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrganizationNameRequired):
			response.Error(w, http.StatusBadRequest, "ORGANIZATION_NAME_REQUIRED", "At least one organization name is required")
		case errors.Is(err, ErrOrganizationLanguageInvalid):
			response.Error(w, http.StatusBadRequest, "ORGANIZATION_LANGUAGE_INVALID", "Default language must be ru, kk, or en")
		case errors.Is(err, ErrAdminEmailRequired):
			response.Error(w, http.StatusBadRequest, "ADMIN_EMAIL_REQUIRED", "Admin email is required")
		case errors.Is(err, ErrAdminPasswordRequired):
			response.Error(w, http.StatusBadRequest, "ADMIN_PASSWORD_REQUIRED", "Admin password is required")
		case errors.Is(err, ErrAdminFullNameRequired):
			response.Error(w, http.StatusBadRequest, "ADMIN_FULL_NAME_REQUIRED", "Admin full name is required")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	response.Success(w, http.StatusCreated, result)
}
