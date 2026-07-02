package lesson

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

func (w *auditResponseWriter) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}

	return w.ResponseWriter.Write(data)
}

func LessonAuditMiddleware(auditService AuditService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recorder := &auditResponseWriter{
				ResponseWriter: w,
				statusCode:     0,
			}

			next.ServeHTTP(recorder, r)

			if auditService == nil {
				return
			}

			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			if statusCode < 200 || statusCode >= 400 {
				return
			}

			action, entityID, description := resolveLessonAuditAction(r.Method, r.URL.Path)
			if action == "" {
				return
			}

			currentUser, ok := usercontext.GetUser(r.Context())
			if !ok {
				return
			}

			organizationID := ""
			if currentUser.OrganizationID != nil {
				organizationID = currentUser.OrganizationID.String()
			}

			_ = auditService.Create(r.Context(), auditmodule.CreateAuditLogInput{
				OrganizationID: organizationID,
				UserID:         currentUser.UserID.String(),
				Action:         action,
				EntityType:     "lesson",
				EntityID:       entityID,
				Description:    description,
				Metadata: map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"status_code": statusCode,
					"roles":       currentUser.Roles,
				},
				IPAddress: getClientIP(r),
				UserAgent: r.UserAgent(),
			})
		})
	}
}

func resolveLessonAuditAction(method string, path string) (string, string, string) {
	parts := splitPath(path)

	lessonsIndex := indexOf(parts, "lessons")
	if lessonsIndex == -1 {
		return "", "", ""
	}

	if method == http.MethodPatch &&
		len(parts) >= lessonsIndex+3 &&
		parts[lessonsIndex+2] == "teacher" {
		lessonID := parts[lessonsIndex+1]
		return "lesson.teacher_replaced", lessonID, "Lesson teacher replaced"
	}

	return "", "", ""
}

func splitPath(path string) []string {
	rawParts := strings.Split(strings.Trim(path, "/"), "/")
	parts := make([]string, 0, len(rawParts))

	for _, part := range rawParts {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}

	return parts
}

func indexOf(values []string, target string) int {
	for index, value := range values {
		if value == target {
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

	return strings.TrimSpace(r.RemoteAddr)
}
