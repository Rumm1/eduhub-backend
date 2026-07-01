package middleware

import (
	"context"
	"net/http"
)

type tenantKey struct{}

const TenantHeader = "X-Tenant-ID"

func Tenant() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tenantID := r.Header.Get(TenantHeader)
			if tenantID == "" {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), tenantKey{}, tenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func TenantFromContext(ctx context.Context) string {
	tenantID, _ := ctx.Value(tenantKey{}).(string)
	return tenantID
}
