package platformdashboard

import (
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.GetDashboard(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
		return
	}

	response.Success(w, http.StatusOK, result)
}
