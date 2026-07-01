package middleware

import (
	"net/http"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func RequirePermissions(required ...string) Middleware {
	requiredSet := make(map[string]struct{}, len(required))
	for _, permission := range required {
		requiredSet[permission] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := usercontext.UserFromContext(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "authentication required")
				return
			}
			if hasPermissions(user.Permissions, requiredSet) {
				next.ServeHTTP(w, r)
				return
			}
			response.Error(w, http.StatusForbidden, "permission denied")
		})
	}
}

func hasPermissions(userPermissions []string, required map[string]struct{}) bool {
	if len(required) == 0 {
		return true
	}

	granted := make(map[string]struct{}, len(userPermissions))
	for _, permission := range userPermissions {
		granted[permission] = struct{}{}
	}
	for permission := range required {
		if _, ok := granted[permission]; !ok {
			return false
		}
	}
	return true
}
