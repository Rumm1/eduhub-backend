package audit

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
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

	format := strings.TrimSpace(r.URL.Query().Get("format"))
	if format == "" || format == "json" {
		response.Success(w, http.StatusOK, result)
		return
	}

	if format != "xlsx" {
		response.Error(w, http.StatusBadRequest, "FORMAT_UNSUPPORTED", "Only xlsx export is supported")
		return
	}

	lang := strings.TrimSpace(r.URL.Query().Get("lang"))
	if lang == "" {
		lang = "ru"
	}

	fileBytes, err := BuildAuditLogsXLSX(result, lang)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "EXPORT_FAILED", "Failed to build audit logs export")
		return
	}

	h.writeAuditExportLog(r, lang, result.Total)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="audit_logs.xlsx"`)
	w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(fileBytes)
}

func (h *Handler) writeAuditExportLog(r *http.Request, lang string, total int) {
	currentUser, ok := usercontext.GetUser(r.Context())
	if !ok {
		return
	}

	organizationID := ""
	if currentUser.OrganizationID != nil {
		organizationID = currentUser.OrganizationID.String()
	}

	_ = h.service.Create(r.Context(), CreateAuditLogInput{
		OrganizationID: organizationID,
		UserID:         currentUser.UserID.String(),
		Action:         "audit_logs.exported",
		EntityType:     "audit_logs",
		EntityID:       "",
		Description:    "Audit logs exported",
		Metadata: map[string]interface{}{
			"format":      "xlsx",
			"lang":        lang,
			"query":       r.URL.RawQuery,
			"path":        r.URL.Path,
			"method":      r.Method,
			"total":       total,
			"roles":       currentUser.Roles,
			"status_code": http.StatusOK,
		},
		IPAddress: getClientIP(r),
		UserAgent: r.UserAgent(),
	})
}

func getClientIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(r.Header.Get("X-Forwarded-For"))
	if forwardedFor != "" {
		parts := strings.Split(forwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	realIP := strings.TrimSpace(r.Header.Get("X-Real-IP"))
	if realIP != "" {
		return realIP
	}

	return strings.TrimSpace(r.RemoteAddr)
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
		fmt.Println("AUDIT ERROR:", err)
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
	}
}
