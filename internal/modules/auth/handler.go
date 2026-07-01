package auth

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

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid email or password")
		case errors.Is(err, ErrUserInactive):
			response.Error(w, http.StatusForbidden, "USER_INACTIVE", "User is inactive")
		default:
			response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
		}
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Me(r.Context())
	if err != nil {
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized")
		return
	}

	response.Success(w, http.StatusOK, result)
}
