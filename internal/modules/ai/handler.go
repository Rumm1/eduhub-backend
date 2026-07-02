package ai

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

func (h *Handler) DashboardInsights(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetDashboardInsights(r.Context())
	if err != nil {
		writeAIError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeAIError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
