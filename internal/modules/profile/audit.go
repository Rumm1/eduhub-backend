package profile

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

func ProfileAuditMiddleware(auditService AuditService) func(http.Handler) http.Handler {
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

			action, entityType, entityID, ok := resolveProfileAuditAction(r.Method, r.URL.Path)
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
				Description:    profileAuditDescription(action),
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

func resolveProfileAuditAction(method string, path string) (string, string, string, bool) {
	parts := splitPath(path)
	profilesIndex := indexOf(parts, "profiles")
	if profilesIndex == -1 {
		return "", "", "", false
	}

	afterProfilesCount := len(parts) - profilesIndex - 1

	if afterProfilesCount == 1 {
		profileID := parts[profilesIndex+1]

		switch method {
		case http.MethodPatch:
			return "profile.updated", "profile", profileID, true
		case http.MethodDelete:
			return "profile.disabled", "profile", profileID, true
		default:
			return "", "", "", false
		}
	}

	if afterProfilesCount == 2 {
		profileID := parts[profilesIndex+1]
		actionPart := parts[profilesIndex+2]

		if actionPart == "set-default" && method == http.MethodPost {
			return "profile.set_default", "profile", profileID, true
		}

		if actionPart == "roles" && method == http.MethodPost {
			return "profile.role_added", "profile", profileID, true
		}

		if actionPart == "branches" && method == http.MethodPost {
			return "profile.branch_added", "profile", profileID, true
		}
	}

	if afterProfilesCount == 3 {
		profileID := parts[profilesIndex+1]
		actionPart := parts[profilesIndex+2]

		if actionPart == "roles" && method == http.MethodDelete {
			return "profile.role_removed", "profile", profileID, true
		}

		if actionPart == "branches" && method == http.MethodDelete {
			return "profile.branch_removed", "profile", profileID, true
		}
	}

	return "", "", "", false
}

func profileAuditDescription(action string) string {
	descriptions := map[string]string{
		"profile.created":        "Profile created",
		"profile.updated":        "Profile updated",
		"profile.disabled":       "Profile disabled",
		"profile.set_default":    "Default profile changed",
		"profile.role_added":     "Role added to profile",
		"profile.role_removed":   "Role removed from profile",
		"profile.branch_added":   "Branch added to profile",
		"profile.branch_removed": "Branch removed from profile",
	}

	if description, ok := descriptions[action]; ok {
		return description
	}

	return "Profile action"
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
