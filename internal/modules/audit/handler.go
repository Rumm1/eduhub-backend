package audit

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil && r.URL.Query().Get("limit") != "" {
		response.Error(w, http.StatusBadRequest, "LIMIT_INVALID", "Limit is invalid")
		return
	}

	offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
	if err != nil && r.URL.Query().Get("offset") != "" {
		response.Error(w, http.StatusBadRequest, "OFFSET_INVALID", "Offset is invalid")
		return
	}

	result, err := h.service.List(
		r.Context(),
		r.URL.Query().Get("user_id"),
		r.URL.Query().Get("action"),
		r.URL.Query().Get("entity_type"),
		r.URL.Query().Get("entity_id"),
		r.URL.Query().Get("from_date"),
		r.URL.Query().Get("to_date"),
		limit,
		offset,
	)
	if err != nil {
		writeAuditError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeAuditError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrUserIDInvalid):
		response.Error(w, http.StatusBadRequest, "USER_ID_INVALID", "User id is invalid")
	case errors.Is(err, ErrEntityIDInvalid):
		response.Error(w, http.StatusBadRequest, "ENTITY_ID_INVALID", "Entity id is invalid")
	case errors.Is(err, ErrFromDateInvalid):
		response.Error(w, http.StatusBadRequest, "FROM_DATE_INVALID", "From date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrToDateInvalid):
		response.Error(w, http.StatusBadRequest, "TO_DATE_INVALID", "To date must be in YYYY-MM-DD format")
	case errors.Is(err, ErrLimitInvalid):
		response.Error(w, http.StatusBadRequest, "LIMIT_INVALID", "Limit must be between 1 and 200")
	case errors.Is(err, ErrOffsetInvalid):
		response.Error(w, http.StatusBadRequest, "OFFSET_INVALID", "Offset is invalid")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
