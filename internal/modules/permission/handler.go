package permission

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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ListGroups(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.ListGroups(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		return
	}

	response.Success(w, http.StatusOK, result)
}
