package branding

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

func (h *Handler) Current(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetCurrentBranding(r.Context())
	if err != nil {
		writeBrandingError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) UpdateMyAvatar(w http.ResponseWriter, r *http.Request) {
	var request UpdateAvatarRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.service.UpdateMyAvatar(r.Context(), request)
	if err != nil {
		writeBrandingError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ClearMyAvatar(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ClearMyAvatar(r.Context())
	if err != nil {
		writeBrandingError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) UpdateOrganizationLogo(w http.ResponseWriter, r *http.Request) {
	var request UpdateLogoRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid request body")
		return
	}

	result, err := h.service.UpdateOrganizationLogo(r.Context(), request)
	if err != nil {
		writeBrandingError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ClearOrganizationLogo(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ClearOrganizationLogo(r.Context())
	if err != nil {
		writeBrandingError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeBrandingError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrAvatarRequired):
		response.Error(w, http.StatusBadRequest, "AVATAR_REQUIRED", "Avatar path is required")
	case errors.Is(err, ErrLogoRequired):
		response.Error(w, http.StatusBadRequest, "LOGO_REQUIRED", "Logo path is required")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
