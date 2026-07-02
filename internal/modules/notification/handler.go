package notification

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
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
		writeNotificationError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Types(w http.ResponseWriter, r *http.Request) {
	response.Success(w, http.StatusOK, h.service.ListTypes())
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateNotificationRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeNotificationError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	notificationID := chi.URLParam(r, "notificationID")

	if err := h.service.MarkRead(r.Context(), notificationID); err != nil {
		writeNotificationError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Notification marked as read")
}

func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.MarkAllRead(r.Context())
	if err != nil {
		writeNotificationError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	notificationID := chi.URLParam(r, "notificationID")

	if err := h.service.Delete(r.Context(), notificationID); err != nil {
		writeNotificationError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Notification deleted")
}

func writeNotificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrNotificationIDInvalid):
		response.Error(w, http.StatusBadRequest, "NOTIFICATION_ID_INVALID", "Notification id is invalid")
	case errors.Is(err, ErrTitleRequired):
		response.Error(w, http.StatusBadRequest, "TITLE_REQUIRED", "Title is required")
	case errors.Is(err, ErrTypeInvalid):
		response.Error(w, http.StatusBadRequest, "NOTIFICATION_TYPE_INVALID", "Notification type is invalid")
	case errors.Is(err, ErrTargetUserInvalid):
		response.Error(w, http.StatusBadRequest, "TARGET_USER_INVALID", "Target user id is invalid")
	case errors.Is(err, ErrTargetUserNotFound):
		response.Error(w, http.StatusBadRequest, "TARGET_USER_NOT_FOUND", "Target user not found in organization")
	case errors.Is(err, ErrNotificationNotFound):
		response.Error(w, http.StatusNotFound, "NOTIFICATION_NOT_FOUND", "Notification not found")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
