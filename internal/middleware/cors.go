package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" && isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Vary", "Origin")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept, Origin, X-Requested-With, Access-Control-Request-Method, Access-Control-Request-Headers")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
			w.Header().Set("Access-Control-Max-Age", "86400")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	allowedOrigins := getAllowedOrigins()

	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}

func getAllowedOrigins() []string {
	value := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if value == "" {
		return []string{
			"http://localhost:5173",
			"http://127.0.0.1:5173",
			"http://localhost:3000",
			"http://127.0.0.1:3000",
		}
	}

	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))

	for _, part := range parts {
		origin := strings.TrimSpace(part)
		if origin != "" {
			origins = append(origins, origin)
		}
	}

	return origins
}
