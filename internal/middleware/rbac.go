package middleware

import (
	"net/http"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func RequirePermission(permission string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := usercontext.GetUser(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found")
				return
			}

			if !usercontext.HasPermission(r.Context(), permission) {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "Permission denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, ok := usercontext.GetUser(r.Context())
			if !ok {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "User context not found")
				return
			}

			if !usercontext.HasRole(r.Context(), role) {
				response.Error(w, http.StatusForbidden, "FORBIDDEN", "Role denied")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
