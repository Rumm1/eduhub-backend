package middleware

import (
	"log"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Printf("panic recovered: %v", recovered)
					response.Error(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
