package payroll

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

func PayrollAuditMiddleware(auditService AuditService) func(http.Handler) http.Handler {
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

			if r.Method != http.MethodPost {
				return
			}

			statusCode := recorder.statusCode
			if statusCode == 0 {
				statusCode = http.StatusOK
			}

			if statusCode < 200 || statusCode >= 400 {
				return
			}

			action, entityID, description := resolvePayrollAuditAction(r.URL.Path)
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
				EntityType:     "payroll",
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

func resolvePayrollAuditAction(path string) (string, string, string) {
	parts := splitPath(path)

	payrollIndex := indexOf(parts, "payroll")
	if payrollIndex == -1 {
		return "", "", ""
	}

	if len(parts) <= payrollIndex+1 {
		return "", "", ""
	}

	resource := parts[payrollIndex+1]

	if resource == "periods" {
		if len(parts) == payrollIndex+2 {
			return "payroll.period_created", "", "Payroll period created"
		}

		periodID := parts[payrollIndex+2]

		if len(parts) >= payrollIndex+4 && parts[payrollIndex+3] == "generate" {
			return "payroll.generated", periodID, "Payroll generated for period"
		}
	}

	if resource == "entries" {
		if len(parts) < payrollIndex+4 {
			return "", "", ""
		}

		entryID := parts[payrollIndex+2]
		actionPart := parts[payrollIndex+3]

		switch actionPart {
		case "adjustments":
			return "payroll.adjustment_added", entryID, "Payroll adjustment added"
		case "send-to-teacher":
			return "payroll.sent_to_teacher", entryID, "Payroll entry sent to teacher"
		case "confirm":
			return "payroll.confirmed_by_teacher", entryID, "Payroll entry confirmed by teacher"
		case "dispute":
			return "payroll.disputed_by_teacher", entryID, "Payroll entry disputed by teacher"
		case "approve":
			return "payroll.approved_by_finance", entryID, "Payroll entry approved by finance"
		case "mark-paid":
			return "payroll.marked_paid", entryID, "Payroll entry marked as paid"
		}
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
