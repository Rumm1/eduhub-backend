package role

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
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")

	result, err := h.service.GetByID(r.Context(), roleID)
	if err != nil {
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), req)
	if err != nil {
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")

	var req UpdateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Update(r.Context(), roleID, req)
	if err != nil {
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")

	if err := h.service.Delete(r.Context(), roleID); err != nil {
		writeRoleError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Role deleted")
}

func (h *Handler) AddPermission(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")

	var req AddPermissionRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.AddPermission(r.Context(), roleID, req)
	if err != nil {
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	roleID := chi.URLParam(r, "roleID")
	permissionCode := chi.URLParam(r, "permissionCode")

	result, err := h.service.RemovePermission(r.Context(), roleID, permissionCode)
	if err != nil {
		writeRoleError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeRoleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrRoleIDInvalid):
		response.Error(w, http.StatusBadRequest, "ROLE_ID_INVALID", "Role id is invalid")
	case errors.Is(err, ErrRoleNotFound):
		response.Error(w, http.StatusNotFound, "ROLE_NOT_FOUND", "Role not found")
	case errors.Is(err, ErrRoleNameRequired):
		response.Error(w, http.StatusBadRequest, "ROLE_NAME_REQUIRED", "Role name is required")
	case errors.Is(err, ErrRoleCodeRequired):
		response.Error(w, http.StatusBadRequest, "ROLE_CODE_REQUIRED", "Role code is required")
	case errors.Is(err, ErrPermissionRequired):
		response.Error(w, http.StatusBadRequest, "PERMISSION_REQUIRED", "Permission is required")
	case errors.Is(err, ErrPermissionNotFound):
		response.Error(w, http.StatusBadRequest, "PERMISSION_NOT_FOUND", "Permission not found")
	case errors.Is(err, ErrSystemRoleReadonly):
		response.Error(w, http.StatusBadRequest, "SYSTEM_ROLE_READONLY", "System role is readonly")
	case errors.Is(err, ErrRoleInUse):
		response.Error(w, http.StatusBadRequest, "ROLE_IN_USE", "Role is in use")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
