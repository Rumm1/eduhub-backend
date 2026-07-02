package report

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

func ReportAuditMiddleware(auditService AuditService) func(http.Handler) http.Handler {
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

			if r.Method != http.MethodGet {
				return
			}

			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			if statusCode < 200 || statusCode >= 400 {
				return
			}

			format := strings.TrimSpace(r.URL.Query().Get("format"))
			if format == "" || format == "json" {
				return
			}

			reportName := resolveReportName(r.URL.Path)
			if reportName == "" {
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

			lang := strings.TrimSpace(r.URL.Query().Get("lang"))
			if lang == "" {
				lang = "ru"
			}

			_ = auditService.Create(r.Context(), auditmodule.CreateAuditLogInput{
				OrganizationID: organizationID,
				UserID:         currentUser.UserID.String(),
				Action:         "report.exported",
				EntityType:     "report",
				EntityID:       "",
				Description:    "Report exported",
				Metadata: map[string]interface{}{
					"report":      reportName,
					"format":      format,
					"lang":        lang,
					"method":      r.Method,
					"path":        r.URL.Path,
					"query":       r.URL.RawQuery,
					"status_code": statusCode,
					"roles":       currentUser.Roles,
				},
				IPAddress: getClientIP(r),
				UserAgent: r.UserAgent(),
			})
		})
	}
}

func resolveReportName(path string) string {
	parts := splitPath(path)

	reportsIndex := indexOf(parts, "reports")
	if reportsIndex == -1 {
		return ""
	}

	if len(parts) <= reportsIndex+1 {
		return ""
	}

	return parts[reportsIndex+1]
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
