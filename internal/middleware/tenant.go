package middleware

import (
	"net/http"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func RequireTenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := usercontext.GetUser(r.Context())
		if !ok {
			response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found")
			return
		}

		if user.OrganizationID == nil {
			response.Error(w, http.StatusForbidden, "TENANT_REQUIRED", "Organization is required")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RequireBranchAccess(branchID string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := usercontext.GetUser(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found")
				return
			}

			for _, id := range user.BranchIDs {
				if id.String() == branchID {
					next.ServeHTTP(w, r)
					return
				}
			}

			response.Error(w, http.StatusForbidden, "BRANCH_ACCESS_DENIED", "Branch access denied")
		})
	}
}
