package middleware

import (
	"net/http"
	"strings"

	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func Auth(jwtManager *platformjwt.Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is required")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Error(w, http.StatusUnauthorized, "INVALID_AUTH_HEADER", "Authorization header must be Bearer token")
				return
			}

			claims, err := jwtManager.ParseAccessToken(parts[1])
			if err != nil {
				response.Error(w, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
				return
			}

			user := usercontext.UserContext{
				UserID:         claims.UserID,
				ProfileID:      claims.ProfileID,
				OrganizationID: claims.OrganizationID,
				Roles:          claims.Roles,
				Permissions:    claims.Permissions,
				BranchIDs:      claims.BranchIDs,
			}

			ctx := usercontext.WithUser(r.Context(), user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
