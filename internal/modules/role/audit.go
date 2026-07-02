package role

import (
	"context"
	"net/http"
	"strings"

	auditmodule "github.com/Rumm1/eduhub-backend/internal/modules/audit"
	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
)

type AuditService interface {
	Create(ctx context.Context, input auditmodule.CreateAuditLogInput) error
}

type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *auditResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func RoleAuditMiddleware(auditService AuditService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if auditService == nil {
				next.ServeHTTP(w, r)
				return
			}

			auditWriter := &auditResponseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(auditWriter, r)

			if auditWriter.statusCode < 200 || auditWriter.statusCode >= 400 {
				return
			}

			action, entityType, entityID, ok := resolveRoleAuditAction(r.Method, r.URL.Path)
			if !ok {
				return
			}

			currentUser, ok := usercontext.GetUser(r.Context())
			if !ok || currentUser.OrganizationID == nil {
				return
			}

			_ = auditService.Create(r.Context(), auditmodule.CreateAuditLogInput{
				OrganizationID: currentUser.OrganizationID.String(),
				UserID:         currentUser.UserID.String(),
				Action:         action,
				EntityType:     entityType,
				EntityID:       entityID,
				Description:    roleAuditDescription(action),
				IPAddress:      getClientIP(r),
				UserAgent:      r.UserAgent(),
				Metadata: map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"query":       r.URL.RawQuery,
					"status_code": auditWriter.statusCode,
					"roles":       currentUser.Roles,
				},
			})
		})
	}
}

func resolveRoleAuditAction(method string, path string) (string, string, string, bool) {
	parts := splitPath(path)
	rolesIndex := indexOf(parts, "roles")
	if rolesIndex == -1 {
		return "", "", "", false
	}

	afterRolesCount := len(parts) - rolesIndex - 1

	if afterRolesCount == 0 && method == http.MethodPost {
		return "role.created", "role", "", true
	}

	if afterRolesCount == 1 {
		roleID := parts[rolesIndex+1]

		switch method {
		case http.MethodPatch:
			return "role.updated", "role", roleID, true
		case http.MethodDelete:
			return "role.deleted", "role", roleID, true
		default:
			return "", "", "", false
		}
	}

	if afterRolesCount == 2 && parts[rolesIndex+2] == "permissions" && method == http.MethodPost {
		roleID := parts[rolesIndex+1]
		return "role.permission_added", "role", roleID, true
	}

	if afterRolesCount == 3 && parts[rolesIndex+2] == "permissions" && method == http.MethodDelete {
		roleID := parts[rolesIndex+1]
		return "role.permission_removed", "role", roleID, true
	}

	return "", "", "", false
}

func roleAuditDescription(action string) string {
	descriptions := map[string]string{
		"role.created":            "Role created",
		"role.updated":            "Role updated",
		"role.deleted":            "Role deleted",
		"role.permission_added":   "Permission added to role",
		"role.permission_removed": "Permission removed from role",
	}

	if description, ok := descriptions[action]; ok {
		return description
	}

	return "Role action"
}

func splitPath(path string) []string {
	rawParts := strings.Split(path, "/")
	parts := make([]string, 0, len(rawParts))

	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parts = append(parts, part)
	}

	return parts
}

func indexOf(items []string, target string) int {
	for index, item := range items {
		if item == target {
			return index
		}
	}

	return -1
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

	return r.RemoteAddr
}
