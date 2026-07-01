package middleware

import (
	"net/http"
	"strings"

	usercontext "github.com/Rumm1/eduhub-backend/internal/shared/context"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func Auth() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
			if token == "" {
				response.Error(w, http.StatusUnauthorized, "missing bearer token")
				return
			}

			user := usercontext.User{ID: token}
			next.ServeHTTP(w, r.WithContext(usercontext.WithUser(r.Context(), user)))
		})
	}
}
