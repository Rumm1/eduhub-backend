package platformuser

import (
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(chi.URLParam(r, "userID"))
	if err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_USER_ID", "Invalid user id")
		return
	}

	result, err := h.service.ResetPassword(r.Context(), userID)
	if err != nil {
		writePlatformUserError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writePlatformUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrUserNotFound):
		response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}
}
