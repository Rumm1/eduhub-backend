package profile

import (
	"encoding/json"
	"errors"
	"net/http"

	auditmodule "github.com/Rumm1/eduhub-backend/internal/modules/audit"
	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service      *Service
	auditService AuditService
}

func NewHandler(service *Service, auditServices ...AuditService) *Handler {
	var auditService AuditService
	if len(auditServices) > 0 {
		auditService = auditServices[0]
	}

	return &Handler{
		service:      service,
		auditService: auditService,
	}
}

func (h *Handler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	result, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	result, err := h.service.GetByID(r.Context(), profileID)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")

	var req CreateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Create(r.Context(), userID, req)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	h.writeProfileCreatedAudit(r, result.ID, userID)

	response.Success(w, http.StatusCreated, result)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	var req UpdateProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.Update(r.Context(), profileID, req)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Disable(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	if err := h.service.Disable(r.Context(), profileID); err != nil {
		writeProfileError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "Profile disabled")
}

func (h *Handler) SetDefault(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	result, err := h.service.SetDefault(r.Context(), profileID)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) AddRole(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	var req AddRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.AddRole(r.Context(), profileID, req)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) RemoveRole(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")
	roleCode := chi.URLParam(r, "roleCode")

	result, err := h.service.RemoveRole(r.Context(), profileID, roleCode)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) AddBranch(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")

	var req AddBranchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	result, err := h.service.AddBranch(r.Context(), profileID, req)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) RemoveBranch(w http.ResponseWriter, r *http.Request) {
	profileID := chi.URLParam(r, "profileID")
	branchID := chi.URLParam(r, "branchID")

	result, err := h.service.RemoveBranch(r.Context(), profileID, branchID)
	if err != nil {
		writeProfileError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) writeProfileCreatedAudit(r *http.Request, profileID string, targetUserID string) {
	if h.auditService == nil {
		return
	}

	currentUser, ok := usercontext.GetUser(r.Context())
	if !ok || currentUser.OrganizationID == nil {
		return
	}

	_ = h.auditService.Create(r.Context(), auditmodule.CreateAuditLogInput{
		OrganizationID: currentUser.OrganizationID.String(),
		UserID:         currentUser.UserID.String(),
		Action:         "profile.created",
		EntityType:     "profile",
		EntityID:       profileID,
		Description:    "Profile created",
		IPAddress:      getClientIP(r),
		UserAgent:      r.UserAgent(),
		Metadata: map[string]interface{}{
			"method":         r.Method,
			"path":           r.URL.Path,
			"query":          r.URL.RawQuery,
			"status_code":    http.StatusCreated,
			"roles":          currentUser.Roles,
			"target_user_id": targetUserID,
		},
	})
}

func writeProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrTenantRequired):
		response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Tenant organization is required")
	case errors.Is(err, ErrUserIDInvalid):
		response.Error(w, http.StatusBadRequest, "USER_ID_INVALID", "User id is invalid")
	case errors.Is(err, ErrProfileIDInvalid):
		response.Error(w, http.StatusBadRequest, "PROFILE_ID_INVALID", "Profile id is invalid")
	case errors.Is(err, ErrBranchIDInvalid):
		response.Error(w, http.StatusBadRequest, "BRANCH_ID_INVALID", "Branch id is invalid")
	case errors.Is(err, ErrRoleRequired):
		response.Error(w, http.StatusBadRequest, "ROLE_REQUIRED", "At least one role is required")
	case errors.Is(err, ErrRoleInvalid):
		response.Error(w, http.StatusBadRequest, "ROLE_INVALID", "Role is invalid")
	case errors.Is(err, ErrUserNotFound):
		response.Error(w, http.StatusNotFound, "USER_NOT_FOUND", "User not found")
	case errors.Is(err, ErrProfileNotFound):
		response.Error(w, http.StatusNotFound, "PROFILE_NOT_FOUND", "Profile not found")
	case errors.Is(err, ErrBranchNotFound):
		response.Error(w, http.StatusBadRequest, "BRANCH_NOT_FOUND", "Branch not found in organization")
	case errors.Is(err, ErrProfileInactive):
		response.Error(w, http.StatusBadRequest, "PROFILE_INACTIVE", "Profile is inactive")
	case errors.Is(err, ErrDefaultProfileRequired):
		response.Error(w, http.StatusBadRequest, "DEFAULT_PROFILE_REQUIRED", "Default profile is required")
	case errors.Is(err, ErrCannotDisableDefaultProfile):
		response.Error(w, http.StatusBadRequest, "CANNOT_DISABLE_DEFAULT_PROFILE", "Cannot disable default profile")
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
