package middleware

import (
	"net/http"
	"strings"
)

type CORSOptions struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

func DefaultCORSOptions() CORSOptions {
	return CORSOptions{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders: []string{"Authorization", "Content-Type", "X-Request-ID", "X-Tenant-ID"},
	}
}

func CORS(options CORSOptions) Middleware {
	if len(options.AllowedOrigins) == 0 {
		options.AllowedOrigins = []string{"*"}
	}
	if len(options.AllowedMethods) == 0 {
		options.AllowedMethods = DefaultCORSOptions().AllowedMethods
	}
	if len(options.AllowedHeaders) == 0 {
		options.AllowedHeaders = DefaultCORSOptions().AllowedHeaders
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", strings.Join(options.AllowedOrigins, ", "))
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ", "))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
